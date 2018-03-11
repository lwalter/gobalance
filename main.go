package main

import (
	"fmt"
	"os"

	"github.com/lwalter/gobalance/config"
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
		fmt.Println(err)
		os.Exit(1)
	}

	//	router := mux.NewRouter()

	// router.Handle("/", http.HandlerFunc(func() {

	// }))
}
