package klondike

import (
	"games-ebiten/data"
	u "games-ebiten/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	_ "image/png"
)

// Card - Implements CasinoCards interface
// RevealedState - if the card's frontFace is visible
// DraggedState - holds the dragged state
// ColNum - refers to what column the card is currently in if it's in any column
// IsActive - every card below the dragged card (to keep the state of all dragged cards, and perform multi cards drag)
type Card struct {
	Img           *ebiten.Image
	BackImg       *ebiten.Image
	Suit          string
	Color         string
	Value         int
	ColNum        int
	X, Y, W, H    int
	ScX           float64
	ScY           float64
	RevealedState bool
	IsActive      bool
	DraggedState  bool
}

func (c *Card) DrawCard(screen *ebiten.Image) {
	cx, cy := ebiten.CursorPosition()

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(c.ScX, c.ScY)
	img := c.Img

	if !c.IsRevealed() {
		img = c.BackImg
	}

	if c.IsRevealed() {
		// drag only clicked revealed cards

		if c.IsHovered(cx, cy) && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			c.SetDraggedState(true)
		}

		// drag and set location
		if inpututil.MouseButtonPressDuration(ebiten.MouseButtonLeft) > 3 && c.IsDragged() {
			c.X = cx - c.W/2
			c.Y = cy - c.H/2
			data.DraggedCard = c
		}

		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			c.SetDraggedState(false)
			data.DraggedCard = nil
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
	if c.IsRevealed() {
		img = c.Img

		// set the clicked card's drag status to true
		if c.IsHovered(cx, cy) && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			c.SetDraggedState(true)
		}

		// drag the stack of cards, and set the cards' IsActive state to true
		if inpututil.MouseButtonPressDuration(ebiten.MouseButtonLeft) > 3 && c.IsDragged() {
			data.DraggedCard = c
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
			data.DraggedCard = nil
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

func (c *Card) IsDragged() bool {
	return c.DraggedState
}

func (c *Card) SetDraggedState(state bool) {
	c.DraggedState = state
}

func (c *Card) IsRevealed() bool {
	return c.RevealedState
}

func (c *Card) SetRevealedState(state bool) {
	c.RevealedState = state
}

// IsHovered - Returns true if the card is hovered
func (c *Card) IsHovered(cx, cy int) bool {
	return image.Pt(cx, cy).In(c.HitBox())
}

func (c *Card) HitBox() image.Rectangle {
	return image.Rect(c.X, c.Y, c.X+c.W, c.Y+c.H)
}
