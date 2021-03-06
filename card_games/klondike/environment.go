package klondike

import (
	"games-ebiten/card_games"
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

func (e *Environment) HitBox(i interface{}) image.Rectangle {
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

func (e *Environment) DrawPlayground(screen *ebiten.Image, th *card_games.Theme) {
	// Draw the BG Image
	opBg := &ebiten.DrawImageOptions{}
	opBg.GeoM.Scale(50, 50)
	screen.DrawImage(e.BgImg, opBg)
	envTh := th.EnvScaleValue[res.ActiveCardsTheme]

	var img *ebiten.Image

	// Draw the Stock Slot
	opStockSlot := &ebiten.DrawImageOptions{}
	opStockSlot.GeoM.Scale(envTh[res.X], envTh[res.Y])

	if res.ActiveCardsTheme == res.PixelatedTheme {
		opStockSlot.GeoM.Translate(float64(res.CenterOnX(e.W, e.Quadrants[0]))+3.5, float64(e.SpacerV+res.MainMenuH)+3.5)
	} else {
		opStockSlot.GeoM.Translate(float64(res.CenterOnX(e.W, e.Quadrants[0])), float64(e.SpacerV+res.MainMenuH))
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
		opFoundationSlot.GeoM.Scale(envTh[res.X], envTh[res.Y])

		if res.ActiveCardsTheme == res.PixelatedTheme {
			opFoundationSlot.GeoM.Translate(float64(res.CenterOnX(e.W, e.Quadrants[3+i]))+3.5, float64(e.SpacerV+res.MainMenuH)+3.5)
		} else {
			opFoundationSlot.GeoM.Translate(float64(res.CenterOnX(e.W, e.Quadrants[3+i])), float64(e.SpacerV+res.MainMenuH))
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
	e.Quadrants = res.FlexRowQuadrants(7)

	e.StockPile.Cards = make([]*Card, 0, 24)
	e.WastePile.Cards = make([]*Card, 0, 24)
	e.FoundationPiles = []FoundationPile{
		{Cards: make([]*Card, 0, 13)},
		{Cards: make([]*Card, 0, 13)},
		{Cards: make([]*Card, 0, 13)},
		{Cards: make([]*Card, 0, 13)},
	}
	e.Columns = []CardColumn{
		{Cards: make([]*Card, 0)},
		{Cards: make([]*Card, 0)},
		{Cards: make([]*Card, 0)},
		{Cards: make([]*Card, 0)},
		{Cards: make([]*Card, 0)},
		{Cards: make([]*Card, 0)},
		{Cards: make([]*Card, 0)},
	}
	e.WastePile.X = res.CenterOnX(e.W, e.Quadrants[1])

	e.StockPile.X = res.CenterOnX(e.W, e.Quadrants[0])
	e.StockPile.Y = e.SpacerV + res.MainMenuH
	e.StockPile.W = e.W
	e.StockPile.H = e.H

	// starts from the fourth quadrant
	for i := range e.FoundationPiles {
		fx := res.CenterOnX(e.W, e.Quadrants[3+i])
		e.FoundationPiles[i].X = fx
		e.FoundationPiles[i].Y = e.SpacerV + res.MainMenuH
		e.FoundationPiles[i].W = e.W
		e.FoundationPiles[i].H = e.H
	}

	// fill every column array with its relative count of cards and save GeomData of columns placeholders
	cardIndex := 0
	for i := range e.Columns {
		// initiate the location of the Card Column placeholders
		colx := res.CenterOnX(e.W, e.Quadrants[0+i])
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

// StockToWaste - contains the Draw Card functionality
func (e *Environment) StockToWaste(cx, cy int) {
	//
	// Handle Stock to Waste functionality
	//
	if len(e.StockPile.Cards) > 0 {
		if e.StockPile.Cards[0].Hovered(cx, cy) {
			last := len(e.StockPile.Cards) - 1
			if len(e.StockPile.Cards) == 1 {
				last = 0
			}
			// append every last card from StockPile to WastePile, then trim last card from StockPile
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				e.WastePile.Cards = append(e.WastePile.Cards, e.StockPile.Cards[last])
				e.StockPile.Cards = e.StockPile.Cards[:last]
			}
		}
	} else {
		// if there are no more cards, clicking the circle will reset the process
		if !e.IsVegas && e.IsStockPileHovered(cx, cy) && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			for i := range e.WastePile.Cards {
				e.StockPile.Cards = append(e.StockPile.Cards, e.WastePile.Cards[i])
				e.StockPile.Cards[i].SetRevealedState(false)
			}
			e.WastePile.Cards = e.WastePile.Cards[:0]

			// reverse order of newly stacked StockPile cards:
			for i, j := 0, len(e.StockPile.Cards)-1; i < j; i, j = i+1, j-1 {
				e.StockPile.Cards[i], e.StockPile.Cards[j] = e.StockPile.Cards[j], e.StockPile.Cards[i]
			}
		}
	}
}

// HandleGameLogic - contains Drag from Waste Pile, valid Column, Foundation Pile to any valid slot; and right click to foundations functionalities
func (e *Environment) HandleGameLogic(cx, cy int) {
	e.RightClickToFoundations(cx, cy)
	e.StockToWaste(cx, cy)

	if res.DraggedCard != nil {
		//
		// drag FROM Waste Pile
		//
		if len(e.WastePile.Cards) > 0 {
			lc := len(e.WastePile.Cards) - 1
			source := e.WastePile.Cards[lc].HitBox()

			// set the prior card's state dragged to false, so it can stick to its location
			if e.WastePile.Cards[lc].Dragged() {
				if len(e.WastePile.Cards) > 1 {
					e.WastePile.Cards[lc-1].SetRevealedState(false)
				}

				// drop ON Column
				for j := range e.Columns {
					if len(e.Columns[j].Cards) == 0 && e.WastePile.Cards[lc].Value == res.CardRanks[res.King] {
						// K card
						target := e.HitBox(e.Columns[j])

						if res.IsCollision(source, target) && inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
							e.MoveFromSrcToTarget(e.WastePile, e.Columns, lc, j, ebiten.MouseButtonLeft)
							res.DraggedCard = nil
							return

						}
					} else if len(e.Columns[j].Cards) > 0 {
						// Other Cards
						lj := len(e.Columns[j].Cards) - 1 // lj = last card in the current context
						target := e.Columns[j].Cards[lj].HitBox()

						if res.IsCollision(source, target) &&
							e.WastePile.Cards[lc].Value == e.Columns[j].Cards[lj].Value-1 &&
							e.WastePile.Cards[lc].Color != e.Columns[j].Cards[lj].Color &&
							inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
							e.MoveFromSrcToTarget(e.WastePile, e.Columns, lc, j, ebiten.MouseButtonLeft)
							res.DraggedCard = nil
							return
						}
					}
				}

				// drop ON Foundation Pile
				for j := range e.FoundationPiles {
					target := e.HitBox(e.FoundationPiles[j])

					if len(e.FoundationPiles[j].Cards) == 0 {
						if e.WastePile.Cards[lc].Value == res.CardRanks[res.Ace] &&
							res.IsCollision(source, target) && inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
							e.MoveFromSrcToTarget(e.WastePile, e.FoundationPiles, lc, j, ebiten.MouseButtonLeft)
							res.DraggedCard = nil
							return
						}

					} else {
						lj := len(e.FoundationPiles[j].Cards) - 1
						if res.IsCollision(source, target) &&
							e.WastePile.Cards[lc].Value > res.CardRanks[res.Ace] &&
							e.WastePile.Cards[lc].Value == e.FoundationPiles[j].Cards[lj].Value+1 &&
							e.WastePile.Cards[lc].Suit == e.FoundationPiles[j].Cards[lj].Suit &&
							inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
							e.MoveFromSrcToTarget(e.WastePile, e.FoundationPiles, lc, j, ebiten.MouseButtonLeft)
							res.DraggedCard = nil
							return
						}
					}
				}
			}
		}

		//
		// drag FROM Column
		//
		for i := range e.Columns {
			if len(e.Columns[i].Cards) > 0 {
				li := len(e.Columns[i].Cards) - 1 // li = last card in the source column
				source := e.Columns[i].Cards[li].HitBox()

				// drop ON Column
				for j := range e.Columns {
					if j != i {
						// K card
						if len(e.Columns[j].Cards) == 0 {
							for _, c := range e.Columns[i].Cards {

								if c.Dragged() && c.Value == res.CardRanks[res.King] {
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
							e.Columns[i].Cards[li].Value == res.CardRanks[res.Ace] {
							e.MoveFromSrcToTarget(e.Columns, e.FoundationPiles, i, j, ebiten.MouseButtonLeft)
							res.DraggedCard = nil
							return
						}
					} else {
						lj := len(e.FoundationPiles[j].Cards) - 1
						if res.IsCollision(source, target) &&
							inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) &&
							e.Columns[i].Cards[li].Value > res.CardRanks[res.Ace] &&
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
								e.FoundationPiles[i].Cards[li].Value == res.CardRanks[res.Ace] {
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

			// reveal the last card from the source column, and revert its height to original
			last := len(e.Columns[i].Cards)
			if last > 0 {
				e.Columns[i].Cards[last-1].SetRevealedState(true)
				e.Columns[i].Cards[last-1].H = e.H
			}

		// move TO Foundation Pile
		case []FoundationPile:
			if draggedIndex == len(e.Columns[i].Cards)-1 || btn == ebiten.MouseButtonRight {
				e.Columns[i].Cards[li].ColNum = 0
				e.FoundationPiles[j].Cards = append(e.FoundationPiles[j].Cards, e.Columns[i].Cards[li])
				e.Columns[i].Cards = e.Columns[i].Cards[:li]

				// reveal the last card from the source column, and revert its height to original
				if len(e.Columns[i].Cards) > 0 {
					e.Columns[i].Cards[li-1].SetRevealedState(true)
					e.Columns[i].Cards[li-1].H = e.H
				}
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

	// move FROM Waste Pile
	case WastePile:
		switch target.(type) {

		// move TO Column
		case []CardColumn:
			e.WastePile.Cards[i].ColNum = j + 1
			e.Columns[j].Cards = append(e.Columns[j].Cards, e.WastePile.Cards[i])
			e.WastePile.Cards = e.WastePile.Cards[:i]

		// move TO Foundation Piles
		case []FoundationPile:
			e.FoundationPiles[j].Cards = append(e.FoundationPiles[j].Cards, e.WastePile.Cards[i])
			e.WastePile.Cards = e.WastePile.Cards[:i]
		}
	}
}

// RightClickToFoundations - moves the card directly to Foundations if any valid spot is available
func (e *Environment) RightClickToFoundations(cx, cy int) {

	// move from Waste Pile to Foundation
	if len(e.WastePile.Cards) > 0 {
		lw := len(e.WastePile.Cards) - 1
		if e.WastePile.Cards[lw].Hovered(cx, cy) && inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
			card := e.WastePile.Cards[lw]
			for j := range e.FoundationPiles {
				if len(e.FoundationPiles[j].Cards) == 0 && card.Value == res.CardRanks[res.Ace] {
					e.MoveFromSrcToTarget(e.WastePile, e.FoundationPiles, lw, j, ebiten.MouseButtonRight)
					return
				} else if len(e.FoundationPiles[j].Cards) > 0 {
					lf := len(e.FoundationPiles[j].Cards) - 1
					fpCard := e.FoundationPiles[j].Cards[lf]

					if fpCard.Value == card.Value-1 && fpCard.Suit == card.Suit {
						e.MoveFromSrcToTarget(e.WastePile, e.FoundationPiles, lw, j, ebiten.MouseButtonRight)
						return
					}
				}
			}
		}
	}

	// move from Column to Foundation
	for i := range e.Columns {
		if len(e.Columns[i].Cards) > 0 {
			lc := len(e.Columns[i].Cards) - 1
			card := e.Columns[i].Cards[lc]
			if card.Hovered(cx, cy) && inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
				for j := range e.FoundationPiles {
					if len(e.FoundationPiles[j].Cards) == 0 && card.Value == res.CardRanks[res.Ace] {
						e.MoveFromSrcToTarget(e.Columns, e.FoundationPiles, i, j, ebiten.MouseButtonRight)
						return
					} else if len(e.FoundationPiles[j].Cards) > 0 {
						lf := len(e.FoundationPiles[j].Cards) - 1
						fpCard := e.FoundationPiles[j].Cards[lf]

						if fpCard.Value == card.Value-1 && fpCard.Suit == card.Suit {
							e.MoveFromSrcToTarget(e.Columns, e.FoundationPiles, i, j, ebiten.MouseButtonRight)
							return
						}
					}
				}
			}
		}
	}
}

func AnimateCardMovement(ticks int, card Card, dest image.Rectangle) {

	ticks++
}
