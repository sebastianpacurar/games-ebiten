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

//type Border struct {
//	X, Y, W, H float64
//	SX, SY     float64
//}

func (ic *Icon) DrawIcon(screen *ebiten.Image) {
	//border
	//border := ebiten.NewImage(int(ic.Border.W), int(ic.Border.H))
	//border.Fill(color.NRGBA{R: 120, G: 175, B: 175, A: 255})
	//opBorder := &ebiten.DrawImageOptions{}
	//opBorder.GeoM.Scale(ic.Border.SX, ic.Border.SY)
	//opBorder.GeoM.Translate(ic.Border.X, ic.Border.Y)
	//screen.DrawImage(border, opBorder)

	img := ebiten.NewImage(int(ic.W), int(ic.H))
	img.Fill(color.NRGBA{R: 120, G: 175, B: 175, A: 255})

	if ic.IsRevealed {
		img = ic.Img
	}

	// sprite image
	opi := &ebiten.DrawImageOptions{}
	opi.GeoM.Scale(1.5, 1.5)
	opi.GeoM.Translate(ic.X, ic.Y)
	screen.DrawImage(img, opi)
}

func (ic *Icon) GetRevealState() bool {
	return true
}

func (ic *Icon) SetRevealState(val bool) {
	ic.IsRevealed = val
}

// GetPosition - returns the icon included in the border area
func (ic *Icon) GetPosition() (float64, float64, float64, float64) {
	return ic.X, ic.Y, ic.W, ic.H
}
