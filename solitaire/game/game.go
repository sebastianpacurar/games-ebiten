package game

import (
	"fmt"
	u "games-ebiten/resources/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"image/color"
)

type Game struct {
	Cards []*Card
	Th    *Theme
}

func NewGame() *Game {
	th := NewTheme()
	return &Game{
		Cards: GenerateDeck(th),
		Th:    th,
	}
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.Key2) && g.Th.Active != ClassicTheme {
		g.Th.Active = ClassicTheme
	} else if ebiten.IsKeyPressed(ebiten.Key1) && g.Th.Active != PixelatedTheme {
		g.Th.Active = PixelatedTheme
	}

	// when Active != LastActive, regenerate the cards, and update LastActive
	if g.Th.IsToggled() {
		g.Cards = GenerateDeck(g.Th)
		g.Th.LastActive = g.Th.Active
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.NRGBA{G: 175, B: 175, A: 255})
	for i := range g.Cards {
		g.Cards[i].DrawStaticSprite(screen)
	}

	ebitenutil.DebugPrint(screen, "Press 1 or 2 to change Themes")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return u.ScreenWidth, u.ScreenHeight
}

// GenerateDeck - returns a []*Card{} in which all elements have the corresponding details and image
func GenerateDeck(th *Theme) []*Card {
	var colStart, colEnd int
	active := th.Active

	frOX, frOY, frW, frH := th.GetFrameBoundaries(active)
	mulX, mulY := th.LocMultiplier[active][u.X], th.LocMultiplier[active][u.Y]
	deck := make([]*Card, 0, 52)

	// this logic is needed due to the discrepancy of the sprite sheets:
	// one Image starts with card Ace as the first Column value, while others start with card number or other value
	if active == PixelatedTheme {
		colStart = 1
		colEnd = 14
	} else if active == ClassicTheme {
		colStart = 0
		colEnd = 13
	}

	// there are 4 suits on the image, and 1 suit consists of 13 cards
	for si, suit := range th.SuitsOrder[active] {
		for i := colStart; i < colEnd; i++ {
			if i == 8 && active == ClassicTheme {
				fmt.Println("a")
			}
			x, y := frOX+i*frW, frOY+si*frH

			// crete card dynamically, based on the Active Theme.
			card := &Card{
				Img:    th.Sources[active].SubImage(image.Rect(x, y, x+frW, y+frH)).(*ebiten.Image),
				Suit:   suit,
				Value:  Translation[active][i],
				ScaleX: th.ScaleValue[active][u.X],
				ScaleY: th.ScaleValue[active][u.Y],
				LX:     float64(frW) * mulX * float64(i),
				LY:     float64(frH) * mulY * float64(si),
				W:      float64(frW) * th.ScaleValue[active][u.X],
				H:      float64(frH) * th.ScaleValue[active][u.Y],
			}

			// append every customized card to the deck
			deck = append(deck, card)
		}
	}
	return deck
}
