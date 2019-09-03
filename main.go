package main

import (
	"fmt"
	"os"

	"github.com/action-averager/pkg/actionaverager"
)

const (
	exitOK  = 0
	exitErr = 1
)

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

	fmt.Println("Single action...")
	action := `{"action":"run","time":75}`
	handleAddAction(averager, action)
	handleGetStats(averager)

	fmt.Println("Multiple actions...")
	actions := []string{
		`{"action":"jump","time":123}`,
		`{"action":"crawl","time":300}`,
		`{"action":"jump","time":456}`,
	}
	for i := range actions {
		handleAddAction(averager, actions[i])
	}
	handleGetStats(averager)

	fmt.Println("Finished example run exiting...")
	os.Exit(exitOK)
}
