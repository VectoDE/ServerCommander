package cmd

import (
	"fmt"

	"servercommander/src/services"
	"servercommander/src/utils"
)

func HelpCommand() {
	fmt.Println(utils.Green, "Available commands:")
	fmt.Println(utils.Blue, "  help  ", utils.White, "- Shows this help")
	fmt.Println(utils.Blue, "  clear ", utils.White, "- Clears the console")
	fmt.Println(utils.Blue, "  exit  ", utils.White, "- Exits the program")

	services.LogToFile("Help command executed")
}
