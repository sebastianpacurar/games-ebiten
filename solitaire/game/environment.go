package game

import (
	u "games-ebiten/resources/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

type (
	Environment struct {
		Quadrants       map[int]image.Rectangle
		Deck            []*Card
		BgImg           *ebiten.Image
		EmptySlotImg    *ebiten.Image
		Columns         []CardColumn
		FoundationPiles []FoundationPile
		WastePile
		StockPile
		SpacerV   int
		IsVegas   bool
		DrawCount int
		//W, H - stand for default width and height
		W, H int
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
		X, Y  int
		Cards []*Card
	}

	CardColumn struct {
		X, Y, W, H int
		Cards      []*Card
	}
)

// IsStockPileHovered - returns true if the StockPile is hovered - used to click on green Circle to redraw
func (e *Environment) IsStockPileHovered(cx, cy int) bool {
	x, y, w, h := e.StockPile.X, e.StockPile.Y, e.StockPile.W, e.StockPile.H
	return image.Pt(cx, cy).In(image.Rect(x, y, x+w, y+h))
}

func (e *Environment) GetGeomData(i interface{}) image.Rectangle {
	rect := image.Rectangle{}
	switch i.(type) {
	case FoundationPile:
		area := i.(FoundationPile)
		rect = image.Rect(area.X, area.Y, area.X+area.W, area.Y+area.H)
	case StockPile:
		area := i.(StockPile)
		rect = image.Rect(area.X, area.Y, area.X+area.W, area.Y+area.H)
	case CardColumn:
		area := i.(CardColumn)
		rect = image.Rect(area.X, area.Y, area.X+area.W, area.Y+area.H)
	}
	return rect
}

func (e *Environment) DrawPlayground(screen *ebiten.Image, th *Theme) {
	// Draw the BG Image
	opBg := &ebiten.DrawImageOptions{}
	opBg.GeoM.Scale(50, 50)
	screen.DrawImage(e.BgImg, opBg)
	envTh := th.EnvScaleValue[th.Active]

	var img *ebiten.Image

	// Draw the Stock Slot
	opStockSlot := &ebiten.DrawImageOptions{}
	opStockSlot.GeoM.Scale(envTh[u.X], envTh[u.Y])

	if th.Active == u.PixelatedTheme {
		opStockSlot.GeoM.Translate(float64(u.CenterItem(e.W, e.Quadrants[0]))+3.5, float64(e.SpacerV)+3.5)
	} else {
		opStockSlot.GeoM.Translate(float64(u.CenterItem(e.W, e.Quadrants[0])), float64(e.SpacerV))
	}

	if e.IsVegas {
		e.DrawCount++
	}

	if e.IsVegas && e.DrawCount == 3 {
		img = e.RedSlotImg
	} else {
		img = e.GreenSlotImg
	}
	screen.DrawImage(img, opStockSlot)

	// Draw the Foundation Slots
	for i := 0; i < 4; i++ {
		opFoundationSlot := &ebiten.DrawImageOptions{}
		opFoundationSlot.GeoM.Scale(envTh[u.X], envTh[u.Y])

		if th.Active == u.PixelatedTheme {
			opFoundationSlot.GeoM.Translate(float64(u.CenterItem(e.W, e.Quadrants[3+i]))+3.5, float64(e.SpacerV)+3.5)
		} else {
			opFoundationSlot.GeoM.Translate(float64(u.CenterItem(e.W, e.Quadrants[3+i])), float64(e.SpacerV))
		}
		screen.DrawImage(e.EmptySlotImg, opFoundationSlot)
	}
}

func (e *Environment) DrawEnding(screen *ebiten.Image) {
	opBg := &ebiten.DrawImageOptions{}
	opBg.GeoM.Scale(50, 50)
	screen.DrawImage(e.BgImg, opBg)
}

func (e *Environment) UpdateEnv() {
	e.W = e.Deck[0].W
	e.H = e.Deck[0].H
	e.Quadrants = u.GetFlexboxQuadrants(7)

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
	e.WastePile.X = u.CenterItem(e.W, e.Quadrants[1])

	e.StockPile.X = u.CenterItem(e.W, e.Quadrants[0])
	e.StockPile.Y = e.SpacerV
	e.StockPile.W = e.W
	e.StockPile.H = e.H

	// starts from the fourth quadrant
	for s := range e.FoundationPiles {
		fx := u.CenterItem(e.W, e.Quadrants[3+s])
		e.FoundationPiles[s].X = fx
		e.FoundationPiles[s].Y = e.SpacerV
		e.FoundationPiles[s].W = e.W
		e.FoundationPiles[s].H = e.H
	}

	// fill every column array with its relative count of cards and save GeomData of columns placeholders
	cardIndex := 0
	for i := range e.Columns {
		// initiate the location of the Card Column placeholders
		colx := u.CenterItem(e.W, e.Quadrants[0+i])
		coly := e.Quadrants[0+i].Max.Y / 3
		e.Columns[i].X = colx
		e.Columns[i].Y = coly
		e.Columns[i].W = e.W
		e.Columns[i].H = e.H

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

	// fill the StockPile array
	for i := range e.Deck[cardIndex:] {
		e.StockPile.Cards = append(e.StockPile.Cards, e.Deck[cardIndex:][i])
	}
}

// IsGameOver - if count of cards in FoundationPiles is 52, return true
func (e *Environment) IsGameOver() bool {
	total := 0
	for _, store := range e.FoundationPiles {
		total += len(store.Cards)
	}
	return total == 52
}
