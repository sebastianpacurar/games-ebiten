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
		Cards []*Card
		Environment
		Th *Theme
	}
)

func NewGame() *Game {
	classicImg := ebiten.NewImageFromImage(u.LoadSpriteImage("resources/images/cards/classic-solitaire.png"))
	th := NewTheme()
	deck := GenerateDeck(th)

	return &Game{
		Cards: deck,
		Th:    th,
		Environment: Environment{
			SpacerV:      50,
			SpacerH:      40,
			BgImg:        classicImg.SubImage(image.Rect(700, 500, 750, 550)).(*ebiten.Image),
			EmptySlotImg: classicImg.SubImage(image.Rect(852, 384, 852+71, 384+96)).(*ebiten.Image),

			DrawCardSlot: DrawCardSlot{
				GreenSlotImg: classicImg.SubImage(image.Rect(710, 384, 710+71, 384+96)).(*ebiten.Image),
				RedSlotImg:   classicImg.SubImage(image.Rect(781, 384, 781+71, 384+96)).(*ebiten.Image),
				Cards:        []*Card{},
			},
			Stores: []CardStore{
				{Cards: []*Card{}},
				{Cards: []*Card{}},
				{Cards: []*Card{}},
				{Cards: []*Card{}},
			},
			Columns: []CardColumn{
				{Cards: []*Card{}},
				{Cards: []*Card{}},
				{Cards: []*Card{}},
				{Cards: []*Card{}},
				{Cards: []*Card{}},
				{Cards: []*Card{}},
				{Cards: []*Card{}},
			},
		},
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.DrawEnvironment(screen, g.Th)

	for _, c := range g.Cards {
		c.Img = nil
	}

	ebitenutil.DebugPrint(screen, "Press 1 or 2 to change Themes")
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.Key1) && g.Th.Active != u.ClassicTheme {
		g.Th.Active = u.ClassicTheme
	} else if ebiten.IsKeyPressed(ebiten.Key2) && g.Th.Active != u.PixelatedTheme {
		g.Th.Active = u.PixelatedTheme
	}

	// when Active != LastActive, regenerate the cards, and update LastActive
	if g.Th.IsToggled() {
		g.Cards = GenerateDeck(g.Th)
		g.Th.LastActive = g.Th.Active
	}

	for _, c := range g.Cards {
		cx, cy := ebiten.CursorPosition()
		if u.IsImgHovered(c, cx, cy) && inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
			c.IsRevealed = !c.IsRevealed
		}
		u.DragAndDrop(c)
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
	//muX, muY := th.LocMultiplier[active][u.X], th.LocMultiplier[active][u.Y]

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
