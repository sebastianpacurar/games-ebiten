package data

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

// Translation - used for acquiring the right card index while getting the card SubImage from the Image
// CardRanks - smallest is "Ace"(0), while highest is "King"(12)
var (
	Translation = map[string]map[int]string{
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

// Theme holds data about the correct frame dimensions, and the correct images to draw.
// It is used more as a helper to handle the Deck recreation quicker
type Theme struct {
	Sources            map[string]*ebiten.Image
	FrontFaceFrameData map[string]map[string]int
	BackFaceFrameData  map[string]map[string][]int
	EmptySlotFrameData map[string][]int
	SuitsOrder         map[string][]string
	CardScaleValue     map[string]map[string]float64
	EnvScaleValue      map[string]map[string]float64

	// Active represents the current theme
	Active string
}

// NewTheme - returns data about the current frame dimensions, related to what Theme is being used
func NewTheme() *Theme {
	return &Theme{

		// The images used for the themes
		Sources: map[string]*ebiten.Image{
			PixelatedTheme: ebiten.NewImageFromImage(LoadSpriteImage("assets/8BitDeckAssets.png")),
			ClassicTheme:   ebiten.NewImageFromImage(LoadSpriteImage("assets/classic-solitaire.png")),
		},

		// The Frame dimensions for the themes. data is stored in this order: FrOX, FrOY, FrameWidth, FrameHeight
		FrontFaceFrameData: map[string]map[string]int{
			PixelatedTheme: {FrOX: 0, FrOY: 0, FrW: 35, FrH: 47},
			ClassicTheme:   {FrOX: 0, FrOY: 0, FrW: 71, FrH: 96},
		},

		// The Frame Dimensions of the available back faces
		// Stored in the form of: FrOX, FrOY, FrW, FrH, FrC
		BackFaceFrameData: map[string]map[string][]int{
			PixelatedTheme: {
				StaticBack1: []int{0, 0, 35, 47, 0},
			},
			ClassicTheme: {
				StaticBack1: []int{0, 384, 71, 96, 0},
				//u.StaticBack2:   []int{0, 480, 71, 96, 0},
				//u.VYnamicCastle: []int{71, 480, 71, 96, 2},
				//u.VYnamicBeach:  []int{213, 480, 71, 96, 3},
			},
		},

		// The Sub Images of the Main Image are different from one theme to another
		SuitsOrder: map[string][]string{
			PixelatedTheme: {Hearts, Clubs, Diamonds, Spades},
			ClassicTheme:   {Spades, Hearts, Clubs, Diamonds},
		},

		CardScaleValue: map[string]map[string]float64{
			PixelatedTheme: {
				X: 2,
				Y: 2,
			},
			ClassicTheme: {
				X: 1,
				Y: 1,
			},
		},

		EnvScaleValue: map[string]map[string]float64{
			PixelatedTheme: {
				X: 0.9,
				Y: 0.9,
			},
			ClassicTheme: {
				X: 1,
				Y: 1,
			},
		},

		// defaults to Classic Theme
		Active: ClassicTheme,
	}
}

// GetFrontFrameGeomData - returns 4 integer values which are: FrOX, FrOY, FrameWidth, FrameHeight
func (th *Theme) GetFrontFrameGeomData(active string) image.Rectangle {
	ath := th.FrontFaceFrameData[active]
	return image.Rect(ath[FrOX], ath[FrOY], ath[FrOX]+ath[FrW], ath[FrOY]+ath[FrH])
}

// GetBackFrameGeomData - returns 4 integer values which are: FrOX, FrOY, FrameWidth, FrameHeight
// TODO: consolidate with FrontFrameGeomData
func (th *Theme) GetBackFrameGeomData(active, backFace string) []int {
	return th.BackFaceFrameData[active][backFace]
}
