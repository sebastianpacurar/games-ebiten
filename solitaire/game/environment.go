package game

import (
	u "games-ebiten/resources/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

type (
	Environment struct {
		BgImg           *ebiten.Image
		EmptySlotImg    *ebiten.Image
		Columns         []CardColumn
		FoundationPiles []FoundationPile
		WastePile
		StockPile
		SpacerV   int
		SpacerH   int
		CardFullH int
	}

	StockPile struct {
		X, Y, W, H   int
		GreenSlotImg *ebiten.Image
		RedSlotImg   *ebiten.Image
		Cards        []*Card
		IsGreen      bool
	}

	FoundationPile struct {
		X, Y, W, H int
		Cards      []*Card
	}

	WastePile struct {
		Cards []*Card
	}

	CardColumn struct {
		X, Y, W, H int
		Cards      []*Card
	}
)

// IsStockPileHovered - returns true if the StockPile is hovered. Applies only when there are no cards in stack
func (e *Environment) IsStockPileHovered(cx, cy int) bool {
	x, y, w, h := e.StockPile.X, e.StockPile.Y, e.StockPile.W, e.StockPile.H
	pt := image.Pt(cx, cy)
	dims := image.Rectangle{
		Min: image.Point{X: x, Y: y},
		Max: image.Point{X: x + w, Y: y + h},
	}
	return pt.In(dims)
}

func (e *Environment) GetGeomData(i interface{}) image.Rectangle {
	data := image.Rectangle{}
	switch i.(type) {
	case FoundationPile:
		area := i.(FoundationPile)
		data = image.Rectangle{
			Min: image.Point{X: area.X, Y: area.Y},
			Max: image.Point{X: area.X + area.W, Y: area.Y + area.H},
		}
	case StockPile:
		area := i.(StockPile)
		data = image.Rectangle{
			Min: image.Point{X: area.X, Y: area.Y},
			Max: image.Point{X: area.X + area.W, Y: area.Y + area.H},
		}
	case CardColumn:
		area := i.(CardColumn)
		data = image.Rectangle{
			Min: image.Point{X: area.X, Y: area.Y},
			Max: image.Point{X: area.X + area.W, Y: area.Y + area.H},
		}
	case Card:
		area := i.(Card)
		data = image.Rectangle{
			Min: image.Point{X: area.X, Y: area.Y},
			Max: image.Point{X: area.X + area.W, Y: area.Y + area.H},
		}
	}
	return data
}

func (e *Environment) DrawPlayground(screen *ebiten.Image, th *Theme) {
	// Draw the BG Image
	opBg := &ebiten.DrawImageOptions{}
	opBg.GeoM.Scale(50, 50)
	screen.DrawImage(e.BgImg, opBg)

	// Draw the Draw Stack and the Drawn Stack
	for i := 1; i <= 2; i++ {
		var img *ebiten.Image
		opCardStack := &ebiten.DrawImageOptions{}
		opCardStack.GeoM.Scale(th.ScaleValue[th.Active][u.X], th.ScaleValue[th.Active][u.Y])

		x := (th.FrontFaceFrameData[th.Active][u.FrW] + e.SpacerH) * i
		y := e.SpacerV
		opCardStack.GeoM.Translate(float64(x), float64(y))

		if i == 2 {
			img = e.EmptySlotImg
		} else {
			img = e.GreenSlotImg
		}
		screen.DrawImage(img, opCardStack)
	}

	// Draw the Card Stores
	for i := 0; i < 4; i++ {
		opStackSlot := &ebiten.DrawImageOptions{}
		opStackSlot.GeoM.Scale(th.ScaleValue[th.Active][u.X], th.ScaleValue[th.Active][u.Y])

		// align Stacked Cards with the left columns of the Column slots
		x := (th.FrontFaceFrameData[th.Active][u.FrW] + e.SpacerH) * (i + 4)

		opStackSlot.GeoM.Translate(float64(x), float64(e.SpacerV))
		screen.DrawImage(e.EmptySlotImg, opStackSlot)
	}

	// Draw the colum slots
	for i := 1; i <= 7; i++ {
		opColumnSlot := &ebiten.DrawImageOptions{}
		opColumnSlot.GeoM.Scale(th.ScaleValue[th.Active][u.X], th.ScaleValue[th.Active][u.Y])

		x := (th.FrontFaceFrameData[th.Active][u.FrW] + e.SpacerH) * i
		y := u.ScreenHeight / 3

		opColumnSlot.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(e.EmptySlotImg, opColumnSlot)
	}
}

func (e *Environment) DrawEnding(screen *ebiten.Image) {
	opBg := &ebiten.DrawImageOptions{}
	opBg.GeoM.Scale(50, 50)
	screen.DrawImage(e.BgImg, opBg)
}

// IsGameOver - if count of cards in FoundationPiles is 52, return true
func (e *Environment) IsGameOver() bool {
	total := 0
	for _, store := range e.FoundationPiles {
		total += len(store.Cards)
	}
	return total == 52
}
