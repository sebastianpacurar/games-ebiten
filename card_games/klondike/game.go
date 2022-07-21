package klondike

import (
	data "games-ebiten/card_games"
	res "games-ebiten/resources"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
)

type (
	Game struct {
		*data.Theme
		*Environment
	}
)

func NewGame() *Game {
	classicImg := ebiten.NewImageFromImage(res.LoadSpriteImage("resources/assets/cards/classic-solitaire.png"))
	th := data.NewTheme()

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
	g.BuildDeck()
	return g
}

// BuildDeck - initiates the Piles and populates them with cards
func (g *Game) BuildDeck() {
	g.Deck = GenerateDeck(g.Theme)
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
			c.Y = g.SpacerV + res.MainMenuH
			c.H = g.H
			c.SetRevealedState(true)

			// draw the prior card as revealed when the current card is dragged
			if i > 0 && g.WastePile.Cards[i].Dragged() {
				g.WastePile.Cards[i-1].SetDragged(false)
				g.WastePile.Cards[i-1].DrawCard(screen)
			}
			c.DrawCard(screen)
		}

		// draw the Foundation Piles
		for i := range g.FoundationPiles {
			for j := range g.FoundationPiles[i].Cards {
				g.FoundationPiles[i].Cards[j].X = g.FoundationPiles[i].X
				g.FoundationPiles[i].Cards[j].Y = g.SpacerV + res.MainMenuH
				g.FoundationPiles[i].Cards[j].H = g.H

				// draw the prior card as revealed when the current card is dragged
				if j > 0 && g.FoundationPiles[i].Cards[j].Dragged() {
					g.FoundationPiles[i].Cards[j-1].SetDragged(false)
					g.FoundationPiles[i].Cards[j-1].DrawCard(screen)
				}

				g.FoundationPiles[i].Cards[j].DrawCard(screen)
			}
		}

		// Draw the Card Columns
		for i := range g.Columns {
			for j, card := range g.Columns[i].Cards {
				card.X = g.Columns[i].X
				card.Y = g.Columns[i].Y + (j * CardsVSpacer) + res.MainMenuH

				// draw the overlapped with the height of the space in which the card is visible
				if j != len(g.Columns[i].Cards)-1 {
					card.H = CardsVSpacer
				} else {
					card.H = g.H
				}
				card.DrawColCard(screen, g.Columns[i].Cards, j, cx, cy)
			}
		}

		// force card or stack of cards image(s) persistence over other cards
		// practically draw the dragged card again. or draw the entire stack again, at the end
		if res.DraggedCard != nil {
			switch res.DraggedCard.(type) {
			case *Card:
				c := res.DraggedCard.(*Card)
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
	ebitenutil.DebugPrintAt(screen, "Right Click to quick move to Foundations", 10, res.ScreenHeight-95)
	ebitenutil.DebugPrintAt(screen, "Press R to deal New Game", 10, res.ScreenHeight-65)
	ebitenutil.DebugPrintAt(screen, "Press 1 or 2 to change Themes", 10, res.ScreenHeight-35)

}

func (g *Game) Update() error {
	cx, cy := ebiten.CursorPosition()

	switch {
	case inpututil.IsKeyJustReleased(ebiten.Key1):
		res.ActiveCardsTheme = res.ClassicTheme
		g.BuildDeck()
	case inpututil.IsKeyJustReleased(ebiten.Key2):
		res.ActiveCardsTheme = res.PixelatedTheme
		g.BuildDeck()
	case inpututil.IsKeyJustReleased(ebiten.KeyR):
		g.BuildDeck()
	}

	if !g.IsGameOver() {
		g.HandleGameLogic(cx, cy)
	}

	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return res.ScreenWidth, res.ScreenHeight
}
