package game

import (
	u "games-ebiten/resources/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
)

type (
	Game struct {
		*Theme
		*Environment
	}
)

func NewGame() *Game {
	classicImg := ebiten.NewImageFromImage(u.LoadSpriteImage("resources/images/cards/classic-solitaire.png"))
	th := NewTheme()
	deck := GenerateDeck(th)

	g := &Game{
		Theme: th,
		Environment: &Environment{
			CardsVSpacer: 30,
			SpacerV:      50,
			SpacerH:      float64(th.FrontFaceFrameData[th.Active][u.FrW]) + 10,
			BgImg:        classicImg.SubImage(image.Rect(700, 500, 750, 550)).(*ebiten.Image),
			EmptySlotImg: classicImg.SubImage(image.Rect(852, 384, 852+71, 384+96)).(*ebiten.Image),

			DrawCardSlot: DrawCardSlot{
				GreenSlotImg: classicImg.SubImage(image.Rect(710, 384, 710+71, 384+96)).(*ebiten.Image),
				RedSlotImg:   classicImg.SubImage(image.Rect(781, 384, 781+71, 384+96)).(*ebiten.Image),
				Cards:        make([]*Card, 0, 24),
			},
			DrawnCardsSlot: DrawnCardsSlot{
				Cards: make([]*Card, 0, 24),
			},
			CardStores: []CardStore{
				{Cards: make([]*Card, 0, 13)},
				{Cards: make([]*Card, 0, 13)},
				{Cards: make([]*Card, 0, 13)},
				{Cards: make([]*Card, 0, 13)},
			},
			Columns: []CardColumn{
				{Cards: make([]*Card, 0, 1)},
				{Cards: make([]*Card, 0, 2)},
				{Cards: make([]*Card, 0, 3)},
				{Cards: make([]*Card, 0, 4)},
				{Cards: make([]*Card, 0, 5)},
				{Cards: make([]*Card, 0, 6)},
				{Cards: make([]*Card, 0, 7)},
			},
		},
	}
	_, _, frW, frH := th.GetFrontFrameGeomData(g.Active)
	x, y := float64(frW)+g.SpacerH, g.SpacerV
	w, h := float64(frW)*th.ScaleValue[g.Active][u.X], float64(frH)*th.ScaleValue[g.Active][u.Y]

	g.DrawCardSlot.X = x
	g.DrawCardSlot.Y = y
	g.DrawCardSlot.W = w
	g.DrawCardSlot.H = h

	for s := range g.CardStores {
		sx := (float64(frW) + g.SpacerH) * (float64(s) + 4)
		g.CardStores[s].X = sx
		g.CardStores[s].Y = y
		g.CardStores[s].W = w
		g.CardStores[s].H = h
	}

	// fill every column array with its relative count of cards and save GeomData of columns placeholders
	cardIndex := 0
	for i := range g.Columns {
		// initiate the location of the Card Column placeholders
		colx := (float64(frW) + g.SpacerH) * float64(i+1)
		coly := float64(u.ScreenHeight / 3)
		g.Columns[i].X = colx
		g.Columns[i].Y = coly
		g.Columns[i].W = w
		g.Columns[i].H = h

		for j := 0; j <= i; j++ {
			deck[cardIndex].ColCount = i + 1
			// keep only the last one revealed
			if j == i {
				deck[cardIndex].IsRevealed = true
			}
			g.Columns[i].Cards = append(g.Columns[i].Cards, deck[cardIndex])
			cardIndex++
		}
	}

	// fill the DrawCard array
	for i := range deck[cardIndex:] {
		g.DrawCardSlot.Cards = append(g.DrawCardSlot.Cards, deck[cardIndex:][i])
	}

	return g
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.DrawEnvironment(screen, g.Theme)

	// Draw the Card Columns
	for i := range g.Columns {
		x := (float64(g.FrontFaceFrameData[g.Active][u.FrW]) + g.SpacerH) * (float64(i) + 1)
		for j := range g.Columns[i].Cards {
			y := float64(u.ScreenHeight/3) + float64(j)*g.CardsVSpacer

			// draw the overlapped with the height of the space in which the card is visible
			if j != len(g.Columns[i].Cards)-1 {
				g.Columns[i].Cards[j].H = g.CardsVSpacer
			}
			g.Columns[i].Cards[j].X = x
			g.Columns[i].Cards[j].Y = y
			g.Columns[i].Cards[j].DrawCardSprite(screen)
		}
	}

	// Draw the Draw Card Slot
	if len(g.DrawCardSlot.Cards) > 0 {
		for i := range g.DrawCardSlot.Cards {
			g.DrawCardSlot.Cards[i].X = g.DrawCardSlot.X
			g.DrawCardSlot.Cards[i].Y = g.DrawCardSlot.Y
			g.DrawCardSlot.Cards[i].DrawCardSprite(screen)
		}
	}

	//secondLast := len(g.DrawnCardsSlot.Cards) - 2
	//if HoveredCard != nil && len(g.DrawnCardsSlot.Cards) > 1 {
	//	g.DrawnCardsSlot.Cards[secondLast].X = (float64(g.FrontFaceFrameData[g.Active][u.FrW]) + g.SpacerH) * 2
	//	g.DrawnCardsSlot.Cards[secondLast].Y = g.SpacerV
	//	g.DrawCardSlot.Cards[secondLast].IsDragged = false
	//	g.DrawnCardsSlot.Cards[secondLast].IsRevealed = true
	//	g.DrawnCardsSlot.Cards[secondLast].DrawCardSprite(screen)
	//}

	// Draw the Drawn Card Slot
	for i := range g.DrawnCardsSlot.Cards {
		x := (float64(g.FrontFaceFrameData[g.Active][u.FrW]) + g.SpacerH) * 2
		y := g.SpacerV
		g.DrawnCardsSlot.Cards[i].X = x
		g.DrawnCardsSlot.Cards[i].Y = y
		g.DrawnCardsSlot.Cards[i].IsRevealed = true
		g.DrawnCardsSlot.Cards[i].DrawCardSprite(screen)

	}

	for i := range g.CardStores {
		x := (float64(g.FrontFaceFrameData[g.Active][u.FrW]) + g.SpacerH) * (float64(i) + 4)
		for j := range g.CardStores[i].Cards {
			y := g.SpacerV
			g.CardStores[i].Cards[j].X = x
			g.CardStores[i].Cards[j].Y = y
			g.CardStores[i].Cards[j].DrawCardSprite(screen)
		}
	}

	// force card image(s) persistence while dragged. append the dragged card + any following cards, if card is part of a column
	if HoveredCard != nil {
		switch HoveredCard.(type) {
		case *Card:
			c := HoveredCard.(*Card)
			if c.ColCount > 0 &&
				g.Columns[c.ColCount-1].Cards[len(g.Columns[c.ColCount-1].Cards)-1] != c {

				/// Drag multiple cards
				cx, cy := ebiten.CursorPosition()

				for i, card := range g.Columns[c.ColCount-1].Cards {
					if i == 0 && int(card.X) <= cx && cx < int(card.X+card.W) && int(card.Y) <= cy && cy < int(card.Y+card.H) && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
						card.IsDragged = true
					}

					// drag and set location
					if inpututil.MouseButtonPressDuration(ebiten.MouseButtonLeft) > 0 {
						if i == 0 {
							HoveredCard = card
						}
						_, _, w, _ := card.GetPosition()
						card.X = float64(cx) - w/2

						if i == len(g.Columns[c.ColCount-1].Cards)-1 {
							card.Y = float64(cy - (30 * i))
						} else {
							//TODO: add this one to Theme (it only appears in Environment as CardsVSpacer
							card.Y = float64(cy - (30 * i))
							card.H = 30
						}

						opc := &ebiten.DrawImageOptions{}
						opc.GeoM.Scale(c.ScaleX, c.ScaleY)
						opc.GeoM.Translate(c.X, c.Y)
						screen.DrawImage(c.Img, opc)
					}
				}
			} else {
				opc := &ebiten.DrawImageOptions{}
				opc.GeoM.Scale(c.ScaleX, c.ScaleY)
				opc.GeoM.Translate(c.X, c.Y)
				screen.DrawImage(c.Img, opc)
			}
		}
	}

	ebitenutil.DebugPrint(screen, "Press 1 or 2 to change Themes")
}

