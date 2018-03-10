package config

import (
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

func TestLoadConfigPanicsWhenFileNotFound(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	path := "../test/asdf.yaml"
	LoadConfig(path)
}

func TestLoadConfigPanicsWhenSelectionNotSet(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("The code did not panic")
		}

		if r != "No selection algorithm was parsed from configuration file" {
			t.Errorf("Incorrect panic message: %s", r)
		}
	}()

	path := "../test/test1.yaml"
	LoadConfig(path)
}

func TestLoadConfigPanicsWhenWorkersNotSet(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("The code did not panic")
		}

		if r != "No workers were parsed from configuration file" {
			t.Errorf("Incorrect panic message: %s", r)
		}
	}()

	path := "../test/test2.yaml"
	LoadConfig(path)
}

func TestLoadConfigPanicsWhenSchemeNotSet(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("The code did not panic")
		}

		if r != "Worker at index 0 did not have a scheme set" {
			t.Errorf("Incorrect panic message: [%s]", r)
		}
	}()

	path := "../test/test3.yaml"
	LoadConfig(path)
}

func TestLoadConfigPanicsWhenHostNotSet(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("The code did not panic")
		}

		if r != "Worker at index 0 did not have a host set" {
			t.Errorf("Incorrect panic message: [%s]", r)
		}
	}()

	path := "../test/test4.yaml"
	LoadConfig(path)
}
