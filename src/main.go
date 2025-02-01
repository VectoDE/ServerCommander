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
		input, _ := reader.ReadString('\n')
		command := strings.TrimSpace(input)

		switch command {
		case "help":
			cmd.HelpCommand()
		case "clear":
			cmd.ClearConsole()
		case "exit":
			cmd.ExitCommand()
		default:
			fmt.Println(utils.Red, "Unknown command. Type 'help'.", utils.Reset)
		}
	}
}
