//go:build desktop
// +build desktop

package ui

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"servercommander/src/utils"
)

type commandExecutor func(string) error

var (
	consoleInstance *DesktopConsole
	consoleMu       sync.RWMutex
)

var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;]*[A-Za-z]`)

const (
	portfolioLink      = "https://uplytech.com/portfolio"
	websiteLink        = "https://uplytech.com"
	commandPlaceholder = "Befehl eingeben und Enter drücken…"
	promptPlaceholder  = "Antwort eingeben und Enter drücken…"
)

// DesktopConsole encapsulates the GUI console experience, piping stdout and
// stderr into a dedicated window that accepts command execution requests.
type DesktopConsole struct {
	executor commandExecutor

	app    fyne.App
	window fyne.Window

	output   *widget.Label
	scroller *container.Scroll
	input    *widget.Entry

	stdoutPipe *os.File
	stderrPipe *os.File
	origStdout *os.File
	origStderr *os.File

	promptBridge    *promptBridge
	promptMu        sync.Mutex
	pendingResponse chan string

	logMu      sync.Mutex
	logBuilder strings.Builder

	cleanupOnce sync.Once
}

// RunStandaloneConsole launches the GUI console and blocks until it closes.
func RunStandaloneConsole(executor func(string) error) error {
	console := &DesktopConsole{executor: executor}

	application := app.New()
	console.app = application

	if err := console.redirectStandardStreams(); err != nil {
		return err
	}

	console.initializePromptBridge()

	consoleMu.Lock()
	consoleInstance = console
	consoleMu.Unlock()

	defer func() {
		console.cleanup()
		consoleMu.Lock()
		consoleInstance = nil
		consoleMu.Unlock()
	}()

	return console.run()
}

// ClearConsole removes the accumulated log content from the GUI window.
func ClearConsole() bool {
	consoleMu.RLock()
	console := consoleInstance
	consoleMu.RUnlock()
	if console == nil {
		return false
	}

	console.clear()
	return true
}

func (c *DesktopConsole) run() error {
	c.window = c.app.NewWindow("ServerCommander")
	c.window.SetMaster()
	c.window.SetFixedSize(true)
	c.window.Resize(fyne.NewSize(960, 600))
	c.window.CenterOnScreen()

	c.window.SetCloseIntercept(func() {
		c.cleanup()
		c.window.Close()
	})

	c.buildUI()

	c.window.Show()
	ApplicationBanner()
	c.app.Run()
	return nil
}

func (c *DesktopConsole) buildUI() {
	c.output = widget.NewLabel("")
	c.output.TextStyle = fyne.TextStyle{Monospace: true}
	c.output.Wrapping = fyne.TextWrapWord
	c.scroller = container.NewVScroll(c.output)
	c.scroller.SetMinSize(fyne.NewSize(0, 400))

	c.input = widget.NewEntry()
	c.input.SetPlaceHolder(commandPlaceholder)
	c.input.OnSubmitted = func(text string) {
		c.handleInput(text)
	}

	runButton := widget.NewButton("Ausführen", func() {
		c.handleInput(c.input.Text)
	})
	runButton.Importance = widget.HighImportance

	clearButton := widget.NewButton("Leeren", func() {
		c.clear()
		ApplicationBanner()
	})

	actions := container.NewHBox(clearButton, layout.NewSpacer(), runButton)
	footer := container.NewBorder(nil, nil, nil, actions, c.input)

	header := container.NewHBox(newLogoWidget(c), layout.NewSpacer())

	c.window.SetContent(container.NewBorder(header, footer, nil, nil, c.scroller))
}

func (c *DesktopConsole) executeCommand(command string) {
	if c.executor == nil {
		return
	}

	if err := c.executor(command); err != nil {
		fmt.Println(utils.Red, err.Error(), utils.Reset)
	}
}

func (c *DesktopConsole) handleInput(input string) {
	trimmed := strings.TrimSpace(input)
	c.input.SetText("")
	if trimmed == "" {
		return
	}

	if c.deliverPromptResponse(trimmed) {
		c.appendLog(trimmed + "\n")
		return
	}

	c.appendLog(">> " + trimmed + "\n")
	go c.executeCommand(trimmed)
}

func (c *DesktopConsole) redirectStandardStreams() error {
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

func (c *DesktopConsole) initializePromptBridge() {
	bridge := newPromptBridge()
	c.promptBridge = bridge
	utils.SetPromptReader(bridge)

	go c.watchPromptRequests()
}

func (c *DesktopConsole) readPipe(pipe *os.File, mirror *os.File) {
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
				c.appendLog(emit)
			}

			pending.Reset()
			if remainder != "" {
				pending.WriteString(remainder)
			}
		}

		if err != nil {
			if pending.Len() > 0 {
				c.appendLog(pending.String())
			}
			return
		}
	}
}

func (c *DesktopConsole) appendLog(text string) {
	cleaned := sanitizeText(text)
	if cleaned == "" {
		return
	}

	c.logMu.Lock()
	c.logBuilder.WriteString(cleaned)
	logSnapshot := c.logBuilder.String()
	c.logMu.Unlock()

	if c.app == nil {
		return
	}

	c.app.QueueUpdate(func() {
		if c.output != nil {
			c.output.SetText(logSnapshot)
		}
		if c.scroller != nil {
			c.scroller.ScrollToBottom()
		}
	})
}

func (c *DesktopConsole) clear() {
	c.logMu.Lock()
	c.logBuilder.Reset()
	c.logMu.Unlock()

	if c.app == nil {
		return
	}

	c.app.QueueUpdate(func() {
		if c.output != nil {
			c.output.SetText("")
		}
		if c.scroller != nil {
			c.scroller.ScrollToTop()
		}
	})
}

func (c *DesktopConsole) cleanup() {
	c.cleanupOnce.Do(func() {
		if c.stdoutPipe != nil {
			c.stdoutPipe.Close()
		}
		if c.stderrPipe != nil {
			c.stderrPipe.Close()
		}
		if c.origStdout != nil {
			os.Stdout = c.origStdout
		}
		if c.origStderr != nil {
			os.Stderr = c.origStderr
		}
		if c.promptBridge != nil {
			c.promptBridge.Close()
		}
		utils.SetPromptReader(os.Stdin)
		c.promptMu.Lock()
		if c.pendingResponse != nil {
			close(c.pendingResponse)
			c.pendingResponse = nil
		}
		c.promptMu.Unlock()
	})
}

func (c *DesktopConsole) watchPromptRequests() {
	if c.promptBridge == nil {
		return
	}

	for responseCh := range c.promptBridge.Requests() {
		c.promptMu.Lock()
		c.pendingResponse = responseCh
		c.promptMu.Unlock()
		c.setPromptPlaceholder(true)
	}
	c.setPromptPlaceholder(false)
}

func (c *DesktopConsole) deliverPromptResponse(response string) bool {
	c.promptMu.Lock()
	responseCh := c.pendingResponse
	if responseCh != nil {
		c.pendingResponse = nil
	}
	c.promptMu.Unlock()

	if responseCh == nil {
		return false
	}

	responseCh <- response
	close(responseCh)
	c.setPromptPlaceholder(false)
	return true
}

func (c *DesktopConsole) setPromptPlaceholder(waiting bool) {
	if c.app == nil {
		return
	}

	placeholder := commandPlaceholder
	if waiting {
		placeholder = promptPlaceholder
	}

	c.app.QueueUpdate(func() {
		if c.input != nil {
			c.input.SetPlaceHolder(placeholder)
			if waiting {
				c.input.Focus()
			}
		}
	})
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

type logoWidget struct {
	widget.BaseWidget

	label  *canvas.Text
	parent *DesktopConsole
}

func newLogoWidget(console *DesktopConsole) *logoWidget {
	widget := &logoWidget{parent: console}
	widget.ExtendBaseWidget(widget)
	return widget
}

func (l *logoWidget) CreateRenderer() fyne.WidgetRenderer {
	l.label = canvas.NewText("ServerCommander", theme.PrimaryColor())
	l.label.TextStyle = fyne.TextStyle{Bold: true}
	l.label.Alignment = fyne.TextAlignCenter

	content := container.NewCenter(l.label)
	return widget.NewSimpleRenderer(content)
}

func (l *logoWidget) MinSize() fyne.Size {
	if l.label == nil {
		return fyne.NewSize(140, 48)
	}

	min := l.label.MinSize()
	return fyne.NewSize(min.Width+32, min.Height+16)
}

func (l *logoWidget) Tapped(_ *fyne.PointEvent) {}

func (l *logoWidget) TappedSecondary(event *fyne.PointEvent) {
	menu := fyne.NewMenu("",
		fyne.NewMenuItem("Portfolio", func() {
			l.openURL(portfolioLink)
		}),
		fyne.NewMenuItem("UplyTech Website", func() {
			l.openURL(websiteLink)
		}),
	)

	if l.parent != nil && l.parent.window != nil {
		pop := widget.NewPopUpMenu(menu, l.parent.window.Canvas())
		pop.ShowAtPosition(event.AbsolutePosition)
	}
}

func (l *logoWidget) openURL(target string) {
	if l.parent == nil || l.parent.app == nil {
		return
	}

	parsed, err := url.Parse(target)
	if err != nil {
		fmt.Println(utils.Red, "Ungültige URL:", target, utils.Reset)
		return
	}

	if err := l.parent.app.OpenURL(parsed); err != nil {
		fmt.Println(utils.Red, "Konnte URL nicht öffnen:", target, utils.Reset)
	}
}

type promptBridge struct {
	requestCh chan chan string
	closeCh   chan struct{}

	mu     sync.Mutex
	buffer []byte
}

func newPromptBridge() *promptBridge {
	return &promptBridge{
		requestCh: make(chan chan string),
		closeCh:   make(chan struct{}),
	}
}

func (p *promptBridge) Read(dst []byte) (int, error) {
	if len(dst) == 0 {
		return 0, nil
	}

	p.mu.Lock()
	if len(p.buffer) > 0 {
		n := copy(dst, p.buffer)
		p.buffer = p.buffer[n:]
		p.mu.Unlock()
		return n, nil
	}
	p.mu.Unlock()

	responseCh := make(chan string, 1)

	select {
	case <-p.closeCh:
		return 0, io.EOF
	case p.requestCh <- responseCh:
	}

	var (
		response string
		ok       bool
	)
	select {
	case <-p.closeCh:
		return 0, io.EOF
	case response, ok = <-responseCh:
	}

	if !ok {
		return 0, io.EOF
	}

	if !strings.HasSuffix(response, "\n") {
		response += "\n"
	}

	data := []byte(response)
	n := copy(dst, data)

	if n < len(data) {
		p.mu.Lock()
		p.buffer = append(p.buffer, data[n:]...)
		p.mu.Unlock()
	}

	return n, nil
}

func (p *promptBridge) Close() {
	p.mu.Lock()
	select {
	case <-p.closeCh:
		p.mu.Unlock()
		return
	default:
	}
	close(p.closeCh)
	close(p.requestCh)
	p.buffer = nil
	p.mu.Unlock()
}

func (p *promptBridge) Requests() <-chan chan string {
	return p.requestCh
}
