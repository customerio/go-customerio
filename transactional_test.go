package customerio_test

import (
	"io/ioutil"
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
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			t.Error(err)
		}
		defer req.Body.Close()

		verify(b)

		w.Write([]byte(`{
			"delivery_id": "` + testDeliveryID + `",
			"queued_at": ` + strconv.Itoa(testQueuedAt) + `
		  }`))
	}))

	api := customerio.NewAPIClient("myKey")
	api.URL = srv.URL

	return api, srv
}
