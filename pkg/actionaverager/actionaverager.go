package actionaverager

import (
	"encoding/json"
	"sync"
)

// ActionAverager averages times for actions
type ActionAverager interface {
	AddAction(string) error
	GetStats() string
}

type inputJSON struct {
	Action string  `json:"action"`
	Time   float64 `json:"time"`
}

type outputJSON struct {
	Action  string  `json:"action"`
	Average float64 `json:"avg"`
}

type actionData struct {
	TotalTime float64
	CallCount float64
}

type safeActionDatastore struct {
	Mux  sync.Mutex
	Data map[string]*actionData
}

// ActionAverage implements the ActionAverager interface
type ActionAverage struct {
	actionData *safeActionDatastore
}

// NewActionAverager creates a new ActionAverager
func NewActionAverager() ActionAverager {
	return &ActionAverage{
		actionData: &safeActionDatastore{
			Data: make(map[string]*actionData),
		},
	}
}

// AddAction takes a json serialized string and adds the action and time to the datastore
func (acav *ActionAverage) AddAction(input string) error {
	inputStruct := &inputJSON{}
	if err := json.Unmarshal([]byte(input), inputStruct); err != nil {
		return err
	}

	acav.actionData.Mux.Lock()
	defer acav.actionData.Mux.Unlock()

	// Check if action is already tracked in datastore if not add an entry for it, otherwise update existing entry
	data, ok := acav.actionData.Data[inputStruct.Action]
	if ok {
		data.TotalTime += inputStruct.Time
		data.CallCount++
	} else {
		ad := &actionData{
			TotalTime: inputStruct.Time,
			CallCount: 1,
		}
		acav.actionData.Data[inputStruct.Action] = ad
	}

	return nil
}

// GetStats computes the average time for each action in the datastore
func (acav *ActionAverage) GetStats() string {
	acav.actionData.Mux.Lock()
	defer acav.actionData.Mux.Unlock()

	var output []*outputJSON
	for action, data := range acav.actionData.Data {
		item := &outputJSON{
			Action:  action,
			Average: data.TotalTime / data.CallCount,
		}
		output = append(output, item)
	}

	// WARNING: should not suppress error, but have to because of assignment constraints
	jsonBytes, _ := json.Marshal(output)
	return string(jsonBytes)
}
