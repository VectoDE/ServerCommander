package cmd

import (
	"fmt"

	"servercommander/src/utils"
)

func helpCommand(args []string) error {
	commands := ListCommands()
	fmt.Println(utils.Green, "Available commands:")
	for _, descriptor := range commands {
		fmt.Printf("%s  %-8s%s - %s\n", utils.Blue, descriptor.Name, utils.Reset, descriptor.Description)
	}
	return nil
}
