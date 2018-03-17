package balancer

import (
	"fmt"
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
			log.Println(fmt.Sprintf("Couldnt get a worker: %s", err))
			http.Error(w, "Internal server error", 500)
			return
		}

		// TODO(lnw) make concurrent, use channels?
		resp, err := wrkr.SendRequest(r)

		if err != nil {
			log.Println(fmt.Sprintf("Upstream worker failed to send request: %s", err))
			http.Error(w, "Internal server error", 500)
			return
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Println(fmt.Sprintf("Could not send request to upstream worker: %s", err))
			http.Error(w, "Internal server error", 500)
			return
		}

		wrkr.MergeHeaders(resp, w)
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
	}
}
