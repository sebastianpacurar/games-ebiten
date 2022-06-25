package game

import (
	cg "games-ebiten/card_games"
	u "games-ebiten/resources/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

type (
	Environment struct {
		Quadrants       map[int]image.Rectangle
		Deck            []*cg.Card
		BgImg           *ebiten.Image
		EmptySlotImg    *ebiten.Image
		Columns         []CardColumn
		FoundationPiles []FoundationPile
		FreeCells       []FreeCell
		SpacerV         int
		W, H            int
	}

	FoundationPile struct {
		X, Y, W, H int
		Cards      []*cg.Card
	}

	FreeCell struct {
		X, Y, W, H int
		Cards      []*cg.Card
	}

	CardColumn struct {
		X, Y, W, H int
		Cards      []*cg.Card
	}
)

func (e *Environment) UpdateEnv() {
	e.W = e.Deck[0].W
	e.H = e.Deck[0].H
	e.Quadrants = u.GetFlexboxQuadrants(8)

	e.FoundationPiles = []FoundationPile{
		{Cards: make([]*cg.Card, 0, 13)},
		{Cards: make([]*cg.Card, 0, 13)},
		{Cards: make([]*cg.Card, 0, 13)},
		{Cards: make([]*cg.Card, 0, 13)},
	}
	e.FreeCells = []FreeCell{
		{Cards: make([]*cg.Card, 0, 13)},
		{Cards: make([]*cg.Card, 0, 13)},
		{Cards: make([]*cg.Card, 0, 13)},
		{Cards: make([]*cg.Card, 0, 13)},
	}
	e.Columns = []CardColumn{
		{Cards: make([]*cg.Card, 0)},
		{Cards: make([]*cg.Card, 0)},
		{Cards: make([]*cg.Card, 0)},
		{Cards: make([]*cg.Card, 0)},
		{Cards: make([]*cg.Card, 0)},
		{Cards: make([]*cg.Card, 0)},
		{Cards: make([]*cg.Card, 0)},
		{Cards: make([]*cg.Card, 0)},
	}

	// start from the first quadrant
	for i := range e.FreeCells {
		fx := u.CenterItem(e.W, e.Quadrants[0+i])
		e.FreeCells[i].X = fx
		e.FreeCells[i].Y = e.SpacerV
		e.FreeCells[i].W = e.W
		e.FreeCells[i].H = e.H
	}

	// starts from the fourth quadrant
	for i := range e.FoundationPiles {
		fx := u.CenterItem(e.W, e.Quadrants[4+i])
		e.FoundationPiles[i].X = fx
		e.FoundationPiles[i].Y = e.SpacerV
		e.FoundationPiles[i].W = e.W
		e.FoundationPiles[i].H = e.H
	}

	// fill every column array with its relative count of cards and save GeomData of columns placeholders
	cardIndex := 0
	for cardIndex < len(e.Deck)-1 {
		for i := range e.Columns {
			if cardIndex < len(e.Deck) {
				// initiate the location of the Card Column placeholders
				colx := u.CenterItem(e.W, e.Quadrants[0+i])
				coly := e.Quadrants[0+i].Max.Y / 3
				e.Columns[i].X = colx
				e.Columns[i].Y = coly
				e.Columns[i].W = e.W
				e.Columns[i].H = e.H
				e.Deck[cardIndex].ColNum = i + 1
				e.Deck[cardIndex].SetRevealedState(true)
				e.Columns[i].Cards = append(e.Columns[i].Cards, e.Deck[cardIndex])
				cardIndex++
			} else {
				break
			}
		}
	}
}

func (e *Environment) HitBox(i interface{}) image.Rectangle {
	rect := image.Rectangle{}
	switch i.(type) {
	case FoundationPile:
		area := i.(FoundationPile)
		rect = image.Rect(area.X, area.Y, area.X+area.W, area.Y+area.H)
	case CardColumn:
		area := i.(CardColumn)
		rect = image.Rect(area.X, area.Y, area.X+area.W, area.Y+area.H)
	}
	return rect
}

func (e *Environment) DrawPlayground(screen *ebiten.Image, th *cg.Theme) {
	// Draw the BG Image
	opBg := &ebiten.DrawImageOptions{}
	opBg.GeoM.Scale(50, 50)
	screen.DrawImage(e.BgImg, opBg)
	envTh := th.EnvScaleValue[th.Active]

	// Draw the FreeCell Slots
	for i := 0; i < 4; i++ {
		opFreeCell := &ebiten.DrawImageOptions{}
		opFreeCell.GeoM.Scale(envTh[u.X], envTh[u.Y])

		if th.Active == u.PixelatedTheme {
			opFreeCell.GeoM.Translate(float64(u.CenterItem(e.W, e.Quadrants[0+i]))+3.5, float64(e.SpacerV)+3.5)
		} else {
			opFreeCell.GeoM.Translate(float64(u.CenterItem(e.W, e.Quadrants[0+i])), float64(e.SpacerV))
		}
		screen.DrawImage(e.EmptySlotImg, opFreeCell)
	}

	// Draw the Foundation Slots
	for i := 0; i < 4; i++ {
		opFoundationSlot := &ebiten.DrawImageOptions{}
		opFoundationSlot.GeoM.Scale(envTh[u.X], envTh[u.Y])

		if th.Active == u.PixelatedTheme {
			opFoundationSlot.GeoM.Translate(float64(u.CenterItem(e.W, e.Quadrants[4+i]))+3.5, float64(e.SpacerV)+3.5)
		} else {
			opFoundationSlot.GeoM.Translate(float64(u.CenterItem(e.W, e.Quadrants[4+i])), float64(e.SpacerV))
		}
		screen.DrawImage(e.EmptySlotImg, opFoundationSlot)
	}
}

// MoveFromSrcToTarget - contains the logic behind moving cards from a specific Pile to another specific Pile
func (e *Environment) MoveFromSrcToTarget(src, target interface{}, i, j int) {
	switch src.(type) {

	// move FROM Column
	case []CardColumn:
		li := len(e.Columns[i].Cards) - 1
		switch target.(type) {

		// move TO Column
		case []CardColumn:
			draggedIndex := 0
			for _, card := range e.Columns[i].Cards {
				if card.IsDragged() {
					break
				}
				draggedIndex++
			}

			e.Columns[j].Cards = append(e.Columns[j].Cards, e.Columns[i].Cards[draggedIndex:]...)
			e.Columns[i].Cards = e.Columns[i].Cards[:draggedIndex]

			for _, card := range e.Columns[j].Cards {
				card.ColNum = j + 1
			}

			// reveal the last card from the source column, and revert its height to original
			last := len(e.Columns[i].Cards)
			if last > 0 {
				e.Columns[i].Cards[last-1].SetRevealedState(true)
				e.Columns[i].Cards[last-1].H = e.H
			}

		// move TO Foundation Pile
		case []FoundationPile:
			e.Columns[i].Cards[li].ColNum = 0
			e.FoundationPiles[j].Cards = append(e.FoundationPiles[j].Cards, e.Columns[i].Cards[li])
			e.Columns[i].Cards = e.Columns[i].Cards[:li]

			// reveal the last card from the source column, and revert its height to original
			if len(e.Columns[i].Cards) > 0 {
				e.Columns[i].Cards[li-1].SetRevealedState(true)
				e.Columns[i].Cards[li-1].H = e.H
			}
		}

	// move FROM Foundation Pile
	case []FoundationPile:
		li := len(e.FoundationPiles[i].Cards) - 1
		switch target.(type) {

		// move TO Column
		case []CardColumn:
			e.FoundationPiles[i].Cards[li].ColNum = j + 1
			e.Columns[j].Cards = append(e.Columns[j].Cards, e.FoundationPiles[i].Cards[li])
			e.FoundationPiles[i].Cards = e.FoundationPiles[i].Cards[:li]

		// move TO Foundation Pile
		case []FoundationPile:
			e.FoundationPiles[j].Cards = append(e.FoundationPiles[j].Cards, e.FoundationPiles[i].Cards[li])
			e.FoundationPiles[i].Cards = e.FoundationPiles[i].Cards[:li]
		}
	}
}
