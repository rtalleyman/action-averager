package main

import (
	"fmt"
	"os"
	"time"

	"github.com/action-averager/pkg/actionaverager"
)

const (
	exitOK        = 0
	exitErr       = 1
	delayDuration = 10

	delay = true
)

func handleMutlipleAdds(averager actionaverager.ActionAverager, actions []string, isDelay bool) {
	for i := range actions {
		handleAddAction(averager, actions[i])
		if isDelay {
			time.Sleep(delayDuration * time.Millisecond)
		}
	}
}

func handleAddAction(averager actionaverager.ActionAverager, action string) {
	fmt.Println("Adding action: ", action)
	if err := averager.AddAction(action); err != nil {
		fmt.Println("Encountered err: ", err)
		os.Exit(exitErr)
	}
}

func handleGetStats(averager actionaverager.ActionAverager) {
	stats := averager.GetStats()
	fmt.Println("Got stats: ", stats)
}

func main() {
	fmt.Println("Starting example run...")
	averager := actionaverager.NewActionAverager()

	action := `{"action":"run","time":50}`
	fmt.Println("Adding single action...")
	handleAddAction(averager, action)
	handleGetStats(averager)

	actions0 := []string{
		`{"action":"jump","time":111}`,
		`{"action":"crawl","time":300}`,
		`{"action":"jump","time":100}`,
	}

	fmt.Println("Adding multiple actions to previous action...")
	handleMutlipleAdds(averager, actions0, !delay)
	handleGetStats(averager)

	actions1 := []string{
		`{"action":"jump","time":87}`,
		`{"action":"run","time":75}`,
		`{"action":"crawl","time":323}`,
	}
	actions2 := []string{
		`{"action":"crawl","time":296}`,
		`{"action":"run","time":67}`,
		`{"action":"jump","time":123}`,
	}

	fmt.Println("Adding multiple concurrent actions to previous actions...")
	// NOTE: done is the sync channel for the concurrent go func
	done := make(chan bool)
	go func() {
		handleMutlipleAdds(averager, actions1, delay)
		done <- true
	}()
	handleMutlipleAdds(averager, actions2, delay)
	// NOTE: block until done is received meaning concurrent go func is finished
	<-done
	handleGetStats(averager)

	fmt.Println("Finished example run exiting...")
	os.Exit(exitOK)
}
