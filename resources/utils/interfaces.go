package utils

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

type (
	CasinoCards interface {
		GetGeomData() image.Rectangle // GetGeomData - returns Min x,y and Max x,y
		IsHovered(int, int) bool      // IsHovered - returns true if cursor is inside the shape
		IsDragged() bool              // IsDragged - returns true if the image is dragged
		SetDraggedState(bool)         // SetDraggedState - sets the Dragged state to the given value
		IsRevealed() bool             // IsRevealed - returns true if the image is revealed
		SetRevealedState(bool)        // SetRevealedState - sets the Revealed state to the given value
		DrawCard(*ebiten.Image)
	}

	MatchIcons interface {
		GetGeomData() image.Rectangle // GetGeomData - returns Min x,y and Max x,y
		GetRevealedState() bool       // GetRevealedState - returns true if an icon is revealed
		SetRevealedState(bool)        // SetRevealedState - sets an icon to be hidden or revealed
		DrawIcon(image *ebiten.Image)
	}

	InteractiveSprite interface {
		GetLocations() (float64, float64) // GetLocations - returns (x, y)
		GetSize() (float64, float64)      // GetSize - returns (w, h)
		SetLocation(string, float64)      // SetLocation - updates X or Y, based on X or Y axis
		SetDelta(string, float64)         // SetDelta - updates DX or DY, based on X or Y axis
		DrawInteractiveSprite(*ebiten.Image)
	}

	StaticSprite interface {
		GetLocations() (float64, float64) // GetLocations - returns (X, Y)
		GetSize() (float64, float64)      // GetSize - returns (width, height)
		SetLocation(string, float64)      // SetLocation - updates X or Y, based on X or Y axis
		DrawStaticSprite(*ebiten.Image)
	}
)
