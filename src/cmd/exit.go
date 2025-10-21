package cmd

import (
	"fmt"
	"os"

	"servercommander/src/console"
	"servercommander/src/utils"
)

func exitCommand(args []string) error {
	console.GoodbyeBanner()
	fmt.Println(utils.Red, "Exiting the program...", utils.Reset)
	os.Exit(0)
	return nil
}
