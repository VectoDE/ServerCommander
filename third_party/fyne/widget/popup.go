package widget

import "fyne.io/fyne/v2"

// PopUpMenu is a stub implementation of a pop-up menu.
type PopUpMenu struct {
	menu   *fyne.Menu
	canvas fyne.Canvas
}

// NewPopUpMenu creates a new pop-up menu instance.
func NewPopUpMenu(menu *fyne.Menu, canvas fyne.Canvas) *PopUpMenu {
	return &PopUpMenu{menu: menu, canvas: canvas}
}

// ShowAtPosition has no visible behaviour in the stub implementation.
func (p *PopUpMenu) ShowAtPosition(fyne.Position) {}

// Dismiss is a no-op for the stub implementation.
func (p *PopUpMenu) Dismiss() {}
