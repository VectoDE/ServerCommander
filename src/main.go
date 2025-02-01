package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/lxn/walk"
	"servercommander/src/cmd"
	"servercommander/src/ui"
	"servercommander/src/utils"
)

func main() {
	ui.ApplicationBanner()
	setAppIcon("src/assets/icon.ico") // Icon setzen

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

// setAppIcon setzt das Fenster-Icon (für Windows GUI-Apps)
func setAppIcon(iconPath string) {
	icon, err := walk.NewIconFromFile(iconPath)
	if err != nil {
		fmt.Println("⚠️ Warning: Could not load icon:", err)
		return
	}

	// Dummy-Fenster erstellen, um das Icon zu setzen
	mw, _ := walk.NewMainWindow()
	mw.SetIcon(icon)

	fmt.Println("✅ Icon set successfully!")
}