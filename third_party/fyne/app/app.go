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
	mu        sync.Mutex
	windows   map[*stubWindow]struct{}
	runOnce   sync.Once
	runCh     chan struct{}
	runClosed bool
}

func (a *stubApp) NewWindow(title string) fyne.Window {
	win := &stubWindow{title: title, app: a}

	a.mu.Lock()
	if a.windows == nil {
		a.windows = make(map[*stubWindow]struct{})
	}
	a.windows[win] = struct{}{}
	a.mu.Unlock()

	a.runChannel()

	return win
}

func (a *stubApp) Run() {
	ch := a.runChannel()

	shouldClose := false
	a.mu.Lock()
	if len(a.windows) == 0 && !a.runClosed {
		a.runClosed = true
		shouldClose = true
	}
	a.mu.Unlock()

	if shouldClose {
		close(ch)
	}

	<-ch
}

func (a *stubApp) runChannel() chan struct{} {
	a.runOnce.Do(func() {
		a.runCh = make(chan struct{})
	})
	return a.runCh
}

func (a *stubApp) windowClosed(win *stubWindow) {
	var ch chan struct{}
	shouldClose := false

	a.mu.Lock()
	if a.windows != nil {
		if _, ok := a.windows[win]; ok {
			delete(a.windows, win)
			if len(a.windows) == 0 && !a.runClosed && a.runCh != nil {
				a.runClosed = true
				ch = a.runCh
				shouldClose = true
			}
		}
	}
	a.mu.Unlock()

	if shouldClose {
		close(ch)
	}
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

	mu     sync.Mutex
	closed bool
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
	w.mu.Lock()
	if w.closed {
		w.mu.Unlock()
		return
	}
	w.closed = true
	intercept := w.closeIntercept
	w.closeIntercept = nil
	w.mu.Unlock()

	if intercept != nil {
		intercept()
	}

	if w.app != nil {
		w.app.windowClosed(w)
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
