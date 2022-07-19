package animation_movement

import (
	res "games-ebiten/resources"
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
	Img        *ebiten.Image
	ImgCount   int
	RowCount   int
	X, Y, W, H int
}

func (i *Item) HitBox() image.Rectangle {
	return image.Rect(i.X, i.Y, i.X+i.W, i.Y+i.H)
}

func (i *Item) SetLocation(axis string, val int) {
	if axis == res.X {
		i.X = val
	} else if axis == res.Y {
		i.Y = val
	}
}

func (i *Item) DrawSprite(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(ItemScale, ItemScale)
	op.GeoM.Translate(float64(i.X), float64(i.Y))

	x, y := ItemFrameOX+i.ImgCount*ItemFrameWidth, ItemFrameOY+i.RowCount*ItemFrameHeight
	screen.DrawImage(i.Img.SubImage(image.Rect(x, y, x+ItemFrameWidth, y+ItemFrameHeight)).(*ebiten.Image), op)
}

// UpdateItemState - mainly used to update location and sub image after every collision with the player or NPCs
func (i *Item) UpdateItemState() {
	i.X, i.Y = res.GenerateRandomPosition(0, 0, res.ScreenWidth-ItemFrameWidth, res.ScreenHeight-ItemFrameHeight)
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
