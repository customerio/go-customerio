package customerio

import (
	"fmt"
	"net/url"
	"os"
	"testing"
	"time"
)

var cio *CustomerIO

func TestMain(m *testing.M) {

	siteID := os.Getenv("CUSTOMERIO_SITE_ID")
	apiKey := os.Getenv("CUSTOMERIO_API_KEY")

	var exitCode int

	if siteID != "" && apiKey != "" {
		cio = NewCustomerIO(siteID, apiKey)
		exitCode = m.Run()
	} else {
		exitCode = 1
		fmt.Println("Must set CUSTOMERIO_SITE_ID and CUSTOMERIO_API_KEY environment variables to test this library")
	}

	os.Exit(exitCode)
}

func TestIdentify(t *testing.T) {

	attributes := map[string]interface{}{}

	err := cio.Identify("golang-test-noattributes", attributes)
	defer cio.Delete("golang-test-noattributes")

	if err != nil {
		t.Error(err.Error())
	}

	attributes["email"] = "golang@customer.io"
	attributes["first_name"] = "golang"
	attributes["last_name"] = "testsuite"

	err = cio.Identify("golang-test-stringattributes", attributes)
	defer cio.Delete("golang-test-stringattributes")

	if err != nil {
		t.Error(err.Error())
	}

	attributes["paid"] = true
	attributes["numUsers"] = 10

	err = cio.Identify("golang-test-mixed-attributes", attributes)
	defer cio.Delete("golang-test-mixed-attributes")

	if err != nil {
		t.Error(err.Error())
	}

	attributes["_last_visit"] = time.Now().Unix()

	err = cio.Identify("golang-test-magic-attributes", attributes)
	defer cio.Delete("golang-test-magic-attributes")

	if err != nil {
		t.Error(err.Error())
	}
}

func TestTrack(t *testing.T) {
	cio.Identify("golang-test-events", map[string]interface{}{})
	defer cio.Delete("golang-test-events")

	err := cio.Track("golang-test-events", "golang-test", map[string]interface{}{})

	if err != nil {
		t.Error(err.Error())
	}

	err = cio.Track("golang-test-events", "golang-test-data", map[string]interface{}{"value": 1, "name": "event"})

	if err != nil {
		t.Error(err.Error())
	}

}

func TestTrackAnonymous(t *testing.T) {
	err := cio.TrackAnonymous("golang-test-anonymous", map[string]interface{}{"recipient": "golang@customer.io"})
	if err != nil {
		t.Error(err.Error())
	}

	err = cio.TrackAnonymous("golang-test-data", map[string]interface{}{"value": 1, "name": "event"})
	if err != nil {
		t.Error(err.Error())
	}
}

func TestDelete(t *testing.T) {
	cio.Identify("golang-test-delete", map[string]interface{}{})

	err := cio.Delete("golang-test-delete")
	if err != nil {
		t.Error(err.Error())
	}
}

func TestAddDevice(t *testing.T) {
	err := cio.Identify("golang-test-adddevice", map[string]interface{}{})
	defer cio.Delete("golang-test-adddevice")

	if err != nil {
		t.Error(err.Error())
	}

	err = cio.AddDevice("golang-test-adddevice", "golang-test-add", "ios", map[string]interface{}{"last_used": time.Now().Unix()})
	defer cio.DeleteDevice("golang-test-adddevice", "golang-test-add")

	if err != nil {
		t.Error(err.Error())
	}
}

func TestDeleteDevice(t *testing.T) {
	cio.Identify("golang-test-deletedevice", map[string]interface{}{})
	defer cio.Delete("golang-test-deletedevice")

	cio.AddDevice("golang-test-deletedevice", "golang-test-delete", "android", map[string]interface{}{"last_used": time.Now().Unix()})

	err := cio.DeleteDevice("golang-test-deletedevice", "golang-test-delete")
	if err != nil {
		t.Error(err.Error())
	}
}

func TestStringEncoding(t *testing.T) {
	encoded := encodeID(url.PathEscape("test path"))
	if encoded != "test%20path" {
		t.Errorf("got: %s, want: %s", encoded, url.PathEscape("test path"))
	}

	encoded = encodeID("test path")
	if encoded != "test%20path" {
		t.Errorf("got: %s, want: %s", encoded, url.PathEscape("test path"))
	}
}
