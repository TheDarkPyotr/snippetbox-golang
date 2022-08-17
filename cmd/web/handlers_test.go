package main

import (
	"net/http"
	"testing"
)

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
