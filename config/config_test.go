package config

import (
	"errors"
	"fmt"
	"testing"
)

// TODO(lnw) will want to introduce file system mocking
func TestLoadConfigReadsFile(t *testing.T) {
	path := "../test/pool.yaml"
	LoadConfig(path)

	if Config.Pool.Selection == "" {
		t.Errorf("The parsed Selection value was empty from config file: %s", path)
	}

	if Config.Pool.Workers == nil || len(Config.Pool.Workers) != 2 {
		t.Errorf("The parsed Workers collection did not have two Workers from config file: %s", path)
	}
}

func TestLoadConfigReturnsErrWhenFileNotFound(t *testing.T) {
	path := "../test/asdf.yaml"
	err := LoadConfig(path)

	if err == nil {
		t.Errorf("Expected error from LoadConfig")
	}

	expected := fmt.Errorf("Config file does not exist: %s", path)
	if err.Error() != expected.Error() {
		t.Errorf("Actual error different than expected.\nActual: [%s]\nExpected: [%s]", err, expected)
	}
}

func TestLoadConfigReturnsErrWhenBalancerHostNotSet(t *testing.T) {
	path := "../test/test5.yaml"
	err := LoadConfig(path)

	if err == nil {
		t.Errorf("Expected error from LoadConfig")
	}

	expected := errors.New("No host specified for load balancer")
	if err.Error() != expected.Error() {
		t.Errorf("Actual error different than expected.\nActual: %s\nExpected: %s", err, expected)
	}
}

func TestLoadConfigReturnsErrWhenSelectionNotSet(t *testing.T) {
	path := "../test/test1.yaml"
	err := LoadConfig(path)

	if err == nil {
		t.Errorf("Expected error from LoadConfig")
	}

	expected := errors.New("No selection algorithm was parsed from configuration file")
	if err.Error() != expected.Error() {
		t.Errorf("Actual error different than expected.\nActual: %s\nExpected: %s", err, expected)
	}
}

func TestLoadConfigReturnsErrWhenWorkersNotSet(t *testing.T) {
	path := "../test/test2.yaml"
	err := LoadConfig(path)

	if err == nil {
		t.Errorf("Expected error from LoadConfig")
	}

	expected := errors.New("No workers were parsed from configuration file")
	if err.Error() != expected.Error() {
		t.Errorf("Actual error different than expected.\nActual: %s\nExpected: %s", err, expected)
	}
}

func TestLoadConfigReturnsErrWhenSchemeNotSet(t *testing.T) {
	path := "../test/test3.yaml"
	err := LoadConfig(path)

	if err == nil {
		t.Errorf("Expected error from LoadConfig")
	}

	expected := errors.New("Worker at index 0 did not have a scheme set")
	if err.Error() != expected.Error() {
		t.Errorf("Actual error different than expected.\nActual: %s\nExpected: %s", err, expected)
	}
}

func TestLoadConfigReturnsErrWhenHostNotSet(t *testing.T) {
	path := "../test/test4.yaml"
	err := LoadConfig(path)

	if err == nil {
		t.Errorf("Expected error from LoadConfig")
	}

	expected := errors.New("Worker at index 0 did not have a host set")
	if err.Error() != expected.Error() {
		t.Errorf("Actual error different than expected.\nActual: %s\nExpected: %s", err, expected)
	}
}
