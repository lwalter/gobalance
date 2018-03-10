package pool

import (
	"errors"
	"fmt"
	"net/http"
)

// Worker describes a member of a given pool which will receive requests
type Worker struct {
	Scheme string
	Host   string
	Port   int
	clnt   *http.Client
	// Statistics struct?
}

// NewWorker creates an instance of a Worker
func NewWorker(scheme string, host string, port int) *Worker {
	return &Worker{Scheme: scheme, Host: host, Port: port, clnt: &http.Client{}}
}

// SendRequest executes an http call to the current worker
func (w *Worker) SendRequest(method string, endpoint string) (*http.Response, error) {
	if w.clnt == nil {
		return nil, errors.New("Worker http client not initialized")
	}

	url := fmt.Sprintf("%s://%s:%d/%s", w.Scheme, w.Host, w.Port, endpoint)
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}

	resp, err := w.clnt.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}
