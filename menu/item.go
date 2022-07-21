package menu

import (
	res "games-ebiten/resources"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"image"
	"image/color"
)

// Item - serves as the selected or available Section Item
type Item struct {
	X, Y, W, H int
	Id         int
	TxtBounds  image.Rectangle
	Img        *ebiten.Image
	Color      color.NRGBA
	TxtColor   color.NRGBA
	Name       string
	IsSelected bool
}

func NewMainMenuItems() []*Item {
	return []*Item{
		{Name: "Free Cell", Id: 1, IsSelected: true, TxtBounds: text.BoundString(res.FontFace, "Free Cell")},
		{Name: "Klondike", Id: 2, IsSelected: false, TxtBounds: text.BoundString(res.FontFace, "Klondike")},
		{Name: "Match Pairs", Id: 3, IsSelected: false, TxtBounds: text.BoundString(res.FontFace, "Match Pairs")},
		{Name: "Anim Move", Id: 4, IsSelected: false, TxtBounds: text.BoundString(res.FontFace, "Anim Move")},
	}
}

func NewCardsThemeItems() []*Item {
	return []*Item{
		{Name: "Classic", Id: 1, IsSelected: true, TxtBounds: text.BoundString(res.FontFace, "Classic")},
		{Name: "8-bit", Id: 2, IsSelected: false, TxtBounds: text.BoundString(res.FontFace, "8-bit")}}
}

func (opt *Item) Draw(i int, screen *ebiten.Image) {
	opo := &ebiten.DrawImageOptions{}
	opo.GeoM.Translate(float64(opt.X), float64(opt.Y))
	opt.Img.Fill(opt.Color)

	txtX := opt.X + Padding
	txtY := (opt.TxtBounds.Dy() + ItemPadding) * (i + 2)

	screen.DrawImage(opt.Img, opo)
	text.Draw(screen, opt.Name, res.FontFace, txtX, txtY, opt.TxtColor)
}

func (opt *Item) HitBox() image.Rectangle {
	return image.Rect(opt.X, opt.Y, opt.X+opt.W, opt.Y+opt.H)
}

func (opt *Item) Hovered(cx, cy int) bool {
	return image.Pt(cx, cy).In(opt.HitBox())
}
