package data

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

const (
	FoodFrameOX     = 0
	FoodFrameOY     = 0
	FoodFrameWidth  = 72
	FoodFrameHeight = 75
	FoodScale       = 0.75
)

type Food struct {
	Img         *ebiten.Image
	ImgNo       int
	LocX, LocY  float64
	W, H        float64
	HitBox      map[string]float64
	IsDisplayed bool
}

func (f *Food) DrawImage(screen *ebiten.Image) {
	opFood := &ebiten.DrawImageOptions{}
	opFood.GeoM.Scale(FoodScale, FoodScale)
	opFood.GeoM.Translate(f.LocX, f.LocY)

	// load every sub image from left to right from the first row
	fx, fy := FoodFrameOX+f.ImgNo*FoodFrameWidth, FoodFrameOY
	screen.DrawImage(f.Img.SubImage(image.Rect(fx, fy, fx+FoodFrameWidth, fy+FoodFrameHeight)).(*ebiten.Image), opFood)
}
