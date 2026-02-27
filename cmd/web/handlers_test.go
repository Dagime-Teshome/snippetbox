package main

import (
	"bytes"
	"net/http"
	"testing"
)

func TestPing(t *testing.T) {

	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")
	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}
	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}

}

func TestShowSnippet(t *testing.T) {

	tests := []struct {
		name       string
		urlPath    string
		wantStatus int
		wantBody   []byte
	}{
		{"Valid Id", "/snippet/1", http.StatusOK, []byte("An old silent pond...")},
		{"Invalid Id", "/snippet/2", http.StatusNotFound, nil},
		{"String Id", "/snippet/rand", http.StatusNotFound, nil},
		{"Empty Id", "/snippet/", http.StatusNotFound, nil},
		{"Trailing Id", "/snippet/1/", http.StatusNotFound, nil},
		{"Negative ID", "/snippet/-1", http.StatusNotFound, nil},
		{"Decimal ID", "/snippet/1.23", http.StatusNotFound, nil},
	}
	testRunner(t, tests)
}

func testRunner(t *testing.T, testList []struct {
	name       string
	urlPath    string
	wantStatus int
	wantBody   []byte
}) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	for _, test := range testList {
		t.Run(test.name, func(t *testing.T) {
			status, _, body := ts.get(t, test.urlPath)
			if status != test.wantStatus {
				t.Errorf("wanted %d got %d", test.wantStatus, status)
			}
			if !bytes.Contains(body, test.wantBody) {
				t.Errorf("wanted %s got %s", string(test.wantBody), string(body))
			}
		})
	}

}
