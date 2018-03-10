package config

import (
	"fmt"
	"io/ioutil"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// Gobalance defines the configuration file structure
type Gobalance struct {
	Pool Pool
}

// Pool defines the configuration node for a Pool
type Pool struct {
	Workers   []Worker
	Selection string
}

// Worker defines the configuration node for a Worker
type Worker struct {
	Scheme string
	Host   string
	Port   int
}

// Config stores application configuration
var Config Gobalance

// LoadConfig reads in gobalance configuration from a yaml file
func LoadConfig(path string) {
	source, err := ioutil.ReadFile(path)

	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(source, &Config)

	if err != nil {
		panic(err)
	}

	validateConfig()
}

func validateConfig() {
	if Config.Pool.Selection == "" {
		panic("No selection algorithm was parsed from configuration file")
	}

	wCount := len(Config.Pool.Workers)
	if wCount == 0 {
		panic("No workers were parsed from configuration file")
	}

	var errs []string
	for i := 0; i < wCount; i++ {
		if Config.Pool.Workers[i].Host == "" {
			errs = append(errs, fmt.Sprintf("Worker at index %d did not have a host set", i))
		}

		if Config.Pool.Workers[i].Scheme == "" {
			errs = append(errs, fmt.Sprintf("Worker at index %d did not have a scheme set", i))
		}
	}

	if len(errs) > 0 {
		err := strings.Join(errs, "\n")
		panic(err)
	}
}
