package pool

import (
	"errors"
	"fmt"
)

// Manager manages the workers
type Manager struct {
	Workers   []*Worker
	Selection string
	curr      *Worker
}

// NewManager creates an instance of a Manager given a selection and a Worker array
func NewManager(selection string, workers []*Worker) (*Manager, error) {
	if len(workers) == 0 {
		return nil, errors.New("Cannot initialize Manager with no workers")
	}

	m := &Manager{
		Workers:   workers,
		Selection: selection,
		curr:      workers[0],
	}

	return m, nil
}

const (
	// RoundRobin defines a selection algorithm in which requests are distributed evenly amongst the workers in a pool
	RoundRobin = "roundrobin"

	// LeastConn defines a selection algorithm in which requests are sent to the worker that has the least amount of connections
	LeastConn = "leastconnections"

	// IPHash defines a selection algorithm in the IP address of the client is used to determine which worker receives the request
	IPHash = "iphash"
)

// GetNextWorker returns the next worker in line to perform work based on the selection algorithm
func (m *Manager) GetNextWorker() (*Worker, error) {
	found := false
	w := m.curr

	switch m.Selection {
	case RoundRobin:
		// Do roundrobin
		// Find the index of the current worker
		wCount := len(m.Workers)
		for i := 0; i < wCount; i++ {
			if m.Workers[i] == w {
				// Go to next worker in array, loop back if last element
				nxt := (i + 1) % wCount

				// Assign new current worker
				m.curr = m.Workers[nxt]

				found = true
				break
			}
		}

	case LeastConn:
		// Do least connection
		return nil, fmt.Errorf("%s not implemented", LeastConn)

	case IPHash:
		// Do IPHash
		return nil, fmt.Errorf("%s not implemented", IPHash)
	}

	if !found {
		return nil, errors.New("Could not retrieve next worker")
	}

	return w, nil
}
