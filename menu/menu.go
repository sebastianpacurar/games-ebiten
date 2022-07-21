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
type Menu struct {
	X, Y, W, H     int
	ContainerImage *ebiten.Image
	Sections       []*Section
}

func NewMenu() *Menu {
	sects := NewSections()
	res.MainMenuH = sects[0].Header.TxtBounds.Dy() + Padding*2
	m := &Menu{
		X:        0,
		Y:        0,
		W:        res.ScreenWidth,
		H:        res.MainMenuH,
		Sections: sects,
	}
	m.ContainerImage = ebiten.NewImage(m.W, m.H)
	m.ContainerImage.Fill(res.White)

	// set menu top item Hitbox
	for i, s := range m.Sections {
		if i == 0 {
			s.Header.X = Padding * 2
		} else {
			s.Header.X = m.StartXFrom(i)
		}
		s.Header.Y = 0
		s.Header.W = s.Header.TxtBounds.Dx()
		s.Header.H = m.H

		s.Header.TxtColor = res.Black
		s.Header.Img = ebiten.NewImage(s.Header.W, s.Header.H)

		if len(s.Items) > 0 {
			// set option HitBox
			for j, item := range s.Items {
				item.X = s.Header.X
				item.Y = s.Header.H*(j+1) + Border
				item.H = s.Header.H
			}
			s.FormatOptsWidth()

			firstOpt := s.Items[0]
			lastOpt := s.Items[len(s.Items)-1]

			// set the dropdown box (usually it's black)
			s.DropArea = image.Rect(firstOpt.X, firstOpt.Y, lastOpt.X+lastOpt.W, lastOpt.Y+lastOpt.H).Inset(-Border)
		}
	}

	// set default game to Free Cell
	res.ActiveGame = *free_cell.NewGame()

	return m
}

// Draw - draws the menu bar along with the Section elements
func (m *Menu) Draw(screen *ebiten.Image) {
	opm := &ebiten.DrawImageOptions{}
	opm.GeoM.Translate(float64(m.X), float64(m.Y))
	screen.DrawImage(m.ContainerImage, opm)

	for _, s := range m.Sections {
		s.Draw(screen)
	}
}

func (m *Menu) HitBox() image.Rectangle {
	return image.Rect(m.X, m.Y, m.X+m.W, m.Y+m.H)
}
