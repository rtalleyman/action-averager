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
	delayDuration = 100
	delay         = true
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
	fmt.Println("Single action...")
	handleAddAction(averager, action)
	handleGetStats(averager)

	actions0 := []string{
		`{"action":"jump","time":111}`,
		`{"action":"crawl","time":300}`,
		`{"action":"jump","time":100}`,
	}

	fmt.Println("Multiple actions...")
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

	fmt.Println("Multiple concurrent actions...")
	go handleMutlipleAdds(averager, actions1, delay)
	handleMutlipleAdds(averager, actions2, delay)
	handleGetStats(averager)

	fmt.Println("Finished example run exiting...")
	os.Exit(exitOK)
}
