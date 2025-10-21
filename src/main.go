package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"servercommander/src/cmd"
	"servercommander/src/ui"
	"servercommander/src/utils"
)

func main() {
	ui.ApplicationBanner()
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(utils.Yellow, "\n>> ", utils.Reset)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(utils.Red, "Failed to read input:", err, utils.Reset)
			continue
		}

		command := strings.TrimSpace(input)
		if command == "" {
			continue
		}

		if err := cmd.Execute(command); err != nil {
			fmt.Println(utils.Red, err, utils.Reset)
		}
	}
}
