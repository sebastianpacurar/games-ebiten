package game

import (
	u "games-ebiten/resources/utils"
	"github.com/hajimehoshi/ebiten/v2"
	_ "image/png"
)

// Translation - translates the value of a card which is located at certain coordinates
var Translation = map[string]map[int]string{
	PixelatedTheme: {
		0: "2", 1: "3", 2: "4", 3: "5", 4: "6", 5: "7",
		6: "8", 7: "9", 8: "10", 9: "J", 10: "Q", 11: "K", 12: "A",
	},
	ClassicTheme: {
		0: "A", 1: "2", 2: "3", 3: "4", 4: "5", 5: "6", 6: "7",
		7: "8", 8: "9", 9: "10", 10: "J", 11: "Q", 12: "K",
	},
}

// Card - Implements StaticSprite interface
type Card struct {
	Img       *ebiten.Image
	FrameInfo map[string]int
	Suit      string
	Value     string
	LX, LY    float64
	W, H      float64
	ScaleX    float64
	ScaleY    float64
}

func (c *Card) GetLocations() (float64, float64) {
	return c.LX, c.LY
}

func (c *Card) GetSize() (float64, float64) {
	return c.W, c.H
}

func (c *Card) GetFramePosition() (int, int, int, int) {
	return c.FrameInfo[u.FrameOX], c.FrameInfo[u.FrameOY], c.FrameInfo[u.FrameW], c.FrameInfo[u.FrameH]
}

func (c *Card) GetImg() *ebiten.Image {
	return c.Img
}

func (c *Card) SetLocation(axis string, val float64) {
	if axis == u.X {
		c.LX = val
	} else if axis == u.Y {
		c.LY = val
	}
}

func (c *Card) DrawStaticSprite(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(c.ScaleX, c.ScaleY)
	op.GeoM.Translate(c.LX, c.LY)
	screen.DrawImage(c.Img, op)
}
