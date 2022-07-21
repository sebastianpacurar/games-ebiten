package resources

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"math/rand"
)

const (
	ScreenWidth  = 960
	ScreenHeight = 960

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
	ActiveGame       interface{}
	DraggedCard      interface{}
	ActiveCardsTheme string
	FontFace         font.Face
	MainMenuH        int
	AnimInProgress   bool
	Black            = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	White            = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	Green            = color.NRGBA{R: 0, G: 255, B: 0, A: 255}
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

// GenerateRandomPosition - generate a random location for X and y
// (x is between 0 and ScreenWidth) and (y is between 0 and ScreenHeight)
func GenerateRandomPosition(minX, minY, maxX, maxY int) (int, int) {
	return rand.Intn(maxX-minX) + minX, rand.Intn(maxY-minY) + minY
}

// IsCollision - returns true if rectangle images overlap in any way
func IsCollision(src, target image.Rectangle) bool {
	return target.Overlaps(src)
}

// BoundaryValidation - prevents characters to move out of the view, in any of the 4 directions
func BoundaryValidation(i interface{}, minX, maxX, minY, maxY int) {
	switch i.(type) {
	case Sprite:
		img := i.(Sprite)
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
	isHovered := false
	switch i.(type) {
	case Hoverable:
		isHovered = i.(Hoverable).Hovered(ebiten.CursorPosition())
	}
	return isHovered
}

// FlexRowQuadrants - splits the screen in x quadrants and returns a rect for each of them (works only on X-Axis)
func FlexRowQuadrants(cols int) map[int]image.Rectangle {
	quads := make(map[int]image.Rectangle)
	unit := ScreenWidth / cols

	for i := 0; i < cols; i++ {
		minX := unit * i
		maxX := unit * (i + 1)
		quads[i] = image.Rect(minX, 0, maxX, ScreenHeight)
	}
	return quads
}

// GridQuadrants - splits the screen into a grid view
func GridQuadrants(rows, cols int) map[int]map[int]image.Rectangle {
	quads := make(map[int]map[int]image.Rectangle)
	xUnit := ScreenWidth / rows
	yUnit := ScreenHeight / cols

	for i := 0; i < rows; i++ {
		minX := xUnit * i
		maxX := xUnit * (i + 1)
		rowQuads := make(map[int]image.Rectangle)
		for j := 0; j < cols; j++ {
			minY := yUnit*j + MainMenuH
			maxY := yUnit * (j + 1)
			rowQuads[j] = image.Rect(minX, minY, maxX, maxY)
		}
		quads[i] = rowQuads
	}

	return quads
}

// CenterOnX - centers the item on the X-axis within the given quadrant
func CenterOnX(width int, quad image.Rectangle) int {
	return (quad.Min.X + quad.Dx()/2) - width/2
}

// CenterOnY - centers the item on the Y-axis within the given quadrant
func CenterOnY(height int, quad image.Rectangle) int {
	return quad.Min.Y + quad.Dy()/2 - height/2
}
