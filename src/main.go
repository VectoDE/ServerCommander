package main

import (
	"log"

	"servercommander/src/cmd"
	"servercommander/src/ui"
)

func main() {
	if err := ui.RunStandaloneConsole(cmd.Execute); err != nil {
		log.Fatal(err)
	}
}
