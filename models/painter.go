package models

import "image/color"

type TextElement struct {
	X     int
	Y     int
	Text  string
	Color color.RGBA
}
