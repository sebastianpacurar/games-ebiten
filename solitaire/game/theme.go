package game

import (
	u "games-ebiten/resources/utils"
	"github.com/hajimehoshi/ebiten/v2"
)

// Theme holds data about the correct frame dimensions, and the correct images to draw.
// It is used more as a helper to handle the Deck recreation quicker
type Theme struct {
	Sources            map[string]*ebiten.Image
	FrontFaceFrameData map[string]map[string]int
	BackFaceFrameData  map[string]map[string][]int
	SuitsOrder         map[string][]string
	ScaleValue         map[string]map[string]float64

	// LocMultiplier is used to properY compute image locations in case of a grid display.
	// It is used onY when GenerateDeck runs, so the images won't be so spaced out between them.
	LocMultiplier map[string]map[string]float64

	// Active and LastActive are used to track down when the Theme Changing gets triggered
	Active     string
	LastActive string
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

		// The Frame Dimensions of the available back faces of the current Theme.
		// Stored in the form of: FrOX, FrOY, FrW, FrH, FrC
		BackFaceFrameData: map[string]map[string][]int{
			u.PixelatedTheme: {
				u.StaticBack1: []int{0, 0, 35, 47, 0},
			},
			u.ClassicTheme: {
				u.StaticBack1:   []int{0, 384, 71, 96, 0},
				u.StaticBack2:   []int{0, 480, 71, 96, 0},
				u.DynamicCastle: []int{71, 480, 71, 96, 2},
				u.DynamicBeach:  []int{213, 480, 71, 96, 3},
			},
		},

		// The Sub Images of the Main Image are different from one theme to another
		SuitsOrder: map[string][]string{
			u.PixelatedTheme: {u.Hearts, u.Clubs, u.Diamonds, u.Spades},
			u.ClassicTheme:   {u.Spades, u.Hearts, u.Clubs, u.Diamonds},
		},

		ScaleValue: map[string]map[string]float64{
			u.PixelatedTheme: {
				u.X: 2,
				u.Y: 2,
			},
			u.ClassicTheme: {
				u.X: 1.5,
				u.Y: 1.5,
			},
			u.SimpleTheme: {
				u.X: 1,
				u.Y: 1,
			},
		},

		// The value which will be multiplied with either X or Y, based on the given scenario
		LocMultiplier: map[string]map[string]float64{
			u.PixelatedTheme: {
				u.X: 3,
				u.Y: 3,
			},
			u.ClassicTheme: {
				u.X: 1.25,
				u.Y: 1.25,
			},
		},

		// defaults to Classic Theme
		Active: u.ClassicTheme,

		// used to see if the state of Active has changed
		LastActive: u.ClassicTheme,
	}
}

// GetFrontFrameGeomData - returns 4 integer values which are: FrOX, FrOY, FrameWidth, FrameHeight
func (th *Theme) GetFrontFrameGeomData(active string) (int, int, int, int) {
	activeTh := th.FrontFaceFrameData[active]
	return activeTh[u.FrOX], activeTh[u.FrOY], activeTh[u.FrW], activeTh[u.FrH]
}

// GetBackFrameGeomData - returns 4 integer values which are: FrOX, FrOY, FrameWidth, FrameHeight
func (th *Theme) GetBackFrameGeomData(active, backFace string) []int {
	return th.BackFaceFrameData[active][backFace]
}

// GetSource - returns the source of the image. Useful to toggle between Themes
func (th *Theme) GetSource(active string) *ebiten.Image {
	return th.Sources[active]
}

func (th *Theme) GetScaleValue(active string) (float64, float64) {
	return th.ScaleValue[active][u.X], th.ScaleValue[active][u.Y]
}

// IsToggled - whenever the "Change Theme" gets triggered, should return true
func (th *Theme) IsToggled() bool {
	return th.LastActive != th.Active
}
