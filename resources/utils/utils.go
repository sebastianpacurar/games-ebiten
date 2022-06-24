package utils

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"io/ioutil"
	"log"
	"math/rand"
	"time"
)

const (
	ScreenWidth  = 840
	ScreenHeight = 840

	// X Y - used as aliases for Main Axis and Cross Axis
	X = "x"
	Y = "y"

	// FrOX FrOY FrW FrH = minX, minY, maxX, maxY, for the area of an Image
	FrOX = "FrameOX"
	FrOY = "FrameOY"
	FrW  = "FrameW"
	FrH  = "FrameH"

	Hearts   = "Hearts"
	Clubs    = "Clubs"
	Diamonds = "Diamonds"
	Spades   = "Spades"

	Ace   = "A"
	Jack  = "J"
	Queen = "Q"
	King  = "K"

	Red   = "Red"
	Black = "Black"

	ClassicTheme   = "classic"
	PixelatedTheme = "8bit"

	StaticBack1 = "StaticBack1"

	CardsVSpacer = 25
)

func LoadSpriteImage(path string) image.Image {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}
	img, _, err := image.Decode(bytes.NewReader(file))
	if err != nil {
		log.Fatalln(err)
	}
	return img
}

// GenerateRandomPosition - generate a random location for X and y
// (x is between 0 and ScreenWidth) and (y is between 0 and ScreenHeight)
func GenerateRandomPosition(minX, minY, maxX, maxY int) (int, int) {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(maxX-minX) + minX, rand.Intn(maxY-minY) + minY
}

// IsCollision - returns true if rectangle images overlap in any way
func IsCollision(src, target image.Rectangle) bool {
	return target.Overlaps(src)
}

// BoundaryValidation - prevents characters to move out of the view, in any of the 4 directions
func BoundaryValidation(i interface{}, minX, maxX, minY, maxY int) {
	switch i.(type) {
	case InteractiveSprite:
		img := i.(InteractiveSprite)
		bounds := img.HitBox()

		if bounds.Min.X <= minX {
			img.SetLocation(X, minX)
			img.SetDelta(X, 0)
		}
		if bounds.Min.X >= maxX-bounds.Dx() {
			img.SetLocation(X, maxX-bounds.Dx())
			img.SetDelta(X, 0)
		}
		if bounds.Min.Y <= minY {
			img.SetLocation(Y, minY)
			img.SetDelta(Y, 0)
		}
		if bounds.Min.Y >= maxY-bounds.Dy() {
			img.SetLocation(Y, maxY-bounds.Dy())
			img.SetDelta(Y, 0)
		}
	}
}

// IsAreaHovered - Returns true if the cursor overlaps the target interface
func IsAreaHovered(i interface{}) bool {
	pt := image.Pt(ebiten.CursorPosition())
	var area image.Rectangle
	switch i.(type) {
	case CasinoCards:
		c := i.(CasinoCards)
		area = c.HitBox()
	case MatchIcons:
		mi := i.(MatchIcons)
		area = mi.HitBox()
	}
	return pt.In(area)
}

// GetFlexboxQuadrants - splits the screen in x quadrants and returns a rect for each of them (works only on X-Axis)
func GetFlexboxQuadrants(cols int) map[int]image.Rectangle {
	quads := make(map[int]image.Rectangle, 0)
	unit := (ScreenWidth) / cols

	for i := 0; i < cols; i++ {
		minX := unit * i
		maxX := unit * (i + 1)
		quads[i] = image.Rect(minX, 0, maxX, ScreenHeight)
	}
	return quads
}

// CenterItem - centers the item within the given quadrant
func CenterItem(width int, quad image.Rectangle) int {
	return (quad.Min.X + quad.Dx()/2) - width/2
}
