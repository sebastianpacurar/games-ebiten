package k_game

import (
	data2 "games-ebiten/card_games/data"
	u "games-ebiten/resources/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
)

type (
	Game struct {
		*data2.Theme
		*Environment
	}
)

func NewGame() *Game {
	classicImg := ebiten.NewImageFromImage(u.LoadSpriteImage("resources/images/cards/classic-solitaire.png"))
	th := data2.NewTheme()

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
func (g *Game) BuildDeck(th *data2.Theme) {
	g.Deck = GenerateDeck(th)
	g.UpdateEnv()
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
			if i > 0 && g.WastePile.Cards[i].IsDragged() {
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
		if data2.DraggedCard != nil {
			switch data2.DraggedCard.(type) {
			case *Card:
				c := data2.DraggedCard.(*Card)
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
	ebitenutil.DebugPrintAt(screen, "Right Click to quick move to Foundations", 10, u.ScreenHeight-95)
	ebitenutil.DebugPrintAt(screen, "Press F2 to start new Game", 10, u.ScreenHeight-65)
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

	if !g.IsGameOver() {
		g.HandleGameLogic(cx, cy)
	}

	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return u.ScreenWidth, u.ScreenHeight
}
