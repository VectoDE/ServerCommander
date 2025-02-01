package cmd

import (
	"fmt"
	"os"

	"servercommander/src/services"
	"servercommander/src/ui"
	"servercommander/src/utils"
)

func ExitCommand() {
	services.LogToFile("Program exited")

	ui.GoodbyeBanner()

	fmt.Println(utils.Red, "Exiting the program...", utils.Reset)

	os.Exit(0)
}
