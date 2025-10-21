package app

import (
	"fmt"
	"net/url"
	"sync"

	"fyne.io/fyne/v2"
)

// New returns a minimal stub implementation of fyne.App so builds can succeed
// in environments where the full dependency graph is unavailable.
func New() fyne.App {
	return &stubApp{}
}

type stubApp struct {
	mu      sync.Mutex
	windows []*stubWindow
}

func (a *stubApp) NewWindow(title string) fyne.Window {
	win := &stubWindow{title: title, app: a}
	a.mu.Lock()
	a.windows = append(a.windows, win)
	a.mu.Unlock()
	return win
}

func (a *stubApp) Run() {
	// No event loop in the stub implementation.
}

func (a *stubApp) QueueUpdate(fn func()) {
	if fn != nil {
		fn()
	}
}

func (a *stubApp) OpenURL(u *url.URL) error {
	if u == nil {
		return fmt.Errorf("nil url provided")
	}
	// The stub just acknowledges the request.
	return nil
}

type stubWindow struct {
	app            *stubApp
	title          string
	size           fyne.Size
	content        fyne.CanvasObject
	closeIntercept func()
}

func (w *stubWindow) SetMaster() {}

func (w *stubWindow) SetFixedSize(bool) {}

func (w *stubWindow) Resize(size fyne.Size) {
	w.size = size
}

func (w *stubWindow) CenterOnScreen() {}

func (w *stubWindow) SetCloseIntercept(fn func()) {
	w.closeIntercept = fn
}

func (w *stubWindow) Show() {}

func (w *stubWindow) Close() {
	if w.closeIntercept != nil {
		w.closeIntercept()
	}
}

func (w *stubWindow) SetContent(obj fyne.CanvasObject) {
	w.content = obj
}

func (w *stubWindow) Canvas() fyne.Canvas {
	return stubCanvas{}
}

type stubCanvas struct{}

func (stubCanvas) MinSize() fyne.Size {
	return fyne.Size{}
}
