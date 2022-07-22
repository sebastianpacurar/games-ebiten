package free_cell

import (
	res "games-ebiten/resources"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
)

type (
	Game struct {
		*res.Theme
		*Environment
	}
)

func NewGame() *Game {
	classicImg := ebiten.NewImageFromImage(res.LoadSpriteImage("resources/assets/cards/classic-solitaire.png"))
	th := res.NewTheme()
	g := &Game{
		Theme: th,
		Environment: &Environment{
			Deck:         make([]*Card, 0, 52),
			Quadrants:    make(map[int]image.Rectangle, 0),
			SpacerV:      50,
			BgImg:        classicImg.SubImage(image.Rect(700, 500, 750, 550)).(*ebiten.Image),
			EmptySlotImg: classicImg.SubImage(image.Rect(852, 384, 852+71, 384+96)).(*ebiten.Image),
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

	g.DrawPlayground(screen, g.Theme)

	for i := range g.FreeCells {
		for j := range g.FreeCells[i].Cards {
			g.FreeCells[i].Cards[j].X = g.FreeCells[i].X
			g.FreeCells[i].Cards[j].Y = g.SpacerV + res.MainMenuH
			g.FreeCells[i].Cards[j].H = g.H
			g.FreeCells[i].Cards[j].DrawCard(screen)
		}
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

	ebitenutil.DebugPrintAt(screen, "Press R to deal New Game", 10, res.ScreenHeight-65)
	ebitenutil.DebugPrintAt(screen, "Press 1 or 2 to change Themes", 10, res.ScreenHeight-35)
}

func (g *Game) Update() error {
	g.HandleGameLogic()
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return res.ScreenWidth, res.ScreenHeight
}
