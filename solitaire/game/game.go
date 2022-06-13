package game

import (
	u "games-ebiten/resources/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	"math/rand"
	"time"
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

	g := &Game{
		Theme: th,
		Environment: &Environment{
			Quadrants:    make(map[int]image.Rectangle, 0),
			SpacerV:      50,
			BgImg:        classicImg.SubImage(image.Rect(700, 500, 750, 550)).(*ebiten.Image),
			EmptySlotImg: classicImg.SubImage(image.Rect(852, 384, 852+71, 384+96)).(*ebiten.Image),

			StockPile: StockPile{
				GreenSlotImg: classicImg.SubImage(image.Rect(710, 384, 710+71, 384+96)).(*ebiten.Image),
				RedSlotImg:   classicImg.SubImage(image.Rect(781, 384, 781+71, 384+96)).(*ebiten.Image),
			},
		},
	}
	g.BuildDeck(th)
	return g
}

// BuildDeck - initiates the Piles and populates them with cards
func (g *Game) BuildDeck(th *Theme) {
	g.Deck = g.GenerateDeck(th)
	g.UpdateEnv()
}

// GenerateDeck - returns a []*Card{} in which all elements have the corresponding details and images
func (g *Game) GenerateDeck(th *Theme) []*Card {
	var colStart, colEnd int
	deck := make([]*Card, 0, 52)
	active := th.Active
	cardSc := th.CardScaleValue[active]

	// set which BackFace the cards have (FrOX, FRoY, FrW, FrH)
	bf := th.GetBackFrameGeomData(active, u.StaticBack1)

	// set which FrontFace the cards have
	frame := th.GetFrontFrameGeomData(active)

	// this logic is needed due to the discrepancy between sprite sheets:
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
			x, y := frame.Min.X+i*frame.Dx(), frame.Min.Y+si*frame.Dy()
			w, h := frame.Dx(), frame.Dy()

			// crete card dynamicalY, based on the Active Theme.
			card := &Card{
				Img:     th.Sources[active].SubImage(image.Rect(x, y, x+w, y+h)).(*ebiten.Image),
				BackImg: th.Sources[active].SubImage(image.Rect(bf[0], bf[1], bf[2]+bf[0], bf[3]+bf[1])).(*ebiten.Image),
				Suit:    suit,
				Value:   CardRanks[Translation[active][i]],
				Color:   color,
				ScX:     cardSc[u.X],
				ScY:     cardSc[u.Y],
				W:       int(float64(w) * cardSc[u.X]),
				H:       int(float64(h) * cardSc[u.Y]),
			}

			// append every customized card to the deck
			deck = append(deck, card)
		}
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(deck), func(i, j int) {
			deck[i], deck[j] = deck[j], deck[i]
		})
	}
	return deck
}

func (g *Game) Draw(screen *ebiten.Image) {
	cx, cy := ebiten.CursorPosition()

	if !g.IsGameOver() {
		g.DrawPlayground(screen, g.Theme)

		// Draw the Stock Pile
		if len(g.StockPile.Cards) > 0 {
			for _, c := range g.StockPile.Cards {
				c.X = g.StockPile.X
				c.Y = g.StockPile.Y
				c.H = g.H
				c.DrawCard(screen)
			}
		}

		// Draw the Waste Pile
		for i, c := range g.WastePile.Cards {
			c.X = g.WastePile.X
			c.Y = g.SpacerV
			c.H = g.H
			c.SetRevealedState(true)

			// draw the prior card as revealed when the current card is dragged
			if i > 0 && g.WastePile.Cards[i].GetDraggedState() {
				g.WastePile.Cards[i-1].SetDraggedState(false)
				g.WastePile.Cards[i-1].DrawCard(screen)
			}
			c.DrawCard(screen)
		}

		// draw the Foundation Piles
		for i := range g.FoundationPiles {
			for j := range g.FoundationPiles[i].Cards {
				g.FoundationPiles[i].Cards[j].X = g.FoundationPiles[i].X
				g.FoundationPiles[i].Cards[j].Y = g.SpacerV
				g.FoundationPiles[i].Cards[j].H = g.H

				// draw the prior card as revealed when the current card is dragged
				if j > 0 && g.FoundationPiles[i].Cards[j].GetDraggedState() {
					g.FoundationPiles[i].Cards[j-1].SetDraggedState(false)
					g.FoundationPiles[i].Cards[j-1].DrawCard(screen)
				}

				g.FoundationPiles[i].Cards[j].DrawCard(screen)
			}
		}

		// Draw the Card Columns
		for i := range g.Columns {
			for j, card := range g.Columns[i].Cards {
				card.X = g.Columns[i].X
				card.Y = g.Columns[i].Y + (j * u.CardsVSpacer)

				// draw the overlapped with the height of the space in which the card is visible
				if j != len(g.Columns[i].Cards)-1 {
					card.H = u.CardsVSpacer
				} else {
					card.H = g.H
				}
				card.DrawColCard(screen, g.Columns[i].Cards, j, cx, cy)
			}
		}

		// force card or stack of cards image(s) persistence over other cards
		// practically draw the dragged card again. or draw the entire stack again, at the end
		if DraggedCard != nil {
			switch DraggedCard.(type) {
			case *Card:
				c := DraggedCard.(*Card)
				if c.ColNum == 0 {
					opc := &ebiten.DrawImageOptions{}
					opc.GeoM.Scale(c.ScX, c.ScY)
					opc.GeoM.Translate(float64(c.X), float64(c.Y))
					screen.DrawImage(c.Img, opc)
				} else {
					for i, card := range g.Columns[c.ColNum-1].Cards {
						card.DrawColCard(screen, g.Columns[c.ColNum-1].Cards, i, cx, cy)
					}
				}
			}
		}
	} else {
		g.DrawEnding(screen)
	}
	ebitenutil.DebugPrintAt(screen, "Press F2 to start new game", 10, u.ScreenHeight-35)
	ebitenutil.DebugPrintAt(screen, "Press 1 or 2 to change Themes", 10, u.ScreenHeight-65)
}

