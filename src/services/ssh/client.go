package ssh

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"servercommander/src/services/config"
)

// Client encapsulates metadata required to spawn SSH processes.
type Client struct {
	session  config.Session
	password string
}

// Connect prepares an SSH client for the provided session. No network
// connection is established at this point; instead we rely on the system's
// native ssh binary to handle the heavy lifting when commands are executed.
func Connect(session config.Session, password string, _ []byte) (*Client, error) {
	if session.Protocol != config.ProtocolSSH && session.Protocol != config.ProtocolSFTP {
		return nil, fmt.Errorf("protocol %s cannot be used with SSH", session.Protocol)
	}

	return &Client{session: session, password: password}, nil
}

// Close is a no-op kept for API compatibility with other services.
func (c *Client) Close() error {
	return nil
}

// InteractiveShell spawns the system ssh command and attaches STDIN/STDOUT to
// provide an interactive shell. Password authentication is delegated to the ssh
// binary which prompts the user as needed.
func (c *Client) InteractiveShell() error {
	args := c.buildBaseArgs()
	cmd := exec.Command("ssh", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Run executes a remote command via ssh and captures its combined output.
func (c *Client) Run(command string) (string, error) {
	args := append(c.buildBaseArgs(), command)
	cmd := exec.Command("ssh", args...)
	var buffer bytes.Buffer
	cmd.Stdout = &buffer
	cmd.Stderr = &buffer

	if err := cmd.Run(); err != nil {
		return buffer.String(), fmt.Errorf("remote command failed: %w", err)
	}

	return buffer.String(), nil
}

// Raw exposes the configuration to enable reuse in other packages. Since the
// implementation relies on external processes the method returns nil and is
// retained purely for API compatibility.
func (c *Client) Raw() interface{} {
	return nil
}

func (c *Client) buildBaseArgs() []string {
	args := []string{"-p", strconv.Itoa(c.session.Port), fmt.Sprintf("%s@%s", c.session.Username, c.session.Host)}
	if c.session.AuthMethod == config.AuthPrivateKey && c.session.KeyPath != "" {
		args = append([]string{"-i", c.session.KeyPath}, args...)
	}
	return args
}
