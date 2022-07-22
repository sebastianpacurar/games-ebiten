package free_cell

import (
	res "games-ebiten/resources"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"math/rand"
)

// GenerateDeck - returns a []*Card{} in which all elements have the corresponding details and images
func GenerateDeck(th *res.Theme) []*Card {
	var colStart, colEnd int
	deck := make([]*Card, 0, 52)

	cardSc := th.CardScaleValue[res.ActiveCardsTheme]

	// set which BackFace the cards have (FrOX, FRoY, FrW, FrH)
	bf := th.GetBackFrameGeomData(res.ActiveCardsTheme, res.StaticBack1)

	// set which FrontFace the cards have
	frame := th.GetFrontFrameGeomData(res.ActiveCardsTheme)

	// this logic is needed due to the discrepancy between sprite sheets:
	// one Image starts with card Ace as the first Column value, while others start with card number or other value
	switch res.ActiveCardsTheme {
	case res.PixelatedTheme:
		colStart = 1
		colEnd = 14
	case res.ClassicTheme:
		colStart = 0
		colEnd = 13
	}

	// there are 4 suits on the image, and 1 suit consists of 13 cards
	for si, suit := range th.SuitsOrder[res.ActiveCardsTheme] {
		color := ""
		switch suit {
		case res.Hearts, res.Diamonds:
			color = res.RED
		case res.Spades, res.Clubs:
			color = res.BLACK
		}

		for i := colStart; i < colEnd; i++ {
			x, y := frame.Min.X+i*frame.Dx(), frame.Min.Y+si*frame.Dy()
			w, h := frame.Dx(), frame.Dy()

			// crete card Dynamically, based on the Active Theme.
			card := &Card{
				Img:     th.Sources[res.ActiveCardsTheme].SubImage(image.Rect(x, y, x+w, y+h)).(*ebiten.Image),
				BackImg: th.Sources[res.ActiveCardsTheme].SubImage(image.Rect(bf[0], bf[1], bf[2]+bf[0], bf[3]+bf[1])).(*ebiten.Image),
				Suit:    suit,
				Value:   res.CardRanks[res.Translation[res.ActiveCardsTheme][i]],
				Color:   color,
				ScX:     cardSc[res.X],
				ScY:     cardSc[res.Y],
				W:       int(float64(w) * cardSc[res.X]),
				H:       int(float64(h) * cardSc[res.Y]),
			}

			// append every customized card to the deck
			deck = append(deck, card)
		}
		rand.Shuffle(len(deck), func(i, j int) {
			deck[i], deck[j] = deck[j], deck[i]
		})
	}
	return deck
}
