package utils

import "github.com/hajimehoshi/ebiten/v2"

type (
	CasinoCards interface {
		GetPosition() (float64, float64, float64, float64) // GetPosition - returns x, y, w, h
		SetLocation(string, float64)                       // SetLocation - updates x or y, based on the axis
		GetDraggedState() bool                             // GetDraggedState - returns true if the image is dragged
		SetDraggedState(bool)                              // SetDraggedState - sets the drag state to the given value
		DrawCard(*ebiten.Image)
	}

	MatchIcons interface {
		GetPosition() (float64, float64, float64, float64) // GetPosition - returns x, y, w, h
		GetRevealState() bool                              // GetRevealState - returns true if an icon is revealed
		SetRevealState(bool)                               // SetRevealState - sets an icon to be hidden or revealed
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
