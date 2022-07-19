package card_games

import (
	res "games-ebiten/resources"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

// Translation - used for acquiring the right card index while getting the card SubImage from the Image
// CardRanks - smallest is "Ace"(0), while highest is "King"(12)
var (
	Translation = map[string]map[int]string{
		res.PixelatedTheme: {
			1: "2", 2: "3", 3: "4", 4: "5", 5: "6", 6: "7",
			7: "8", 8: "9", 9: "10", 10: "J", 11: "Q", 12: "K", 13: "A",
		},
		res.ClassicTheme: {
			0: "A", 1: "2", 2: "3", 3: "4", 4: "5", 5: "6", 6: "7",
			7: "8", 8: "9", 9: "10", 10: "J", 11: "Q", 12: "K",
		},
	}
	CardRanks = map[string]int{
		res.Ace: 0, "2": 1, "3": 2, "4": 3, "5": 4, "6": 5, "7": 6,
		"8": 7, "9": 8, "10": 9, res.Jack: 10, res.Queen: 11, res.King: 12,
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
			res.PixelatedTheme: ebiten.NewImageFromImage(res.LoadSpriteImage("resources/assets/cards/8BitDeckAssets.png")),
			res.ClassicTheme:   ebiten.NewImageFromImage(res.LoadSpriteImage("resources/assets/cards/classic-solitaire.png")),
		},

		// The Frame dimensions for the themes. data is stored in this order: FrOX, FrOY, FrameWidth, FrameHeight
		FrontFaceFrameData: map[string]map[string]int{
			res.PixelatedTheme: {res.FrOX: 0, res.FrOY: 0, res.FrW: 35, res.FrH: 47},
			res.ClassicTheme:   {res.FrOX: 0, res.FrOY: 0, res.FrW: 71, res.FrH: 96},
		},

		// The Frame Dimensions of the available back faces
		// Stored in the form of: FrOX, FrOY, FrW, FrH, FrC
		BackFaceFrameData: map[string]map[string][]int{
			res.PixelatedTheme: {
				res.StaticBack1: []int{0, 0, 35, 47, 0},
			},
			res.ClassicTheme: {
				res.StaticBack1: []int{0, 384, 71, 96, 0},
				//res.StaticBack2:   []int{0, 480, 71, 96, 0},
				//res.VYnamicCastle: []int{71, 480, 71, 96, 2},
				//res.VYnamicBeach:  []int{213, 480, 71, 96, 3},
			},
		},

		// The Sub Images of the Main Image are different from one theme to another
		SuitsOrder: map[string][]string{
			res.PixelatedTheme: {res.Hearts, res.Clubs, res.Diamonds, res.Spades},
			res.ClassicTheme:   {res.Spades, res.Hearts, res.Clubs, res.Diamonds},
		},

		CardScaleValue: map[string]map[string]float64{
			res.PixelatedTheme: {
				res.X: 2,
				res.Y: 2,
			},
			res.ClassicTheme: {
				res.X: 1,
				res.Y: 1,
			},
		},

		EnvScaleValue: map[string]map[string]float64{
			res.PixelatedTheme: {
				res.X: 0.9,
				res.Y: 0.9,
			},
			res.ClassicTheme: {
				res.X: 1,
				res.Y: 1,
			},
		},

		// defaults to Classic Theme
		Active: res.ClassicTheme,
	}
}

// GetFrontFrameGeomData - returns 4 integer values which are: FrOX, FrOY, FrameWidth, FrameHeight
func (th *Theme) GetFrontFrameGeomData(active string) image.Rectangle {
	ath := th.FrontFaceFrameData[active]
	return image.Rect(ath[res.FrOX], ath[res.FrOY], ath[res.FrOX]+ath[res.FrW], ath[res.FrOY]+ath[res.FrH])
}

// GetBackFrameGeomData - returns 4 integer values which are: FrOX, FrOY, FrameWidth, FrameHeight
// TODO: consolidate with FrontFrameGeomData
func (th *Theme) GetBackFrameGeomData(active, backFace string) []int {
	return th.BackFaceFrameData[active][backFace]
}
