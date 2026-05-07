package customerio_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/customerio/go-customerio/v3"
)

var (
	testDeliveryID = "ABCDEFG"
	testQueuedAt   = 1500111111
)

func transactionalServer(t *testing.T, verify func(request []byte)) (*customerio.APIClient, *httptest.Server) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		b, err := io.ReadAll(req.Body)
		if err != nil {
			t.Error(err)
		}
		defer func() {
			_ = req.Body.Close()
		}()

		verify(b)

		if _, err := w.Write([]byte(`{
				"delivery_id": "` + testDeliveryID + `",
				"queued_at": ` + strconv.Itoa(testQueuedAt) + `
			  }`)); err != nil {
			t.Error(err)
		}
	}))

	api := customerio.NewAPIClient("myKey")
	api.URL = srv.URL

	return api, srv
}
