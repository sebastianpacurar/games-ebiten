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
	ScreenWidth  = 1280
	ScreenHeight = 860

	// X Y - used as aliases for Main Axis and Cross Axis
	X = "x"
	Y = "y"

	// MinX MinY MaxX MaxY - represent the vertices points of an Image
	MinX = "minX"
	MinY = "minY"
	MaxX = "maxX"
	MaxY = "maxY"

	// FrOX FrOY FrW FrH = minX, minY, maxX, maxY, for the area of an Image
	FrOX = "FrameOX"
	FrOY = "FrameOY"
	FrW  = "FrameW"
	FrH  = "FrameH"

	NPC1 = "npc1"
	NPC2 = "npc2"
	NPC3 = "npc3"
	NPC4 = "npc4"
	NPC5 = "npc5"

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
	AbstractTheme  = "abstract"
	SimpleTheme    = "simple"

	StaticBack1   = "StaticBack1"
	StaticBack2   = "StaticBack2"
	DynamicRobot  = "DynamicRobot"
	DynamicCastle = "DynamicCastle"
	DynamicBeach  = "DynamicBeach"
	DynamicSleeve = "DynamicSleeve"
)

var ScreenDims = map[string]float64{MinX: 0, MaxX: ScreenWidth, MinY: 0, MaxY: ScreenHeight}

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

// GenerateRandomLocation - generate a random location for X and y
// (x is between 0 and ScreenWidth) and (y is between 0 and ScreenHeight)
func GenerateRandomLocation(minX, maxX, minY, maxY float64) (float64, float64) {
	rand.Seed(time.Now().UnixNano())
	return (rand.Float64() * (maxX - minX)) + minX, (rand.Float64() * (maxY - minY)) + minY
}

// IsCollision - returns true if rectangle images overlap in any way
func IsCollision(x1, y1, w1, h1, x2, y2, w2, h2 float64) bool {
	return x1 < x2+w2 && x1+w1 > x2 && y1 < y2+h2 && h1+y1 > y2
}

// HitBox - generate the shape's hitbox (minX, maxX, minY, maxY)
func HitBox(x, y, w, h float64) map[string]float64 {
	return map[string]float64{
		MinX: x,
		MaxX: x + w,
		MinY: y,
		MaxY: y + h,
	}
}

// BoundaryValidation - prevents characters to move out of the view, in any of the 4 directions
func BoundaryValidation(i interface{}, minX, maxX, minY, maxY float64) {
	switch i.(type) {
	case InteractiveSprite:
		img := i.(InteractiveSprite)

		locX, locY := img.GetLocations()
		w, h := img.GetSize()

		if locX <= minX {
			img.SetLocation(X, minX)
			img.SetDelta(X, 0)
		}
		if locX >= maxX-w {
			img.SetLocation(X, maxX-w)
			img.SetDelta(X, 0)
		}
		if locY <= minY {
			img.SetLocation(Y, minY)
			img.SetDelta(Y, 0)
		}
		if locY >= maxY-h {
			img.SetLocation(Y, maxY-h)
			img.SetDelta(Y, 0)
		}
	}
}

// IsAreaHovered - Returns true if the cursor overlaps the target interface
func IsAreaHovered(i interface{}) bool {
	cx, cy := ebiten.CursorPosition()
	var x, y, w, h float64

	switch i.(type) {
	case CasinoCards:
		c := i.(CasinoCards)
		x, y, w, h = c.GetPosition()
	case MatchIcons:
		mi := i.(MatchIcons)
		x, y, w, h = mi.GetPosition()
	}
	return int(x) <= cx && cx < int(x+w) && int(y) <= cy && cy < int(y+h)
}

// DragAndDrop - Drags the image
//func DragAndDrop(i interface{}) {
//	switch i.(type) {
//	case CasinoCards:
//		c := i.(CasinoCards)
//
//		if IsAreaHovered(c) && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
//			c.SetDraggedState(true)`
//		}
//
//		// drag and set location
//		if inpututil.MouseButtonPressDuration(ebiten.MouseButtonLeft) > 0 && c.GetDraggedState() {
//			cx, cy := ebiten.CursorPosition()
//			_, _, w, h := c.GetPosition()
//			c.SetLocation(X, float64(cx)-w/2)
//			c.SetLocation(Y, float64(cy)-h/2)
//		}
//
//		// release
//		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
//			c.SetDraggedState(false)
//		}
//	}
//}
