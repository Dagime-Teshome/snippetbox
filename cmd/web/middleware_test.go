package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecureHeader(t *testing.T) {
	rr := httptest.NewRecorder()

	req := httptest.NewRequest("GET", "/", nil)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	securityMiddleware(handler).ServeHTTP(rr, req)

	result := rr.Result()
	frameOpts := result.Header.Get("X-Frame-Options")
	if frameOpts != "deny" {
		t.Errorf("expected %s got %s ", "deny", frameOpts)
	}
	xssProtection := result.Header.Get("X-XSS-Protection")
	if xssProtection != "1; mode=block" {
		t.Errorf("expected %s got %s", "1; mode=block", xssProtection)
	}

	statusCode := result.StatusCode

	if statusCode != http.StatusOK {
		t.Errorf("expected %d got %d", http.StatusOK, statusCode)
	}
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatal(err)
	}
	bodyString := string(body)
	if bodyString != "OK" {
		t.Errorf("expected %s got %s", "OK", bodyString)
	}
}
