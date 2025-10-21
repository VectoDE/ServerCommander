package main

import (
	"log"

	"servercommander/src/cmd"
	"servercommander/src/console"
)

func main() {
	if err := console.Run(cmd.Execute); err != nil {
		log.Fatal(err)
	}
}
