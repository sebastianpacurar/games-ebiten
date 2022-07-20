package match_pairs

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"image/color"
)

type Icon struct {
	Img           *ebiten.Image
	X, Y, W, H    int
	ID            string
	revealedState bool
	removedState  bool
}

func (ic *Icon) DrawIcon(screen *ebiten.Image) {
	img := ic.Img

	// if the image is hidden
	if !ic.Revealed() {
		img = ebiten.NewImage(img.Size())
		img.Fill(color.NRGBA{R: 120, G: 175, B: 175, A: 255})
	}
	if ic.Removed() {
		img = ebiten.NewImage(img.Size())
		img.Fill(color.NRGBA{R: 240, G: 240, B: 240, A: 255})
	}

	// sprite image
	opi := &ebiten.DrawImageOptions{}
	if !ic.Removed() {
		opi.GeoM.Scale(ScX, ScY)
		opi.GeoM.Translate(float64(ic.X), float64(ic.Y))
	} else {
		opi.GeoM.Scale(0, 0)
		opi.GeoM.Translate(0, 0)
	}

	screen.DrawImage(img, opi)
}

func (ic *Icon) Revealed() bool {
	return ic.revealedState
}

func (ic *Icon) Removed() bool {
	return ic.removedState
}

func (ic *Icon) SetRevealed(val bool) {
	ic.revealedState = val
}

func (ic *Icon) SetRemoved(val bool) {
	ic.removedState = val
}

// HitBox - returns the shape in image.Rectangle format
func (ic *Icon) HitBox() image.Rectangle {
	return image.Rect(ic.X, ic.Y, ic.X+ic.W, ic.Y+ic.H)
}

func (ic *Icon) Hovered(cx, cy int) bool {
	return image.Pt(cx, cy).In(ic.HitBox())
}
