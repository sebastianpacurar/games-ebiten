package game

import (
	u "games-ebiten/resources/utils"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	// ClassicTheme - represents a Theme reference key
	// PixelatedTheme - represents a Theme reference key
	ClassicTheme   = "classic"
	PixelatedTheme = "8bit"
)

// Theme holds data about the correct frame dimensions, and the correct images to draw.
// It is used more as a helper to handle the Deck recreation quicker
type Theme struct {
	Sources         map[string]*ebiten.Image
	FrameBoundaries map[string]map[string]int
	SuitsOrder      map[string][]string
	ScaleValue      map[string]map[string]float64

	// LocMultiplier is used to properly compute image locations in case of a grid display.
	// It is used only when GenerateDeck runs, so the images won't be so spaced out between them.
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
			PixelatedTheme: ebiten.NewImageFromImage(u.LoadSpriteImage("resources/images/cards/8BitDeckAssets.png")),
			ClassicTheme:   ebiten.NewImageFromImage(u.LoadSpriteImage("resources/images/cards/classic-solitaire.png")),
		},

		// The Frame dimensions for the themes. data is stored in this order: FrameOX, FrameOY, FrameWidth, FrameHeight
		FrameBoundaries: map[string]map[string]int{
			PixelatedTheme: {u.FrameOX: 0, u.FrameOY: 0, u.FrameW: 35, u.FrameH: 47},
			ClassicTheme:   {u.FrameOX: 0, u.FrameOY: 0, u.FrameW: 71, u.FrameH: 96},
		},

		// The Sub Images of the Main Image are different from one theme to another
		SuitsOrder: map[string][]string{
			PixelatedTheme: {u.Hearts, u.Clubs, u.Diamonds, u.Spades},
			ClassicTheme:   {u.Spades, u.Hearts, u.Clubs, u.Diamonds},
		},

		ScaleValue: map[string]map[string]float64{
			PixelatedTheme: {
				u.X: 2,
				u.Y: 2,
			},
			ClassicTheme: {
				u.X: 1,
				u.Y: 1,
			},
		},

		// The value which will be multiplied with either LX or LY, based on the given scenario
		LocMultiplier: map[string]map[string]float64{
			PixelatedTheme: {
				u.X: 3,
				u.Y: 3,
			},
			ClassicTheme: {
				u.X: 1.5,
				u.Y: 1.5,
			},
		},

		// defaults to Pixelated Theme
		Active: PixelatedTheme,

		// used to see if the state of Active has changed
		LastActive: PixelatedTheme,
	}
}

// GetFrameBoundaries - returns 4 integer values which are: FrameOX, FrameOY, FrameWidth, FrameHeight
func (th *Theme) GetFrameBoundaries(active string) (int, int, int, int) {
	activeTh := th.FrameBoundaries[active]
	return activeTh[u.FrameOX], activeTh[u.FrameOY], activeTh[u.FrameW], activeTh[u.FrameH]
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
