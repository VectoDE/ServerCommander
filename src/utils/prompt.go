package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

var (
	readerMu sync.RWMutex
	reader   = bufio.NewReader(os.Stdin)
)

// SetPromptReader allows callers to override the input source used by the
// prompt helpers. Passing nil keeps the existing reader untouched.
func SetPromptReader(r io.Reader) {
	if r == nil {
		return
	}

	readerMu.Lock()
	reader = bufio.NewReader(r)
	readerMu.Unlock()
}

func readLine() (string, error) {
	readerMu.RLock()
	active := reader
	readerMu.RUnlock()
	return active.ReadString('\n')
}

// Prompt requests free-form input from the user. The default value is used when
// the user submits an empty string.
func Prompt(question, defaultValue string) (string, error) {
	fmt.Printf("%s%s%s", Cyan, question, Reset)
	if defaultValue != "" {
		fmt.Printf(" [%s]", defaultValue)
	}
	fmt.Print(": ")

	input, err := readLine()
	if err != nil {
		return "", err
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return defaultValue, nil
	}

	return input, nil
}

// PromptPassword reads sensitive input from the terminal without echoing the
// characters. It returns an error when the stdin file descriptor is not a
// terminal.
func PromptPassword(question string) (string, error) {
	fmt.Printf("%s%s (input hidden not supported): %s", Cyan, question, Reset)
	value, err := readLine()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(value), nil
}

// PromptBool converts user input into a boolean. Accepted inputs are "y",
// "yes", "n", "no" (case insensitive). The default value is returned when the
// user submits an empty string.
func PromptBool(question string, defaultValue bool) (bool, error) {
	def := "n"
	if defaultValue {
		def = "y"
	}

	answer, err := Prompt(fmt.Sprintf("%s (y/n)", question), def)
	if err != nil {
		return false, err
	}

	switch strings.ToLower(answer) {
	case "y", "yes":
		return true, nil
	case "n", "no":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean value: %s", answer)
	}
}
