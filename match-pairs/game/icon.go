package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

type Icon struct {
	Img        *ebiten.Image
	X, Y, W, H float64
}

func (ic *Icon) DrawIcon(screen *ebiten.Image) {
	//border
	border := ebiten.NewImage(int(ic.W)+10, int(ic.H)+10)
	border.Fill(color.NRGBA{R: 120, G: 175, B: 175, A: 255})
	opBorder := &ebiten.DrawImageOptions{}
	opBorder.GeoM.Scale(1.5, 1.5)
	opBorder.GeoM.Translate(ic.X-7.5, ic.Y-7.5)
	screen.DrawImage(border, opBorder)

	// sprite image
	opi := &ebiten.DrawImageOptions{}
	opi.GeoM.Scale(1.5, 1.5)
	opi.GeoM.Translate(ic.X, ic.Y)
	screen.DrawImage(ic.Img, opi)
}
