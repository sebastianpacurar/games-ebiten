package resources

import (
	"image"
)

type (
	RectArea interface {
		HitBox() image.Rectangle // HitBox - returns Rect{} props
	}

	Sprite interface {
		RectArea
		SetLocation(string, int)  // SetLocation - updates X or Y, based on X or Y axis
		SetDelta(string, float64) // SetDelta - updates VX or VY, based on X or Y axis
	}

	Hoverable interface {
		RectArea
		Hovered(cx, cy int) bool // returns true if HitBox is hovered
	}
)
