package free_cell

import (
	data "games-ebiten/card_games"
	res "games-ebiten/resources"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
		FreeCells       []FreeCell
		SpacerV         int
		W, H            int
	}

	FoundationPile struct {
		X, Y, W, H int
		Cards      []*Card
	}

	FreeCell struct {
		X, Y, W, H int
		Cards      []*Card
	}

	CardColumn struct {
		X, Y, W, H int
		Cards      []*Card
		IsOpen     bool
	}
)

func (e *Environment) UpdateEnv() {
	e.W = e.Deck[0].W
	e.H = e.Deck[0].H
	e.Quadrants = res.GetFlexboxQuadrants(8)

	e.FoundationPiles = []FoundationPile{
		{Cards: make([]*Card, 0, 13)},
		{Cards: make([]*Card, 0, 13)},
		{Cards: make([]*Card, 0, 13)},
		{Cards: make([]*Card, 0, 13)},
	}
	e.FreeCells = []FreeCell{
		{Cards: make([]*Card, 0, 1)},
		{Cards: make([]*Card, 0, 1)},
		{Cards: make([]*Card, 0, 1)},
		{Cards: make([]*Card, 0, 1)},
	}
	e.Columns = []CardColumn{
		{Cards: make([]*Card, 0)},
		{Cards: make([]*Card, 0)},
		{Cards: make([]*Card, 0)},
		{Cards: make([]*Card, 0)},
		{Cards: make([]*Card, 0)},
		{Cards: make([]*Card, 0)},
		{Cards: make([]*Card, 0)},
		{Cards: make([]*Card, 0)},
	}

	// start from the first quadrant
	for i := range e.FreeCells {
		fx := res.CenterItem(e.W, e.Quadrants[0+i])
		e.FreeCells[i].X = fx
		e.FreeCells[i].Y = e.SpacerV
		e.FreeCells[i].W = e.W
		e.FreeCells[i].H = e.H
	}

	// starts from the fourth quadrant
	for i := range e.FoundationPiles {
		fx := res.CenterItem(e.W, e.Quadrants[4+i])
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
				colx := res.CenterItem(e.W, e.Quadrants[0+i])
				coly := e.Quadrants[0+i].Max.Y / 3
				e.Columns[i].X = colx
				e.Columns[i].Y = coly
				e.Columns[i].W = e.W
				e.Columns[i].H = e.H
				e.Deck[cardIndex].ColNum = i + 1
				e.Columns[i].Cards = append(e.Columns[i].Cards, e.Deck[cardIndex])
				cardIndex++
			} else {
				break
			}
		}
	}
	e.SetColDraggableCards()
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
	case FreeCell:
		area := i.(FreeCell)
		rect = image.Rect(area.X, area.Y, area.X+area.W, area.Y+area.H)
	}
	return rect
}

func (e *Environment) DrawPlayground(screen *ebiten.Image, th *data.Theme) {
	// Draw the BG Image
	opBg := &ebiten.DrawImageOptions{}
	opBg.GeoM.Scale(50, 50)
	screen.DrawImage(e.BgImg, opBg)
	envTh := th.EnvScaleValue[th.Active]

	// Draw the FreeCell Slots
	for i := 0; i < 4; i++ {
		opFreeCell := &ebiten.DrawImageOptions{}
		opFreeCell.GeoM.Scale(envTh[res.X], envTh[res.Y])

		if th.Active == res.PixelatedTheme {
			opFreeCell.GeoM.Translate(float64(res.CenterItem(e.W, e.Quadrants[0+i]))+3.5, float64(e.SpacerV)+3.5)
		} else {
			opFreeCell.GeoM.Translate(float64(res.CenterItem(e.W, e.Quadrants[0+i])), float64(e.SpacerV))
		}
		screen.DrawImage(e.EmptySlotImg, opFreeCell)
	}

	// Draw the Foundation Slots
	for i := 0; i < 4; i++ {
		opFoundationSlot := &ebiten.DrawImageOptions{}
		opFoundationSlot.GeoM.Scale(envTh[res.X], envTh[res.Y])

		if th.Active == res.PixelatedTheme {
			opFoundationSlot.GeoM.Translate(float64(res.CenterItem(e.W, e.Quadrants[4+i]))+3.5, float64(e.SpacerV)+3.5)
		} else {
			opFoundationSlot.GeoM.Translate(float64(res.CenterItem(e.W, e.Quadrants[4+i])), float64(e.SpacerV))
		}
		screen.DrawImage(e.EmptySlotImg, opFoundationSlot)
	}
}

