package customerio_test

import (
	"context"
	"testing"

	"github.com/customerio/go-customerio/v3"
)

func TestDesignStudioFolders(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.CreateDesignStudioFolder(context.Background(), customerio.DesignStudioFolderInput{}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "name")
	}
	if _, err := api.GetDesignStudioFolder(context.Background(), ""); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "folderID")
	}

	api, rec := apiServer(t, 200, map[string]any{"folders": []map[string]any{{"id": "f1", "name": "Newsletters"}}})
	got, err := api.ListDesignStudioFolders(context.Background(), customerio.DesignStudioListOptions{
		ParentFolderID: "root", DirectDescendantsOnly: boolPtr(true), SortBy: customerio.DesignStudioSortByName, Limit: 10,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/design_studio/folders?direct_descendants_only=true&limit=10&parent_folder_id=root&sort_by=name")
	if len(got.Folders) != 1 || got.Folders[0].Name != "Newsletters" {
		t.Errorf("unexpected response: %#v", got)
	}

	if _, err := api.CreateDesignStudioFolder(context.Background(), customerio.DesignStudioFolderInput{Name: "New"}); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "POST", "/v1/design_studio/folders")

	if _, err := api.GetDesignStudioFolder(context.Background(), "f1"); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/design_studio/folders/f1")

	data := map[string]any{"name": "Renamed"}
	if _, err := api.UpdateDesignStudioFolder(context.Background(), "f1", data); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "PUT", "/v1/design_studio/folders/f1")
	assertJSONEqual(t, rec.body, data)

	if err := api.DeleteDesignStudioFolder(context.Background(), "f1"); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "DELETE", "/v1/design_studio/folders/f1")
}

func TestDesignStudioEmails(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.CreateDesignStudioEmail(context.Background(), customerio.DesignStudioEmailInput{}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "name")
	}

	api, rec := apiServer(t, 200, map[string]any{"emails": []map[string]any{{"id": "e1", "name": "Welcome"}}})
	got, err := api.ListDesignStudioEmails(context.Background(), customerio.ListDesignStudioEmailsOptions{
		IsTemplate: customerio.DesignStudioFilterAny,
		IsLinked:   customerio.DesignStudioFilterFalse,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/design_studio/emails?is_linked=false&is_template=any")
	if len(got.Emails) != 1 {
		t.Errorf("unexpected response: %#v", got)
	}

	input := customerio.DesignStudioEmailInput{
		Name:    "Welcome",
		Content: &customerio.DesignStudioEmailContent{Subject: "Hi"},
	}
	if _, err := api.CreateDesignStudioEmail(context.Background(), input); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "POST", "/v1/design_studio/emails")
	assertJSONEqual(t, rec.body, input)

	if _, err := api.GetDesignStudioEmail(context.Background(), "e1"); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/design_studio/emails/e1")

	updateData := map[string]any{"name": "Renamed"}
	if _, err := api.UpdateDesignStudioEmail(context.Background(), "e1", updateData); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "PUT", "/v1/design_studio/emails/e1")

	if err := api.DeleteDesignStudioEmail(context.Background(), "e1"); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "DELETE", "/v1/design_studio/emails/e1")
}

func TestDesignStudioEmailLanguages(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.CreateDesignStudioEmailLanguage(context.Background(), "e1", customerio.DesignStudioEmailTranslationInput{}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "language")
	}

	api, rec := apiServer(t, 200, map[string]any{"email_translations": []map[string]any{{"language": "en"}}})
	got, err := api.ListDesignStudioEmailLanguages(context.Background(), "e1")
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/design_studio/emails/e1/languages")
	if len(got.EmailTranslations) != 1 {
		t.Errorf("unexpected response: %#v", got)
	}

	input := customerio.DesignStudioEmailTranslationInput{Language: "fr"}
	if _, err := api.CreateDesignStudioEmailLanguage(context.Background(), "e1", input); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "POST", "/v1/design_studio/emails/e1/languages")

	if _, err := api.GetDesignStudioEmailLanguage(context.Background(), "e1", "fr"); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/design_studio/emails/e1/languages/fr")

	data := map[string]any{"content": map[string]any{"subject": "Bonjour"}}
	if _, err := api.UpdateDesignStudioEmailLanguage(context.Background(), "e1", "fr", data); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "PUT", "/v1/design_studio/emails/e1/languages/fr")

	if err := api.DeleteDesignStudioEmailLanguage(context.Background(), "e1", "fr"); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "DELETE", "/v1/design_studio/emails/e1/languages/fr")
}

func TestDesignStudioComponents(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.CreateDesignStudioComponent(context.Background(), customerio.DesignStudioComponentInput{Name: "Header"}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "tag")
	}

	api, rec := apiServer(t, 200, map[string]any{"components": []map[string]any{{"id": "c1", "tag": "header"}}})
	got, err := api.ListDesignStudioComponents(context.Background(), customerio.ListDesignStudioComponentsOptions{Tag: "header"})
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/design_studio/components?tag=header")
	if len(got.Components) != 1 {
		t.Errorf("unexpected response: %#v", got)
	}

	input := customerio.DesignStudioComponentInput{Name: "Header", Tag: "header"}
	if _, err := api.CreateDesignStudioComponent(context.Background(), input); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "POST", "/v1/design_studio/components")

	if _, err := api.GetDesignStudioComponent(context.Background(), "c1"); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/design_studio/components/c1")

	if _, err := api.UpdateDesignStudioComponent(context.Background(), "c1", map[string]any{"name": "Header 2"}); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "PUT", "/v1/design_studio/components/c1")

	if err := api.DeleteDesignStudioComponent(context.Background(), "c1"); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "DELETE", "/v1/design_studio/components/c1")
}
