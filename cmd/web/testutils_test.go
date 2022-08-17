package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"snippetbox/pkg/models/mock"
	"testing"
	"time"

	"github.com/golangcollege/sessions"
)

func newTestApplication(t *testing.T) *application {

	templateCache, err := newTemplateCache("./../../ui/html/")
	if err != nil {
		t.Fatal(err)
	}

	// Create a session manager instance, with the same settings as production.
	session := sessions.New([]byte("3dSm5MnygFHh7XidAtbskXrjbwfoJcbJ"))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	return &application{
		errorLog:      log.New(ioutil.Discard, "", 0),
		infoLog:       log.New(ioutil.Discard, "", 0),
		session:       session,
		database:      &mock.SnippetModel{},
		templateCache: templateCache,
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {

	tServer := httptest.NewTLSServer(h)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	// Add the cookie jar to the client, so that response cookies are stored
	tServer.Client().Jar = jar
	tServer.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return &testServer{tServer}

}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, []byte) {

	result, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}
	defer result.Body.Close()

	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		t.Fatal(err)
	}

	return result.StatusCode, result.Header, body

}
