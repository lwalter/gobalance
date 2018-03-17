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
func (w *Worker) SendRequest(r *http.Request) (*http.Response, error) {
	if w.clnt == nil {
		return nil, errors.New("Worker http client not initialized")
	}

	endpoint := r.URL.RequestURI()
	url := fmt.Sprintf("%s://%s:%d%s", w.Scheme, w.Host, w.Port, endpoint)
	log.Printf(fmt.Sprintf("Sending request to %s", url))

	u, err := http.NewRequest(r.Method, url, r.Body)

	if err != nil {
		return nil, err
	}

	w.AdjustHeaders(r, u)
	resp, err := w.clnt.Do(u)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// AdjustHeaders modifies http headers for a given request prior to sending upstream
func (w *Worker) AdjustHeaders(in *http.Request, out *http.Request) {
	log.Printf(fmt.Sprintf("remote: %s scheme: %s host: %s port: %s",
		in.RemoteAddr,
		in.URL.Scheme,
		in.Host,
		in.URL.Port()))

	if in.RemoteAddr != "" {
		out.Header.Set("X-Forwarded-For", in.RemoteAddr)
	}

	if in.URL.Scheme != "" {
		out.Header.Set("X-Forwarded-Scheme", in.URL.Scheme)
	}

	if in.Host != "" {
		out.Header.Set("X-Forwarded-Host", in.Host)
	}

	if in.URL.Port() != "" {
		out.Header.Set("X-Forwarded-Port", in.URL.Port())
	}

	out.Header.Set("Connection", "close")
}

// MergeHeaders asdf
func (w *Worker) MergeHeaders(r *http.Response, wr http.ResponseWriter) {
	for n, hs := range r.Header {
		for _, h := range hs {
			log.Printf(fmt.Sprintf("setting: [%s] to [%s]", n, h))
			wr.Header().Set(n, h)
		}
	}
}
