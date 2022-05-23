package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

type Icon struct {
	Img        *ebiten.Image
	X, Y, W, H float64
	ID         string
	IsRevealed bool
}

func (ic *Icon) DrawIcon(screen *ebiten.Image) {
	img := ic.Img

	// if the image is hidden
	if !ic.IsRevealed {
		img = ebiten.NewImage(img.Size())
		img.Fill(color.NRGBA{R: 120, G: 175, B: 175, A: 255})
	}

	// sprite image
	opi := &ebiten.DrawImageOptions{}
	opi.GeoM.Scale(ScX, ScY)
	opi.GeoM.Translate(ic.X, ic.Y)
	screen.DrawImage(img, opi)
}

func (ic *Icon) GetRevealState() bool {
	return true
}

func (ic *Icon) SetRevealState(val bool) {
	ic.IsRevealed = val
}

// GetPosition - returns the icon position (x, y, w, h)
func (ic *Icon) GetPosition() (float64, float64, float64, float64) {
	return ic.X, ic.Y, ic.W, ic.H
}