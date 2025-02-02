package ui

import (
	"fmt"
	"servercommander/src/services"

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

func ApplicationBanner() {
	version := "v1.0"

	s, _ := pterm.DefaultBigText.WithLetters(putils.LettersFromString("Server")).Srender()
	c, _ := pterm.DefaultBigText.WithLetters(putils.LettersFromString("Commander")).Srender()

	pterm.DefaultCenter.Println(s)
	pterm.DefaultCenter.Println(c)

	pterm.DefaultCenter.WithCenterEachLineSeparately().Println(fmt.Sprintf("Version: %s\nDeveloper: Tim Hauke (VectoDE)\n© 2025 VectoDE", version))

	services.LogToFile("Displayed Application Banner", "INFO")
}
