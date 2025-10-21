package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"servercommander/src/services/config"
	"servercommander/src/utils"
)

func init() {
	RegisterCommand("session", "Manage saved server sessions", sessionCommand)
}

func sessionCommand(args []string) error {
	if len(args) == 0 {
		return errors.New(utils.FormatUsageError("session <add|list|remove|show> [alias]"))
	}

	action := strings.ToLower(args[0])
	switch action {
	case "add":
		if err := ensureUsage(args[1:], 1, 1, "session add <alias>"); err != nil {
			return err
		}
		return sessionAdd(args[1])
	case "list":
		return sessionList()
	case "remove":
		if err := ensureUsage(args[1:], 1, 1, "session remove <alias>"); err != nil {
			return err
		}
		return sessionRemove(args[1])
	case "show":
		if err := ensureUsage(args[1:], 1, 1, "session show <alias>"); err != nil {
			return err
		}
		return sessionShow(args[1])
	default:
		return fmt.Errorf("unknown session action '%s'", action)
	}
}

func sessionAdd(alias string) error {
	store, err := config.LoadSessions()
	if err != nil {
		return err
	}

	existing, exists := store.Get(alias)
	session, err := promptSessionDetails(alias, existing, exists)
	if err != nil {
		return err
	}

	store.Upsert(session)
	if err := store.Save(); err != nil {
		return err
	}

	fmt.Printf("%sSession '%s' saved.%s\n", utils.Green, session.Alias, utils.Reset)
	return nil
}

func sessionList() error {
	store, err := config.LoadSessions()
	if err != nil {
		return err
	}

	sessions := store.List()
	if len(sessions) == 0 {
		fmt.Println(utils.Yellow, "No sessions stored.", utils.Reset)
		return nil
	}

	fmt.Printf("%s%-15s %-8s %-25s %-10s %-10s%s\n", utils.Cyan, "Alias", "Protocol", "Host", "User", "Auth", utils.Reset)
	for _, session := range sessions {
		fmt.Printf("%-15s %-8s %-25s %-10s %-10s\n",
			session.Alias,
			session.Protocol,
			fmt.Sprintf("%s:%d", session.Host, session.Port),
			session.Username,
			session.AuthMethod,
		)
	}

	return nil
}

func sessionRemove(alias string) error {
	store, err := config.LoadSessions()
	if err != nil {
		return err
	}

	if err := store.Remove(alias); err != nil {
		return err
	}

	if err := store.Save(); err != nil {
		return err
	}

	fmt.Printf("%sSession '%s' removed.%s\n", utils.Green, strings.ToLower(alias), utils.Reset)
	return nil
}

func sessionShow(alias string) error {
	session, err := loadSession(alias)
	if err != nil {
		return err
	}

	fmt.Printf("%sAlias:%s        %s\n", utils.Blue, utils.Reset, session.Alias)
	fmt.Printf("%sProtocol:%s     %s\n", utils.Blue, utils.Reset, session.Protocol)
	fmt.Printf("%sHost:%s         %s:%d\n", utils.Blue, utils.Reset, session.Host, session.Port)
	fmt.Printf("%sUsername:%s     %s\n", utils.Blue, utils.Reset, session.Username)
	fmt.Printf("%sAuth Method:%s  %s\n", utils.Blue, utils.Reset, session.AuthMethod)
	if session.KeyPath != "" {
		fmt.Printf("%sKey Path:%s     %s\n", utils.Blue, utils.Reset, session.KeyPath)
	}
	if session.Protocol == config.ProtocolFTP {
		fmt.Printf("%sTLS Enabled:%s %t\n", utils.Blue, utils.Reset, session.UseTLS)
	}
	if session.Description != "" {
		fmt.Printf("%sDescription:%s %s\n", utils.Blue, utils.Reset, session.Description)
	}
	fmt.Printf("%sRequires Pass:%s %t\n", utils.Blue, utils.Reset, session.RequiresPass)
	fmt.Printf("%sCreated:%s      %s\n", utils.Blue, utils.Reset, session.CreatedAt.Format(time.RFC3339))
	fmt.Printf("%sUpdated:%s      %s\n", utils.Blue, utils.Reset, session.UpdatedAt.Format(time.RFC3339))
	return nil
}