func (g *Game) Update() error {
	cx, cy := ebiten.CursorPosition()

	switch {
	case inpututil.IsKeyJustReleased(ebiten.Key1):
		g.Active = u.ClassicTheme
		g.BuildDeck(g.Theme)
	case inpututil.IsKeyJustReleased(ebiten.Key2):
		g.Active = u.PixelatedTheme
		g.BuildDeck(g.Theme)
	case inpututil.IsKeyJustReleased(ebiten.KeyF2):
		g.BuildDeck(g.Theme)
	}

	if !g.IsGameOver() {

		//
		// Handle Stock to Waste functionality
		//
		if len(g.StockPile.Cards) > 0 {
			if g.StockPile.Cards[0].IsHovered(cx, cy) {
				last := len(g.StockPile.Cards) - 1
				if len(g.StockPile.Cards) == 1 {
					last = 0
				}
				// append every last card from StockPile to WastePile, then trim last card from StockPile
				if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
					g.WastePile.Cards = append(g.WastePile.Cards, g.StockPile.Cards[last])
					g.StockPile.Cards = g.StockPile.Cards[:last]
				}
			}
		} else {
			// if there are no more cards, clicking the circle will reset the process
			if !g.IsVegas && g.IsStockPileHovered(cx, cy) && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				for i := range g.WastePile.Cards {
					g.StockPile.Cards = append(g.StockPile.Cards, g.WastePile.Cards[i])
					g.StockPile.Cards[i].SetRevealedState(false)
				}
				g.WastePile.Cards = g.WastePile.Cards[:0]

				// reverse order of newly stacked StockPile cards:
				for i, j := 0, len(g.StockPile.Cards)-1; i < j; i, j = i+1, j-1 {
					g.StockPile.Cards[i], g.StockPile.Cards[j] = g.StockPile.Cards[j], g.StockPile.Cards[i]
				}
			}
		}

		if DraggedCard != nil {
			//
			// drag from Waste Pile to valid Column or FoundationPile slot
			//
			if len(g.WastePile.Cards) > 0 {
				lc := len(g.WastePile.Cards) - 1
				source := g.WastePile.Cards[lc].GetGeomData()

				// set the prior card's state dragged to false, so it can stick to its location
				if g.WastePile.Cards[lc].GetDraggedState() {
					if len(g.WastePile.Cards) > 1 {
						g.WastePile.Cards[lc-1].SetRevealedState(false)
					}

					// drag from Waste Pile to Column Slot
					for j := range g.Columns {
						if len(g.Columns[j].Cards) == 0 && g.WastePile.Cards[lc].Value == CardRanks[u.King] {
							// if there are no cards on the Column Slot and the source card is a King
							target := g.GetGeomData(g.Columns[j])

							if u.IsCollision(source, target) && inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
								g.WastePile.Cards[lc].ColNum = j + 1
								g.Columns[j].Cards = append(g.Columns[j].Cards, g.WastePile.Cards[lc])
								g.WastePile.Cards = g.WastePile.Cards[:lc]

								// exit entirely to prevent redundant iterations
								DraggedCard = nil
								return nil
							}
						} else if len(g.Columns[j].Cards) > 0 {
							// applies for any other card than K, also prevents iteration over the empty slots if there are any
							lj := len(g.Columns[j].Cards) - 1 // lj = last card in the current context
							target := g.Columns[j].Cards[lj].GetGeomData()

							if u.IsCollision(source, target) &&
								g.WastePile.Cards[lc].Value == g.Columns[j].Cards[lj].Value-1 &&
								g.WastePile.Cards[lc].Color != g.Columns[j].Cards[lj].Color &&
								inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {

								g.WastePile.Cards[lc].ColNum = j + 1
								g.Columns[j].Cards = append(g.Columns[j].Cards, g.WastePile.Cards[lc])
								g.WastePile.Cards = g.WastePile.Cards[:lc]

								// exit entirely to prevent redundant iterations
								DraggedCard = nil
								return nil
							}
						}
					}

					// draw from Waste Pile to Foundation Piles
					for j := range g.FoundationPiles {
						target := g.GetGeomData(g.FoundationPiles[j])

						if len(g.FoundationPiles[j].Cards) == 0 {
							if g.WastePile.Cards[lc].Value == CardRanks[u.Ace] &&
								u.IsCollision(source, target) && inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {

								g.WastePile.Cards[lc].ColNum = 0
								g.FoundationPiles[j].Cards = append(g.FoundationPiles[j].Cards, g.WastePile.Cards[lc])
								g.WastePile.Cards = g.WastePile.Cards[:lc]

								// exit entirely to prevent redundant iterations
								DraggedCard = nil
								return nil
							}

						} else {
							lj := len(g.FoundationPiles[j].Cards) - 1
							if u.IsCollision(source, target) &&
								g.WastePile.Cards[lc].Value > CardRanks[u.Ace] &&
								g.WastePile.Cards[lc].Value == g.FoundationPiles[j].Cards[lj].Value+1 &&
								g.WastePile.Cards[lc].Suit == g.FoundationPiles[j].Cards[lj].Suit &&
								inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {

								g.WastePile.Cards[lc].ColNum = 0
								g.FoundationPiles[j].Cards = append(g.FoundationPiles[j].Cards, g.WastePile.Cards[lc])
								g.WastePile.Cards = g.WastePile.Cards[:lc]

								// exit entirely to prevent redundant iterations
								DraggedCard = nil
								return nil
							}
						}
					}
				}
			}

			//
			// drag from Column to valid Column or Foundation Pile
			//
			for i := range g.Columns {
				if len(g.Columns[i].Cards) > 0 {
					li := len(g.Columns[i].Cards) - 1 // li = last card in the source column
					source := g.Columns[i].Cards[li].GetGeomData()

					// drag card(s) from Column to Column
					for j := range g.Columns {

						// avoid iteration over the same column
						if j != i {

							// handle moving K Column card or stack on empty Column Slot
							if len(g.Columns[j].Cards) == 0 {
								for _, c := range g.Columns[i].Cards {

									if c.Value == CardRanks[u.King] {
										target := g.GetGeomData(g.Columns[j])

										if u.IsCollision(c.GetGeomData(), target) &&
											inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {

											g.Columns[j].Cards = append(g.Columns[j].Cards, g.Columns[i].Cards[g.GetIndexOfDraggedColCard(c.ColNum-1):]...)
											g.Columns[i].Cards = g.Columns[i].Cards[:g.GetIndexOfDraggedColCard(c.ColNum-1)]

											for _, card := range g.Columns[j].Cards {
												card.ColNum = j + 1
											}

											// reveal the last card from the source column, and revert its height to original
											last := len(g.Columns[i].Cards)
											if last > 0 {
												g.Columns[i].Cards[last-1].SetRevealedState(true)
												g.Columns[i].Cards[last-1].H = g.H
											}

											// exit entirely to prevent redundant iterations
											DraggedCard = nil
											return nil
										}
									}
								}

							} else {
								// handle all cases except K
								if len(g.Columns[j].Cards) > 0 {
									lj := len(g.Columns[j].Cards) - 1 // lj = last card in the current context (target)
									target := g.Columns[j].Cards[lj].GetGeomData()

									for _, c := range g.Columns[i].Cards {
										if u.IsCollision(c.GetGeomData(), target) && c.GetDraggedState() &&
											c.Value == g.Columns[j].Cards[lj].Value-1 &&
											c.Color != g.Columns[j].Cards[lj].Color {

											if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {

												g.Columns[j].Cards = append(g.Columns[j].Cards, g.Columns[i].Cards[g.GetIndexOfDraggedColCard(c.ColNum-1):]...)
												g.Columns[i].Cards = g.Columns[i].Cards[:g.GetIndexOfDraggedColCard(c.ColNum-1)]

												for _, card := range g.Columns[j].Cards {
													card.ColNum = j + 1
												}

												// reveal the last card from the source column, and revert its height to original
												last := len(g.Columns[i].Cards)
												if last > 0 {
													g.Columns[i].Cards[last-1].SetRevealedState(true)
													g.Columns[i].Cards[last-1].H = g.H
												}

												// exit entirely to prevent redundant iterations
												DraggedCard = nil
												return nil
											}
										}
									}
								}
							}
						}
					}

					// loop over the Foundation Piles
					for j := range g.FoundationPiles {
						target := g.GetGeomData(g.FoundationPiles[j])

						if len(g.FoundationPiles[j].Cards) == 0 {
							if u.IsCollision(source, target) &&
								g.Columns[i].Cards[li].GetDraggedState() &&
								g.Columns[i].Cards[li].Value == CardRanks[u.Ace] &&
								inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {

								g.Columns[i].Cards[li].ColNum = 0
								g.FoundationPiles[j].Cards = append(g.FoundationPiles[j].Cards, g.Columns[i].Cards[li])
								g.Columns[i].Cards = g.Columns[i].Cards[:li]

								// reveal the last card from the source column, and revert its height to original
								if len(g.Columns[i].Cards) > 0 {
									g.Columns[i].Cards[li-1].SetRevealedState(true)
									g.Columns[i].Cards[li-1].H = g.H
								}

								// exit entirely to prevent redundant iterations
								DraggedCard = nil
								return nil
							}
						} else {
							lj := len(g.FoundationPiles[j].Cards) - 1
							if u.IsCollision(source, target) &&
								g.Columns[i].Cards[li].GetDraggedState() &&
								g.Columns[i].Cards[li].Value > CardRanks[u.Ace] &&
								g.Columns[i].Cards[li].Value == g.FoundationPiles[j].Cards[lj].Value+1 &&
								g.Columns[i].Cards[li].Suit == g.FoundationPiles[j].Cards[lj].Suit &&
								inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {

								g.Columns[i].Cards[li].ColNum = 0
								g.FoundationPiles[j].Cards = append(g.FoundationPiles[j].Cards, g.Columns[i].Cards[li])
								g.Columns[i].Cards = g.Columns[i].Cards[:li]

								// reveal the last card from the source column, and revert its height to original
								if len(g.Columns[i].Cards) > 0 {
									g.Columns[i].Cards[li-1].SetRevealedState(true)
									g.Columns[i].Cards[li-1].H = g.H
								}

								DraggedCard = nil
								return nil
							}
						}
					}
				}
			}

			//
			// drag from Foundation Pile to valid Column or another Foundation Pile
			//
			for i := range g.FoundationPiles {
				if len(g.FoundationPiles[i].Cards) > 0 {
					li := len(g.FoundationPiles[i].Cards) - 1 // li = last card in the current context
					if g.FoundationPiles[i].Cards[li].GetDraggedState() {
						source := g.FoundationPiles[i].Cards[li].GetGeomData()

						// loop over all the columns
						for j := range g.Columns {
							if len(g.Columns[j].Cards) > 0 {
								lj := len(g.Columns[j].Cards) - 1 // lj = last card in the current context
								target := g.Columns[j].Cards[lj].GetGeomData()

								if u.IsCollision(source, target) &&
									g.FoundationPiles[i].Cards[li].Value == g.Columns[j].Cards[lj].Value-1 &&
									g.Columns[j].Cards[lj].Color != g.FoundationPiles[i].Cards[li].Color {

									if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
										g.FoundationPiles[i].Cards[li].ColNum = j + 1
										g.Columns[j].Cards = append(g.Columns[j].Cards, g.FoundationPiles[i].Cards[li])
										g.FoundationPiles[i].Cards = g.FoundationPiles[i].Cards[:li]

										// exit entirely to prevent redundant iterations
										DraggedCard = nil
										return nil
									}
								}
							}
						}

						// loop over the Other Foundation Piles (this applies only to the Ace card being moved from a Store to another)
						for j := range g.FoundationPiles {
							if i != j {
								target := g.GetGeomData(g.FoundationPiles[j])

								if len(g.FoundationPiles[j].Cards) == 0 {
									if g.FoundationPiles[i].Cards[li].Value == CardRanks[u.Ace] && u.IsCollision(source, target) {
										if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
											g.FoundationPiles[j].Cards = append(g.FoundationPiles[j].Cards, g.FoundationPiles[i].Cards[li])
											g.FoundationPiles[i].Cards = g.FoundationPiles[i].Cards[:li]

											DraggedCard = nil
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
	}

	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return u.ScreenWidth, u.ScreenHeight
}

func (g *Game) GetIndexOfDraggedColCard(col int) int {
	count := 0
	for _, c := range g.Columns[col].Cards {
		if c.GetDraggedState() {
			break
		}
		count++
	}
	return count
}
