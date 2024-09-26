package status

import (
	"errors"
	"fmt"
	"math/rand/v2"
)

// MockDB is a mock database that randomly fails so we can demonstrate the retry/queue logic.
type MockDBClient struct {
	FailureRate float64 // failureRate is a value (0.0, 1.0) defining the rate of simulated failures.
	Storage     map[string]any

	disabled bool // Set true to simulate the remote being unavailable, so an error is returned each time.
}

func NewMockDBClient(failureRate float64) *MockDBClient {
	return &MockDBClient{
		FailureRate: failureRate,
		Storage:     make(map[string]any),
	}
}

func (mock *MockDBClient) Stop() {
	mock.disabled = true
}

func (mock *MockDBClient) Start() {
	mock.disabled = false
}

func (mock *MockDBClient) Write(key string, value any) error {

	// Randomly simulate some network failure.
	rng := rand.Float64()
	if rng < mock.FailureRate {
		return errors.New("simulating unreachable database")
	}

	fmt.Println(value)
	mock.Storage[key] = value

	return nil
}
