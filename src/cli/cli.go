package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"servercommander/src/cli/cmd"
	"servercommander/src/ui"
	"servercommander/src/utils"
)

func StartCLI() {
	ui.WelcomeCLIBanner()

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
		case "htop":
			cmd.HtopCommand()
		case "exit":
			cmd.ExitCommand()
			return
		default:
			fmt.Println(utils.Red, "Unknown command. Type 'help' to learn more.", utils.Reset)
		}
	}
}
