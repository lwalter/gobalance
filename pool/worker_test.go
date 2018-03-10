package pool

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

func createExampleServer() *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))

	return ts
}

func TestSendRequest(t *testing.T) {
	ts := createExampleServer()
	defer ts.Close()

	url, err := url.Parse(ts.URL)
	if err != nil {
		t.Errorf("Could not parse url from test server: %s", err)
	}

	port, err := strconv.Atoi(url.Port())

	if err != nil {
		t.Errorf("Could not parse host: %s", err)
	}

	wrkr := NewWorker(url.Scheme, url.Hostname(), port)

	resp, err := wrkr.SendRequest("GET", "/")

	if err != nil {
		t.Errorf("Could not send request: %s", err)
	}

	if resp == nil {
		t.Error("Response was nil")
	}

	if resp.StatusCode != 200 {
		t.Errorf("Invalid response from test server: %d", resp.StatusCode)
	}
}
