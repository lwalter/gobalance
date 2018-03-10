package pool

import (
	"testing"
)

func generateTestWorkers() []*Worker {
	workers := []*Worker{
		NewWorker("http", "localhost", 8080),
		NewWorker("http", "localhost", 8081),
	}
	return workers
}

func TestGetNextWorkerReturnsFirstWithRoundRobin(t *testing.T) {
	workers := generateTestWorkers()
	manager, err := NewManager(RoundRobin, workers)

	if err != nil {
		t.Error("Could not create a new Manager")
	}

	curr, err := manager.GetNextWorker()

	if err != nil {
		t.Errorf("Could not get next worker: %s", err)
	}

	if curr != workers[0] {
		t.Errorf("Expected the first worker, expected: %p actual: %p", workers[0], curr)
	}
}

func TestGetNextWorkerReturnsNextWithRoundRobin(t *testing.T) {
	workers := generateTestWorkers()
	manager, err := NewManager(RoundRobin, workers)

	if err != nil {
		t.Error("Could not create a new Manager")
	}

	// Move past intial
	manager.GetNextWorker()

	curr, err := manager.GetNextWorker()
	if err != nil {
		t.Errorf("Could not get next worker: %s", err)
	}

	if curr != workers[1] {
		t.Errorf("Expected the second worker, expected: %p actual: %p", workers[0], curr)
	}
}

func TestGetNextWorkerWrapsWithRoundRobin(t *testing.T) {
	workers := generateTestWorkers()
	manager, err := NewManager(RoundRobin, workers)

	if err != nil {
		t.Error("Could not create a new Manager")
	}

	// Move so that next wraps
	manager.GetNextWorker()
	manager.GetNextWorker()

	curr, err := manager.GetNextWorker()
	if err != nil {
		t.Errorf("Could not get next worker: %s", err)
	}

	if curr != workers[0] {
		t.Errorf("Expected the first worker, expected: %p actual: %p", workers[0], curr)
	}
}
