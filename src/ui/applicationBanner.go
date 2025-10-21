package ui

import (
	"fmt"
	"servercommander/src/utils"
)

func ApplicationBanner() {
	fmt.Println(utils.Cyan, "==============================")
	fmt.Println(utils.Green, "     Server Commander v1.0.1")
	fmt.Println(utils.Cyan, "==============================", utils.Reset)
}
