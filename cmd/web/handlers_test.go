package main

import (
	"bytes"
	"net/http"
	"testing"
)

func TestShowSnipept(t *testing.T) {
	app := newTestApplication(t)
	testServer := newTestServer(t, app.routes())
	defer testServer.Close()

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			code, _, body := testServer.get(t, tt.urlPath)
			if code != tt.wantCode {

				t.Errorf("want %d got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("body not contains the wanted body %q", body)
			}

		})
	}
}

func TestPing(t *testing.T) {

	app := newTestApplication(t)
	testServer := newTestServer(t, app.routes())
	defer testServer.Close()

	code, _, body := testServer.get(t, "/ping")

	if code != http.StatusOK {
		t.Errorf("want %d got %d", code, http.StatusOK)

	}

	if string(body) != "OK" {
		t.Errorf("want body equal to %q", "OK")
	}

}
