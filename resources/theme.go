package resources

import (
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
}

// NewTheme - returns data about the current frame dimensions, related to what Theme is being used
func NewTheme() *Theme {
	th := &Theme{

		// The images used for the themes
		Sources: map[string]*ebiten.Image{
			PixelatedTheme: ebiten.NewImageFromImage(LoadSpriteImage("resources/assets/cards/8BitDeckAssets.png")),
			ClassicTheme:   ebiten.NewImageFromImage(LoadSpriteImage("resources/assets/cards/classic-solitaire.png")),
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
				//res.StaticBack2:   []int{0, 480, 71, 96, 0},
				//res.VYnamicCastle: []int{71, 480, 71, 96, 2},
				//res.VYnamicBeach:  []int{213, 480, 71, 96, 3},
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
	}
	return th
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