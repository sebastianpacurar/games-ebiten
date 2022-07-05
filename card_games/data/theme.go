package data

import (
	u "games-ebiten/resources/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
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
			u.PixelatedTheme: ebiten.NewImageFromImage(u.LoadSpriteImage("resources/images/cards/8BitDeckAssets.png")),
			u.ClassicTheme:   ebiten.NewImageFromImage(u.LoadSpriteImage("resources/images/cards/classic-solitaire.png")),
		},

		// The Frame dimensions for the themes. data is stored in this order: FrOX, FrOY, FrameWidth, FrameHeight
		FrontFaceFrameData: map[string]map[string]int{
			u.PixelatedTheme: {u.FrOX: 0, u.FrOY: 0, u.FrW: 35, u.FrH: 47},
			u.ClassicTheme:   {u.FrOX: 0, u.FrOY: 0, u.FrW: 71, u.FrH: 96},
		},

		// The Frame Dimensions of the available back faces
		// Stored in the form of: FrOX, FrOY, FrW, FrH, FrC
		BackFaceFrameData: map[string]map[string][]int{
			u.PixelatedTheme: {
				u.StaticBack1: []int{0, 0, 35, 47, 0},
			},
			u.ClassicTheme: {
				u.StaticBack1: []int{0, 384, 71, 96, 0},
				//u.StaticBack2:   []int{0, 480, 71, 96, 0},
				//u.VYnamicCastle: []int{71, 480, 71, 96, 2},
				//u.VYnamicBeach:  []int{213, 480, 71, 96, 3},
			},
		},

		// The Sub Images of the Main Image are different from one theme to another
		SuitsOrder: map[string][]string{
			u.PixelatedTheme: {u.Hearts, u.Clubs, u.Diamonds, u.Spades},
			u.ClassicTheme:   {u.Spades, u.Hearts, u.Clubs, u.Diamonds},
		},

		CardScaleValue: map[string]map[string]float64{
			u.PixelatedTheme: {
				u.X: 2,
				u.Y: 2,
			},
			u.ClassicTheme: {
				u.X: 1,
				u.Y: 1,
			},
		},

		EnvScaleValue: map[string]map[string]float64{
			u.PixelatedTheme: {
				u.X: 0.9,
				u.Y: 0.9,
			},
			u.ClassicTheme: {
				u.X: 1,
				u.Y: 1,
			},
		},

		// defaults to Classic Theme
		Active: u.ClassicTheme,
	}
}

// GetFrontFrameGeomData - returns 4 integer values which are: FrOX, FrOY, FrameWidth, FrameHeight
func (th *Theme) GetFrontFrameGeomData(active string) image.Rectangle {
	ath := th.FrontFaceFrameData[active]
	return image.Rect(ath[u.FrOX], ath[u.FrOY], ath[u.FrOX]+ath[u.FrW], ath[u.FrOY]+ath[u.FrH])
}

// GetBackFrameGeomData - returns 4 integer values which are: FrOX, FrOY, FrameWidth, FrameHeight
// TODO: consolidate with FrontFrameGeomData
func (th *Theme) GetBackFrameGeomData(active, backFace string) []int {
	return th.BackFaceFrameData[active][backFace]
}
