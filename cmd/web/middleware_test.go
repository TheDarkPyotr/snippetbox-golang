package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecureHeaders(t *testing.T) {
	// Initialize a new httptest.ResponseRecorder and dummy http.Request.
	rr := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Create a mock HTTP handler that we can pass to our secureHeaders
	// middleware, which writes a 200 status code and "OK" response body.
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	// Pass the mock HTTP handler to our secureHeaders middleware. Because
	// secureHeaders *returns* a http.Handler we can call its ServeHTTP()
	// method, passing in the http.ResponseRecorder and dummy http.Request to
	// execute it.
	secureHeaders(next).ServeHTTP(rr, r)
	// Call the Result() method on the http.ResponseRecorder to get the results
	// of the test.
	rs := rr.Result()
	// Check that the middleware has correctly set the X-Frame-Options header
	// on the response.
	frameOptions := rs.Header.Get("X-Frame-Options")
	if frameOptions != "deny" {
		t.Errorf("want %q; got %q", "deny", frameOptions)
	}
	// Check that the middleware has correctly set the X-XSS-Protection header
	// on the response.
	xssProtection := rs.Header.Get("X-XSS-Protection")
	if xssProtection != "1; mode=block" {
		t.Errorf("want %q; got %q", "1; mode=block", xssProtection)
	}

	// Check that the middleware has correctly called the next handler in line
	// and the response status code and body are as expected.
	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}
	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}

func TestShowSnippet(t *testing.T) {

	app := newTestApplication(t)
	testServer := newTestServer(t, app.routes())
	defer testServer.Close()

	// Set up some table-driven tests to check the responses sent by our
	// application for different URLs.
	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody []byte
	}{
		{"Valid ID", "/snippet/1", http.StatusOK, []byte("An old silent pond...")},
		{"Non-existent ID", "/snippet/2", http.StatusNotFound, nil},
		{"Negative ID", "/snippet/-1", http.StatusNotFound, nil},
		{"Decimal ID", "/snippet/1.23", http.StatusNotFound, nil},
		{"String ID", "/snippet/foo", http.StatusNotFound, nil},
		{"Empty ID", "/snippet/", http.StatusNotFound, nil},
		{"Trailing slash", "/snippet/1/", http.StatusNotFound, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			code, _, body := testServer.get(t, tt.urlPath)
			if code != tt.wantCode {
				t.Errorf("want %d, got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body contain %q", tt.wantBody)
			}

		})
	}

}