// HandleGameLogic - contains Drag and Drop functionality and cards' state updates
func (e *Environment) HandleGameLogic() {
	if res.DraggedCard != nil {
		//
		// drag FROM Column
		//
		for i := range e.Columns {
			if len(e.Columns[i].Cards) > 0 {
				li := len(e.Columns[i].Cards) - 1 // li = last card in the source column
				source := e.Columns[i].Cards[li].HitBox()

				// drop ON Free Cell
				for j := range e.FreeCells {
					target := e.HitBox(e.FreeCells[j])
					if res.IsCollision(source, target) {
						if len(e.FreeCells[j].Cards) == 0 &&
							inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
							e.MoveFromSrcToTarget(e.Columns, e.FreeCells, i, j, ebiten.MouseButtonLeft)
							res.DraggedCard = nil
							return
						}
					}
				}

				// drop ON Column
				for j := range e.Columns {
					if j != i {
						// K card
						if len(e.Columns[j].Cards) == 0 {
							for _, c := range e.Columns[i].Cards {

								if c.Dragged() && c.Value == data.CardRanks[res.King] {
									target := e.HitBox(e.Columns[j])
									if res.IsCollision(source, target) &&
										inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
										e.MoveFromSrcToTarget(e.Columns, e.Columns, i, j, ebiten.MouseButtonLeft)
										res.DraggedCard = nil
										return
									}
								}
							}

						} else {
							// Other Cards
							if len(e.Columns[j].Cards) > 0 {
								lj := len(e.Columns[j].Cards) - 1 // lj = last card in the current context (target)
								target := e.Columns[j].Cards[lj].HitBox()

								for _, c := range e.Columns[i].Cards {
									if c.Dragged() {
										if res.IsCollision(source, target) &&
											inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) &&
											c.Value == e.Columns[j].Cards[lj].Value-1 &&
											c.Color != e.Columns[j].Cards[lj].Color {
											e.MoveFromSrcToTarget(e.Columns, e.Columns, i, j, ebiten.MouseButtonLeft)
											res.DraggedCard = nil
											return
										}
									}
								}
							}
						}
					}
				}

				// drop ON Foundation Pile
				for j := range e.FoundationPiles {
					target := e.HitBox(e.FoundationPiles[j])

					if len(e.FoundationPiles[j].Cards) == 0 {
						if res.IsCollision(source, target) &&
							inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) &&
							e.Columns[i].Cards[li].Value == data.CardRanks[res.Ace] {
							e.MoveFromSrcToTarget(e.Columns, e.FoundationPiles, i, j, ebiten.MouseButtonLeft)
							res.DraggedCard = nil
							return
						}
					} else {
						lj := len(e.FoundationPiles[j].Cards) - 1
						if res.IsCollision(source, target) &&
							inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) &&
							e.Columns[i].Cards[li].Value > data.CardRanks[res.Ace] &&
							e.Columns[i].Cards[li].Value == e.FoundationPiles[j].Cards[lj].Value+1 &&
							e.Columns[i].Cards[li].Suit == e.FoundationPiles[j].Cards[lj].Suit {
							e.MoveFromSrcToTarget(e.Columns, e.FoundationPiles, i, j, ebiten.MouseButtonLeft)
							res.DraggedCard = nil
							return
						}
					}
				}
			}
		}

		//
		// drag FROM Free Cell
		//
		for i := range e.FreeCells {
			if len(e.FreeCells[i].Cards) == 1 {
				source := e.FreeCells[i].Cards[0].HitBox()
				// drop ON Column
				for j := range e.Columns {

					// K card
					if len(e.Columns[j].Cards) == 0 {
						if e.FreeCells[i].Cards[0].Value == data.CardRanks[res.King] {
							target := e.HitBox(e.Columns[j])
							if res.IsCollision(source, target) &&
								inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
								e.MoveFromSrcToTarget(e.FreeCells, e.Columns, i, j, ebiten.MouseButtonLeft)
								res.DraggedCard = nil
								return
							}
						}
					} else {
						// Other cases
						lj := len(e.Columns[j].Cards) - 1
						target := e.Columns[j].Cards[lj].HitBox()

						if res.IsCollision(source, target) {
							if e.FreeCells[i].Cards[0].Value+1 == e.Columns[j].Cards[lj].Value &&
								e.FreeCells[i].Cards[0].Color != e.Columns[j].Cards[lj].Color &&
								inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
								e.MoveFromSrcToTarget(e.FreeCells, e.Columns, i, j, ebiten.MouseButtonLeft)
								res.DraggedCard = nil
								return
							}
						}
					}
				}

				// drop ON Free Cell
				for j := range e.FreeCells {
					if i != j {
						target := e.HitBox(e.FreeCells[j])
						if res.IsCollision(source, target) {
							if len(e.FreeCells[j].Cards) == 0 &&
								inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
								e.MoveFromSrcToTarget(e.FreeCells, e.FreeCells, i, j, ebiten.MouseButtonLeft)
								res.DraggedCard = nil
								return
							}
						}
					}
				}

				// drop ON Foundation Pile
				for j := range e.FoundationPiles {
					target := e.HitBox(e.FoundationPiles[j])

					if len(e.FoundationPiles[j].Cards) == 0 {
						if res.IsCollision(source, target) &&
							inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) &&
							e.FreeCells[i].Cards[0].Value == data.CardRanks[res.Ace] {
							e.MoveFromSrcToTarget(e.FreeCells, e.FoundationPiles, i, j, ebiten.MouseButtonLeft)
							res.DraggedCard = nil
							return
						}
					} else {
						lj := len(e.FoundationPiles[j].Cards) - 1
						if res.IsCollision(source, target) &&
							inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) &&
							e.FreeCells[i].Cards[0].Value > data.CardRanks[res.Ace] &&
							e.FreeCells[i].Cards[0].Value == e.FoundationPiles[j].Cards[lj].Value+1 &&
							e.FreeCells[i].Cards[0].Suit == e.FoundationPiles[j].Cards[lj].Suit {
							e.MoveFromSrcToTarget(e.Columns, e.FoundationPiles, i, j, ebiten.MouseButtonLeft)
							res.DraggedCard = nil
							return
						}
					}
				}
			}
		}

		//
		// drag FROM Foundation Pile
		//
		for i := range e.FoundationPiles {
			if len(e.FoundationPiles[i].Cards) > 0 {
				li := len(e.FoundationPiles[i].Cards) - 1 // li = last card in the current context
				if e.FoundationPiles[i].Cards[li].Dragged() {
					source := e.FoundationPiles[i].Cards[li].HitBox()

					// drop ON Column
					for j := range e.Columns {
						if len(e.Columns[j].Cards) > 0 {
							lj := len(e.Columns[j].Cards) - 1 // lj = last card in the current context
							target := e.Columns[j].Cards[lj].HitBox()

							if res.IsCollision(source, target) && inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) &&
								e.FoundationPiles[i].Cards[li].Value == e.Columns[j].Cards[lj].Value-1 &&
								e.Columns[j].Cards[lj].Color != e.FoundationPiles[i].Cards[li].Color {
								e.MoveFromSrcToTarget(e.FoundationPiles, e.Columns, i, j, ebiten.MouseButtonLeft)
								res.DraggedCard = nil
								return
							}

						}
					}

					// drop ON Foundation Pile
					for j := range e.FoundationPiles {
						if i != j {
							target := e.HitBox(e.FoundationPiles[j])

							if len(e.FoundationPiles[j].Cards) == 0 && res.IsCollision(source, target) &&
								inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) &&
								e.FoundationPiles[i].Cards[li].Value == data.CardRanks[res.Ace] {
								e.MoveFromSrcToTarget(e.FoundationPiles, e.FoundationPiles, i, j, ebiten.MouseButtonLeft)
								res.DraggedCard = nil
								return
							}
						}
					}
				}
			}
		}
	}
}

