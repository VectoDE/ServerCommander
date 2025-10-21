package cmd

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"servercommander/src/services"
	"servercommander/src/utils"
)

// CommandHandler represents the business logic executed for a command.
type CommandHandler func(args []string) error

// commandDescriptor keeps metadata for a registered command.
type commandDescriptor struct {
	Name        string
	Description string
	Handler     CommandHandler
}

var commandRegistry = map[string]commandDescriptor{}

// RegisterCommand adds a new command to the registry. It will panic if the
// command name collides with an existing entry because this indicates a
// developer error that should be caught during tests.
func RegisterCommand(name, description string, handler CommandHandler) {
	if name == "" {
		panic("command name cannot be empty")
	}

	key := strings.ToLower(name)
	if _, exists := commandRegistry[key]; exists {
		panic(fmt.Sprintf("command %s already registered", name))
	}

	commandRegistry[key] = commandDescriptor{
		Name:        key,
		Description: description,
		Handler:     handler,
	}
}

// Execute resolves the command within the registry and runs the attached
// handler. It returns an error if the command does not exist or if the
// handler reports an error.
func Execute(input string) error {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return nil
	}

	key := strings.ToLower(parts[0])
	descriptor, exists := commandRegistry[key]
	if !exists {
		return fmt.Errorf("unknown command '%s'. Type 'help' to list available commands", parts[0])
	}

	if err := descriptor.Handler(parts[1:]); err != nil {
		services.LogToFile(fmt.Sprintf("command '%s' failed: %v", descriptor.Name, err))
		return err
	}

	services.LogToFile(fmt.Sprintf("command '%s' executed successfully", descriptor.Name))
	return nil
}

// ListCommands returns a deterministic, alphabetically sorted slice of
// descriptors. This is primarily used by the help command but also enables
// other commands to query the available functionality.
func ListCommands() []commandDescriptor {
	commands := make([]commandDescriptor, 0, len(commandRegistry))
	for _, descriptor := range commandRegistry {
		commands = append(commands, descriptor)
	}

	sort.Slice(commands, func(i, j int) bool {
		return commands[i].Name < commands[j].Name
	})

	return commands
}

// init registers the built-in commands to guarantee they are always
// available even when the application is extended with plugins.
func init() {
	RegisterCommand("help", "Shows this help", helpCommand)
	RegisterCommand("clear", "Clears the console", clearCommand)
	RegisterCommand("exit", "Exits the program", exitCommand)
}

// ensureUsage enforces the expected number of positional arguments for a
// command. It simplifies handler implementations by providing a consistent
// user-facing error message. Some commands use advanced parsing and therefore
// do not rely on this helper.
func ensureUsage(args []string, min, max int, usage string) error {
	if len(args) < min || (max >= 0 && len(args) > max) {
		return errors.New(utils.FormatUsageError(usage))
	}
	return nil
}
