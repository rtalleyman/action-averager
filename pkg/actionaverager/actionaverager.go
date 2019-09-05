package actionaverager

import (
	"encoding/json"
	"fmt"
	"sync"
)

const (
	expInputLen = 2

	actKey         = "action"
	timeKey        = "time"
	emptyArrayJSON = "[]"
)

// ActionAverager averages times for actions
type ActionAverager interface {
	AddAction(string) error
	GetStats() string
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
	// NOTE: not doing an unmarshal to an explicit struct here, since input like
	// {"action":"run","random":"random"} will give {action:"run",time:0} and
	// {"action":"run","time":20,"random":"randon"} will give {action:"run",time:20}
	// unmarshaling to an interface allows explicit verification of fields
	var inInterface interface{}
	if err := json.Unmarshal([]byte(input), &inInterface); err != nil {
		return err
	}

	inMap, ok := inInterface.(map[string]interface{})
	if !ok {
		return fmt.Errorf("unable to convert input %s to internal data, rejecting", input)
	}

	numKeys := len(inMap)
	if numKeys != expInputLen {
		return fmt.Errorf("unexpected number of fields, %d, in input %s, expect %d, rejecting", numKeys, input, expInputLen)
	}

	action, ok := inMap[actKey]
	if !ok {
		return fmt.Errorf(`input %s is missing "action" field, rejecting`, input)
	}
	actStr, ok := action.(string)
	if !ok {
		return fmt.Errorf("action field is not a string in input %s, rejecting", input)
	}

	time, ok := inMap[timeKey]
	if !ok {
		return fmt.Errorf(`input %s is missing "time" field, rejecting`, input)
	}
	timeFlt, ok := time.(float64)
	if !ok {
		return fmt.Errorf("time field is not a number in input %s, rejecting", input)
	}
	if timeFlt < 0 {
		return fmt.Errorf("negative time value for input %s, rejecting", input)
	}

	acav.actionData.Mux.Lock()
	defer acav.actionData.Mux.Unlock()

	// Check if action is already tracked in datastore if not add an entry for it, otherwise update existing entry
	data, ok := acav.actionData.Data[actStr]
	if ok {
		// NOTE: data is a pointer to an actionData object so this will update the underlying object
		data.TotalTime += timeFlt
		data.CallCount++
	} else {
		ad := &actionData{
			TotalTime: timeFlt,
			CallCount: 1,
		}
		acav.actionData.Data[actStr] = ad
	}

	return nil
}

// GetStats computes the average time for each action in the datastore
func (acav *ActionAverage) GetStats() string {
	// NOTE: the defer unlock could be moved to after the for loop for performance, but is here for organization
	acav.actionData.Mux.Lock()
	defer acav.actionData.Mux.Unlock()

	var output []*outputJSON
	for action, data := range acav.actionData.Data {
		// NOTE: future proofing against possible divide by 0, currently this should not be possible
		if data.CallCount <= 0 {
			continue
		}
		// NOTE: divisions are expensive so only calculate averages when asked not as actions are added
		item := &outputJSON{
			Action:  action,
			Average: data.TotalTime / data.CallCount,
		}
		output = append(output, item)
	}

	// Return an empty json array if output is empty
	if len(output) == 0 {
		return emptyArrayJSON
	}

	// WARNING: should not suppress error, but have to because of assignment constraints.
	// However, this should not be an issue, since proper formatting is handled on our end.
	jsonBytes, _ := json.Marshal(output)
	return string(jsonBytes)
}
