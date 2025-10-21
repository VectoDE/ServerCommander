package container

import "fyne.io/fyne/v2"

// Container is a simple grouping of canvas objects.
type Container struct {
	Objects []fyne.CanvasObject
	size    fyne.Size
}

// MinSize returns a naive size estimation based on contained objects.
func (c *Container) MinSize() fyne.Size {
	if c == nil {
		return fyne.Size{}
	}
	if c.size != (fyne.Size{}) {
		return c.size
	}
	var width, height float32
	for _, obj := range c.Objects {
		if obj == nil {
			continue
		}
		sz := obj.MinSize()
		if sz.Width > width {
			width = sz.Width
		}
		height += sz.Height
	}
	return fyne.NewSize(width, height)
}

// SetMinSize hints the container size.
func (c *Container) SetMinSize(size fyne.Size) {
	if c != nil {
		c.size = size
	}
}

// NewVScroll wraps content inside a scroll container.
func NewVScroll(content fyne.CanvasObject) *Scroll {
	return &Scroll{Content: content}
}

// NewHBox returns a container laid out horizontally.
func NewHBox(objects ...fyne.CanvasObject) *Container {
	return &Container{Objects: append([]fyne.CanvasObject(nil), objects...)}
}

// NewBorder arranges objects around a center element.
func NewBorder(top, bottom, left, right, center fyne.CanvasObject) *Container {
	objs := []fyne.CanvasObject{center, top, bottom, left, right}
	return &Container{Objects: objs}
}

// NewCenter centres a single object.
func NewCenter(object fyne.CanvasObject) *Container {
	return &Container{Objects: []fyne.CanvasObject{object}}
}

// Scroll represents a vertical scroll container.
type Scroll struct {
	Content fyne.CanvasObject
	size    fyne.Size
}

// MinSize returns the stored minimum size.
func (s *Scroll) MinSize() fyne.Size {
	if s == nil {
		return fyne.Size{}
	}
	if s.size != (fyne.Size{}) {
		return s.size
	}
	if s.Content == nil {
		return fyne.Size{}
	}
	return s.Content.MinSize()
}

// SetMinSize stores a requested minimum size.
func (s *Scroll) SetMinSize(size fyne.Size) {
	if s != nil {
		s.size = size
	}
}

// ScrollToBottom is a no-op in the stub implementation.
func (s *Scroll) ScrollToBottom() {}

// ScrollToTop is a no-op in the stub implementation.
func (s *Scroll) ScrollToTop() {}
