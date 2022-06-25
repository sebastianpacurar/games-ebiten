package game

import (
	cg "games-ebiten/card_games"
	u "games-ebiten/resources/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
)

type (
	Game struct {
		*cg.Theme
		*Environment
	}
)

func NewGame() *Game {
	classicImg := ebiten.NewImageFromImage(u.LoadSpriteImage("resources/images/cards/classic-solitaire.png"))
	th := cg.NewTheme()
	g := &Game{
		Theme: th,
		Environment: &Environment{
			Quadrants:    make(map[int]image.Rectangle, 0),
			SpacerV:      50,
			BgImg:        classicImg.SubImage(image.Rect(700, 500, 750, 550)).(*ebiten.Image),
			EmptySlotImg: classicImg.SubImage(image.Rect(852, 384, 852+71, 384+96)).(*ebiten.Image),
		},
	}
	g.BuildDeck(th)
	return g
}

// BuildDeck - initiates the Piles and populates them with cards
func (g *Game) BuildDeck(th *cg.Theme) {
	g.Deck = cg.GenerateDeck(th)
	g.UpdateEnv()
}

func (g *Game) Draw(screen *ebiten.Image) {
	cx, cy := ebiten.CursorPosition()

	g.DrawPlayground(screen, g.Theme)

	for i := range g.FreeCells {
		for j := range g.FreeCells[i].Cards {
			g.FreeCells[i].Cards[j].X = g.FreeCells[i].X
			g.FreeCells[i].Cards[j].Y = g.SpacerV
			g.FreeCells[i].Cards[j].H = g.H

			// draw the prior card as revealed when the current card is dragged
			if j > 0 && g.FoundationPiles[i].Cards[j].IsDragged() {
				g.FoundationPiles[i].Cards[j-1].SetDraggedState(false)
				g.FoundationPiles[i].Cards[j-1].DrawCard(screen)
			}

			g.FoundationPiles[i].Cards[j].DrawCard(screen)
		}
	}

	// draw the Foundation Piles
	for i := range g.FoundationPiles {
		for j := range g.FoundationPiles[i].Cards {
			g.FoundationPiles[i].Cards[j].X = g.FoundationPiles[i].X
			g.FoundationPiles[i].Cards[j].Y = g.SpacerV
			g.FoundationPiles[i].Cards[j].H = g.H

			// draw the prior card as revealed when the current card is dragged
			if j > 0 && g.FoundationPiles[i].Cards[j].IsDragged() {
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
	if cg.DraggedCard != nil {
		switch cg.DraggedCard.(type) {
		case *cg.Card:
			c := cg.DraggedCard.(*cg.Card)
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
	ebitenutil.DebugPrintAt(screen, "Right Click to quick move to Foundations", 10, u.ScreenHeight-95)
	ebitenutil.DebugPrintAt(screen, "Press F2 to start new game", 10, u.ScreenHeight-65)
	ebitenutil.DebugPrintAt(screen, "Press 1 or 2 to change Themes", 10, u.ScreenHeight-35)
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

	g.HandleGameLogic(cx, cy)

	return nil
}

// HandleGameLogic - contains Drag from Waste Pile, valid Column, Foundation Pile to any valid slot; and right click to foundations functionalities
func (e *Environment) HandleGameLogic(cx, cy int) {
	//e.RightClickToFoundations(cx, cy)

	if cg.DraggedCard != nil {
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

								if c.IsDragged() && c.Value == cg.CardRanks[u.King] {
									target := e.HitBox(e.Columns[j])
									if u.IsCollision(source, target) &&
										inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
										e.MoveFromSrcToTarget(e.Columns, e.Columns, i, j)
										cg.DraggedCard = nil
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
									if c.IsDragged() {
										if u.IsCollision(source, target) &&
											inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) &&
											c.Value == e.Columns[j].Cards[lj].Value-1 &&
											c.Color != e.Columns[j].Cards[lj].Color {
											e.MoveFromSrcToTarget(e.Columns, e.Columns, i, j)
											cg.DraggedCard = nil
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
						if u.IsCollision(source, target) &&
							inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) &&
							e.Columns[i].Cards[li].Value == cg.CardRanks[u.Ace] {
							e.MoveFromSrcToTarget(e.Columns, e.FoundationPiles, i, j)
							cg.DraggedCard = nil
							return
						}
					} else {
						lj := len(e.FoundationPiles[j].Cards) - 1
						if u.IsCollision(source, target) &&
							inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) &&
							e.Columns[i].Cards[li].Value > cg.CardRanks[u.Ace] &&
							e.Columns[i].Cards[li].Value == e.FoundationPiles[j].Cards[lj].Value+1 &&
							e.Columns[i].Cards[li].Suit == e.FoundationPiles[j].Cards[lj].Suit {
							e.MoveFromSrcToTarget(e.Columns, e.FoundationPiles, i, j)
							cg.DraggedCard = nil
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
				if e.FoundationPiles[i].Cards[li].IsDragged() {
					source := e.FoundationPiles[i].Cards[li].HitBox()

					// drop ON Column
					for j := range e.Columns {
						if len(e.Columns[j].Cards) > 0 {
							lj := len(e.Columns[j].Cards) - 1 // lj = last card in the current context
							target := e.Columns[j].Cards[lj].HitBox()

							if u.IsCollision(source, target) && inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) &&
								e.FoundationPiles[i].Cards[li].Value == e.Columns[j].Cards[lj].Value-1 &&
								e.Columns[j].Cards[lj].Color != e.FoundationPiles[i].Cards[li].Color {
								e.MoveFromSrcToTarget(e.FoundationPiles, e.Columns, i, j)
								cg.DraggedCard = nil
								return
							}

						}
					}

					// drop ON Foundation Pile
					for j := range e.FoundationPiles {
						if i != j {
							target := e.HitBox(e.FoundationPiles[j])

							if len(e.FoundationPiles[j].Cards) == 0 && u.IsCollision(source, target) &&
								inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) &&
								e.FoundationPiles[i].Cards[li].Value == cg.CardRanks[u.Ace] {
								e.MoveFromSrcToTarget(e.FoundationPiles, e.FoundationPiles, i, j)
								cg.DraggedCard = nil
								return
							}
						}
					}
				}
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return u.ScreenWidth, u.ScreenHeight
}
