package utils

import "fmt"

// FormatUsageError standardises the error message returned when the user
// provides invalid arguments for a command. Keeping this logic in a single
// location ensures consistency across the CLI.
func FormatUsageError(usage string) string {
	if usage == "" {
		return "invalid command usage"
	}
	return fmt.Sprintf("invalid command usage. expected: %s", usage)
}
