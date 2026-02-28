package main

import (
	"html"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/Dagime-Teshome/snippetbox/pkg/models/mock"
	"github.com/golangcollege/sessions"
)

var csrfTokenRX = regexp.MustCompile(
	`name=['"]csrf_token['"][^>]*value=['"]([^'"]+)['"]`,
)

func extractCSRFToken(t *testing.T, body []byte) string {

	matches := csrfTokenRX.FindSubmatch(body)

	if len(matches) < 2 {
		t.Fatal("no csrf token found in body")
	}

	return html.UnescapeString(string(matches[1]))
}

func newTestApplication(t *testing.T) *application {

	templateCache, err := newTemplateCache("./../../ui/html/")
	if err != nil {
		t.Fatal(err)
	}
	secret := "xK7vQp3Lz9R4mN8tY2cHf6Bq1Ws5Dj0U"
	session := sessions.New([]byte(secret))
	session.Lifetime = time.Minute * 30
	session.SameSite = http.SameSiteStrictMode
	session.Secure = true

	return &application{
		errorLog:      log.New(io.Discard, "", 0),
		infoLog:       log.New(io.Discard, "", 0),
		Session:       session,
		snippet:       &mock.SnippetModel{},
		user:          &mock.UserModel{},
		templateCache: templateCache,
	}
}

type testServer struct {
	*httptest.Server
	client *http.Client
}

// func (ts *testServer) postForm(t *testing.T, urlPath string, form url.Values) (int, http.Header, []byte) {
// 	u, _ := url.Parse(ts.URL + urlPath)
// 	t.Logf("jar cookies before POST: %v", ts.client.Jar.Cookies(u))
// 	rs, err := ts.client.PostForm(ts.URL+urlPath, form)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	defer rs.Body.Close()

// 	body, err := io.ReadAll(rs.Body)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

//		return rs.StatusCode, rs.Header, body
//	}
func (ts *testServer) postForm(t *testing.T, urlPath string, form url.Values) (int, http.Header, []byte) {

	req, err := http.NewRequest("POST", ts.URL+urlPath, strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", ts.URL)

	rs, err := ts.client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, body
}
func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	client := ts.Client()
	client.Jar = jar
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts, client}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, []byte) {
	rs, err := ts.client.Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	return rs.StatusCode, rs.Header, body
}
