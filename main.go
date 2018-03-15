package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lwalter/gobalance/balancer"
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

func main() {
	parseFlags()
	err := config.LoadConfig(cfgPath)

	if err != nil {
		log.Fatal(err)
	}

	wrkrs := pool.CreateWorkersFromConfig(config.Config.Pool.Workers)
	m, err := pool.NewManager(config.Config.Pool.Selection, wrkrs)

	if err != nil {
		log.Fatal("Could not instantiate a manager")
	}

	router := mux.NewRouter()
	catchAll := router.PathPrefix("/")
	catchAll.Handler(http.HandlerFunc(balancer.CatchAllHandler(m)))

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
