package game

import (
	u "games-ebiten/resources/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	_ "image/png"
)

// Translation - used for acquiring the right card index while getting the card SubImage from the Image
// DraggedCard - used to force the hovered card to overlap other images, while dragged
// CardRanks - smallest is "Ace"(0), while highest is "King"(13)
var (
	DraggedCard interface{}
	Translation = map[string]map[int]string{
		u.PixelatedTheme: {
			1: "2", 2: "3", 3: "4", 4: "5", 5: "6", 6: "7",
			7: "8", 8: "9", 9: "10", 10: "J", 11: "Q", 12: "K", 13: "A",
		},
		u.ClassicTheme: {
			0: "A", 1: "2", 2: "3", 3: "4", 4: "5", 5: "6", 6: "7",
			7: "8", 8: "9", 9: "10", 10: "J", 11: "Q", 12: "K",
		},
		//u.AbstractTheme: {
		//	0: "2", 1: "3", 2: "4", 3: "5", 4: "6", 5: "7",
		//	6: "8", 7: "9", 8: "10", 9: "Jack", 10: "Queen", 11: "King", 12: "Ace",
		//},
	}
	CardRanks = map[string]int{
		u.Ace:   0,
		"2":     1,
		"3":     2,
		"4":     3,
		"5":     4,
		"6":     5,
		"7":     6,
		"8":     7,
		"9":     8,
		"10":    9,
		u.Jack:  10,
		u.Queen: 11,
		u.King:  12,
	}
)

// Card - Implements CasinoCards interface
// IsRevealed - if the card's frontFace is visible
// IsDragged - holds the dragged state
// ColNum - refers to what column the card is currently in if it's in any column
// IsActive - every card below the dragged card (to keep the state of all dragged cards, and perform multi cards drag)
type Card struct {
	Img        *ebiten.Image
	BackImg    *ebiten.Image
	Suit       string
	Color      string
	Value      int
	ColNum     int
	X, Y, W, H int
	ScX        float64
	ScY        float64
	IsRevealed bool
	IsActive   bool
	IsDragged  bool
}

func (c *Card) DrawCard(screen *ebiten.Image) {
	cx, cy := ebiten.CursorPosition()

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(c.ScX, c.ScY)
	img := c.Img

	if !c.GetRevealedState() {
		img = c.BackImg
	}

	if c.GetRevealedState() {
		// drag only clicked revealed cards

		if c.IsHovered(cx, cy) && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			c.SetDraggedState(true)
		}

		// drag and set location
		if inpututil.MouseButtonPressDuration(ebiten.MouseButtonLeft) > 2 && c.GetDraggedState() {
			c.X = cx - c.W/2
			c.Y = cy - c.H/2
			DraggedCard = c
		}

		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			c.SetDraggedState(false)
			DraggedCard = nil
		}
	}

	op.GeoM.Translate(float64(c.X), float64(c.Y))
	screen.DrawImage(img, op)
}

// DrawColCard - handles the drag multiple revealed cards functionality
// cards param = cards array of which the DraggedCard is part
// ci param = index of the DraggedCard
// cx, cy params = cursor position
func (c *Card) DrawColCard(screen *ebiten.Image, cards []*Card, ci, cx, cy int) {
	var img *ebiten.Image

	// if card is revealed
	if c.GetRevealedState() {
		img = c.Img

		// set the clicked card's drag status to true
		if c.IsHovered(cx, cy) && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			c.SetDraggedState(true)
		}

		// drag the stack of cards, and set the cards' IsActive state to true
		if inpututil.MouseButtonPressDuration(ebiten.MouseButtonLeft) > 1 && c.GetDraggedState() {
			DraggedCard = c
			c.IsActive = true
			c.X = cx - c.W/2
			c.Y = cy - c.H/2

			// if not the last card then height is 35, and c.Y is 35/2
			if len(cards[ci:]) > 1 && c != cards[len(cards[ci:])-1] {
				c.H = u.CardsVSpacer
				c.Y = cy
			}

			// draw the dragged card first
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(c.ScX, c.ScY)
			op.GeoM.Translate(float64(c.X), float64(c.Y))
			screen.DrawImage(c.Img, op)

			// draw the rest of the cards
			for i, card := range cards[ci+1:] {
				if i != len(cards[ci+1:])-1 {
					card.H = u.CardsVSpacer
				}
				card.X = cx - card.W/2
				card.Y = cy + ((i + 1) * u.CardsVSpacer)
				card.IsActive = true
				opc := &ebiten.DrawImageOptions{}
				opc.GeoM.Scale(card.ScX, card.ScY)
				opc.GeoM.Translate(float64(card.X), float64(card.Y))
				screen.DrawImage(card.Img, opc)
			}
		}

		// upon release, set dragged and active states to false
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			c.SetDraggedState(false)
			c.IsActive = false
			DraggedCard = nil
		}
	} else {
		// if card is not revealed, draw back face image
		img = c.BackImg
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(c.ScX, c.ScY)
		op.GeoM.Translate(float64(c.X), float64(c.Y))
		screen.DrawImage(img, op)
	}

	// draw the revealed cards above the DraggedCard
	if !c.IsActive {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(c.ScX, c.ScY)
		op.GeoM.Translate(float64(c.X), float64(c.Y))
		screen.DrawImage(img, op)
	}
}

func (c *Card) GetDraggedState() bool {
	return c.IsDragged
}

func (c *Card) SetDraggedState(state bool) {
	c.IsDragged = state
}

func (c *Card) GetRevealedState() bool {
	return c.IsRevealed
}

func (c *Card) SetRevealedState(state bool) {
	c.IsRevealed = state
}

// IsHovered - Returns true if the card is hovered
func (c *Card) IsHovered(cx, cy int) bool {
	return image.Pt(cx, cy).In(c.GetGeomData())
}

func (c *Card) GetGeomData() image.Rectangle {
	return image.Rect(c.X, c.Y, c.X+c.W, c.Y+c.H)
}
