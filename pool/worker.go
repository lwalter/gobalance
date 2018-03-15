package pool

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/lwalter/gobalance/config"
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

// CreateWorkersFromConfig generates a Worker array from configuration
func CreateWorkersFromConfig(cfg []config.Worker) ([]*Worker, error) {
	if len(cfg) == 0 {
		return nil, errors.New("Worker config is empty")
	}

	var wrkrs []*Worker
	for _, wc := range cfg {
		wrkr := NewWorker(wc.Scheme, wc.Host, wc.Port)
		wrkrs = append(wrkrs, wrkr)
	}

	return wrkrs, nil
}

// SendRequest executes an http call to the current worker
func (w *Worker) SendRequest(method string, endpoint string) (*http.Response, error) {
	if w.clnt == nil {
		return nil, errors.New("Worker http client not initialized")
	}

	url := fmt.Sprintf("%s://%s:%d%s", w.Scheme, w.Host, w.Port, endpoint)
	log.Printf(fmt.Sprintf("Sending request to %s", url))
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
