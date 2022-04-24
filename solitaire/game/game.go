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

	// force card image persistence while dragged
	if HoveredCard != nil {
		switch HoveredCard.(type) {
		case *Card:
			c := HoveredCard.(*Card)
			opc := &ebiten.DrawImageOptions{}
			opc.GeoM.Scale(c.ScaleX, c.ScaleY)
			opc.GeoM.Translate(c.X, c.Y)
			screen.DrawImage(c.Img, opc)
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
				count := 0
				//revealed := g.Columns[i].GetCountOfRevealed()

				// iterate through all the visible cards, for multiple grab
				for x, _ := range g.Columns[i].Cards[hidden:] {
					count++
					//li := len(g.Columns[i].Cards) - 1 // li = last card in the current context
					if g.Columns[i].Cards[x].IsDragged {

						ix, iy := g.Columns[i].Cards[x].X, g.Columns[i].Cards[x].Y
						iw, ih := g.Columns[i].Cards[x].W, g.Columns[i].Cards[x].H

						// loop over all the other columns
						for j := range g.Columns {
							if j != i {
								if len(g.Columns[j].Cards) == 0 && g.Columns[i].Cards[x].Value == CardRanks[u.King] {
									jx, jy, jw, jh := g.Columns[j].GetColumnGeoMData()
									if u.IsCollision(ix, iy, iw, ih, jx, jy, jw, jh) &&
										inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {

										g.Columns[j].Cards = append(g.Columns[j].Cards, g.Columns[i].Cards[hidden:hidden+count]...)
										g.Columns[i].Cards = g.Columns[i].Cards[hidden : hidden+count]

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
										g.Columns[i].Cards[x].Value == g.Columns[j].Cards[lj].Value-1 &&
										g.Columns[i].Cards[x].Color != g.Columns[j].Cards[lj].Color &&
										inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {

										g.Columns[j].Cards = append(g.Columns[j].Cards, g.Columns[i].Cards[hidden:hidden+count]...)
										g.Columns[i].Cards = g.Columns[i].Cards[hidden : hidden+count]

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

						// loop over the Store slots
						for j := range g.CardStores {
							jx, jy, jw, jh := g.CardStores[j].GetStoreGeomData()

							if len(g.CardStores[j].Cards) == 0 {
								if u.IsCollision(ix, iy, iw, ih, jx, jy, jw, jh) &&
									g.Columns[i].Cards[x].Value == CardRanks[u.Ace] &&
									// check if it's the source card hovered is the last from the columns
									g.Columns[i].Cards[x].Value == g.Columns[i].Cards[len(g.Columns[i].Cards)-1].Value &&
									inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
									g.CardStores[j].Cards = append(g.CardStores[j].Cards, g.Columns[i].Cards[x])
									g.Columns[i].Cards = g.Columns[i].Cards[:x]

									// reveal the last card from the source column, and revert its height to original
									if len(g.Columns[i].Cards) > 0 {
										g.Columns[i].Cards[x-1].IsRevealed = true
										g.Columns[i].Cards[x-1].H = g.DrawCardSlot.H
									}

									// exit entirely to prevent redundant iterations
									HoveredCard = nil
									return nil
								}
							} else {
								lj := len(g.CardStores[j].Cards) - 1
								if u.IsCollision(ix, iy, iw, ih, jx, jy, jw, jh) &&
									g.Columns[i].Cards[x].Value > CardRanks[u.Ace] &&
									g.Columns[i].Cards[x].Value == g.CardStores[j].Cards[lj].Value+1 &&
									g.Columns[i].Cards[x].Suit == g.CardStores[j].Cards[lj].Suit {

									g.CardStores[j].Cards = append(g.CardStores[j].Cards, g.Columns[i].Cards[x])
									g.Columns[i].Cards = g.Columns[i].Cards[:x]

									HoveredCard = nil
									return nil
								}
							}
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

										// make the card under the current Revealed TODO: seems redundant
										last := len(g.CardStores[i].Cards)
										if last > 0 {
											g.CardStores[i].Cards[last-1].IsRevealed = true
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