func promptSessionDetails(alias string, existing config.Session, exists bool) (config.Session, error) {
	protocolDefault := string(config.ProtocolSSH)
	if exists {
		protocolDefault = string(existing.Protocol)
	}

	protocolInput, err := utils.Prompt("Protocol (ssh/sftp/ftp)", protocolDefault)
	if err != nil {
		return config.Session{}, err
	}

	protocol := config.Protocol(strings.ToLower(protocolInput))
	switch protocol {
	case config.ProtocolSSH, config.ProtocolSFTP, config.ProtocolFTP:
	default:
		return config.Session{}, fmt.Errorf("unsupported protocol '%s'", protocol)
	}

	host, err := utils.Prompt("Host", existing.Host)
	if err != nil {
		return config.Session{}, err
	}
	if host == "" {
		return config.Session{}, errors.New("host cannot be empty")
	}

	portDefault := defaultPort(protocol)
	if exists && existing.Port != 0 {
		portDefault = existing.Port
	}

	portInput, err := utils.Prompt("Port", strconv.Itoa(portDefault))
	if err != nil {
		return config.Session{}, err
	}

	port, err := strconv.Atoi(portInput)
	if err != nil || port <= 0 {
		return config.Session{}, fmt.Errorf("invalid port: %s", portInput)
	}

	username, err := utils.Prompt("Username", existing.Username)
	if err != nil {
		return config.Session{}, err
	}
	if username == "" {
		return config.Session{}, errors.New("username cannot be empty")
	}

	authDefault := string(config.AuthPassword)
	if exists {
		authDefault = string(existing.AuthMethod)
	}

	authMethod := config.AuthMethod(authDefault)
	if protocol != config.ProtocolFTP {
		authInput := authDefault
		for {
			authInput, err = utils.Prompt("Authentication method (password/private_key)", authInput)
			if err != nil {
				return config.Session{}, err
			}

			authMethod = config.AuthMethod(strings.ToLower(authInput))
			if authMethod == config.AuthPassword || authMethod == config.AuthPrivateKey {
				break
			}

			fmt.Printf("%sUnsupported value. Enter 'password' to supply the password when connecting or 'private_key' to use a key file.%s\n", utils.Yellow, utils.Reset)
		}
	} else {
		authMethod = config.AuthPassword
	}

	requiresPass := authMethod == config.AuthPassword

	keyPath := existing.KeyPath
	if authMethod == config.AuthPrivateKey {
		keyPath, err = utils.Prompt("Private key path", existing.KeyPath)
		if err != nil {
			return config.Session{}, err
		}
		if keyPath == "" {
			return config.Session{}, errors.New("private key path cannot be empty when using key authentication")
		}
	} else {
		keyPath = ""
	}

	description, err := utils.Prompt("Description", existing.Description)
	if err != nil {
		return config.Session{}, err
	}

	useTLS := existing.UseTLS
	if protocol == config.ProtocolFTP {
		useTLS, err = utils.PromptBool("Use explicit TLS", existing.UseTLS)
		if err != nil {
			return config.Session{}, err
		}
	}

	return config.Session{
		Alias:        alias,
		Protocol:     protocol,
		Host:         host,
		Port:         port,
		Username:     username,
		AuthMethod:   authMethod,
		KeyPath:      keyPath,
		UseTLS:       useTLS,
		Description:  description,
		RequiresPass: requiresPass,
	}, nil
}

func defaultPort(protocol config.Protocol) int {
	switch protocol {
	case config.ProtocolSSH, config.ProtocolSFTP:
		return 22
	case config.ProtocolFTP:
		return 21
	default:
		return 0
	}
}

func loadSession(alias string) (config.Session, error) {
	store, err := config.LoadSessions()
	if err != nil {
		return config.Session{}, err
	}

	session, ok := store.Get(alias)
	if !ok {
		return config.Session{}, fmt.Errorf("session '%s' not found", alias)
	}

	return session, nil
}
