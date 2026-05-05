package customerio_test

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/customerio/go-customerio/v3"
)

type stubRoundTripper struct{}

func (stubRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       http.NoBody,
	}, nil
}

func TestAPIOptions(t *testing.T) {

	client := customerio.NewAPIClient("mykey")
	if client.URL != customerio.RegionUS.APIURL() {
		t.Errorf("wrong default url. got: %s, want: %s", client.URL, customerio.RegionUS.APIURL())
	}
	if client.Client != http.DefaultClient {
		t.Errorf("wrong default http client. got: %#v, want: %#v", client.Client, http.DefaultClient)
	}

	client = customerio.NewAPIClient("mykey", customerio.WithRegion(customerio.RegionEU))
	if client.URL != customerio.RegionEU.APIURL() {
		t.Errorf("wrong url. got: %s, want: %s", client.URL, customerio.RegionEU.APIURL())
	}

	hc := &http.Client{}
	client = customerio.NewAPIClient("mykey", customerio.WithHTTPClient(hc))
	if !reflect.DeepEqual(client.Client, hc) {
		t.Errorf("wrong http client. got: %#v, want: %#v", client.Client, hc)
	}

	customUserAgent := "Customer.io"
	client = customerio.NewAPIClient("mykey", customerio.WithUserAgent(customUserAgent))
	if client.UserAgent != customUserAgent {
		t.Errorf("wrong user-agent. got: %s, want: %s", client.UserAgent, customUserAgent)
	}
}

func TestTrackOptions(t *testing.T) {

	client := customerio.NewTrackClient("site_id", "api_key")
	if client.URL != customerio.RegionUS.TrackURL() {
		t.Errorf("wrong default url. got: %s, want: %s", client.URL, customerio.RegionUS.TrackURL())
	}
	defaultHTTPClient, ok := client.Client.(*http.Client)
	if !ok {
		t.Fatalf("expected default HTTP client to be *http.Client, got %T", client.Client)
	}
	if defaultHTTPClient.Timeout != 0 {
		t.Errorf("wrong default timeout. got: %s, want: 0s", defaultHTTPClient.Timeout)
	}
	defaultTransport, ok := defaultHTTPClient.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("expected default transport to be *http.Transport, got %T", defaultHTTPClient.Transport)
	}
	if defaultTransport.MaxIdleConnsPerHost != 100 {
		t.Errorf("wrong default max idle conns per host. got: %d, want: 100", defaultTransport.MaxIdleConnsPerHost)
	}

	client = customerio.NewTrackClient("site_id", "api_key", customerio.WithRegion(customerio.RegionEU))
	if client.URL != customerio.RegionEU.TrackURL() {
		t.Errorf("wrong url. got: %s, want: %s", client.URL, customerio.RegionEU.TrackURL())
	}

	hc := &http.Client{}
	client = customerio.NewTrackClient("site_id", "api_key", customerio.WithHTTPClient(hc))
	if !reflect.DeepEqual(client.Client, hc) {
		t.Errorf("wrong http client. got: %#v, want: %#v", client.Client, hc)
	}

	customUserAgent := "Customer.io"
	client = customerio.NewTrackClient("site_id", "api_key", customerio.WithUserAgent(customUserAgent))
	if client.UserAgent != customUserAgent {
		t.Errorf("wrong user-agent. got: %s, want: %s", client.UserAgent, customUserAgent)
	}
}

func TestNilOptionIsIgnored(t *testing.T) {
	var opt customerio.Option

	api := customerio.NewAPIClient("mykey", opt)
	if api.URL != customerio.RegionUS.APIURL() {
		t.Errorf("wrong default api url. got: %s, want: %s", api.URL, customerio.RegionUS.APIURL())
	}

	track := customerio.NewTrackClient("site_id", "api_key", opt)
	if track.URL != customerio.RegionUS.TrackURL() {
		t.Errorf("wrong default track url. got: %s, want: %s", track.URL, customerio.RegionUS.TrackURL())
	}
}

func TestTrackDefaultClientAcceptsInstrumentedDefaultTransport(t *testing.T) {
	original := http.DefaultTransport
	rt := stubRoundTripper{}
	http.DefaultTransport = rt
	t.Cleanup(func() {
		http.DefaultTransport = original
	})

	client := customerio.NewTrackClient("site_id", "api_key")
	defaultHTTPClient, ok := client.Client.(*http.Client)
	if !ok {
		t.Fatalf("expected default HTTP client to be *http.Client, got %T", client.Client)
	}
	if defaultHTTPClient.Transport != rt {
		t.Errorf("wrong default transport. got: %#v, want: %#v", defaultHTTPClient.Transport, rt)
	}
}
