package widget

import "fyne.io/fyne/v2"

// BaseWidget provides helper methods for building custom widgets.
type BaseWidget struct{}

// ExtendBaseWidget is a no-op for the stub implementation.
func (b *BaseWidget) ExtendBaseWidget(fyne.Widget) {}

// Refresh is a no-op for the stub implementation.
func (b *BaseWidget) Refresh() {}

// stubRenderer implements fyne.WidgetRenderer with no behaviour.
type stubRenderer struct {
	object fyne.CanvasObject
}

func (r *stubRenderer) Layout(fyne.Size)             {}
func (r *stubRenderer) MinSize() fyne.Size           { return r.object.MinSize() }
func (r *stubRenderer) Refresh()                     {}
func (r *stubRenderer) Destroy()                     {}
func (r *stubRenderer) Objects() []fyne.CanvasObject { return []fyne.CanvasObject{r.object} }

// NewSimpleRenderer creates a renderer that simply exposes the object.
func NewSimpleRenderer(obj fyne.CanvasObject) fyne.WidgetRenderer {
	return &stubRenderer{object: obj}
}

// Label represents simple text output.
type Label struct {
	BaseWidget
	Text      string
	TextStyle fyne.TextStyle
	Wrapping  fyne.TextWrap
}

// NewLabel creates a new label instance.
func NewLabel(text string) *Label {
	l := &Label{Text: text}
	l.ExtendBaseWidget(l)
	return l
}

// SetText updates the label contents.
func (l *Label) SetText(text string) {
	l.Text = text
}

// MinSize returns a naive size estimation.
func (l *Label) MinSize() fyne.Size {
	width := float32(len(l.Text)) * 7
	if width == 0 {
		width = 1
	}
	height := float32(14)
	return fyne.NewSize(width, height)
}

// CreateRenderer returns a minimal renderer for compatibility.
func (l *Label) CreateRenderer() fyne.WidgetRenderer {
	return &stubRenderer{object: l}
}

// Entry represents text input.
type Entry struct {
	BaseWidget
	Text        string
	placeholder string
	OnSubmitted func(string)
}

// NewEntry constructs an empty entry widget.
func NewEntry() *Entry {
	e := &Entry{}
	e.ExtendBaseWidget(e)
	return e
}

// SetPlaceHolder stores the placeholder text.
func (e *Entry) SetPlaceHolder(text string) {
	e.placeholder = text
}

// SetText updates the entry value.
func (e *Entry) SetText(text string) {
	e.Text = text
}

// Focus is a no-op for the stub implementation.
func (e *Entry) Focus() {}

// MinSize returns a naive entry size.
func (e *Entry) MinSize() fyne.Size {
	return fyne.NewSize(120, 20)
}

// CreateRenderer returns a stub renderer.
func (e *Entry) CreateRenderer() fyne.WidgetRenderer {
	return &stubRenderer{object: e}
}

// Importance controls button emphasis.
type Importance int

const (
	LowImportance Importance = iota
	MediumImportance
	HighImportance
)

// Button is a clickable element.
type Button struct {
	BaseWidget
	Label      string
	Importance Importance
	onTapped   func()
}

// NewButton creates a new button instance.
func NewButton(label string, tapped func()) *Button {
	b := &Button{Label: label, onTapped: tapped}
	b.ExtendBaseWidget(b)
	return b
}

// MinSize returns a naive button size.
func (b *Button) MinSize() fyne.Size {
	width := float32(len(b.Label))*7 + 16
	if width < 40 {
		width = 40
	}
	return fyne.NewSize(width, 24)
}

// CreateRenderer returns a stub renderer.
func (b *Button) CreateRenderer() fyne.WidgetRenderer {
	return &stubRenderer{object: b}
}

// Tapped triggers the button callback.
func (b *Button) Tapped() {
	if b.onTapped != nil {
		b.onTapped()
	}
}