// SetColDraggableCards - Sets which cards are in the proper order to be dragged as stack
func (e *Environment) SetColDraggableCards() {
	for _, col := range e.Columns {
		lc := len(col.Cards) - 1
		if lc >= 0 {
			// if the first undraggable card is found
			isFoundAt := -1

			// the last card will always be draggable
			col.Cards[lc].SetDraggableState(true)

			for i := lc - 1; i >= 0; i-- {
				if isFoundAt < 0 {
					if col.Cards[i].Color != col.Cards[i+1].Color && col.Cards[i].Value-1 == col.Cards[i+1].Value {
						col.Cards[i].SetDraggableState(true)
					} else {
						isFoundAt = i
						break
					}
				}
			}

			// after the index of first unmatch is found, set all the cards' DraggableState, until that index, to false
			if isFoundAt >= 0 {
				for _, card := range col.Cards[:isFoundAt] {
					card.SetDraggableState(false)
				}
			}
		}
	}
}

// MoveFromSrcToTarget - contains the logic behind moving cards from a specific Pile to another specific Pile
func (e *Environment) MoveFromSrcToTarget(src, target interface{}, i, j int, btn ebiten.MouseButton) {
	switch src.(type) {

	// move FROM Column
	case []CardColumn:
		li := len(e.Columns[i].Cards) - 1
		draggedIndex := 0
		for _, card := range e.Columns[i].Cards {
			if card.Dragged() {
				break
			}
			draggedIndex++
		}
		switch target.(type) {

		// move TO Column
		case []CardColumn:
			e.Columns[j].Cards = append(e.Columns[j].Cards, e.Columns[i].Cards[draggedIndex:]...)
			e.Columns[i].Cards = e.Columns[i].Cards[:draggedIndex]

			for _, card := range e.Columns[j].Cards {
				card.ColNum = j + 1
			}

			// revert last card's height to original
			last := len(e.Columns[i].Cards)
			if last > 0 {
				e.Columns[i].Cards[last-1].H = e.H
			}
			e.SetColDraggableCards()

		// move TO Foundation Pile
		case []FoundationPile:
			if draggedIndex == len(e.Columns[i].Cards)-1 || btn == ebiten.MouseButtonRight {
				e.Columns[i].Cards[li].ColNum = 0
				e.FoundationPiles[j].Cards = append(e.FoundationPiles[j].Cards, e.Columns[i].Cards[li])
				e.Columns[i].Cards = e.Columns[i].Cards[:li]
				e.SetColDraggableCards()
			}

		// move TO Free Cell
		case []FreeCell:
			if draggedIndex == len(e.Columns[i].Cards)-1 {
				e.Columns[i].Cards[li].ColNum = 0
				e.FreeCells[j].Cards = append(e.FreeCells[j].Cards, e.Columns[i].Cards[li])
				e.Columns[i].Cards = e.Columns[i].Cards[:li]
				e.SetColDraggableCards()
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
			e.SetColDraggableCards()

		// move TO Foundation Pile
		case []FoundationPile:
			e.FoundationPiles[j].Cards = append(e.FoundationPiles[j].Cards, e.FoundationPiles[i].Cards[li])
			e.FoundationPiles[i].Cards = e.FoundationPiles[i].Cards[:li]
		}

	// move FROM Free Cell
	case []FreeCell:
		switch target.(type) {

		// move TO Column
		case []CardColumn:
			e.FreeCells[i].Cards[0].ColNum = j + 1
			e.Columns[j].Cards = append(e.Columns[j].Cards, e.FreeCells[i].Cards[0])
			e.FreeCells[i].Cards = nil
			e.SetColDraggableCards()

		// move TO Free Cell
		case []FreeCell:
			e.FreeCells[j].Cards = append(e.FreeCells[j].Cards, e.FreeCells[i].Cards[0])
			e.FreeCells[i].Cards = nil

		// move TO Foundations Pile
		case []FoundationPile:
			e.FoundationPiles[j].Cards = append(e.FoundationPiles[j].Cards, e.FreeCells[i].Cards[0])
			e.FreeCells[i].Cards = nil
		}
	}
}
