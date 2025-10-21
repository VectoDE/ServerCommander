package ui

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"sync"

	"servercommander/src/utils"
)

type commandExecutor func(string) error

type eventPayload struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

type ConsoleServer struct {
	executor  commandExecutor
	listener  net.Listener
	clients   map[int]chan eventPayload
	clientsMu sync.Mutex
	history   []eventPayload
	historyMu sync.Mutex
	errCh     chan error

	stdoutPipe *os.File
	stderrPipe *os.File
	origStdout *os.File
	origStderr *os.File
}

var (
	consoleInstance *ConsoleServer
	consoleMu       sync.RWMutex
)

const (
	controlClearEvent = "clear"
	appendEvent       = "append"
)

var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;]*[A-Za-z]`)

func RunStandaloneConsole(executor func(string) error) error {
	server, err := newConsoleServer(executor)
	if err != nil {
		return err
	}

	consoleMu.Lock()
	consoleInstance = server
	consoleMu.Unlock()

	return server.run()
}

func ClearConsole() bool {
	consoleMu.RLock()
	server := consoleInstance
	consoleMu.RUnlock()
	if server == nil {
		return false
	}
	server.clear()
	return true
}

func newConsoleServer(executor commandExecutor) (*ConsoleServer, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("failed to allocate listener: %w", err)
	}

	server := &ConsoleServer{
		executor: executor,
		listener: listener,
		clients:  make(map[int]chan eventPayload),
		errCh:    make(chan error, 1),
	}

	if err := server.redirectStandardStreams(); err != nil {
		listener.Close()
		return nil, err
	}

	return server, nil
}

func (c *ConsoleServer) run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", c.handleIndex)
	mux.HandleFunc("/execute", c.handleExecute)
	mux.HandleFunc("/events", c.handleEvents)

	go func() {
		c.errCh <- http.Serve(c.listener, mux)
	}()

	url := "http://" + c.listener.Addr().String() + "/"
	if err := openBrowser(url); err != nil {
		fmt.Println(utils.Yellow, "Unable to automatically open the interface. Visit", url, "in your browser.", utils.Reset)
	} else {
		fmt.Println(utils.Green, "ServerCommander interface started at", url, utils.Reset)
	}

	ApplicationBanner()

	err := <-c.errCh
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (c *ConsoleServer) redirectStandardStreams() error {
	stdoutReader, stdoutWriter, err := os.Pipe()
	if err != nil {
		return fmt.Errorf("failed to redirect stdout: %w", err)
	}

	stderrReader, stderrWriter, err := os.Pipe()
	if err != nil {
		stdoutReader.Close()
		stdoutWriter.Close()
		return fmt.Errorf("failed to redirect stderr: %w", err)
	}

	c.stdoutPipe = stdoutWriter
	c.stderrPipe = stderrWriter
	c.origStdout = os.Stdout
	c.origStderr = os.Stderr

	os.Stdout = stdoutWriter
	os.Stderr = stderrWriter

	go c.readPipe(stdoutReader, c.origStdout)
	go c.readPipe(stderrReader, c.origStderr)

	return nil
}

func (c *ConsoleServer) readPipe(pipe *os.File, mirror *os.File) {
	reader := bufio.NewReader(pipe)
	buffer := make([]byte, 1024)
	var pending strings.Builder

	for {
		n, err := reader.Read(buffer)
		if n > 0 {
			chunk := buffer[:n]
			if mirror != nil {
				_, _ = mirror.Write(chunk)
			}

			pending.Write(chunk)
			emit, remainder := splitANSISequences(pending.String())
			if emit != "" {
				c.broadcastAppend(emit)
			}

			pending.Reset()
			if remainder != "" {
				pending.WriteString(remainder)
			}
		}

		if err != nil {
			if pending.Len() > 0 {
				c.broadcastAppend(pending.String())
			}
			return
		}
	}
}

func splitANSISequences(input string) (string, string) {
	last := strings.LastIndexByte(input, 0x1b)
	if last == -1 {
		return input, ""
	}

	fragment := input[last:]
	if fragment == "" {
		return input, ""
	}

	if ansiSequenceComplete(fragment) {
		return input, ""
	}

	return input[:last], fragment
}

func ansiSequenceComplete(fragment string) bool {
	if len(fragment) < 2 {
		return false
	}

	if fragment[0] != 0x1b {
		return true
	}

	if fragment[1] != '[' {
		return true
	}

	for i := 2; i < len(fragment); i++ {
		b := fragment[i]
		if (b >= '0' && b <= '9') || b == ';' {
			continue
		}
		if b >= '@' && b <= '~' {
			return true
		}
		return false
	}

	return false
}

func (c *ConsoleServer) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, consoleHTML)
}

func (c *ConsoleServer) handleExecute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload struct {
		Command string `json:"command"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	command := strings.TrimSpace(payload.Command)
	if command == "" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	c.broadcastAppend(">> " + command + "\n")
	go c.executeCommand(command)
	w.WriteHeader(http.StatusAccepted)
}

func (c *ConsoleServer) executeCommand(command string) {
	if c.executor == nil {
		return
	}

	if err := c.executor(command); err != nil {
		fmt.Println(utils.Red, err.Error(), utils.Reset)
	}
}

