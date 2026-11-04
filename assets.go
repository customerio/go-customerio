package customerio

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"
)

// AssetListOptions filters/paginates the asset and asset-folder list endpoints.
type AssetListOptions struct {
	ParentFolderID        int
	DirectDescendantsOnly *bool
	Page                  int
	Limit                 int
}

func (o AssetListOptions) apply(q *queryBuilder) *queryBuilder {
	return q.setInt("parent_folder_id", o.ParentFolderID).
		setBool("direct_descendants_only", o.DirectDescendantsOnly).
		setInt("page", o.Page).
		setInt("limit", o.Limit)
}

// Asset is a file in the workspace's asset library. Path is where the file
// is hosted; the API doesn't document a content-type or dimensions field on
// the asset object itself (only as upload constraints).
type Asset struct {
	ID             int    `json:"id"`
	Name           string `json:"name,omitempty"`
	ParentFolderID *int   `json:"parent_folder_id,omitempty"`
	Path           string `json:"path,omitempty"`
	Size           int64  `json:"size,omitempty"`
	Created        int64  `json:"created,omitempty"`
	Updated        int64  `json:"updated,omitempty"`
}

// ListAssetsResponse is the decoded shape of GET /v1/assets.
type ListAssetsResponse struct {
	Assets []Asset `json:"assets"`
}

