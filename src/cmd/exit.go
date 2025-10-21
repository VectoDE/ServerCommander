package cmd

import (
	"fmt"
	"os"

	"servercommander/src/ui"
	"servercommander/src/utils"
)

func exitCommand(args []string) error {
	ui.GoodbyeBanner()
	fmt.Println(utils.Red, "Exiting the program...", utils.Reset)
	os.Exit(0)
	return nil
}
