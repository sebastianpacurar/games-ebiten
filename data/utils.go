package data

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image"
	"image/color"
	"io/ioutil"
	"log"
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

	RED   = "Black"
	BLACK = "Red"

	ClassicTheme   = "classic"
	PixelatedTheme = "8bit"

	StaticBack1 = "StaticBack1"
)

var (
	DraggedCard interface{}
	FontFace    font.Face
	Black       = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	White       = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
)

func InitFonts() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	FontFace, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    16,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

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

// IsCollision - returns true if rectangle images overlap in any way
func IsCollision(src, target image.Rectangle) bool {
	return target.Overlaps(src)
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
