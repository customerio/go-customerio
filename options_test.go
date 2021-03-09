package customerio

import (
	"net/http"
	"reflect"
	"testing"
)

func TestAPIOptions(t *testing.T) {

	client := NewAPIClient("mykey")
	if client.URL != RegionUS.apiURL {
		t.Errorf("wrong default url. got: %s, want: %s", client.URL, RegionUS.apiURL)
	}

	client = NewAPIClient("mykey", WithRegion(RegionEU))
	if client.URL != RegionEU.apiURL {
		t.Errorf("wrong url. got: %s, want: %s", client.URL, RegionEU.apiURL)
	}

	hc := &http.Client{}
	client = NewAPIClient("mykey", WithHTTPClient(hc))
	if !reflect.DeepEqual(client.Client, hc) {
		t.Errorf("wrong http client. got: %#v, want: %#v", client.Client, hc)
	}
}

func TestTrackOptions(t *testing.T) {

	client := NewTrackClient("site_id", "api_key")
	if client.URL != RegionUS.trackURL {
		t.Errorf("wrong default url. got: %s, want: %s", client.URL, RegionUS.trackURL)
	}

	client = NewTrackClient("site_id", "api_key", WithRegion(RegionEU))
	if client.URL != RegionEU.trackURL {
		t.Errorf("wrong url. got: %s, want: %s", client.URL, RegionEU.trackURL)
	}

	hc := &http.Client{}
	client = NewTrackClient("site_id", "api_key", WithHTTPClient(hc))
	if !reflect.DeepEqual(client.Client, hc) {
		t.Errorf("wrong http client. got: %#v, want: %#v", client.Client, hc)
	}
}
