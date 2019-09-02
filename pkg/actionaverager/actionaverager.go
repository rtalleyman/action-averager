package actionaverager

import (
	"sync"
)

// ActionAverager averages times for actions
type ActionAverager interface {
	AddAction(string) error
	GetStats() string
}

type actionData struct {
	TotalTime uint64
	CallCount uint64
}

type safeActionDatastore struct {
	Mux  sync.Mutex
	Data map[string]actionData
}

// ActionAverage implements the ActionAverager interface
type ActionAverage struct {
	actionData *safeActionDatastore
}

// NewActionAverager creates a new ActionAverager
func NewActionAverager() ActionAverager {
	return &ActionAverage{
		actionData: &safeActionDatastore{
			Data: make(map[string]actionData),
		},
	}
}

// AddAction takes a json serialized string and adds the action and time to the datastore
func (acav *ActionAverage) AddAction(input string) error {
	return nil
}

// GetStats computes the average time for each action in the datastore
func (acav *ActionAverage) GetStats() string {
	return ""
}