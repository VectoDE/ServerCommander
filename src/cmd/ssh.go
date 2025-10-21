package cmd

import (
	"errors"
	"fmt"
	"strings"

	"servercommander/src/services/config"
	sshservice "servercommander/src/services/ssh"
	"servercommander/src/utils"
)

func init() {
	RegisterCommand("ssh", "Execute SSH operations", sshCommand)
	RegisterCommand("connect", "Open an interactive SSH session", connectCommand)
}

func sshCommand(args []string) error {
	if len(args) == 0 {
		return errors.New(utils.FormatUsageError("ssh <connect|exec> <alias> [command]"))
	}

	action := strings.ToLower(args[0])
	switch action {
	case "connect":
		if err := ensureUsage(args[1:], 1, 1, "ssh connect <alias>"); err != nil {
			return err
		}
		session, err := loadSession(args[1])
		if err != nil {
			return err
		}
		if session.Protocol != config.ProtocolSSH {
			return fmt.Errorf("session '%s' is not an SSH session", session.Alias)
		}
		return startInteractiveSSH(session)
	case "exec":
		if err := ensureUsage(args[1:], 2, -1, "ssh exec <alias> <command>"); err != nil {
			return err
		}
		session, err := loadSession(args[1])
		if err != nil {
			return err
		}
		if session.Protocol != config.ProtocolSSH {
			return fmt.Errorf("session '%s' is not an SSH session", session.Alias)
		}
		command := strings.Join(args[2:], " ")
		return executeRemoteCommand(session, command)
	default:
		return fmt.Errorf("unknown ssh action '%s'", action)
	}
}

func connectCommand(args []string) error {
	if err := ensureUsage(args, 1, 1, "connect <alias>"); err != nil {
		return err
	}
	session, err := loadSession(args[0])
	if err != nil {
		return err
	}

	switch session.Protocol {
	case config.ProtocolSSH:
		return startInteractiveSSH(session)
	case config.ProtocolSFTP:
		return fmt.Errorf("session '%s' is configured for SFTP. Use the sftp commands for file operations", session.Alias)
	case config.ProtocolFTP:
		return fmt.Errorf("session '%s' is configured for FTP. Use the ftp commands for file operations", session.Alias)
	default:
		return fmt.Errorf("protocol '%s' is not supported", session.Protocol)
	}
}

func startInteractiveSSH(session config.Session) error {
	password, err := promptPassword(session)
	if err != nil {
		return err
	}

	client, _, err := dialSSHWithRetry(session, password)
	if err != nil {
		return err
	}
	defer client.Close()

	if err := client.InteractiveShell(); err != nil {
		return err
	}

	return nil
}

func executeRemoteCommand(session config.Session, command string) error {
	password, err := promptPassword(session)
	if err != nil {
		return err
	}

	client, _, err := dialSSHWithRetry(session, password)
	if err != nil {
		return err
	}
	defer client.Close()

	output, err := client.Run(command)
	if output != "" {
		fmt.Println(output)
	}
	if err != nil {
		return err
	}

	return nil
}

func dialSSHWithRetry(session config.Session, password string) (*sshservice.Client, []byte, error) {
	client, err := sshservice.Connect(session, password, nil)
	if err != nil {
		return nil, nil, err
	}
	return client, nil, nil
}

func promptPassword(session config.Session) (string, error) {
	if !session.RequiresPass {
		return "", nil
	}
	return utils.PromptPassword(fmt.Sprintf("Password for %s@%s", session.Username, session.Host))
}
