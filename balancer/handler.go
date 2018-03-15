package balancer

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/lwalter/gobalance/pool"
)

// CatchAllHandler handles the incoming request to the load balancer
func CatchAllHandler(m *pool.Manager) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		wrkr, err := m.GetNextWorker()

		if err != nil {
			log.Println("Couldnt get a worker")
			http.Error(w, "Internal server error", 500)
			return
		}

		// TODO(lnw) make concurrent, use channels?
		resp, err := wrkr.SendRequest(r.Method, r.URL.RequestURI())

		if err != nil {
			log.Println("Upstream worker failed to send request")
			http.Error(w, "Internal server error", 500)
			return
		}

		w.WriteHeader(resp.StatusCode)
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Println("Could not send request to upstream worker")
			http.Error(w, "Internal server error", 500)
			return
		}

		w.Write(body)
	}
}
