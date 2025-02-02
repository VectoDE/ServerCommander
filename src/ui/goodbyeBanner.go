package ui

import (
	"fmt"
	"servercommander/src/services"
	"servercommander/src/utils"
	"time"

	"github.com/pterm/pterm"
)

func GoodbyeBanner(duration int) {
	pterm.DefaultCenter.Println(fmt.Sprintf("%sGoodbye! The program will close in %d seconds...", utils.Yellow, duration))

	multi := pterm.DefaultMultiPrinter

	pb1, _ := pterm.DefaultProgressbar.WithTotal(duration).WithWriter(multi.NewWriter()).Start("Shutting down...")

	multi.Start()

	for i := duration; i > 0; i-- {
		pb1.Increment()

		time.Sleep(1 * time.Second)
	}

	multi.Stop()

	services.LogToFile("Displayed Goodbye Banner", "INFO")
}
