package fyne

import "net/url"

// App represents a graphical application instance.
type App interface {
	NewWindow(title string) Window
	Run()
	QueueUpdate(func())
	OpenURL(*url.URL) error
}

// Window describes a top-level user interface window.
type Window interface {
	SetMaster()
	SetFixedSize(bool)
	Resize(Size)
	CenterOnScreen()
	SetCloseIntercept(func())
	Show()
	Close()
	SetContent(CanvasObject)
	Canvas() Canvas
}

// Canvas represents a drawable surface.
type Canvas interface{}

// CanvasObject is any drawable UI element.
type CanvasObject interface {
	MinSize() Size
}

// Widget is any interactive canvas object that can render itself.
type Widget interface {
	CanvasObject
	CreateRenderer() WidgetRenderer
}

// WidgetRenderer draws a widget.
type WidgetRenderer interface {
	Layout(Size)
	MinSize() Size
	Refresh()
	Destroy()
	Objects() []CanvasObject
}

// Size represents a width/height pair.
type Size struct {
	Width  float32
	Height float32
}

// NewSize constructs a Size instance.
func NewSize(width, height float32) Size {
	return Size{Width: width, Height: height}
}

// Position represents a point in 2D space.
type Position struct {
	X float32
	Y float32
}

// PointEvent describes a pointer interaction event.
type PointEvent struct {
	AbsolutePosition Position
}

// TextStyle configures textual rendering options.
type TextStyle struct {
	Bold      bool
	Italic    bool
	Monospace bool
}

// TextWrap controls wrapping behaviour.
type TextWrap int

const (
	TextWrapOff TextWrap = iota
	TextWrapBreak
	TextWrapWord
)

// TextAlign controls textual alignment.
type TextAlign int

const (
	TextAlignLeading TextAlign = iota
	TextAlignCenter
	TextAlignTrailing
)

// Menu represents a pop-up menu structure.
type Menu struct {
	Label string
	Items []*MenuItem
}

// MenuItem is a single selectable menu item.
type MenuItem struct {
	Label  string
	Action func()
}

// NewMenu constructs a new menu.
func NewMenu(label string, items ...*MenuItem) *Menu {
	return &Menu{Label: label, Items: items}
}

// NewMenuItem constructs a new menu item.
func NewMenuItem(label string, action func()) *MenuItem {
	return &MenuItem{Label: label, Action: action}
}

// Activate triggers the menu item's action.
func (m *MenuItem) Activate() {
	if m != nil && m.Action != nil {
		m.Action()
	}
}
