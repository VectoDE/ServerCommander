package canvas

import (
	"image/color"

	"fyne.io/fyne/v2"
)

// Text is a minimal representation of a canvas text element.
type Text struct {
	Text      string
	Color     color.Color
	TextSize  float32
	TextStyle fyne.TextStyle
	Alignment fyne.TextAlign
}

// NewText creates a new canvas text instance.
func NewText(text string, clr color.Color) *Text {
	return &Text{Text: text, Color: clr, TextSize: 12}
}

// MinSize returns a naive size estimation for the text.
func (t *Text) MinSize() fyne.Size {
	width := float32(len(t.Text)) * 7
	if width == 0 {
		width = 1
	}
	height := t.TextSize
	if height == 0 {
		height = 12
	}
	return fyne.NewSize(width, height)
}

// Refresh does nothing in the stub implementation.
func (t *Text) Refresh() {}
