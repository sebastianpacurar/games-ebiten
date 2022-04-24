package game

import (
	u "games-ebiten/resources/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	_ "image/png"
)

// Translation - used for acquiring the right card index while getting the card SubImage from the Image
// HoveredCard - used to force the hovered card to overlap other images, while dragged
// CardRanks - smallest is "Ace"(0), while highest is "King"(13)
var (
	HoveredCard interface{}
	Translation = map[string]map[int]string{
		u.PixelatedTheme: {
			0: "2", 1: "3", 2: "4", 3: "5", 4: "6", 5: "7",
			6: "8", 7: "9", 8: "10", 9: "J", 10: "Q", 11: "K", 12: "A",
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
// ColCount - refers to what column the card is currently in if it's in any column
type Card struct {
	Img        *ebiten.Image
	BackImg    *ebiten.Image
	ColCount   int
	Suit       string
	Value      int
	Color      string
	X, Y, W, H float64
	ScaleX     float64
	ScaleY     float64
	IsRevealed bool
	IsDragged  bool
}

func (c *Card) GetPosition() (float64, float64, float64, float64) {
	return c.X, c.Y, c.W, c.H
}

func (c *Card) SetLocation(axis string, val float64) {
	if axis == u.X {
		c.X = val
	} else if axis == u.Y {
		c.Y = val
	}
}

func (c *Card) DrawCardSprite(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(c.ScaleX, c.ScaleY)
	img := c.Img

	if !c.IsRevealed {
		img = c.BackImg
	}

	// only revealed cards can be hovered
	if c.IsRevealed {
		c.DragAndDropCard()
		// release
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			c.IsDragged = false
			HoveredCard = nil
		}

		// used to force to draw the card over the other images while being dragged
		if c.IsDragged {
			HoveredCard = c
		}
	}
	op.GeoM.Translate(c.X, c.Y)
	screen.DrawImage(img, op)
}

func (c *Card) GetDraggedState() bool {
	return c.IsDragged
}

func (c *Card) SetDraggedState(state bool) {
	c.IsDragged = state
}

// DragAndDropCard - Drops the card on the valid card from a valid slot
func (c *Card) DragAndDropCard() {
	cx, cy := ebiten.CursorPosition()

	if int(c.X) <= cx && cx < int(c.X+c.W) && int(c.Y) <= cy && cy < int(c.Y+c.H) && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		c.IsDragged = true
	}

	// drag and set location
	if inpututil.MouseButtonPressDuration(ebiten.MouseButtonLeft) > 0 && c.GetDraggedState() {
		_, _, w, h := c.GetPosition()
		c.X = float64(cx) - w/2
		c.Y = float64(cy) - h/2
		HoveredCard = c
	}
}

// IsCardHovered - Returns true if the image is hovered
func (c *Card) IsCardHovered(cx, cy int) bool {
	x, y, w, h := c.GetPosition()
	return int(x) <= cx && cx < int(x+w) && int(y) <= cy && cy < int(y+h)
}
