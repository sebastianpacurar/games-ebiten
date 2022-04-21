package game

import (
	u "games-ebiten/resources/utils"
	"github.com/hajimehoshi/ebiten/v2"
	_ "image/png"
)

// Translation - translates the value of a card which is located at certain coordinates
var Translation = map[string]map[int]string{
	u.PixelatedTheme: {
		0: "2", 1: "3", 2: "4", 3: "5", 4: "6", 5: "7",
		6: "8", 7: "9", 8: "10", 9: "J", 10: "Q", 11: "K", 12: "A",
	},
	u.ClassicTheme: {
		0: "A", 1: "2", 2: "3", 3: "4", 4: "5", 5: "6", 6: "7",
		7: "8", 8: "9", 9: "10", 10: "J", 11: "Q", 12: "K",
	},
	//u.AbstractTheme: {
	//	0: "2", 1: "3", 2: "4", 3: "5", 4: "6", 5: "7",
	//	6: "8", 7: "9", 8: "10", 9: "J", 10: "Q", 11: "K", 12: "A",
	//},
}

// Card - Implements CasinoCards interface
// IsRevealed - if the card's frontFace is visible
// IsLocked - if the card is in the stack slot (and has a card over it) or if it's in the draw stack
// IsDragged - holds the dragged state
type Card struct {
	Img        *ebiten.Image
	BackImg    *ebiten.Image
	Suit       string
	Value      string
	X, Y, W, H float64
	ScaleX     float64
	ScaleY     float64
	IsRevealed bool
	IsLocked   bool
	IsDragged  bool
}

func (c *Card) GetPosition() (float64, float64, float64, float64) {
	return c.X, c.Y, c.W, c.H
}

func (c *Card) SetLocation(axis string, val float64) {
	if axis == u.X {
		c.X = val
	} else if axis == u.Y {
		c.Y = val
	}
}

func (c *Card) DrawCardSprite(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(c.ScaleX, c.ScaleY)
	cx, cy := ebiten.CursorPosition()
	img := c.Img

	if !c.IsRevealed {
		img = c.BackImg
	}
	// if hovered, lower the opacity level of this card's RGBA
	if u.IsImgHovered(c, cx, cy) {
		op.ColorM.Scale(1, 1, 1, 0.85)
	}
	op.GeoM.Translate(c.X, c.Y)
	screen.DrawImage(img, op)
}

func (c *Card) GetDraggedState() bool {
	return c.IsDragged
}

func (c *Card) SetDraggedState(state bool) {
	c.IsDragged = state
}
