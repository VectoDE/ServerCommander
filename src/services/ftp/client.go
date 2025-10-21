package ftp

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/textproto"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"servercommander/src/services/config"
)

// Entry describes a file or directory returned by the FTP server.
type Entry struct {
	Name    string
	Size    int64
	ModTime time.Time
	IsDir   bool
	Raw     string
}

// Client implements a minimal FTP/FTPS client with passive mode support.
type Client struct {
	session   config.Session
	control   *textproto.Conn
	conn      net.Conn
	tlsConfig *tls.Config
}

// Connect establishes a control connection and authenticates the user.
func Connect(session config.Session, password string) (*Client, error) {
	address := net.JoinHostPort(session.Host, strconv.Itoa(session.Port))
	conn, err := net.DialTimeout("tcp", address, 10*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", address, err)
	}

	client := &Client{
		session: session,
		conn:    conn,
		control: textproto.NewConn(conn),
	}

	if _, _, err := client.read(220); err != nil {
		client.Close()
		return nil, err
	}

	if session.UseTLS {
		if err := client.startTLS(); err != nil {
			client.Close()
			return nil, err
		}
	}

	code, err := client.sendUser(session.Username)
	if err != nil {
		client.Close()
		return nil, err
	}
	if code != 230 {
		if err := client.sendPass(password); err != nil {
			client.Close()
			return nil, err
		}
	}

	return client, nil
}

// Close terminates the control connection.
func (c *Client) Close() error {
	if c.control != nil {
		c.control.Cmd("QUIT")
		c.control.Close()
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Upload stores a local file on the remote server.
func (c *Client) Upload(localPath, remotePath string) error {
	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open local file: %w", err)
	}
	defer file.Close()

	if err := c.ensureRemoteDir(path.Dir(remotePath)); err != nil {
		return err
	}

	dataConn, err := c.openDataConnection("STOR " + remotePath)
	if err != nil {
		return err
	}
	defer dataConn.Close()

	if _, err := io.Copy(dataConn, file); err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return c.readTransferResponse()
}

// Download retrieves a remote file and stores it locally.
func (c *Client) Download(remotePath, localPath string) error {
	dataConn, err := c.openDataConnection("RETR " + remotePath)
	if err != nil {
		return err
	}
	defer dataConn.Close()

	if err := os.MkdirAll(path.Dir(localPath), 0750); err != nil {
		return fmt.Errorf("failed to create local directories: %w", err)
	}

	file, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, dataConn); err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}

	return c.readTransferResponse()
}

// List returns a directory listing.
func (c *Client) List(path string) ([]Entry, error) {
	dataConn, err := c.openDataConnection("LIST " + path)
	if err != nil {
		return nil, err
	}
	defer dataConn.Close()

	scanner := bufio.NewScanner(dataConn)
	entries := []Entry{}
	for scanner.Scan() {
		raw := scanner.Text()
		entry := parseListLine(raw)
		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read directory listing: %w", err)
	}

	if err := c.readTransferResponse(); err != nil {
		return nil, err
	}

	return entries, nil
}

func (c *Client) startTLS() error {
	if err := c.control.PrintfLine("AUTH TLS"); err != nil {
		return err
	}
	if _, _, err := c.read(234); err != nil {
		return err
	}

	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	tlsConn := tls.Client(c.conn, tlsConfig)
	if err := tlsConn.Handshake(); err != nil {
		return err
	}

	c.conn = tlsConn
	c.control = textproto.NewConn(tlsConn)
	c.tlsConfig = tlsConfig
	return nil
}

func (c *Client) sendUser(username string) (int, error) {
	if err := c.control.PrintfLine("USER %s", username); err != nil {
		return 0, err
	}
	code, _, err := c.read(331, 230)
	return code, err
}

func (c *Client) sendPass(password string) error {
	if password == "" {
		return fmt.Errorf("password is required for FTP session '%s'", c.session.Alias)
	}
	if err := c.control.PrintfLine("PASS %s", password); err != nil {
		return err
	}
	_, _, err := c.read(230)
	return err
}

