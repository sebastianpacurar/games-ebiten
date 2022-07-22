package menu

import (
	"games-ebiten/animation_movement"
	"games-ebiten/card_games/free_cell"
	"games-ebiten/card_games/klondike"
	"games-ebiten/match_pairs"
	res "games-ebiten/resources"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"image"
	"image/color"
)

// Section - consists of a rectangle formed from the Section's Name.
// every Section has a Header. but may or may not have Items
type Section struct {
	IsCurrent bool
	Header    *Header
	Items     []*Item
	DropArea  image.Rectangle
}

// Header - is the actual text Rect{} of the section.
// TxtBounds - represents the Bounds() returned by text.BoundString.
// IsDropped - applies only to droppable headers.
type Header struct {
	X, Y, W, H  int
	Id          int
	Name        string
	IsDropped   bool
	IsDroppable bool
	Img         *ebiten.Image
	TxtBounds   image.Rectangle
	TxtColor    color.NRGBA
}

// NewSections - returns the menu from left to right:
// the first item refers to all games and appears on all games,
// the following ones are game-specifics
func NewSections() []*Section {
	mis := []*Section{
		{
			IsCurrent: true,
			Header: &Header{
				Name:        "Games",
				TxtBounds:   text.BoundString(res.FontFace, "Games"),
				IsDroppable: true,
			},
			Items: NewMainMenuItems(),
		},

		{
			IsCurrent: true,
			Header: &Header{
				Name:        "New Game",
				TxtBounds:   text.BoundString(res.FontFace, "New Game"),
				IsDroppable: false,
			},
			Items: nil,
		},

		// KLONDIKE and FREE CELL
		{
			IsCurrent: true,
			Header: &Header{
				Name:        "Themes",
				TxtBounds:   text.BoundString(res.FontFace, "Themes"),
				IsDroppable: true,
			},
			Items: NewCardsThemeItems(),
		},
	}

	return mis
	//{
	//	Name: "Themes",
	//	Items: []*Item{
	//		{Name: "Classic", IsSelected: true, TxtBounds: text.BoundString(res.FontFace, "Classic")},
	//		{Name: "8 bit", IsSelected: false, TxtBounds: text.BoundString(res.FontFace, "8 bit")},
	//	},
	//	TxtBounds: text.BoundString(res.FontFace, "Themes"),
	//},
}

func (s *Section) Draw(screen *ebiten.Image) {
	// draw the text for every menu item. Note that Y = text's actual height (plus Padding)
	opc := &ebiten.DrawImageOptions{}
	opc.GeoM.Translate(float64(s.Header.X), float64(s.Header.Y))
	screen.DrawImage(s.Header.Img, opc)

	text.Draw(screen, s.Header.Name, res.FontFace, s.Header.X, s.Header.TxtBounds.Dy()+Padding, s.Header.TxtColor)

	// draw the Section dropdown
	if s.Header.IsDropped {
		x, y := float64(s.DropArea.Min.X), float64(s.DropArea.Min.Y)
		border := ebiten.NewImage(s.DropArea.Dx(), s.DropArea.Dy())

		// draw border (it's a black image behind the dropdown)
		// pls refer to DropArea for Rect{} properties
		opb := &ebiten.DrawImageOptions{}
		opb.GeoM.Translate(x, y)
		border.Fill(res.Black)
		screen.DrawImage(border, opb)

		for i, item := range s.Items {
			item.Draw(i, screen)
		}
	}
}

func (s *Section) SwitchToGame(gameId int) {
	// set current option to active true, and others to false
	for _, v := range s.Items {
		if v.Id != gameId {
			v.IsSelected = false
		} else {
			v.IsSelected = true
		}
	}

	// pass the game to the interface
	switch gameId {
	case 1:
		res.ActiveGame = *free_cell.NewGame()
	case 2:
		res.ActiveGame = *klondike.NewGame()
	case 3:
		res.ActiveGame = *match_pairs.NewGame()
	case 4:
		res.ActiveGame = *animation_movement.NewGame()
	}
}

func (s *Section) SwitchToCardTheme(themeId int) {
	for _, v := range s.Items {
		if v.Id != themeId {
			v.IsSelected = false
		} else {
			v.IsSelected = true
		}
	}

	switch themeId {
	case 1:
		res.ActiveCardsTheme = res.ClassicTheme
	case 2:
		res.ActiveCardsTheme = res.PixelatedTheme
	}

	switch res.ActiveGame.(type) {
	case free_cell.Game:
		g := res.ActiveGame.(free_cell.Game)
		g.BuildDeck()
	case klondike.Game:
		g := res.ActiveGame.(klondike.Game)
		g.BuildDeck()
	}
}

// FormatWidthsBasedOnTxtSize - sets the width of the opts to the highest one from the list
func (s *Section) FormatWidthsBasedOnTxtSize() {
	largest := 0
	for _, item := range s.Items {
		if largest < item.TxtBounds.Dx() {
			largest = item.TxtBounds.Dx()
		}
	}

	// set the highest width, image, and colors, for every option in the section
	for _, item := range s.Items {
		item.W = largest + Padding*2
		item.Img = ebiten.NewImage(item.W, item.H)
		item.TxtColor = res.Black
		item.Color = res.White
	}
}

// StartXFromLeft - returns the starting position based on the trailing items' text widths, from Left to Right
func (m *Menu) StartXFromLeft(i int) int {
	size := Padding * 2 // first item starts from Padding*2
	for _, s := range m.Sections[:i] {
		size += s.Header.TxtBounds.Dx() + Padding*2
	}
	return size
}

func (h *Header) HitBox() image.Rectangle {
	return image.Rect(h.X, h.Y, h.X+h.W, h.Y+h.H)
}

func (h *Header) Hovered(cx, cy int) bool {
	return image.Pt(cx, cy).In(h.HitBox())
}
