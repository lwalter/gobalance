package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lwalter/gobalance/config"
	"github.com/lwalter/gobalance/pool"
	"github.com/ogier/pflag"
)

var (
	cfgPath string
)

func init() {
	pflag.StringVarP(&cfgPath, "config", "c", "", "Configuration file path")
}

func parseFlags() {
	pflag.Parse()

	if pflag.NFlag() == 0 {
		fmt.Printf("Usage: %s [options]\n", os.Args[0])
		fmt.Println("Options:")
		pflag.PrintDefaults()
		os.Exit(1)
	}
}

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
			log.Println("Upstream worker failed to send request")
			http.Error(w, "Internal server error", 500)
			return
		}

		w.Write(body)
	}
}

func main() {
	parseFlags()
	err := config.LoadConfig(cfgPath)

	if err != nil {
		log.Fatal(err)
	}

	// TODO(lnw) this should probably be in the pool package
	var wrkrs []*pool.Worker
	for _, wc := range config.Config.Pool.Workers {
		wrkr := pool.NewWorker(wc.Scheme, wc.Host, wc.Port)
		wrkrs = append(wrkrs, wrkr)
	}

	m, err := pool.NewManager(config.Config.Pool.Selection, wrkrs)

	if err != nil {
		log.Fatal("Could not instantiate a manager")
	}

	router := mux.NewRouter()
	catchAll := router.PathPrefix("/")
	catchAll.Handler(http.HandlerFunc(CatchAllHandler(m)))

	p := strconv.Itoa(config.Config.Balancer.Port)
	h := config.Config.Balancer.Host
	addr := h + ":" + p
	srv := &http.Server{
		Handler: router,
		Addr:    addr,
	}

	log.Println(fmt.Sprintf("Starting gobalance at %s", srv.Addr))
	log.Fatal(srv.ListenAndServe())
}
