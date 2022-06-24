package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"image/color"
)

type Icon struct {
	Img           *ebiten.Image
	X, Y, W, H    int
	ID            string
	RevealedState bool
	RemovedState  bool
}

func (ic *Icon) DrawIcon(screen *ebiten.Image) {
	img := ic.Img

	// if the image is hidden
	if !ic.IsRevealed() {
		img = ebiten.NewImage(img.Size())
		img.Fill(color.NRGBA{R: 120, G: 175, B: 175, A: 255})
	}
	if ic.IsRemoved() {
		img = ebiten.NewImage(img.Size())
		img.Fill(color.NRGBA{R: 240, G: 240, B: 240, A: 255})
	}

	// sprite image
	opi := &ebiten.DrawImageOptions{}
	if !ic.IsRemoved() {
		opi.GeoM.Scale(ScX, ScY)
		opi.GeoM.Translate(float64(ic.X), float64(ic.Y))
	} else {
		opi.GeoM.Scale(0, 0)
		opi.GeoM.Translate(0, 0)
	}

	screen.DrawImage(img, opi)
}

func (ic *Icon) IsRevealed() bool {
	return ic.RevealedState
}

func (ic *Icon) IsRemoved() bool {
	return ic.RemovedState
}

func (ic *Icon) SetRevealedState(val bool) {
	ic.RevealedState = val
}

func (ic *Icon) SetRemovedState(val bool) {
	ic.RemovedState = val
}

// HitBox - returns the shape in image.Rectangle format
func (ic *Icon) HitBox() image.Rectangle {
	return image.Rect(ic.X, ic.Y, ic.X+ic.W, ic.Y+ic.H)
}
