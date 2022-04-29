package game

import (
	u "games-ebiten/resources/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

const (
	ItemFrameOX     = 4
	ItemFrameOY     = 0
	ItemFrameWidth  = 79
	ItemFrameHeight = 80
	ItemScale       = 0.65
)

// Item - implements the StaticSprite interface
type Item struct {
	Img      *ebiten.Image
	ImgCount int
	RowCount int
	X, Y     float64
	W, H     float64
	HitBox   map[string]float64
}

func (i *Item) GetLocations() (float64, float64) {
	return i.X, i.Y
}

func (i *Item) GetSize() (float64, float64) {
	return i.W, i.H
}

func (i *Item) SetLocation(axis string, val float64) {
	if axis == u.X {
		i.X = val
	} else if axis == u.Y {
		i.Y = val
	}
}

func (i *Item) DrawStaticSprite(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(ItemScale, ItemScale)
	op.GeoM.Translate(i.X, i.Y)

	x, y := ItemFrameOX+i.ImgCount*ItemFrameWidth, ItemFrameOY+i.RowCount*ItemFrameHeight
	screen.DrawImage(i.Img.SubImage(image.Rect(x, y, x+ItemFrameWidth, y+ItemFrameHeight)).(*ebiten.Image), op)
}

// UpdateItemState - mainly used to update location and sub image after every collision with the player or NPCs
func (i *Item) UpdateItemState() {
	i.X, i.Y = u.GenerateRandomLocation(u.ScreenDims[u.MinX], u.ScreenDims[u.MaxX]-ItemFrameWidth, u.ScreenDims[u.MinY], u.ScreenDims[u.MaxY]-ItemFrameHeight)
	i.ImgCount++

	if i.RowCount == 2 && i.ImgCount == 6 {
		// if the last image has been rendered, reset RowCount and ImgCount
		i.RowCount = 0
		i.ImgCount = 0
	} else if i.ImgCount == 6 {
		// else if ImgCount reaches limit+1, increase RowCount and reset ImgCount
		i.RowCount++
		i.ImgCount = 0
	}
}
