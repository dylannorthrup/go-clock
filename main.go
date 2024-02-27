package main

import (
	"fmt"
	"os"
)

func exit() {
	fmt.Println("")
	fmt.Println("Program execution complete")
	os.Exit(0)
}

type EnabledFunction struct {
	enabled  bool
	function func()
}

var items map[string]EnabledFunction

func populateItems() {
	items = map[string]EnabledFunction{
		"Clock": EnabledFunction{enabled: true, function: RunClock},
	}
}

func main() {
	populateItems()
	fmt.Printf("Running all enabled items\n")
	for k, v := range items {
		if v.enabled {
			fmt.Printf(":: Running %s\n", k)
			v.function()
		}
	}
	fmt.Printf("\n\nDone running all enabled items\n")
	exit()
}
