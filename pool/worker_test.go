package pool

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/lwalter/gobalance/config"
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

func TestCreateWorkersFromConfigGeneratesArray(t *testing.T) {
	cfg := []config.Worker{
		config.Worker{
			Host:   "localhost",
			Port:   8080,
			Scheme: "http",
		},
		config.Worker{
			Host:   "localhost",
			Port:   8081,
			Scheme: "http",
		},
	}

	wrkrs, err := CreateWorkersFromConfig(cfg)

	if err != nil {
		t.Errorf("Err returned from CreateWorkersFromConfig: %s", err)
	}

	cfgLen := len(cfg)
	wrkrLen := len(wrkrs)
	if len(wrkrs) != len(cfg) {
		t.Errorf("Length of created workers (%d) does not match config input (%d).", wrkrLen, cfgLen)
	}
}

func TestCreateWorkersFromConfigReturnsErrForEmptyInput(t *testing.T) {
	cfg := []config.Worker{}

	_, err := CreateWorkersFromConfig(cfg)

	if err == nil {
		t.Error("Err was not returned from CreateWorkersFromConfig")
	}
}
