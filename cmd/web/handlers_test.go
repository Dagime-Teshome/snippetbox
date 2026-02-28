package main

import (
	"bytes"
	"net/http"
	"net/url"
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

func TestSignUpUser(t *testing.T) {

	app := newTestApplication(t)
	testServer := newTestServer(t, app.routes())
	defer testServer.Close()

	tests := []struct {
		name            string
		userName        string
		userEmail       string
		userPassword    string
		useInvalidToken bool
		wantCode        int
		wantBody        []byte
	}{
		{"Empty name", "", "bob@example.com", "validPa$$word", false, http.StatusOK, []byte("This field cannot be blank")},
		{"Empty email", "Bob", "", "validPa$$word", false, http.StatusOK, []byte("This field cannot be blank")},
		{"Empty password", "Bob", "bob@example.com", "", false, http.StatusOK, []byte("This field cannot be blank")},
		{"Invalid email (incomplete domain)", "Bob", "bob@example.", "validPa$$word", false, http.StatusOK, []byte("Enter a valid email address")},
		{"Invalid email (missing @)", "Bob", "bobexample.com", "validPa$$word", false, http.StatusOK, []byte("Enter a valid email address")},
		{"Invalid email (missing local part)", "Bob", "@example.com", "validPa$$word", false, http.StatusOK, []byte("Enter a valid email address")},
		{"Short password", "Bob", "bob@example.com", "pa$$word", false, http.StatusOK, []byte("Password is too short")},
		{"Duplicate email", "Bob", "dupe@example.com", "validPa$$word", false, http.StatusOK, []byte("Email address is already in use")},
		{"Invalid CSRF Token", "", "", "", true, http.StatusBadRequest, []byte("Bad Request")},
		{"Valid submission", "Bob", "bob@example.com", "validPa$$word", false, http.StatusSeeOther, []byte("")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, body := testServer.get(t, "/user/signup")
			csrfToken := extractCSRFToken(t, body)
			if tt.useInvalidToken {
				csrfToken = "wrongToken"
			}
			form := url.Values{}
			form.Add("name", tt.userName)
			form.Add("email", tt.userEmail)
			form.Add("password", tt.userPassword)
			form.Add("csrf_token", csrfToken)

			code, _, body := testServer.postForm(t, "/user/signup", form)
			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body %q to contain %q", body, tt.wantBody)
			}
		})
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
