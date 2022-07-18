package router

import (
	res "games-ebiten/resources"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"image"
	"image/color"
)

const (
	Padding     = 10
	ItemPadding = 15
	Border      = 1
)

// Menu - global menu which overlaps all game screens.
type Menu struct {
	X, Y, W, H     int
	ContainerImage *ebiten.Image
	MenuItems      []*MenuItem
}

// MenuItem - consists of a rectangle formed from the MenuItem's Name
type MenuItem struct {
	X, Y, W, H int
	Name       string
	IsDropped  bool
	Img        *ebiten.Image
	TxtBounds  image.Rectangle
	TxtColor   color.NRGBA
	Options    []*Option
	dropArea   image.Rectangle
}

// Option - serves as the selected or available MenuItem option
type Option struct {
	X, Y, W, H int
	TxtBounds  image.Rectangle
	Img        *ebiten.Image
	OptColor   color.NRGBA
	Name       string
	Active     bool
}

func NewMenu() *Menu {
	mis := []*MenuItem{
		{
			Name: "Games",
			Options: []*Option{
				{Name: "Free Cell", Active: true, TxtBounds: text.BoundString(res.FontFace, "Free Cell")},
				{Name: "Klondike", Active: false, TxtBounds: text.BoundString(res.FontFace, "Klondike")},
			},
			TxtBounds: text.BoundString(res.FontFace, "Game"),
		},
		{
			Name: "Themes",
			Options: []*Option{
				{Name: "Classic", Active: true, TxtBounds: text.BoundString(res.FontFace, "Classic")},
				{Name: "8 bit", Active: false, TxtBounds: text.BoundString(res.FontFace, "8 bit")},
			},
			TxtBounds: text.BoundString(res.FontFace, "Themes"),
		},
	}

	m := &Menu{
		X: 0,
		Y: 0,
		W: res.ScreenWidth,
		H: mis[0].TxtBounds.Dy() + Padding*2,
	}
	m.ContainerImage = ebiten.NewImage(m.W, m.H)
	m.ContainerImage.Fill(res.White)

	m.MenuItems = mis

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

		// set item option HitBox
		for j, op := range mi.Options {
			op.X = mi.X
			op.Y = mi.H*(j+1) + Border
			op.W = mi.W * 2
			op.H = mi.H

			mi.dropArea.Min.X = mi.X
			mi.dropArea.Min.Y = mi.H + Border
			mi.dropArea.Max.X = mi.W * 2
			mi.dropArea.Max.Y += op.Y

			op.Img = ebiten.NewImage(op.W, op.H)
			op.Img.Fill(res.White)
		}
	}

	return m
}

// DrawMenu - draws the menu bar along with the MenuItem elements
func (m *Menu) DrawMenu(screen *ebiten.Image) {
	opm := &ebiten.DrawImageOptions{}
	opm.GeoM.Translate(float64(m.X), float64(m.Y))
	screen.DrawImage(m.ContainerImage, opm)

	m.DrawMenuItems(screen)
}

func (m *Menu) DrawMenuItems(screen *ebiten.Image) {

	// draw the text for every menu item. Note that Y = text's actual height (plus Padding)
	for _, mi := range m.MenuItems {

		opc := &ebiten.DrawImageOptions{}
		opc.GeoM.Translate(float64(mi.X), float64(mi.Y))
		screen.DrawImage(mi.Img, opc)

		text.Draw(screen, mi.Name, res.FontFace, mi.X, mi.TxtBounds.Dy()+Padding, mi.TxtColor)

		// TODO: finish, and add borders
		// draw the MenuItem dropdown
		if mi.IsDropped {
			x, y := float64(mi.dropArea.Min.X), float64(mi.dropArea.Min.Y)
			border := ebiten.NewImage(mi.dropArea.Dx(), mi.dropArea.Dy())

			// draw border
			opb := &ebiten.DrawImageOptions{}
			opb.GeoM.Translate(x, y)
			border.Fill(res.Black)
			screen.DrawImage(border, opb)

			for oi, opt := range mi.Options {
				opo := &ebiten.DrawImageOptions{}
				opo.GeoM.Translate(float64(opt.X), float64(opt.Y))
				opt.Img.Fill(res.White)

				txtX := opt.X + Padding
				txtY := (opt.TxtBounds.Dy()+ItemPadding)*(oi+2) + Border

				screen.DrawImage(opt.Img, opo)
				text.Draw(screen, opt.Name, res.FontFace, txtX, txtY, res.Black)
			}
		}
	}
}

func (m *Menu) GetHitBox() image.Rectangle {
	return image.Rect(m.X, m.Y, m.X+m.W, m.Y+m.H)
}

func GetItemHitBox(mi *MenuItem) image.Rectangle {
	return image.Rect(mi.X, mi.Y, mi.X+mi.W, mi.Y+mi.H)
}

func GetOptionHitBox(opt *Option) image.Rectangle {
	return image.Rect(opt.X, opt.Y, opt.X+opt.W, opt.Y+opt.H)
}
