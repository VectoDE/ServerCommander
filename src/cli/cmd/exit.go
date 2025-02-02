package cmd

import (
	"fmt"
	"os"

	"servercommander/src/services"
	"servercommander/src/ui"
	"servercommander/src/utils"
)

func ExitCommand() {
	services.LogToFile("Program exited", "INFO")

	ui.GoodbyeBanner(3)

	fmt.Println(utils.Red, "Exiting the program...", utils.Reset)

	os.Exit(0)
}
