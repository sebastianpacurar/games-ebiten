package resources

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

type (
	CasinoCards interface {
		HitBox() image.Rectangle // HitBox - returns Min x,y and Max x,y
		Hovered(int, int) bool   // IsHovered - returns true if cursor is inside the shape
		Dragged() bool           // IsDragged - returns true if the image is dragged
		SetDragged(bool)         // SetDraggedState - sets the Dragged state to the given value
		DrawCard(*ebiten.Image)
	}

	MatchIcons interface {
		HitBox() image.Rectangle // HitBox - returns Min x,y and Max x,y
		Revealed() bool          // IsRevealed - returns true if an icon is revealed
		SetRevealed(bool)        // SetRevealedState - sets an icon to be hidden or revealed
		SetRemoved(bool)         // SetRemovedState - sets an icon to be visible or hidden
		Removed() bool           // IsRemoved - returns true
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
