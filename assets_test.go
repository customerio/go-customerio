package customerio_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/customerio/go-customerio/v3"
)

func TestListAssets(t *testing.T) {
	api, rec := apiServer(t, 200, map[string]any{"assets": []map[string]any{{"id": 1, "name": "logo.png"}}})
	got, err := api.ListAssets(context.Background(), customerio.AssetListOptions{ParentFolderID: 5, Limit: 20})
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/assets?limit=20&parent_folder_id=5")
	if len(got.Assets) != 1 || got.Assets[0].Name != "logo.png" {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestCreateAsset(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.CreateAsset(context.Background(), customerio.CreateAssetInput{Filename: "a.png"}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "data")
	}
	if _, err := api.CreateAsset(context.Background(), customerio.CreateAssetInput{Data: strings.NewReader("x")}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "filename")
	}

	var gotContentType, gotFileContentType, gotFileContent, gotName, gotParentFolderID, gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		gotContentType = req.Header.Get("Content-Type")
		gotAuth = req.Header.Get("Authorization")
		if err := req.ParseMultipartForm(1 << 20); err != nil {
			t.Fatal(err)
		}
		file, header, err := req.FormFile("file")
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = file.Close() }()
		b, err := io.ReadAll(file)
		if err != nil {
			t.Fatal(err)
		}
		gotFileContent = string(b)
		gotFileContentType = header.Header.Get("Content-Type")
		gotName = req.FormValue("name")
		gotParentFolderID = req.FormValue("parent_folder_id")

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"asset":{"id":42,"name":"logo.png"}}`))
	}))
	defer srv.Close()

	api = customerio.NewAPIClient("myKey")
	api.URL = srv.URL

	got, err := api.CreateAsset(context.Background(), customerio.CreateAssetInput{
		Data:           strings.NewReader("fake png bytes"),
		Filename:       "logo.png",
		Name:           "Logo",
		ParentFolderID: 7,
	})
	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasPrefix(gotContentType, "multipart/form-data; boundary=") {
		t.Errorf("expected multipart content-type, got %q", gotContentType)
	}
	if gotAuth != "Bearer myKey" {
		t.Errorf("expected bearer auth, got %q", gotAuth)
	}
	if gotFileContent != "fake png bytes" {
		t.Errorf("expected file content %q, got %q", "fake png bytes", gotFileContent)
	}
	if gotFileContentType != "image/png" {
		t.Errorf("expected derived content-type image/png, got %q", gotFileContentType)
	}
	if gotName != "Logo" {
		t.Errorf("expected name field Logo, got %q", gotName)
	}
	if gotParentFolderID != "7" {
		t.Errorf("expected parent_folder_id field 7, got %q", gotParentFolderID)
	}
	if got.Asset.ID != 42 {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestCreateAssetExplicitContentType(t *testing.T) {
	var gotFileContentType string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if err := req.ParseMultipartForm(1 << 20); err != nil {
			t.Fatal(err)
		}
		file, header, err := req.FormFile("file")
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = file.Close() }()
		gotFileContentType = header.Header.Get("Content-Type")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"asset":{"id":1}}`))
	}))
	defer srv.Close()

	api := customerio.NewAPIClient("myKey")
	api.URL = srv.URL

	if _, err := api.CreateAsset(context.Background(), customerio.CreateAssetInput{
		Data:        strings.NewReader("data"),
		Filename:    "report",
		ContentType: "application/pdf",
	}); err != nil {
		t.Fatal(err)
	}
	if gotFileContentType != "application/pdf" {
		t.Errorf("expected explicit content-type application/pdf, got %q", gotFileContentType)
	}
}

func TestGetUpdateDeleteAsset(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetAsset(context.Background(), ""); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "assetID")
	}

	api, rec := apiServer(t, 200, map[string]any{"asset": map[string]any{"id": 1}})
	if _, err := api.GetAsset(context.Background(), "1"); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/assets/files/1")

	data := map[string]any{"name": "renamed.png"}
	if _, err := api.UpdateAsset(context.Background(), "1", data); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "PUT", "/v1/assets/files/1")
	assertJSONEqual(t, rec.body, data)

	if err := api.DeleteAsset(context.Background(), "1"); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "DELETE", "/v1/assets/files/1")
}

func TestAssetFolders(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.CreateAssetFolder(context.Background(), customerio.AssetFolderInput{}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "name")
	}

	api, rec := apiServer(t, 200, map[string]any{"folders": []map[string]any{{"id": 1, "name": "Logos"}}})
	got, err := api.ListAssetFolders(context.Background(), customerio.AssetListOptions{})
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/assets/folders")
	if len(got.Folders) != 1 {
		t.Errorf("unexpected response: %#v", got)
	}

	input := customerio.AssetFolderInput{Name: "Logos"}
	if _, err := api.CreateAssetFolder(context.Background(), input); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "POST", "/v1/assets/folders")

	if _, err := api.GetAssetFolder(context.Background(), "1"); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/assets/folders/1")

	if _, err := api.UpdateAssetFolder(context.Background(), "1", map[string]any{"name": "Brand"}); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "PUT", "/v1/assets/folders/1")

	if err := api.DeleteAssetFolder(context.Background(), "1"); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "DELETE", "/v1/assets/folders/1")
}
