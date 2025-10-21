package layout

import "fyne.io/fyne/v2"

type spacer struct{}

func (s *spacer) MinSize() fyne.Size { return fyne.Size{} }

// NewSpacer returns a placeholder object used to align widgets.
func NewSpacer() fyne.CanvasObject { return &spacer{} }
