package utils

import "github.com/hajimehoshi/ebiten/v2"

type CasinoCards interface {
	// GetPosition - returns x, y, w, h
	GetPosition() (float64, float64, float64, float64)

	// SetLocation - updates x or y, based on the axis
	SetLocation(string, float64)

	// DrawCardSprite - draws the image
	DrawCardSprite(*ebiten.Image)

	// GetDraggedState - returns true if the image is dragged
	GetDraggedState() bool

	// SetDraggedState - sets the drag state to the given value
	SetDraggedState(bool)
}
