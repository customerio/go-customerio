package customerio

import "context"

// DesignStudioFolder organizes Design Studio emails/components into a tree.
type DesignStudioFolder struct {
	ID             string `json:"id"`
	Name           string `json:"name,omitempty"`
	ParentFolderID string `json:"parent_folder_id,omitempty"`
	Created        int64  `json:"created,omitempty"`
	Updated        int64  `json:"updated,omitempty"`
}

// ListDesignStudioFoldersResponse is the decoded shape of GET /v1/design_studio/folders.
type ListDesignStudioFoldersResponse struct {
	Folders []DesignStudioFolder `json:"folders"`
}

// ListDesignStudioFolders returns Design Studio folders, filtered/sorted/
// paginated by opts.
// See https://docs.customer.io/api/app/#operation/listDesignStudioFolders
func (c *APIClient) ListDesignStudioFolders(ctx context.Context, opts DesignStudioListOptions) (*ListDesignStudioFoldersResponse, error) {
	requestPath := "/v1/design_studio/folders" + opts.apply(newQuery()).String()

	var resp ListDesignStudioFoldersResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DesignStudioFolderResponse wraps a single DesignStudioFolder.
type DesignStudioFolderResponse struct {
	Folder DesignStudioFolder `json:"folder"`
}

// DesignStudioFolderInput creates a Design Studio folder.
type DesignStudioFolderInput struct {
	Name           string `json:"name"`
	ParentFolderID string `json:"parent_folder_id,omitempty"`
}

// CreateDesignStudioFolder creates a Design Studio folder.
// See https://docs.customer.io/api/app/#operation/createDesignStudioFolder
func (c *APIClient) CreateDesignStudioFolder(ctx context.Context, input DesignStudioFolderInput) (*DesignStudioFolderResponse, error) {
	if input.Name == "" {
		return nil, ParamError{Param: "name"}
	}

	var resp DesignStudioFolderResponse
	if err := c.doJSON(ctx, "POST", "/v1/design_studio/folders", input, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetDesignStudioFolder returns one Design Studio folder by id.
// See https://docs.customer.io/api/app/#operation/getDesignStudioFolder
func (c *APIClient) GetDesignStudioFolder(ctx context.Context, folderID string) (*DesignStudioFolderResponse, error) {
	if folderID == "" {
		return nil, ParamError{Param: "folderID"}
	}

	var resp DesignStudioFolderResponse
	if err := c.doJSON(ctx, "GET", formatPath("/v1/design_studio/folders/%s", folderID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateDesignStudioFolder updates a Design Studio folder. data is
// free-form (name/parent_folder_id, whichever you want to change). The API
// returns no body on success; the response struct is populated only if the
// server does return one.
// See https://docs.customer.io/api/app/#operation/updateDesignStudioFolder
func (c *APIClient) UpdateDesignStudioFolder(ctx context.Context, folderID string, data map[string]any) (*DesignStudioFolderResponse, error) {
	if folderID == "" {
		return nil, ParamError{Param: "folderID"}
	}

	var resp DesignStudioFolderResponse
	if err := c.doJSON(ctx, "PUT", formatPath("/v1/design_studio/folders/%s", folderID), data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteDesignStudioFolder deletes a Design Studio folder.
// See https://docs.customer.io/api/app/#operation/deleteDesignStudioFolder
func (c *APIClient) DeleteDesignStudioFolder(ctx context.Context, folderID string) error {
	if folderID == "" {
		return ParamError{Param: "folderID"}
	}

	return c.doJSON(ctx, "DELETE", formatPath("/v1/design_studio/folders/%s", folderID), nil, nil, 200, 204)
}

// DesignStudioEmailHeader is one custom header on a Design Studio email's envelope.
type DesignStudioEmailHeader struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// DesignStudioEmailContent is a Design Studio email's rendered content.
type DesignStudioEmailContent struct {
	Subject       string `json:"subject,omitempty"`
	PreheaderText string `json:"preheader_text,omitempty"`
	HTML          string `json:"html,omitempty"`
	AMP           string `json:"amp,omitempty"`
	Text          string `json:"text,omitempty"`
}

// DesignStudioEmailEnvelope is a Design Studio email's sender/recipient
// configuration. CC is accepted by the App API's request body (per the
// Node SDK) but isn't part of the documented response schema — kept here
// for round-tripping requests, may always be empty on responses.
type DesignStudioEmailEnvelope struct {
	FromID    *int                      `json:"from_id,omitempty"`
	From      string                    `json:"from,omitempty"`
	ReplyToID *int                      `json:"reply_to_id,omitempty"`
	ReplyTo   string                    `json:"reply_to,omitempty"`
	Recipient string                    `json:"recipient,omitempty"`
	CC        string                    `json:"cc,omitempty"`
	BCC       string                    `json:"bcc,omitempty"`
	FakeBCC   *bool                     `json:"fake_bcc,omitempty"`
	Headers   []DesignStudioEmailHeader `json:"headers,omitempty"`
}

// DesignStudioEmail is a reusable email template/design managed in Design Studio.
type DesignStudioEmail struct {
	ID                 string                     `json:"id"`
	Name               string                     `json:"name,omitempty"`
	IsTemplate         bool                       `json:"is_template,omitempty"`
	IsLinked           bool                       `json:"is_linked,omitempty"`
	HasTranslations    bool                       `json:"has_translations,omitempty"`
	ParentFolderID     string                     `json:"parent_folder_id,omitempty"`
	AvailableLanguages []string                   `json:"available_languages,omitempty"`
	Created            int64                      `json:"created,omitempty"`
	Updated            int64                      `json:"updated,omitempty"`
	Content            *DesignStudioEmailContent  `json:"content,omitempty"`
	Envelope           *DesignStudioEmailEnvelope `json:"envelope,omitempty"`
	Transformers       map[string]any             `json:"transformers,omitempty"`
}

// ListDesignStudioEmailsOptions filters/sorts/paginates ListDesignStudioEmails.
type ListDesignStudioEmailsOptions struct {
	DesignStudioListOptions
	IsTemplate      DesignStudioFilter
	HasTranslations DesignStudioFilter
	IsLinked        DesignStudioFilter
}

// ListDesignStudioEmailsResponse is the decoded shape of GET /v1/design_studio/emails.
// List entries are a reduced shape — Content/Envelope/Transformers are only
// populated by GetDesignStudioEmail, not this list endpoint.
type ListDesignStudioEmailsResponse struct {
	Emails []DesignStudioEmail `json:"emails"`
}

// ListDesignStudioEmails returns Design Studio emails.
// See https://docs.customer.io/api/app/#operation/listDesignStudioEmails
func (c *APIClient) ListDesignStudioEmails(ctx context.Context, opts ListDesignStudioEmailsOptions) (*ListDesignStudioEmailsResponse, error) {
	q := opts.DesignStudioListOptions.apply(newQuery()).
		setString("is_template", string(opts.IsTemplate)).
		setString("has_translations", string(opts.HasTranslations)).
		setString("is_linked", string(opts.IsLinked))
	requestPath := "/v1/design_studio/emails" + q.String()

	var resp ListDesignStudioEmailsResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DesignStudioEmailResponse wraps a single DesignStudioEmail.
type DesignStudioEmailResponse struct {
	Email DesignStudioEmail `json:"email"`
}

// DesignStudioEmailInput creates a Design Studio email.
type DesignStudioEmailInput struct {
	Name           string                     `json:"name"`
	ParentFolderID string                     `json:"parent_folder_id,omitempty"`
	IsTemplate     *bool                      `json:"is_template,omitempty"`
	Content        *DesignStudioEmailContent  `json:"content,omitempty"`
	Envelope       *DesignStudioEmailEnvelope `json:"envelope,omitempty"`
	Transformers   map[string]any             `json:"transformers,omitempty"`
}

// CreateDesignStudioEmail creates a Design Studio email.
// See https://docs.customer.io/api/app/#operation/createDesignStudioEmail
func (c *APIClient) CreateDesignStudioEmail(ctx context.Context, input DesignStudioEmailInput) (*DesignStudioEmailResponse, error) {
	if input.Name == "" {
		return nil, ParamError{Param: "name"}
	}

	var resp DesignStudioEmailResponse
	if err := c.doJSON(ctx, "POST", "/v1/design_studio/emails", input, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetDesignStudioEmail returns one Design Studio email by id, including its
// full content/envelope/transformers.
// See https://docs.customer.io/api/app/#operation/getDesignStudioEmail
func (c *APIClient) GetDesignStudioEmail(ctx context.Context, emailID string) (*DesignStudioEmailResponse, error) {
	if emailID == "" {
		return nil, ParamError{Param: "emailID"}
	}

	var resp DesignStudioEmailResponse
	if err := c.doJSON(ctx, "GET", formatPath("/v1/design_studio/emails/%s", emailID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateDesignStudioEmail updates a Design Studio email. data is free-form.
// See https://docs.customer.io/api/app/#operation/updateDesignStudioEmail
func (c *APIClient) UpdateDesignStudioEmail(ctx context.Context, emailID string, data map[string]any) (*DesignStudioEmailResponse, error) {
	if emailID == "" {
		return nil, ParamError{Param: "emailID"}
	}

	var resp DesignStudioEmailResponse
	if err := c.doJSON(ctx, "PUT", formatPath("/v1/design_studio/emails/%s", emailID), data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteDesignStudioEmail deletes a Design Studio email.
// See https://docs.customer.io/api/app/#operation/deleteDesignStudioEmail
func (c *APIClient) DeleteDesignStudioEmail(ctx context.Context, emailID string) error {
	if emailID == "" {
		return ParamError{Param: "emailID"}
	}

	return c.doJSON(ctx, "DELETE", formatPath("/v1/design_studio/emails/%s", emailID), nil, nil, 200, 204)
}

// DesignStudioEmailTranslation is one language variant of a Design Studio email.
type DesignStudioEmailTranslation struct {
	Language string                     `json:"language,omitempty"`
	Content  *DesignStudioEmailContent  `json:"content,omitempty"`
	Envelope *DesignStudioEmailEnvelope `json:"envelope,omitempty"`
}

// ListDesignStudioEmailLanguagesResponse is the decoded shape of
// GET /v1/design_studio/emails/{id}/languages.
type ListDesignStudioEmailLanguagesResponse struct {
	EmailTranslations []DesignStudioEmailTranslation `json:"email_translations"`
}

// ListDesignStudioEmailLanguages returns every language variant of a Design
// Studio email.
// See https://docs.customer.io/api/app/#operation/listDesignStudioEmailLanguages
func (c *APIClient) ListDesignStudioEmailLanguages(ctx context.Context, emailID string) (*ListDesignStudioEmailLanguagesResponse, error) {
	if emailID == "" {
		return nil, ParamError{Param: "emailID"}
	}

	var resp ListDesignStudioEmailLanguagesResponse
	if err := c.doJSON(ctx, "GET", formatPath("/v1/design_studio/emails/%s/languages", emailID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DesignStudioEmailTranslationResponse wraps a single DesignStudioEmailTranslation.
type DesignStudioEmailTranslationResponse struct {
	EmailTranslation DesignStudioEmailTranslation `json:"email_translation"`
}

// DesignStudioEmailTranslationInput creates a language variant of a Design
// Studio email.
type DesignStudioEmailTranslationInput struct {
	Language     string                     `json:"language"`
	Content      *DesignStudioEmailContent  `json:"content,omitempty"`
	Envelope     *DesignStudioEmailEnvelope `json:"envelope,omitempty"`
	Transformers map[string]any             `json:"transformers,omitempty"`
}

// CreateDesignStudioEmailLanguage adds a language variant to a Design
// Studio email.
// See https://docs.customer.io/api/app/#operation/createDesignStudioEmailLanguage
func (c *APIClient) CreateDesignStudioEmailLanguage(ctx context.Context, emailID string, input DesignStudioEmailTranslationInput) (*DesignStudioEmailTranslationResponse, error) {
	if emailID == "" {
		return nil, ParamError{Param: "emailID"}
	}
	if input.Language == "" {
		return nil, ParamError{Param: "language"}
	}

	requestPath := formatPath("/v1/design_studio/emails/%s/languages", emailID)

	var resp DesignStudioEmailTranslationResponse
	if err := c.doJSON(ctx, "POST", requestPath, input, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetDesignStudioEmailLanguage returns one language variant of a Design
// Studio email.
// See https://docs.customer.io/api/app/#operation/getDesignStudioEmailLanguage
func (c *APIClient) GetDesignStudioEmailLanguage(ctx context.Context, emailID, language string) (*DesignStudioEmailTranslationResponse, error) {
	if emailID == "" {
		return nil, ParamError{Param: "emailID"}
	}
	if language == "" {
		return nil, ParamError{Param: "language"}
	}

	requestPath := formatPath("/v1/design_studio/emails/%s/languages/%s", emailID, language)

	var resp DesignStudioEmailTranslationResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateDesignStudioEmailLanguage updates a language variant of a Design
// Studio email. data is free-form.
// See https://docs.customer.io/api/app/#operation/updateDesignStudioEmailLanguage
func (c *APIClient) UpdateDesignStudioEmailLanguage(ctx context.Context, emailID, language string, data map[string]any) (*DesignStudioEmailTranslationResponse, error) {
	if emailID == "" {
		return nil, ParamError{Param: "emailID"}
	}
	if language == "" {
		return nil, ParamError{Param: "language"}
	}

	requestPath := formatPath("/v1/design_studio/emails/%s/languages/%s", emailID, language)

	var resp DesignStudioEmailTranslationResponse
	if err := c.doJSON(ctx, "PUT", requestPath, data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteDesignStudioEmailLanguage removes a language variant from a Design
// Studio email.
// See https://docs.customer.io/api/app/#operation/deleteDesignStudioEmailLanguage
func (c *APIClient) DeleteDesignStudioEmailLanguage(ctx context.Context, emailID, language string) error {
	if emailID == "" {
		return ParamError{Param: "emailID"}
	}
	if language == "" {
		return ParamError{Param: "language"}
	}

	requestPath := formatPath("/v1/design_studio/emails/%s/languages/%s", emailID, language)
	return c.doJSON(ctx, "DELETE", requestPath, nil, nil, 200, 204)
}

// DesignStudioComponent is a reusable content block managed in Design
// Studio. Content isn't part of the documented list/get response schema
// (despite being accepted on create/update) — kept here for round-tripping
// requests, may always be empty on responses.
type DesignStudioComponent struct {
	ID             string `json:"id"`
	Name           string `json:"name,omitempty"`
	Tag            string `json:"tag,omitempty"`
	ParentFolderID string `json:"parent_folder_id,omitempty"`
	Content        string `json:"content,omitempty"`
	Created        int64  `json:"created,omitempty"`
	Updated        int64  `json:"updated,omitempty"`
}

// ListDesignStudioComponentsOptions filters/sorts/paginates ListDesignStudioComponents.
type ListDesignStudioComponentsOptions struct {
	DesignStudioListOptions
	Tag string
}

// ListDesignStudioComponentsResponse is the decoded shape of GET /v1/design_studio/components.
type ListDesignStudioComponentsResponse struct {
	Components []DesignStudioComponent `json:"components"`
}

// ListDesignStudioComponents returns Design Studio components.
// See https://docs.customer.io/api/app/#operation/listDesignStudioComponents
func (c *APIClient) ListDesignStudioComponents(ctx context.Context, opts ListDesignStudioComponentsOptions) (*ListDesignStudioComponentsResponse, error) {
	q := opts.DesignStudioListOptions.apply(newQuery()).setString("tag", opts.Tag)
	requestPath := "/v1/design_studio/components" + q.String()

	var resp ListDesignStudioComponentsResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DesignStudioComponentResponse wraps a single DesignStudioComponent.
type DesignStudioComponentResponse struct {
	Component DesignStudioComponent `json:"component"`
}

// DesignStudioComponentInput creates a Design Studio component.
type DesignStudioComponentInput struct {
	Name           string `json:"name"`
	Tag            string `json:"tag"`
	ParentFolderID string `json:"parent_folder_id,omitempty"`
	Content        string `json:"content,omitempty"`
}

// CreateDesignStudioComponent creates a Design Studio component.
// See https://docs.customer.io/api/app/#operation/createDesignStudioComponent
func (c *APIClient) CreateDesignStudioComponent(ctx context.Context, input DesignStudioComponentInput) (*DesignStudioComponentResponse, error) {
	if input.Name == "" {
		return nil, ParamError{Param: "name"}
	}
	if input.Tag == "" {
		return nil, ParamError{Param: "tag"}
	}

	var resp DesignStudioComponentResponse
	if err := c.doJSON(ctx, "POST", "/v1/design_studio/components", input, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetDesignStudioComponent returns one Design Studio component by id.
// See https://docs.customer.io/api/app/#operation/getDesignStudioComponent
func (c *APIClient) GetDesignStudioComponent(ctx context.Context, componentID string) (*DesignStudioComponentResponse, error) {
	if componentID == "" {
		return nil, ParamError{Param: "componentID"}
	}

	var resp DesignStudioComponentResponse
	if err := c.doJSON(ctx, "GET", formatPath("/v1/design_studio/components/%s", componentID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateDesignStudioComponent updates a Design Studio component. data is free-form.
// See https://docs.customer.io/api/app/#operation/updateDesignStudioComponent
func (c *APIClient) UpdateDesignStudioComponent(ctx context.Context, componentID string, data map[string]any) (*DesignStudioComponentResponse, error) {
	if componentID == "" {
		return nil, ParamError{Param: "componentID"}
	}

	var resp DesignStudioComponentResponse
	if err := c.doJSON(ctx, "PUT", formatPath("/v1/design_studio/components/%s", componentID), data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteDesignStudioComponent deletes a Design Studio component.
// See https://docs.customer.io/api/app/#operation/deleteDesignStudioComponent
func (c *APIClient) DeleteDesignStudioComponent(ctx context.Context, componentID string) error {
	if componentID == "" {
		return ParamError{Param: "componentID"}
	}

	return c.doJSON(ctx, "DELETE", formatPath("/v1/design_studio/components/%s", componentID), nil, nil, 200, 204)
}
