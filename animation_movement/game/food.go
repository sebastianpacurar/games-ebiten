package game

import (
	u "games-ebiten/resources/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

const (
	FoodFrameOX     = 4
	FoodFrameOY     = 0
	FoodFrameWidth  = 79
	FoodFrameHeight = 80
	FoodScale       = 0.65
)

type Food struct {
	Img      *ebiten.Image
	ImgCount int
	RowCount int
	LX, LY   float64
	W, H     float64
	HitBox   map[string]float64
}

func (f *Food) DrawImage(screen *ebiten.Image) {
	opFood := &ebiten.DrawImageOptions{}
	opFood.GeoM.Scale(FoodScale, FoodScale)
	opFood.GeoM.Translate(f.LX, f.LY)

	// load every sub image from left to right from the first row
	x, y := FoodFrameOX+f.ImgCount*FoodFrameWidth, FoodFrameOY+f.RowCount*FoodFrameHeight
	screen.DrawImage(f.Img.SubImage(image.Rect(x, y, x+FoodFrameWidth, y+FoodFrameHeight)).(*ebiten.Image), opFood)
}

// UpdateFoodState - mainly used to update location and sub image after every collision with the player or NPCs
func (f *Food) UpdateFoodState() {
	f.LX, f.LY = u.GenerateRandomLocation(Bounds[u.MinX], Bounds[u.MaxX]-FoodFrameWidth, Bounds[u.MinY], Bounds[u.MaxY]-FoodFrameHeight)
	f.ImgCount++

	if f.RowCount == 2 && f.ImgCount == 6 {
		// if the last image has been rendered, reset RowCount and ImgCount
		f.RowCount = 0
		f.ImgCount = 0
	} else if f.ImgCount == 6 {
		// else if ImgCount reaches limit+1, increase RowCount and reset ImgCount
		f.RowCount++
		f.ImgCount = 0
	}
}