func (c *Client) openDataConnection(command string) (net.Conn, error) {
	host, port, err := c.enterPassiveMode()
	if err != nil {
		return nil, err
	}

	dataConn, err := net.DialTimeout("tcp", net.JoinHostPort(host, strconv.Itoa(port)), 10*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to open data connection: %w", err)
	}

	if c.tlsConfig != nil {
		dataConn = tls.Client(dataConn, c.tlsConfig)
		if err := dataConn.(*tls.Conn).Handshake(); err != nil {
			dataConn.Close()
			return nil, fmt.Errorf("failed to establish TLS data connection: %w", err)
		}
	}

	if err := c.control.PrintfLine(command); err != nil {
		dataConn.Close()
		return nil, err
	}

	if _, _, err := c.read(125, 150); err != nil {
		dataConn.Close()
		return nil, err
	}

	return dataConn, nil
}

func (c *Client) readTransferResponse() error {
	_, _, err := c.read(226, 250)
	return err
}

func (c *Client) enterPassiveMode() (string, int, error) {
	if err := c.control.PrintfLine("PASV"); err != nil {
		return "", 0, err
	}
	_, message, err := c.read(227)
	if err != nil {
		return "", 0, err
	}

	start := strings.Index(message, "(")
	end := strings.Index(message, ")")
	if start == -1 || end == -1 || end <= start+1 {
		return "", 0, fmt.Errorf("invalid PASV response: %s", message)
	}

	parts := strings.Split(message[start+1:end], ",")
	if len(parts) != 6 {
		return "", 0, fmt.Errorf("unexpected PASV response: %s", message)
	}

	host := strings.Join(parts[0:4], ".")
	p1, err := strconv.Atoi(parts[4])
	if err != nil {
		return "", 0, fmt.Errorf("invalid PASV port: %w", err)
	}
	p2, err := strconv.Atoi(parts[5])
	if err != nil {
		return "", 0, fmt.Errorf("invalid PASV port: %w", err)
	}

	return host, p1*256 + p2, nil
}

func (c *Client) ensureRemoteDir(dir string) error {
	if dir == "." || dir == "/" || dir == "" {
		return nil
	}

	segments := strings.Split(dir, "/")
	path := ""
	for _, segment := range segments {
		if segment == "" {
			continue
		}
		path += "/" + segment
		if err := c.control.PrintfLine("MKD %s", path); err != nil {
			return err
		}
		if _, _, err := c.read(257, 550); err != nil {
			return err
		}
	}
	return nil
}

func parseListLine(raw string) Entry {
	entry := Entry{Raw: raw}
	fields := strings.Fields(raw)
	if len(fields) >= 9 {
		entry.Name = strings.Join(fields[8:], " ")
		entry.IsDir = strings.HasPrefix(fields[0], "d")
		size, _ := strconv.ParseInt(fields[4], 10, 64)
		entry.Size = size
		dateParts := fields[5:8]
		parsedTime := parseListTime(dateParts)
		if !parsedTime.IsZero() {
			entry.ModTime = parsedTime
		}
	} else {
		entry.Name = raw
	}
	return entry
}

func parseListTime(parts []string) time.Time {
	if len(parts) != 3 {
		return time.Time{}
	}
	layout := "Jan 2 15:04"
	value := strings.Join(parts, " ")
	if strings.Contains(parts[2], ":") {
		t, err := time.Parse(layout, value)
		if err == nil {
			return t
		}
	} else {
		layout = "Jan 2 2006"
		t, err := time.Parse(layout, value)
		if err == nil {
			return t
		}
	}
	return time.Time{}
}

func (c *Client) read(expected ...int) (int, string, error) {
	if len(expected) == 0 {
		expected = []int{200}
	}

	code, message, err := c.control.ReadResponse(expected[0])
	if err == nil {
		return code, message, nil
	}

	if protoErr, ok := err.(*textproto.Error); ok {
		for _, exp := range expected {
			if protoErr.Code == exp {
				return protoErr.Code, protoErr.Msg, nil
			}
		}
	}

	return code, message, err
}
