package free_cell

import (
	res "games-ebiten/resources"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	_ "image/png"
)

var CardsVSpacer = 25

// Card
// DraggableState - refers to a Column card or stack of cards, which can be dragged
// DraggedState - holds the dragged state
// ColNum - refers to what column the card is currently in if it's in any column
// IsActive - every card below the dragged card (to keep the state of all dragged cards, and perform multi cards drag)
type Card struct {
	Img            *ebiten.Image
	BackImg        *ebiten.Image
	Suit           string
	Color          string
	Value          int
	ColNum         int
	X, Y, W, H     int
	ScX            float64
	ScY            float64
	IsActive       bool
	DraggedState   bool
	DraggableState bool
}

func (c *Card) DrawCard(screen *ebiten.Image) {
	cx, cy := ebiten.CursorPosition()

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(c.ScX, c.ScY)

	if c.Hovered(cx, cy) && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		c.SetDragged(true)
	}

	// drag and set location
	if inpututil.MouseButtonPressDuration(ebiten.MouseButtonLeft) > 3 && c.Dragged() {
		c.X = cx - c.W/2
		c.Y = cy - c.H/2
		res.DraggedCard = c
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		c.SetDragged(false)
		res.DraggedCard = nil
	}

	op.GeoM.Translate(float64(c.X), float64(c.Y))
	screen.DrawImage(c.Img, op)
}

// DrawColCard - handles the drag multiple revealed cards functionality
// cards param = cards array of which the DraggedCard is part
// ci param = index of the DraggedCard
// cx, cy params = cursor position
func (c *Card) DrawColCard(screen *ebiten.Image, cards []*Card, ci, cx, cy int) {

	// set the clicked card's drag status to true
	if c.Hovered(cx, cy) && c.IsDraggable() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		c.SetDragged(true)
	}

	// drag the stack of cards, and set the cards' IsActive state to true
	if inpututil.MouseButtonPressDuration(ebiten.MouseButtonLeft) > 3 && c.Dragged() {
		res.DraggedCard = c
		c.IsActive = true
		c.X = cx - c.W/2
		c.Y = cy - c.H/2

		// if not the last card then height is 35, and c.Y is 35/2
		if len(cards[ci:]) > 1 && c != cards[len(cards[ci:])-1] {
			c.H = CardsVSpacer
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
				card.H = CardsVSpacer
			}
			card.X = cx - card.W/2
			card.Y = cy + ((i + 1) * CardsVSpacer)
			card.IsActive = true
			opc := &ebiten.DrawImageOptions{}
			opc.GeoM.Scale(card.ScX, card.ScY)
			opc.GeoM.Translate(float64(card.X), float64(card.Y))
			screen.DrawImage(card.Img, opc)
		}
	}

	// upon release, set dragged and active states to false
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		c.SetDragged(false)
		c.IsActive = false
		res.DraggedCard = nil
	}

	// draw the cards below the Dragged Card
	if !c.IsActive {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(c.ScX, c.ScY)
		op.GeoM.Translate(float64(c.X), float64(c.Y))
		screen.DrawImage(c.Img, op)
	}
}

func (c *Card) Dragged() bool {
	return c.DraggedState
}

func (c *Card) SetDragged(state bool) {
	c.DraggedState = state
}

func (c *Card) IsDraggable() bool {
	return c.DraggableState
}

func (c *Card) SetDraggableState(state bool) {
	c.DraggableState = state
}

// Hovered - Returns true if the card is hovered
func (c *Card) Hovered(cx, cy int) bool {
	return image.Pt(cx, cy).In(c.HitBox())
}

func (c *Card) HitBox() image.Rectangle {
	return image.Rect(c.X, c.Y, c.X+c.W, c.Y+c.H)
}
