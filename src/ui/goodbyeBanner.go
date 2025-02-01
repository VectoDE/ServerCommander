package ui

import (
	"fmt"
	"servercommander/src/utils"
	"time"
)

func GoodbyeBanner() {
	fmt.Println(utils.Yellow, "Goodbye! The program will close in 3 seconds...", utils.Reset)
	time.Sleep(1 * time.Second)

	for i := 3; i > 0; i-- {
		fmt.Printf("%sClosing in %d seconds...\r", utils.Cyan, i)
		time.Sleep(1 * time.Second)
	}

	time.Sleep(2 * time.Second)
}
