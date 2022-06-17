package game

import (
	u "games-ebiten/resources/utils"
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

// RightClickToFoundations - moves the card directly to Foundations if any valid spot is available
func (e *Environment) RightClickToFoundations(cx, cy int) {

	// move from Waste Pile to Foundation
	if len(e.WastePile.Cards) > 0 {
		lw := len(e.WastePile.Cards) - 1
		if e.WastePile.Cards[lw].IsHovered(cx, cy) && inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
			card := e.WastePile.Cards[lw]
			for j := range e.FoundationPiles {
				if len(e.FoundationPiles[j].Cards) == 0 && card.Value == CardRanks[u.Ace] {
					e.FoundationPiles[j].Cards = append(e.FoundationPiles[j].Cards, card)
					e.WastePile.Cards = e.WastePile.Cards[:lw]

					return
				} else if len(e.FoundationPiles[j].Cards) > 0 {
					lf := len(e.FoundationPiles[j].Cards) - 1
					fpCard := e.FoundationPiles[j].Cards[lf]

					if fpCard.Color == card.Color && fpCard.Value == card.Value-1 && fpCard.Suit == card.Suit {
						e.FoundationPiles[j].Cards = append(e.FoundationPiles[j].Cards, card)
						e.WastePile.Cards = e.WastePile.Cards[:lw]

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
			if card.IsHovered(cx, cy) && inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
				for j := range e.FoundationPiles {
					if len(e.FoundationPiles[j].Cards) == 0 && card.Value == CardRanks[u.Ace] {
						card.ColNum = 0
						e.FoundationPiles[j].Cards = append(e.FoundationPiles[j].Cards, card)
						e.Columns[i].Cards = e.Columns[i].Cards[:lc]

						last := len(e.Columns[i].Cards)
						if last > 0 {
							e.Columns[i].Cards[last-1].SetRevealedState(true)
							e.Columns[i].Cards[last-1].H = e.H
						}

						return
					} else if len(e.FoundationPiles[j].Cards) > 0 {
						lf := len(e.FoundationPiles[j].Cards) - 1
						fpCard := e.FoundationPiles[j].Cards[lf]

						if fpCard.Color == card.Color && fpCard.Value == card.Value-1 && fpCard.Suit == card.Suit {
							card.ColNum = 0
							e.FoundationPiles[j].Cards = append(e.FoundationPiles[j].Cards, card)
							e.Columns[i].Cards = e.Columns[i].Cards[:lc]

							last := len(e.Columns[i].Cards)
							if last > 0 {
								e.Columns[i].Cards[last-1].SetRevealedState(true)
								e.Columns[i].Cards[last-1].H = e.H
							}

							return
						}
					}
				}
			}
		}
	}
}

// StockToWaste - contains the Draw Card functionality
func (e *Environment) StockToWaste(cx, cy int) {
	//
	// Handle Stock to Waste functionality
	//
	if len(e.StockPile.Cards) > 0 {
		if e.StockPile.Cards[0].IsHovered(cx, cy) {
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

	//
	// drag from Waste Pile to valid Column or FoundationPile slot
	//
	if len(e.WastePile.Cards) > 0 {
		lc := len(e.WastePile.Cards) - 1
		source := e.WastePile.Cards[lc].GetGeomData()

		// set the prior card's state dragged to false, so it can stick to its location
		if e.WastePile.Cards[lc].IsDragged() {
			if len(e.WastePile.Cards) > 1 {
				e.WastePile.Cards[lc-1].SetRevealedState(false)
			}

			// drag from Waste Pile to Column Slot
			for j := range e.Columns {
				if len(e.Columns[j].Cards) == 0 && e.WastePile.Cards[lc].Value == CardRanks[u.King] {
					// if there are no cards on the Column Slot and the source card is a King
					target := e.GetGeomData(e.Columns[j])

					if u.IsCollision(source, target) && inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
						e.WastePile.Cards[lc].ColNum = j + 1
						e.Columns[j].Cards = append(e.Columns[j].Cards, e.WastePile.Cards[lc])
						e.WastePile.Cards = e.WastePile.Cards[:lc]

						// exit entirely to prevent redundant iterations
						DraggedCard = nil
						return
					}
				} else if len(e.Columns[j].Cards) > 0 {
					// applies for any other card than K, also prevents iteration over the empty slots if there are any
					lj := len(e.Columns[j].Cards) - 1 // lj = last card in the current context
					target := e.Columns[j].Cards[lj].GetGeomData()

					if u.IsCollision(source, target) &&
						e.WastePile.Cards[lc].Value == e.Columns[j].Cards[lj].Value-1 &&
						e.WastePile.Cards[lc].Color != e.Columns[j].Cards[lj].Color &&
						inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {

						e.WastePile.Cards[lc].ColNum = j + 1
						e.Columns[j].Cards = append(e.Columns[j].Cards, e.WastePile.Cards[lc])
						e.WastePile.Cards = e.WastePile.Cards[:lc]

						// exit entirely to prevent redundant iterations
						DraggedCard = nil
						return
					}
				}
			}

			// draw from Waste Pile to Foundation Piles
			for j := range e.FoundationPiles {
				target := e.GetGeomData(e.FoundationPiles[j])

				if len(e.FoundationPiles[j].Cards) == 0 {
					if e.WastePile.Cards[lc].Value == CardRanks[u.Ace] &&
						u.IsCollision(source, target) && inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {

						e.WastePile.Cards[lc].ColNum = 0
						e.FoundationPiles[j].Cards = append(e.FoundationPiles[j].Cards, e.WastePile.Cards[lc])
						e.WastePile.Cards = e.WastePile.Cards[:lc]

						// exit entirely to prevent redundant iterations
						DraggedCard = nil
						return
					}

				} else {
					lj := len(e.FoundationPiles[j].Cards) - 1
					if u.IsCollision(source, target) &&
						e.WastePile.Cards[lc].Value > CardRanks[u.Ace] &&
						e.WastePile.Cards[lc].Value == e.FoundationPiles[j].Cards[lj].Value+1 &&
						e.WastePile.Cards[lc].Suit == e.FoundationPiles[j].Cards[lj].Suit &&
						inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {

						e.WastePile.Cards[lc].ColNum = 0
						e.FoundationPiles[j].Cards = append(e.FoundationPiles[j].Cards, e.WastePile.Cards[lc])
						e.WastePile.Cards = e.WastePile.Cards[:lc]

						// exit entirely to prevent redundant iterations
						DraggedCard = nil
						return
					}
				}
			}
		}
	}

	//
	// drag from Column to valid Column or Foundation Pile
	//
	for i := range e.Columns {
		if len(e.Columns[i].Cards) > 0 {
			li := len(e.Columns[i].Cards) - 1 // li = last card in the source column
			source := e.Columns[i].Cards[li].GetGeomData()

			// drag card(s) from Column to Column
			for j := range e.Columns {

				// avoid iteration over the same column
				if j != i {

					// handle moving K Column card or stack on empty Column Slot
					if len(e.Columns[j].Cards) == 0 {
						for _, c := range e.Columns[i].Cards {

							if c.Value == CardRanks[u.King] {
								target := e.GetGeomData(e.Columns[j])

								if u.IsCollision(c.GetGeomData(), target) &&
									inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {

									e.Columns[j].Cards = append(e.Columns[j].Cards, e.Columns[i].Cards[e.GetIndexOfDraggedColCard(c.ColNum-1):]...)
									e.Columns[i].Cards = e.Columns[i].Cards[:e.GetIndexOfDraggedColCard(c.ColNum-1)]

									for _, card := range e.Columns[j].Cards {
										card.ColNum = j + 1
									}

									// reveal the last card from the source column, and revert its height to original
									last := len(e.Columns[i].Cards)
									if last > 0 {
										e.Columns[i].Cards[last-1].SetRevealedState(true)
										e.Columns[i].Cards[last-1].H = e.H
									}

									// exit entirely to prevent redundant iterations
									DraggedCard = nil
									return
								}
							}
						}

					} else {
						// handle all cases except K
						if len(e.Columns[j].Cards) > 0 {
							lj := len(e.Columns[j].Cards) - 1 // lj = last card in the current context (target)
							target := e.Columns[j].Cards[lj].GetGeomData()

							for _, c := range e.Columns[i].Cards {
								if u.IsCollision(c.GetGeomData(), target) && c.IsDragged() &&
									c.Value == e.Columns[j].Cards[lj].Value-1 &&
									c.Color != e.Columns[j].Cards[lj].Color {

									if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {

										e.Columns[j].Cards = append(e.Columns[j].Cards, e.Columns[i].Cards[e.GetIndexOfDraggedColCard(c.ColNum-1):]...)
										e.Columns[i].Cards = e.Columns[i].Cards[:e.GetIndexOfDraggedColCard(c.ColNum-1)]

										for _, card := range e.Columns[j].Cards {
											card.ColNum = j + 1
										}

										// reveal the last card from the source column, and revert its height to original
										last := len(e.Columns[i].Cards)
										if last > 0 {
											e.Columns[i].Cards[last-1].SetRevealedState(true)
											e.Columns[i].Cards[last-1].H = e.H
										}

										// exit entirely to prevent redundant iterations
										DraggedCard = nil
										return
									}
								}
							}
						}
					}
				}
			}

			// loop over the Foundation Piles
			for j := range e.FoundationPiles {
				target := e.GetGeomData(e.FoundationPiles[j])

				if len(e.FoundationPiles[j].Cards) == 0 {
					if u.IsCollision(source, target) &&
						e.Columns[i].Cards[li].IsDragged() &&
						e.Columns[i].Cards[li].Value == CardRanks[u.Ace] &&
						inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {

						e.Columns[i].Cards[li].ColNum = 0
						e.FoundationPiles[j].Cards = append(e.FoundationPiles[j].Cards, e.Columns[i].Cards[li])
						e.Columns[i].Cards = e.Columns[i].Cards[:li]

						// reveal the last card from the source column, and revert its height to original
						if len(e.Columns[i].Cards) > 0 {
							e.Columns[i].Cards[li-1].SetRevealedState(true)
							e.Columns[i].Cards[li-1].H = e.H
						}

						// exit entirely to prevent redundant iterations
						DraggedCard = nil
						return
					}
				} else {
					lj := len(e.FoundationPiles[j].Cards) - 1
					if u.IsCollision(source, target) &&
						e.Columns[i].Cards[li].IsDragged() &&
						e.Columns[i].Cards[li].Value > CardRanks[u.Ace] &&
						e.Columns[i].Cards[li].Value == e.FoundationPiles[j].Cards[lj].Value+1 &&
						e.Columns[i].Cards[li].Suit == e.FoundationPiles[j].Cards[lj].Suit &&
						inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {

						e.Columns[i].Cards[li].ColNum = 0
						e.FoundationPiles[j].Cards = append(e.FoundationPiles[j].Cards, e.Columns[i].Cards[li])
						e.Columns[i].Cards = e.Columns[i].Cards[:li]

						// reveal the last card from the source column, and revert its height to original
						if len(e.Columns[i].Cards) > 0 {
							e.Columns[i].Cards[li-1].SetRevealedState(true)
							e.Columns[i].Cards[li-1].H = e.H
						}

						DraggedCard = nil
						return
					}
				}
			}
		}
	}

	//
	// drag from Foundation Pile to valid Column or another Foundation Pile
	//
	for i := range e.FoundationPiles {
		if len(e.FoundationPiles[i].Cards) > 0 {
			li := len(e.FoundationPiles[i].Cards) - 1 // li = last card in the current context
			if e.FoundationPiles[i].Cards[li].IsDragged() {
				source := e.FoundationPiles[i].Cards[li].GetGeomData()

				// loop over all the columns
				for j := range e.Columns {
					if len(e.Columns[j].Cards) > 0 {
						lj := len(e.Columns[j].Cards) - 1 // lj = last card in the current context
						target := e.Columns[j].Cards[lj].GetGeomData()

						if u.IsCollision(source, target) &&
							e.FoundationPiles[i].Cards[li].Value == e.Columns[j].Cards[lj].Value-1 &&
							e.Columns[j].Cards[lj].Color != e.FoundationPiles[i].Cards[li].Color {

							if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
								e.FoundationPiles[i].Cards[li].ColNum = j + 1
								e.Columns[j].Cards = append(e.Columns[j].Cards, e.FoundationPiles[i].Cards[li])
								e.FoundationPiles[i].Cards = e.FoundationPiles[i].Cards[:li]

								// exit entirely to prevent redundant iterations
								DraggedCard = nil
								return
							}
						}
					}
				}

				// loop over the Other Foundation Piles (this applies only to the Ace card being moved from a Store to another)
				for j := range e.FoundationPiles {
					if i != j {
						target := e.GetGeomData(e.FoundationPiles[j])

						if len(e.FoundationPiles[j].Cards) == 0 {
							if e.FoundationPiles[i].Cards[li].Value == CardRanks[u.Ace] && u.IsCollision(source, target) {
								if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
									e.FoundationPiles[j].Cards = append(e.FoundationPiles[j].Cards, e.FoundationPiles[i].Cards[li])
									e.FoundationPiles[i].Cards = e.FoundationPiles[i].Cards[:li]

									DraggedCard = nil
									return
								}
							}
						}
					}
				}
			}
		}
	}
}

func (e *Environment) GetIndexOfDraggedColCard(col int) int {
	count := 0
	for _, c := range e.Columns[col].Cards {
		if c.IsDragged() {
			break
		}
		count++
	}
	return count
}
