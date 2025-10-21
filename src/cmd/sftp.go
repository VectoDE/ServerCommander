package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"servercommander/src/services/config"
	"servercommander/src/utils"
)

func init() {
	RegisterCommand("sftp", "Perform SFTP file operations", sftpCommand)
}

func sftpCommand(args []string) error {
	if len(args) == 0 {
		return errors.New(utils.FormatUsageError("sftp <list|upload|download> <alias> [paths]"))
	}

	action := strings.ToLower(args[0])
	switch action {
	case "list":
		if err := ensureUsage(args[1:], 1, 2, "sftp list <alias> [remote-path]"); err != nil {
			return err
		}
		remotePath := "."
		if len(args) == 3 {
			remotePath = args[2]
		}
		return runSFTPBatch(args[1], []string{fmt.Sprintf("ls %s", remotePath)}, renderSFTPListing)
	case "upload":
		if err := ensureUsage(args[1:], 3, 3, "sftp upload <alias> <local> <remote>"); err != nil {
			return err
		}
		alias := args[1]
		local := args[2]
		remote := args[3]
		if strings.HasSuffix(remote, "/") {
			remote = path.Join(remote, filepath.Base(local))
		}
		return runSFTPBatch(alias, []string{fmt.Sprintf("put %s %s", local, remote)}, nil)
	case "download":
		if err := ensureUsage(args[1:], 3, 3, "sftp download <alias> <remote> <local>"); err != nil {
			return err
		}
		alias := args[1]
		remote := args[2]
		local := args[3]
		if strings.HasSuffix(local, string(filepath.Separator)) {
			local = filepath.Join(local, filepath.Base(remote))
		}
		return runSFTPBatch(alias, []string{fmt.Sprintf("get %s %s", remote, local)}, nil)
	default:
		return fmt.Errorf("unknown sftp action '%s'", action)
	}
}

func runSFTPBatch(alias string, commands []string, postProcess func(string) error) error {
	session, err := loadSession(alias)
	if err != nil {
		return err
	}
	if session.Protocol != config.ProtocolSFTP {
		return fmt.Errorf("session '%s' is not configured for SFTP", session.Alias)
	}

	password, err := promptPassword(session)
	if err != nil {
		return err
	}

	batch := strings.Join(commands, "\n") + "\n"
	args := buildSFTPArgs(session)
	cmd := exec.Command("sftp", args...)
	cmd.Stdin = strings.NewReader(batch)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if password != "" {
		fmt.Println(utils.Yellow, "Note: Enter the password when prompted by sftp.", utils.Reset)
	}

	if err := cmd.Run(); err != nil {
		output := stderr.String()
		if output == "" {
			output = stdout.String()
		}
		return fmt.Errorf("sftp command failed: %s", strings.TrimSpace(output))
	}

	if postProcess != nil {
		return postProcess(stdout.String())
	}

	fmt.Println(strings.TrimSpace(stdout.String()))
	return nil
}

func renderSFTPListing(output string) error {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) == 0 || output == "" {
		fmt.Println("(empty)")
		return nil
	}

	fmt.Printf("%s%-30s %-12s %-20s%s\n", utils.Cyan, "NAME", "SIZE", "MODIFIED", utils.Reset)
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 9 {
			name := strings.Join(fields[8:], " ")
			size := fields[4]
			date := strings.Join(fields[5:8], " ")
			fmt.Printf("%-30s %-12s %-20s\n", name, size, normaliseDate(date))
			continue
		}
		fmt.Println(line)
	}
	return nil
}

func buildSFTPArgs(session config.Session) []string {
	args := []string{"-b", "-"}
	args = append(args, "-P", fmt.Sprintf("%d", session.Port))
	if session.AuthMethod == config.AuthPrivateKey && session.KeyPath != "" {
		args = append(args, "-i", session.KeyPath)
	}
	args = append(args, fmt.Sprintf("%s@%s", session.Username, session.Host))
	return args
}

func normaliseDate(value string) string {
	// sftp output differs across implementations. Try to parse RFC3339 and
	// fall back to the raw string when parsing fails.
	if parsed, err := time.Parse(time.RFC3339, value); err == nil {
		return parsed.UTC().Format(time.RFC3339)
	}
	return value
}
