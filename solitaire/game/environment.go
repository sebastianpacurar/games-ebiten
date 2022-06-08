package game

import (
	u "games-ebiten/resources/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

type (
	Environment struct {
		Deck            []*Card
		BgImg           *ebiten.Image
		EmptySlotImg    *ebiten.Image
		Columns         []CardColumn
		FoundationPiles []FoundationPile
		WastePile
		StockPile
		SpacerV   int
		SpacerH   int
		CardFullH int
		IsVegas   bool
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
		x := (int(float64(th.FrontFaceFrameData[th.Active][u.FrW])*th.ScaleValue[th.Active][u.X]) + e.SpacerH) * (i + 4)

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

func (e *Environment) UpdateEnv(th *Theme) {
	//c := e.Deck[0]
	e.StockPile.Cards = make([]*Card, 0, 24)
	e.WastePile.Cards = make([]*Card, 0, 24)
	e.FoundationPiles = []FoundationPile{
		{Cards: make([]*Card, 0, 13)},
		{Cards: make([]*Card, 0, 13)},
		{Cards: make([]*Card, 0, 13)},
		{Cards: make([]*Card, 0, 13)},
	}
	e.Columns = []CardColumn{
		{Cards: make([]*Card, 0, 1)},
		{Cards: make([]*Card, 0, 2)},
		{Cards: make([]*Card, 0, 3)},
		{Cards: make([]*Card, 0, 4)},
		{Cards: make([]*Card, 0, 5)},
		{Cards: make([]*Card, 0, 6)},
		{Cards: make([]*Card, 0, 7)},
	}
	frame := th.GetFrontFrameGeomData(th.Active)
	x, y := frame.Dx()+e.SpacerH, e.SpacerV
	w, h := int(float64(frame.Dx())*th.ScaleValue[th.Active][u.X]), int(float64(frame.Dy())*th.ScaleValue[th.Active][u.Y])

	e.CardFullH = h

	e.StockPile.X = x
	e.StockPile.Y = y
	e.StockPile.W = w
	e.StockPile.H = h

	for s := range e.FoundationPiles {
		sx := (frame.Dx() + e.SpacerH) * (s + 4)
		e.FoundationPiles[s].X = sx
		e.FoundationPiles[s].Y = y
		e.FoundationPiles[s].W = w
		e.FoundationPiles[s].H = h
	}

	// fill every column array with its relative count of cards and save GeomData of columns placeholders
	cardIndex := 0
	for i := range e.Columns {
		// initiate the location of the Card Column placeholders
		colx := (frame.Dx() + e.SpacerH) * (i + 1)
		coly := u.ScreenHeight / 3
		e.Columns[i].X = colx
		e.Columns[i].Y = coly
		e.Columns[i].W = w
		e.Columns[i].H = h

		for j := 0; j <= i; j++ {
			// keep only the last one revealed
			if j == i {
				e.Deck[cardIndex].SetRevealedState(true)
			}
			e.Deck[cardIndex].ColNum = i + 1
			e.Columns[i].Cards = append(e.Columns[i].Cards, e.Deck[cardIndex])
			cardIndex++
		}
	}

	// fill the DrawCard array
	for i := range e.Deck[cardIndex:] {
		e.StockPile.Cards = append(e.StockPile.Cards, e.Deck[cardIndex:][i])
	}
}
