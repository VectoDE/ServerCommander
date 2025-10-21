package cmd

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"time"

	"servercommander/src/services/config"
	ftpservice "servercommander/src/services/ftp"
	"servercommander/src/utils"
)

func init() {
	RegisterCommand("ftp", "Perform FTP/FTPS file operations", ftpCommand)
}

func ftpCommand(args []string) error {
	if len(args) == 0 {
		return errors.New(utils.FormatUsageError("ftp <list|upload|download> <alias> [paths]"))
	}

	action := strings.ToLower(args[0])
	switch action {
	case "list":
		if err := ensureUsage(args[1:], 1, 2, "ftp list <alias> [remote-path]"); err != nil {
			return err
		}
		remotePath := "."
		if len(args) == 3 {
			remotePath = args[2]
		}
		return withFTPClient(args[1], func(client *ftpservice.Client) error {
			return renderFTPListing(client, remotePath)
		})
	case "upload":
		if err := ensureUsage(args[1:], 3, 3, "ftp upload <alias> <local> <remote>"); err != nil {
			return err
		}
		local := args[2]
		remote := args[3]
		return withFTPClient(args[1], func(client *ftpservice.Client) error {
			if strings.HasSuffix(remote, "/") {
				remote = path.Join(remote, filepath.Base(local))
			}
			return client.Upload(local, remote)
		})
	case "download":
		if err := ensureUsage(args[1:], 3, 3, "ftp download <alias> <remote> <local>"); err != nil {
			return err
		}
		remote := args[2]
		local := args[3]
		return withFTPClient(args[1], func(client *ftpservice.Client) error {
			if strings.HasSuffix(local, string(filepath.Separator)) {
				local = filepath.Join(local, filepath.Base(remote))
			}
			return client.Download(remote, local)
		})
	default:
		return fmt.Errorf("unknown ftp action '%s'", action)
	}
}

func withFTPClient(alias string, fn func(*ftpservice.Client) error) error {
	session, err := loadSession(alias)
	if err != nil {
		return err
	}
	if session.Protocol != config.ProtocolFTP {
		return fmt.Errorf("session '%s' is not configured for FTP", session.Alias)
	}

	password, err := promptPassword(session)
	if err != nil {
		return err
	}

	client, err := ftpservice.Connect(session, password)
	if err != nil {
		return err
	}
	defer client.Close()

	return fn(client)
}

func renderFTPListing(client *ftpservice.Client, path string) error {
	entries, err := client.List(path)
	if err != nil {
		return err
	}

	fmt.Printf("%s%-30s %-12s %-20s%s\n", utils.Cyan, "NAME", "SIZE", "MODIFIED", utils.Reset)
	for _, entry := range entries {
		name := entry.Name
		if entry.IsDir {
			name += "/"
		}
		modTime := ""
		if !entry.ModTime.IsZero() {
			modTime = entry.ModTime.UTC().Format(time.RFC3339)
		}
		fmt.Printf("%-30s %-12d %-20s\n", name, entry.Size, modTime)
	}
	return nil
}
