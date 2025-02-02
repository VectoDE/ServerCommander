package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"servercommander/src/cli"
	"servercommander/src/ui"
	"servercommander/src/utils"
)

func main() {
	ui.ApplicationBanner()

	for {
		fmt.Println(utils.Cyan, "Choose Mode:", utils.Reset)
		fmt.Println("[1] Start CLI")
		fmt.Println("[2] Upcoming Features")
		fmt.Println("[3] Exit")
		fmt.Print(utils.Yellow, ">> ", utils.Reset)

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		choice := strings.TrimSpace(input)

		if choice == "1" {
			cli.StartCLI()
			break
		} else if choice == "2" {
			ui.UpcomingFeaturesBanner()
		} else if choice == "3" {
			fmt.Println(utils.Red, "Program is exiting.", utils.Reset)
			return
		} else {
			fmt.Println(utils.Red, "Invalid choice. Please choose '1', '2', or '3'.", utils.Reset)
		}
	}
}