func (g *Game) Update() error {
	cx, cy := ebiten.CursorPosition()
	if ebiten.IsKeyPressed(ebiten.Key1) && g.Active != u.ClassicTheme {
		g.Active = u.ClassicTheme
	} else if ebiten.IsKeyPressed(ebiten.Key2) && g.Active != u.PixelatedTheme {
		g.Active = u.PixelatedTheme
	}

	//
	// handle the Draw Card functionality
	//
	if len(g.DrawCardSlot.Cards) > 0 {
		if g.DrawCardSlot.Cards[0].IsCardHovered(cx, cy) {
			last := len(g.DrawCardSlot.Cards) - 1
			if len(g.DrawCardSlot.Cards) == 1 {
				last = 0
			}
			// append every last card from DrawCardSlot to DrawnCardsSlot, then trim last card from DrawCardSlot
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				g.DrawnCardsSlot.Cards = append(g.DrawnCardsSlot.Cards, g.DrawCardSlot.Cards[last])
				g.DrawCardSlot.Cards = g.DrawCardSlot.Cards[:last]
			}
		}
	} else {
		// if there are no more cards, clicking the circle will reset the process
		if g.IsDrawCardHovered(cx, cy, g.Theme) && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			for i := range g.DrawnCardsSlot.Cards {
				g.DrawCardSlot.Cards = append(g.DrawCardSlot.Cards, g.DrawnCardsSlot.Cards[i])
				g.DrawCardSlot.Cards[i].IsRevealed = false
			}
			g.DrawnCardsSlot.Cards = g.DrawnCardsSlot.Cards[:0]

			// reverse order of newly stacked DrawCardSlot cards:
			for i, j := 0, len(g.DrawCardSlot.Cards)-1; i < j; i, j = i+1, j-1 {
				g.DrawCardSlot.Cards[i], g.DrawCardSlot.Cards[j] = g.DrawCardSlot.Cards[j], g.DrawCardSlot.Cards[i]
			}
		}
	}

	if HoveredCard != nil {
		//
		// drag from Drawn Stack to valid Column or Store slot
		//
		if len(g.DrawnCardsSlot.Cards) > 0 {
			lc := len(g.DrawnCardsSlot.Cards) - 1
			ix, iy := g.DrawnCardsSlot.Cards[lc].X, g.DrawnCardsSlot.Cards[lc].Y
			iw, ih := g.DrawnCardsSlot.Cards[lc].W, g.DrawnCardsSlot.Cards[lc].H

			if g.DrawnCardsSlot.Cards[lc].IsDragged {
				for j := range g.Columns {

					// if there are no cards on the Column Slot and the source card is a King
					if len(g.Columns[j].Cards) == 0 && g.DrawnCardsSlot.Cards[lc].Value == CardRanks[u.King] {
						jx, jy, jw, jh := g.Columns[j].GetColumnGeoMData()
						if u.IsCollision(ix, iy, iw, ih, jx, jy, jw, jh) &&
							inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {

							g.DrawnCardsSlot.Cards[lc].ColCount = j + 1
							g.Columns[j].Cards = append(g.Columns[j].Cards, g.DrawnCardsSlot.Cards[lc])
							g.DrawnCardsSlot.Cards = g.DrawnCardsSlot.Cards[:lc]

							HoveredCard = nil
							return nil
						}
					} else {
						lj := len(g.Columns[j].Cards) - 1 // lj = last card in the current context
						jx, jy := g.Columns[j].Cards[lj].X, g.Columns[j].Cards[lj].Y
						jw, jh := g.Columns[j].Cards[lj].W, g.Columns[j].Cards[lj].H

						if u.IsCollision(ix, iy, iw, ih, jx, jy, jw, jh) &&
							g.DrawnCardsSlot.Cards[lc].Value == g.Columns[j].Cards[lj].Value-1 &&
							g.DrawnCardsSlot.Cards[lc].Color != g.Columns[j].Cards[lj].Color &&
							inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {

							g.DrawnCardsSlot.Cards[lc].ColCount = j + 1
							g.Columns[j].Cards = append(g.Columns[j].Cards, g.DrawnCardsSlot.Cards[lc])
							g.DrawnCardsSlot.Cards = g.DrawnCardsSlot.Cards[:lc]

							// exit entirely to prevent redundant iterations
							HoveredCard = nil
							return nil
						}
					}
				}
				for j := range g.CardStores {
					jx, jy, jw, jh := g.CardStores[j].GetStoreGeomData()

					if len(g.CardStores[j].Cards) == 0 {
						if g.DrawnCardsSlot.Cards[lc].Value == CardRanks[u.Ace] &&
							u.IsCollision(ix, iy, iw, ih, jx, jy, jw, jh) &&
							inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {

							g.CardStores[j].Cards = append(g.CardStores[j].Cards, g.DrawnCardsSlot.Cards[lc])
							g.DrawnCardsSlot.Cards = g.DrawnCardsSlot.Cards[:lc]

							// exit entirely to prevent redundant iterations
							HoveredCard = nil
							return nil
						}

					} else {
						lj := len(g.CardStores[j].Cards) - 1
						if u.IsCollision(ix, iy, iw, ih, jx, jy, jw, jh) &&
							g.DrawnCardsSlot.Cards[lc].Value > CardRanks[u.Ace] &&
							g.DrawnCardsSlot.Cards[lc].Value == g.CardStores[j].Cards[lj].Value+1 &&
							g.DrawnCardsSlot.Cards[lc].Suit == g.CardStores[j].Cards[lj].Suit &&
							inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {

							g.CardStores[j].Cards = append(g.CardStores[j].Cards, g.DrawnCardsSlot.Cards[lc])
							g.DrawnCardsSlot.Cards = g.DrawnCardsSlot.Cards[:lc]

							// exit entirely to prevent redundant iterations
							HoveredCard = nil
							return nil
						}
					}
				}
			}
		}

		//
		// drag from Column to valid Column or Store slot
		//
		for i := range g.Columns {
			if len(g.Columns[i].Cards) > 0 {
				hidden := g.Columns[i].GetCountOfHidden()
				li := len(g.Columns[i].Cards) - 1 // li = last card in the current context

				ix, iy := g.Columns[i].Cards[li].X, g.Columns[i].Cards[li].Y
				iw, ih := g.Columns[i].Cards[li].W, g.Columns[i].Cards[li].H

				// iterate through all the visible cards, for multiple grab
				for ci, card := range g.Columns[i].Cards[hidden:] {
					//li := len(g.Columns[i].Cards) - 1 // li = last card in the current context
					if card.IsDragged {

						ix, iy = card.X, card.Y
						iw, ih = card.W, card.H

						// loop over all the other columns
						for j := range g.Columns {
							if j != i {

								// handle K case
								if len(g.Columns[j].Cards) == 0 && card.Value == CardRanks[u.King] {
									jx, jy, jw, jh := g.Columns[j].GetColumnGeoMData()
									if u.IsCollision(ix, iy, iw, ih, jx, jy, jw, jh) &&
										inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {

										// append dragged card + any other cards below it. (all cards below hidden ones + card index + 1)
										for _, c := range g.Columns[i].Cards[hidden : hidden+ci+1] {
											c.ColCount = j + 1
										}
										g.Columns[j].Cards = append(g.Columns[j].Cards, g.Columns[i].Cards[hidden:hidden+ci+1]...)

										// source column = all existing ones + the number of dragged ones
										g.Columns[i].Cards = g.Columns[i].Cards[:hidden+ci]

										// reveal the last card from the source column, and revert its height to original
										last := len(g.Columns[i].Cards)
										if last > 0 {
											g.Columns[i].Cards[last-1].IsRevealed = true
											g.Columns[i].Cards[last-1].H = g.DrawCardSlot.H
										}

										// exit entirely to prevent redundant iterations
										HoveredCard = nil
										return nil
									}
								}

								// handle all cases except K
								if len(g.Columns[j].Cards) > 0 {

									lj := len(g.Columns[j].Cards) - 1 // lj = last card in the current context
									jx, jy := g.Columns[j].Cards[lj].X, g.Columns[j].Cards[lj].Y
									jw, jh := g.Columns[j].Cards[lj].W, g.Columns[j].Cards[lj].H

									if u.IsCollision(ix, iy, iw, ih, jx, jy, jw, jh) &&
										card.Value == g.Columns[j].Cards[lj].Value-1 &&
										card.Color != g.Columns[j].Cards[lj].Color {

										if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {

											// append dragged card + any other cards below it. (all cards below hidden ones + card index + 1)
											for _, c := range g.Columns[i].Cards[hidden : hidden+ci+1] {
												c.ColCount = j + 1
											}
											g.Columns[j].Cards = append(g.Columns[j].Cards, g.Columns[i].Cards[hidden:hidden+ci+1]...)

											// source column = all existing ones + the number of dragged ones
											g.Columns[i].Cards = g.Columns[i].Cards[:hidden+ci]

											// reveal the last card from the source column, and revert its height to original
											last := len(g.Columns[i].Cards)
											if last > 0 {
												g.Columns[i].Cards[last-1].IsRevealed = true
												g.Columns[i].Cards[last-1].H = g.DrawCardSlot.H
											}

											// exit entirely to prevent redundant iterations
											HoveredCard = nil
											return nil
										}
									}
								}
							}
						}
					}
				}

				// loop over the Store slots
				for j := range g.CardStores {
					jx, jy, jw, jh := g.CardStores[j].GetStoreGeomData()

					if len(g.CardStores[j].Cards) == 0 {
						if u.IsCollision(ix, iy, iw, ih, jx, jy, jw, jh) &&
							g.Columns[i].Cards[li].IsDragged &&
							g.Columns[i].Cards[li].Value == CardRanks[u.Ace] &&
							inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {

							g.Columns[i].Cards[li].ColCount = 0
							g.CardStores[j].Cards = append(g.CardStores[j].Cards, g.Columns[i].Cards[li])
							g.Columns[i].Cards = g.Columns[i].Cards[:li]

							// reveal the last card from the source column, and revert its height to original
							if len(g.Columns[i].Cards) > 0 {
								g.Columns[i].Cards[li-1].IsRevealed = true
								g.Columns[i].Cards[li-1].H = g.DrawCardSlot.H
							}

							// exit entirely to prevent redundant iterations
							HoveredCard = nil
							return nil
						}
					} else {
						lj := len(g.CardStores[j].Cards) - 1
						if u.IsCollision(ix, iy, iw, ih, jx, jy, jw, jh) &&
							g.Columns[i].Cards[li].IsDragged &&
							g.Columns[i].Cards[li].Value > CardRanks[u.Ace] &&
							g.Columns[i].Cards[li].Value == g.CardStores[j].Cards[lj].Value+1 &&
							g.Columns[i].Cards[li].Suit == g.CardStores[j].Cards[lj].Suit &&
							inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {

							g.Columns[i].Cards[li].ColCount = 0
							g.CardStores[j].Cards = append(g.CardStores[j].Cards, g.Columns[i].Cards[li])
							g.Columns[i].Cards = g.Columns[i].Cards[:li]

							HoveredCard = nil
							return nil
						}
					}
				}
			}
		}

		//
		// drag from Store to valid Column or another available Store
		//
		for i := range g.CardStores {
			if len(g.CardStores[i].Cards) > 0 {
				li := len(g.CardStores[i].Cards) - 1 // li = last card in the current context
				if g.CardStores[i].Cards[li].IsDragged {
					ix, iy := g.CardStores[i].Cards[li].X, g.CardStores[i].Cards[li].Y
					iw, ih := g.CardStores[i].Cards[li].W, g.CardStores[i].Cards[li].H

					// loop over all the columns
					for j := range g.Columns {
						if len(g.Columns[j].Cards) > 0 {
							lj := len(g.Columns[j].Cards) - 1 // lj = last card in the current context
							jx, jy := g.Columns[j].Cards[lj].X, g.Columns[j].Cards[lj].Y
							jw, jh := g.Columns[j].Cards[lj].W, g.Columns[j].Cards[lj].H

							if u.IsCollision(ix, iy, iw, ih, jx, jy, jw, jh) &&
								g.CardStores[i].Cards[li].Value == g.Columns[j].Cards[lj].Value-1 &&
								g.Columns[j].Cards[lj].Color != g.CardStores[i].Cards[li].Color {

								if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
									g.CardStores[i].Cards[li].ColCount = j + 1
									g.Columns[j].Cards = append(g.Columns[j].Cards, g.CardStores[i].Cards[li])
									g.CardStores[i].Cards = g.CardStores[i].Cards[:li]

									// exit entirely to prevent redundant iterations
									HoveredCard = nil
									return nil
								}
							}
						}
					}

					// loop over the Other Store slots (this applies only to the Ace card being moved from a Store to another)
					for j := range g.CardStores {
						if i != j {
							jx, jy, jw, jh := g.CardStores[j].GetStoreGeomData()

							if len(g.CardStores[j].Cards) == 0 {
								if g.CardStores[i].Cards[li].Value == CardRanks[u.Ace] && u.IsCollision(ix, iy, iw, ih, jx, jy, jw, jh) {
									if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
										g.CardStores[j].Cards = append(g.CardStores[j].Cards, g.CardStores[i].Cards[li])
										g.CardStores[i].Cards = g.CardStores[i].Cards[:li]

										//last := len(g.CardStores[i].Cards)
										//if last > 0 {
										//	g.CardStores[i].Cards[last-1].IsRevealed = true
										//}

										HoveredCard = nil
										return nil
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return u.ScreenWidth, u.ScreenHeight
}

// GenerateDeck - returns a []*Card{} in which all elements have the corresponding details and images
func GenerateDeck(th *Theme) []*Card {
	var colStart, colEnd int
	deck := make([]*Card, 0, 52)
	active := th.Active

	// set which BackFace the cards have (FrOX, FRoY, FrW, FrH)
	bf := th.GetBackFrameGeomData(active, u.StaticBack1)

	// set which FrontFace the cards have
	frOX, frOY, frW, frH := th.GetFrontFrameGeomData(active)

	// this logic is needed due to the discrepancy of the sprite sheets:
	// one Image starts with card Ace as the first Column value, while others start with card number or other value
	switch active {
	case u.PixelatedTheme:
		colStart = 1
		colEnd = 14
	case u.ClassicTheme:
		colStart = 0
		colEnd = 13
	}

	// there are 4 suits on the image, and 1 suit consists of 13 cards
	for si, suit := range th.SuitsOrder[active] {
		color := ""
		switch suit {
		case u.Hearts, u.Diamonds:
			color = u.Red
		case u.Spades, u.Clubs:
			color = u.Black
		}

		for i := colStart; i < colEnd; i++ {
			x, y := frOX+i*frW, frOY+si*frH
			w, h := float64(frW)*th.ScaleValue[active][u.X], float64(frH)*th.ScaleValue[active][u.Y]

			// crete card dynamicalY, based on the Active Theme.
			card := &Card{
				Img:     th.Sources[active].SubImage(image.Rect(x, y, x+frW, y+frH)).(*ebiten.Image),
				BackImg: th.Sources[active].SubImage(image.Rect(bf[0], bf[1], bf[2]+bf[0], bf[3]+bf[1])).(*ebiten.Image),
				Suit:    suit,
				Value:   CardRanks[Translation[active][i]],
				Color:   color,
				ScaleX:  th.ScaleValue[active][u.X],
				ScaleY:  th.ScaleValue[active][u.Y],
				W:       w,
				H:       h,
			}

			// append every customized card to the deck
			deck = append(deck, card)
		}
	}
	return deck
}
