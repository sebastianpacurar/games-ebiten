package data

import "github.com/hajimehoshi/ebiten/v2"

type Food struct {
	Img         *ebiten.Image
	ImgNo       int
	LocX, LocY  float64
	HitBox      map[string]float64
	IsDisplayed bool
}