// ListAssets returns files in the workspace's asset library.
// See https://docs.customer.io/api/app/#operation/listAssets
func (c *APIClient) ListAssets(ctx context.Context, opts AssetListOptions) (*ListAssetsResponse, error) {
	requestPath := "/v1/assets" + opts.apply(newQuery()).String()

	var resp ListAssetsResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// AssetResponse wraps a single Asset.
type AssetResponse struct {
	Asset Asset `json:"asset"`
}

// assetContentTypeByExtension is the same extension -> MIME lookup the
// Node SDK uses when CreateAssetInput.ContentType isn't set.
var assetContentTypeByExtension = map[string]string{
	"bmp":  "image/bmp",
	"jpg":  "image/jpeg",
	"jpeg": "image/jpeg",
	"png":  "image/png",
	"gif":  "image/gif",
	"pdf":  "application/pdf",
}

// CreateAssetInput uploads a file to the workspace's asset library. The API
// accepts images (image/bmp, image/jpeg, image/png, image/gif) and
// application/pdf, up to 2MB (images up to 4096px per side). If ContentType
// is empty, it's derived from Filename's extension.
type CreateAssetInput struct {
	Data           io.Reader
	Filename       string
	ContentType    string
	Name           string
	ParentFolderID int
}

// multipartQuoteEscaper matches the escaping mime/multipart.Writer's own
// CreateFormFile uses for filenames in a Content-Disposition header — this
// client can't reuse that helper directly since it's unexported, and a
// custom Content-Type per part (set below) requires CreatePart instead of
// CreateFormFile anyway.
var multipartQuoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

// CreateAsset uploads a file as multipart/form-data.
// See https://docs.customer.io/api/app/#operation/createAsset
func (c *APIClient) CreateAsset(ctx context.Context, input CreateAssetInput) (*AssetResponse, error) {
	if input.Data == nil {
		return nil, ParamError{Param: "data"}
	}
	if input.Filename == "" {
		return nil, ParamError{Param: "filename"}
	}

	contentType := input.ContentType
	if contentType == "" {
		ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(input.Filename), "."))
		contentType = assetContentTypeByExtension[ext]
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	header := make(map[string][]string)
	header["Content-Disposition"] = []string{`form-data; name="file"; filename="` + multipartQuoteEscaper.Replace(input.Filename) + `"`}
	if contentType != "" {
		header["Content-Type"] = []string{contentType}
	}
	part, err := writer.CreatePart(header)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, input.Data); err != nil {
		return nil, err
	}

	if input.Name != "" {
		if err := writer.WriteField("name", input.Name); err != nil {
			return nil, err
		}
	}
	if input.ParentFolderID != 0 {
		if err := writer.WriteField("parent_folder_id", strconv.Itoa(input.ParentFolderID)); err != nil {
			return nil, err
		}
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	var resp AssetResponse
	if err := c.doMultipartJSON(ctx, "POST", "/v1/assets/files", writer.FormDataContentType(), &body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetAsset returns one asset by id.
// See https://docs.customer.io/api/app/#operation/getAsset
func (c *APIClient) GetAsset(ctx context.Context, assetID string) (*AssetResponse, error) {
	if assetID == "" {
		return nil, ParamError{Param: "assetID"}
	}

	var resp AssetResponse
	if err := c.doJSON(ctx, "GET", formatPath("/v1/assets/files/%s", assetID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateAsset updates an asset's name or parent folder. data is free-form.
// See https://docs.customer.io/api/app/#operation/updateAsset
func (c *APIClient) UpdateAsset(ctx context.Context, assetID string, data map[string]any) (*AssetResponse, error) {
	if assetID == "" {
		return nil, ParamError{Param: "assetID"}
	}

	var resp AssetResponse
	if err := c.doJSON(ctx, "PUT", formatPath("/v1/assets/files/%s", assetID), data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteAsset deletes an asset.
// See https://docs.customer.io/api/app/#operation/deleteAsset
func (c *APIClient) DeleteAsset(ctx context.Context, assetID string) error {
	if assetID == "" {
		return ParamError{Param: "assetID"}
	}

	return c.doJSON(ctx, "DELETE", formatPath("/v1/assets/files/%s", assetID), nil, nil, 200, 204)
}

// AssetFolder organizes assets into a tree.
type AssetFolder struct {
	ID             int    `json:"id"`
	Name           string `json:"name,omitempty"`
	ParentFolderID *int   `json:"parent_folder_id,omitempty"`
	Created        int64  `json:"created,omitempty"`
	Updated        int64  `json:"updated,omitempty"`
}

// ListAssetFoldersResponse is the decoded shape of GET /v1/assets/folders.
type ListAssetFoldersResponse struct {
	Folders []AssetFolder `json:"folders"`
}

// ListAssetFolders returns asset folders.
// See https://docs.customer.io/api/app/#operation/listAssetFolders
func (c *APIClient) ListAssetFolders(ctx context.Context, opts AssetListOptions) (*ListAssetFoldersResponse, error) {
	requestPath := "/v1/assets/folders" + opts.apply(newQuery()).String()

	var resp ListAssetFoldersResponse
	if err := c.doJSON(ctx, "GET", requestPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// AssetFolderResponse wraps a single AssetFolder.
type AssetFolderResponse struct {
	Folder AssetFolder `json:"folder"`
}

// AssetFolderInput creates an asset folder.
type AssetFolderInput struct {
	Name           string `json:"name"`
	ParentFolderID int    `json:"parent_folder_id,omitempty"`
}

// CreateAssetFolder creates an asset folder.
// See https://docs.customer.io/api/app/#operation/createAssetFolder
func (c *APIClient) CreateAssetFolder(ctx context.Context, input AssetFolderInput) (*AssetFolderResponse, error) {
	if input.Name == "" {
		return nil, ParamError{Param: "name"}
	}

	var resp AssetFolderResponse
	if err := c.doJSON(ctx, "POST", "/v1/assets/folders", input, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetAssetFolder returns one asset folder by id.
// See https://docs.customer.io/api/app/#operation/getAssetFolder
func (c *APIClient) GetAssetFolder(ctx context.Context, folderID string) (*AssetFolderResponse, error) {
	if folderID == "" {
		return nil, ParamError{Param: "folderID"}
	}

	var resp AssetFolderResponse
	if err := c.doJSON(ctx, "GET", formatPath("/v1/assets/folders/%s", folderID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateAssetFolder updates an asset folder. data is free-form.
// See https://docs.customer.io/api/app/#operation/updateAssetFolder
func (c *APIClient) UpdateAssetFolder(ctx context.Context, folderID string, data map[string]any) (*AssetFolderResponse, error) {
	if folderID == "" {
		return nil, ParamError{Param: "folderID"}
	}

	var resp AssetFolderResponse
	if err := c.doJSON(ctx, "PUT", formatPath("/v1/assets/folders/%s", folderID), data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteAssetFolder deletes an asset folder.
// See https://docs.customer.io/api/app/#operation/deleteAssetFolder
func (c *APIClient) DeleteAssetFolder(ctx context.Context, folderID string) error {
	if folderID == "" {
		return ParamError{Param: "folderID"}
	}

	return c.doJSON(ctx, "DELETE", formatPath("/v1/assets/folders/%s", folderID), nil, nil, 200, 204)
}
