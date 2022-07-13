package utils

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

type (
	CasinoCards interface {
		HitBox() image.Rectangle // HitBox - returns Min x,y and Max x,y
		IsHovered(int, int) bool // IsHovered - returns true if cursor is inside the shape
		IsDragged() bool         // IsDragged - returns true if the image is dragged
		SetDraggedState(bool)    // SetDraggedState - sets the Dragged state to the given value
		DrawCard(*ebiten.Image)
	}
)
