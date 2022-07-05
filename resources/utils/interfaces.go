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

	SolitaireCard interface {
		IsRevealed() bool      // IsRevealed - returns true if the image is revealed
		SetRevealedState(bool) // SetRevealedState - sets the Revealed state to the given value
	}

	FreeCellCard interface {
		IsDraggable() bool      // IsDraggable - returns true if the card from column is draggable
		SetDraggableState(bool) // SetDraggableState - sets the Draggable state to the given value
	}

	MatchIcons interface {
		HitBox() image.Rectangle // HitBox - returns Min x,y and Max x,y
		IsRevealed() bool        // IsRevealed - returns true if an icon is revealed
		SetRevealedState(bool)   // SetRevealedState - sets an icon to be hidden or revealed
		SetRemovedState(bool)    // SetRemovedState - sets an icon to be visible or hidden
		IsRemoved() bool         // IsRemoved - returns true
		DrawIcon(image *ebiten.Image)
	}

	InteractiveSprite interface {
		HitBox() image.Rectangle  // HitBox - returns Min x,y and Max x,y
		SetLocation(string, int)  // SetLocation - updates X or Y, based on X or Y axis
		SetDelta(string, float64) // SetDelta - updates VX or VY, based on X or Y axis
		DrawInteractiveSprite(*ebiten.Image)
	}

	StaticSprite interface {
		HitBox() image.Rectangle // HitBox - returns Min x,y and Max x,y
		SetLocation(string, int) // SetLocation - updates X or Y, based on X or Y axis
		DrawStaticSprite(*ebiten.Image)
	}
)
