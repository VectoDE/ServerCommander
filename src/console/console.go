package console

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"servercommander/src/utils"
)

// Run starts the interactive console loop using the provided executor to
// handle individual command lines.
func Run(executor func(string) error) error {
	ApplicationBanner()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("%s>>>%s ", utils.Cyan, utils.Reset)

		if !scanner.Scan() {
			fmt.Println()
			return scanner.Err()
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if err := executor(line); err != nil {
			fmt.Println(utils.Red, err.Error(), utils.Reset)
		}
	}
}

// ApplicationBanner prints the program banner in the console.
func ApplicationBanner() {
	fmt.Println(utils.Cyan, "==============================")
	fmt.Println(utils.Green, "     Server Commander v1.0.1")
	fmt.Println(utils.Cyan, "==============================", utils.Reset)
}

// GoodbyeBanner displays a small exit animation before terminating.
func GoodbyeBanner() {
	fmt.Println(utils.Yellow, "Goodbye! The program will close in 3 seconds...", utils.Reset)
	time.Sleep(time.Second)

	for i := 3; i > 0; i-- {
		fmt.Printf("%sClosing in %d seconds...\r", utils.Cyan, i)
		time.Sleep(time.Second)
	}

	time.Sleep(2 * time.Second)
}

// ClearConsole tries to clear the terminal using ANSI escape sequences. It
// returns true when handled internally so callers can skip external commands.
func ClearConsole() bool {
	if runtime.GOOS == "windows" {
		return false
	}

	fmt.Print("\033[H\033[2J")
	return true
}
