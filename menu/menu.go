package menu

import (
	"games-ebiten/card_games/free_cell"
	res "games-ebiten/resources"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

const (
	Padding     = 10
	ItemPadding = 18
	Border      = 2
)

// Menu - global menu which overlaps all game screens.
// IsActive - refers to the current displayed game
type Menu struct {
	X, Y, W, H     int
	ContainerImage *ebiten.Image
	MenuItems      []*Item
	IsActive       bool
}

func NewMenu() *Menu {
	mis := NewMenuItems()
	res.MainMenuH = mis[0].TxtBounds.Dy() + Padding*2
	m := &Menu{
		X:         0,
		Y:         0,
		W:         res.ScreenWidth,
		H:         res.MainMenuH,
		MenuItems: mis,
	}
	m.ContainerImage = ebiten.NewImage(m.W, m.H)
	m.ContainerImage.Fill(res.White)

	// set menu item Hitbox
	for i, mi := range m.MenuItems {
		if i == 0 {
			mi.X = Padding * 2
		} else {
			mi.X = mi.TxtBounds.Dx() * (i + 1)
		}
		mi.Y = 0
		mi.W = mi.TxtBounds.Dx()
		mi.H = m.H

		mi.TxtColor = res.Black
		mi.Img = ebiten.NewImage(mi.W, mi.H)

		// set option HitBox
		for j, op := range mi.Options {
			op.X = mi.X
			op.Y = mi.H*(j+1) + Border
			op.W = mi.W * 3
			op.H = mi.H

			op.Img = ebiten.NewImage(op.W, op.H)
			op.Color = res.White
			op.TxtColor = res.Black
		}
		firstOpt := mi.Options[0]
		lastOpt := mi.Options[len(mi.Options)-1]

		// set the dropdown box (usually it's black)
		mi.dropArea.Min.X = firstOpt.X - Border
		mi.dropArea.Min.Y = firstOpt.Y - Border
		mi.dropArea.Max.X = lastOpt.X + lastOpt.W + Border
		mi.dropArea.Max.Y = lastOpt.Y + lastOpt.H + Border
	}

	// set default game to Free Cell
	res.ActiveGame = *free_cell.NewGame()

	return m
}

// Draw - draws the menu bar along with the Item elements
func (m *Menu) Draw(screen *ebiten.Image) {
	opm := &ebiten.DrawImageOptions{}
	opm.GeoM.Translate(float64(m.X), float64(m.Y))
	screen.DrawImage(m.ContainerImage, opm)

	for _, mi := range m.MenuItems {
		mi.Draw(screen)
	}
}

func (m *Menu) HitBox() image.Rectangle {
	return image.Rect(m.X, m.Y, m.X+m.W, m.Y+m.H)
}
