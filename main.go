package main

import (
	"github.com/lwalter/gobalance/config"
)

func main() {
	config.LoadConfig("test/pool.yaml")

	//	router := mux.NewRouter()

	// router.Handle("/", http.HandlerFunc(func() {

	// }))
}
