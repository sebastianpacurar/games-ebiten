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

	game := &Game{
		Theme: th,
		Environment: &Environment{
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
			Stores: []CardStore{
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
	_, _, frW, frH := th.GetFrontFrameGeomData(game.Active)
	x, y := float64(frW)+game.SpacerH, game.SpacerV
	w, h := float64(frW)*th.ScaleValue[game.Active][u.X], float64(frH)*th.ScaleValue[game.Active][u.Y]

	game.DrawCardSlot.X = x
	game.DrawCardSlot.Y = y
	game.DrawCardSlot.W = w
	game.DrawCardSlot.H = h

	// fill every column array with its relative count of cards
	cardIndex := 0
	for i := range game.Columns {
		for j := 0; j <= i; j++ {
			// keep only the last one revealed
			if j == i {
				deck[cardIndex].IsRevealed = true
			}
			game.Columns[i].Cards = append(game.Columns[i].Cards, deck[cardIndex])
			cardIndex++
		}
	}

	// fill the DrawCard array
	for i := range deck[cardIndex:] {
		game.DrawCardSlot.Cards = append(game.DrawCardSlot.Cards, deck[cardIndex:][i])
	}

	return game
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.DrawEnvironment(screen, g.Theme)

	// Draw the Card Columns
	for i := range g.Columns {
		x := (float64(g.FrontFaceFrameData[g.Active][u.FrW]) + g.SpacerH) * (float64(i) + 1)
		for j := range g.Columns[i].Cards {
			y := float64(u.ScreenHeight/3) + float64(j*20)
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

	// Draw the Drawn Card Slot
	for i := range g.DrawnCardsSlot.Cards {
		x := (float64(g.FrontFaceFrameData[g.Active][u.FrW]) + g.SpacerH) * 2
		y := g.SpacerV
		g.DrawnCardsSlot.Cards[i].X = x
		g.DrawnCardsSlot.Cards[i].Y = y
		g.DrawnCardsSlot.Cards[i].IsRevealed = true
		g.DrawnCardsSlot.Cards[i].DrawCardSprite(screen)
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
	//
	// handle the Draw Card functionality
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
		for i := colStart; i < colEnd; i++ {

			x, y := frOX+i*frW, frOY+si*frH
			w, h := float64(frW)*th.ScaleValue[active][u.X], float64(frH)*th.ScaleValue[active][u.Y]

			// crete card dynamicalY, based on the Active Theme.
			card := &Card{
				Img:     th.Sources[active].SubImage(image.Rect(x, y, x+frW, y+frH)).(*ebiten.Image),
				BackImg: th.Sources[active].SubImage(image.Rect(bf[0], bf[1], bf[2]+bf[0], bf[3]+bf[1])).(*ebiten.Image),
				Suit:    suit,
				Value:   Translation[active][i],
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
