package customerio

import (
	"net/http"
	"reflect"
	"testing"
)

func TestAPIOptions(t *testing.T) {

	client := NewAPIClient("mykey")
	if client.URL != apiURL_US {
		t.Errorf("wrong default url. got: %s, want: %s", client.URL, apiURL_US)
	}

	client = NewAPIClient("mykey", WithURL("http://example.com"))
	if client.URL != "http://example.com" {
		t.Errorf("wrong url. got: %s, want: http://example.com", client.URL)
	}

	client = NewAPIClient("mykey", WithRegion(Region_EU))
	if client.URL != apiURL_EU {
		t.Errorf("wrong url. got: %s, want: %s", client.URL, apiURL_EU)
	}

	hc := &http.Client{}
	client = NewAPIClient("mykey", WithHTTPClient(hc))
	if !reflect.DeepEqual(client.Client, hc) {
		t.Errorf("wrong http client. got: %#v, want: %#v", client.Client, hc)
	}

	client = NewAPIClient("mykey", WithURL("http://example.com"), WithRegion(Region_EU))
	if client.URL != apiURL_EU {
		t.Errorf("wrong url. got: %s, want: %s", client.URL, apiURL_EU)
	}
}

func TestTrackOptions(t *testing.T) {

	client := NewTrackClient("site_id", "api_key")
	if client.URL != trackURL_US {
		t.Errorf("wrong default url. got: %s, want: %s", client.URL, trackURL_US)
	}

	client = NewTrackClient("site_id", "api_key", WithURL("http://example.com"))
	if client.URL != "http://example.com" {
		t.Errorf("wrong url. got: %s, want: http://example.com", client.URL)
	}

	client = NewTrackClient("site_id", "api_key", WithRegion(Region_EU))
	if client.URL != trackURL_EU {
		t.Errorf("wrong url. got: %s, want: %s", client.URL, trackURL_EU)
	}

	hc := &http.Client{}
	client = NewTrackClient("site_id", "api_key", WithHTTPClient(hc))
	if !reflect.DeepEqual(client.Client, hc) {
		t.Errorf("wrong http client. got: %#v, want: %#v", client.Client, hc)
	}

	client = NewTrackClient("site_id", "api_key", WithURL("http://example.com"), WithRegion(Region_EU))
	if client.URL != trackURL_EU {
		t.Errorf("wrong url. got: %s, want: %s", client.URL, trackURL_EU)
	}
}
