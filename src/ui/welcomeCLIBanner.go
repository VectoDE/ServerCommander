package ui

import (
	"fmt"

	"servercommander/src/utils"

	"github.com/pterm/pterm"
)

func WelcomeCLIBanner() {
	welcomeText := fmt.Sprintf(
		"Welcome to the Command Line Application!\n\n"+
			"%sHave fun using the CLI!\nType 'help' to learn more.%s",
		utils.Yellow, utils.Reset,
	)

	pterm.DefaultCenter.Println(pterm.DefaultBox.WithTitle(" Welcome! ").Sprint(welcomeText))

	pterm.DefaultCenter.Println("---------------------------------------------------")
}
