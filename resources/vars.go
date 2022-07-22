package resources

import (
	"golang.org/x/image/font"
	"image/color"
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

// Translation - used for acquiring the right card index while getting the card SubImage from the Image
// CardRanks - smallest is "Ace"(0), while highest is "King"(12)
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
	Translation      = map[string]map[int]string{
		PixelatedTheme: {
			1: "2", 2: "3", 3: "4", 4: "5", 5: "6", 6: "7",
			7: "8", 8: "9", 9: "10", 10: "J", 11: "Q", 12: "K", 13: "A",
		},
		ClassicTheme: {
			0: "A", 1: "2", 2: "3", 3: "4", 4: "5", 5: "6", 6: "7",
			7: "8", 8: "9", 9: "10", 10: "J", 11: "Q", 12: "K",
		},
	}
	CardRanks = map[string]int{
		Ace: 0, "2": 1, "3": 2, "4": 3, "5": 4, "6": 5, "7": 6,
		"8": 7, "9": 8, "10": 9, Jack: 10, Queen: 11, King: 12,
	}
)
