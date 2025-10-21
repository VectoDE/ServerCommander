package theme

import "image/color"

// ForegroundColor returns a generic foreground colour.
func ForegroundColor() color.Color {
	return color.NRGBA{R: 0x22, G: 0x22, B: 0x22, A: 0xff}
}

// PrimaryColor returns a highlight colour.
func PrimaryColor() color.Color {
	return color.NRGBA{R: 0x21, G: 0x6c, B: 0xff, A: 0xff}
}

// BackgroundColor returns a generic background colour.
func BackgroundColor() color.Color {
	return color.NRGBA{R: 0xf0, G: 0xf0, B: 0xf0, A: 0xff}
}
