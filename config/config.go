package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// Gobalance defines the configuration file structure
type Gobalance struct {
	Pool     Pool
	Balancer Balancer
}

// Pool defines the configuration node for a Pool
type Pool struct {
	Workers   []Worker
	Selection string
}

// Balancer defines the configuraiton node for the http load balancer
type Balancer struct {
	Port int
	Host string
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
func LoadConfig(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("Config file does not exist: %s", path)
		}

		return fmt.Errorf("Could not read file at: %s\n%s", path, err)
	}

	source, err := ioutil.ReadFile(path)

	if err != nil {
		return fmt.Errorf("Could not read file: %s\n%s", path, err)
	}

	err = yaml.Unmarshal(source, &Config)

	if err != nil {
		return fmt.Errorf("Could not parse config file at: %s\n%s", path, err)
	}

	err = validateConfig()

	if err != nil {
		return err
	}

	return nil
}

func validateConfig() error {
	if Config.Balancer.Host == "" {
		return errors.New("No host specified for load balancer")
	}

	if Config.Pool.Selection == "" {
		return errors.New("No selection algorithm was parsed from configuration file")
	}

	wCount := len(Config.Pool.Workers)
	if wCount == 0 {
		return errors.New("No workers were parsed from configuration file")
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
		return errors.New(err)
	}

	return nil
}