func (c *ConsoleServer) handleEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	clientChan := make(chan eventPayload, 16)
	clientID := c.addClient(clientChan)
	defer c.removeClient(clientID)

	history := c.snapshotHistory()
	for _, event := range history {
		if err := writeSSE(w, event); err != nil {
			return
		}
	}
	flusher.Flush()

	for {
		select {
		case <-r.Context().Done():
			return
		case event := <-clientChan:
			if err := writeSSE(w, event); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}

func (c *ConsoleServer) addClient(ch chan eventPayload) int {
	c.clientsMu.Lock()
	defer c.clientsMu.Unlock()

	id := len(c.clients) + 1
	c.clients[id] = ch
	return id
}

func (c *ConsoleServer) removeClient(id int) {
	c.clientsMu.Lock()
	defer c.clientsMu.Unlock()

	if ch, exists := c.clients[id]; exists {
		close(ch)
		delete(c.clients, id)
	}
}

func (c *ConsoleServer) snapshotHistory() []eventPayload {
	c.historyMu.Lock()
	defer c.historyMu.Unlock()

	snapshot := make([]eventPayload, len(c.history))
	copy(snapshot, c.history)
	return snapshot
}

func (c *ConsoleServer) broadcastAppend(text string) {
	cleaned := sanitizeText(text)
	if cleaned == "" {
		return
	}

	event := eventPayload{Type: appendEvent, Text: cleaned}

	c.historyMu.Lock()
	c.history = append(c.history, event)
	c.historyMu.Unlock()

	c.dispatch(event)
}

func (c *ConsoleServer) clear() {
	c.historyMu.Lock()
	c.history = nil
	c.historyMu.Unlock()

	c.dispatch(eventPayload{Type: controlClearEvent})
}

func (c *ConsoleServer) dispatch(event eventPayload) {
	c.clientsMu.Lock()
	defer c.clientsMu.Unlock()

	for id, ch := range c.clients {
		select {
		case ch <- event:
		default:
			go func(target chan eventPayload, identifier int) {
				target <- event
			}(ch, id)
		}
	}
}

func sanitizeText(text string) string {
	if text == "" {
		return ""
	}

	cleaned := ansiPattern.ReplaceAllString(text, "")
	cleaned = strings.ReplaceAll(cleaned, "\r\n", "\n")
	cleaned = strings.ReplaceAll(cleaned, "\r", "\n")
	return cleaned
}

func writeSSE(w http.ResponseWriter, event eventPayload) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintf(w, "data: %s\n\n", data); err != nil {
		return err
	}
	return nil
}

const consoleHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8" />
<title>ServerCommander</title>
<style>
body {
  margin: 0;
  font-family: "Fira Code", "Consolas", "Courier New", monospace;
  background-color: #111;
  color: #e8e8e8;
  display: flex;
  flex-direction: column;
  height: 100vh;
}
#console {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  white-space: pre-wrap;
}
#input-area {
  display: flex;
  padding: 12px 16px;
  background-color: #1b1b1b;
  border-top: 1px solid #333;
}
#prompt {
  margin-right: 8px;
  color: #4caf50;
  font-weight: bold;
}
#command-input {
  flex: 1;
  background-color: #111;
  color: #e8e8e8;
  border: 1px solid #333;
  padding: 8px;
  font-family: inherit;
}
#command-input:focus {
  outline: none;
  border-color: #4caf50;
}
#execute {
  margin-left: 12px;
  padding: 8px 12px;
  background-color: #4caf50;
  color: #111;
  border: none;
  cursor: pointer;
  font-weight: bold;
}
#execute:hover {
  background-color: #66bb6a;
}
</style>
</head>
<body>
<div id="console"></div>
<form id="input-area">
  <span id="prompt">&gt;&gt;</span>
  <input id="command-input" type="text" autocomplete="off" placeholder="Enter command" />
  <button id="execute" type="submit">Execute</button>
</form>
<script>
const consoleElement = document.getElementById('console');
const form = document.getElementById('input-area');
const input = document.getElementById('command-input');

function appendText(text) {
  if (!text) {
    return;
  }
  consoleElement.textContent += text;
  consoleElement.scrollTop = consoleElement.scrollHeight;
}

function clearConsole() {
  consoleElement.textContent = '';
}

const events = new EventSource('/events');
events.onmessage = function (event) {
  const payload = JSON.parse(event.data);
  if (payload.type === 'clear') {
    clearConsole();
    return;
  }
  if (payload.type === 'append') {
    appendText(payload.text);
  }
};

events.onerror = function () {
  appendText('\n[Connection lost. Retrying...]\n');
};

form.addEventListener('submit', function (e) {
  e.preventDefault();
  const value = input.value.trim();
  if (!value) {
    input.value = '';
    return;
  }
  fetch('/execute', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ command: value })
  }).catch(function () {
    appendText('\n[Failed to send command]\n');
  });
  input.value = '';
});

window.addEventListener('load', function () {
  input.focus();
});
</script>
</body>
</html>`

func openBrowser(url string) error {
	var command string
	var args []string

	switch runtime.GOOS {
	case "windows":
		command = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	case "darwin":
		command = "open"
		args = []string{url}
	default:
		command = "xdg-open"
		args = []string{url}
	}

	executable, err := exec.LookPath(command)
	if err != nil {
		return err
	}

	cmd := exec.Command(executable, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Start()
}
