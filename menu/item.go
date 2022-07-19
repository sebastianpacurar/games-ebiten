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

// Item - consists of a rectangle formed from the Item's Name
type Item struct {
	X, Y, W, H int
	Name       string
	IsDropped  bool
	Img        *ebiten.Image
	TxtBounds  image.Rectangle
	TxtColor   color.NRGBA
	Options    []*Option
	dropArea   image.Rectangle
}

func NewMenuItems() []*Item {
	return []*Item{
		{
			Name:      "Games",
			Options:   NewOptions(),
			TxtBounds: text.BoundString(res.FontFace, "Games"),
		},
		//{
		//	Name: "Themes",
		//	Options: []*Option{
		//		{Name: "Classic", IsSelected: true, TxtBounds: text.BoundString(res.FontFace, "Classic")},
		//		{Name: "8 bit", IsSelected: false, TxtBounds: text.BoundString(res.FontFace, "8 bit")},
		//	},
		//	TxtBounds: text.BoundString(res.FontFace, "Themes"),
		//},
	}
}

func (mi *Item) Draw(screen *ebiten.Image) {
	// draw the text for every menu item. Note that Y = text's actual height (plus Padding)
	opc := &ebiten.DrawImageOptions{}
	opc.GeoM.Translate(float64(mi.X), float64(mi.Y))
	screen.DrawImage(mi.Img, opc)

	text.Draw(screen, mi.Name, res.FontFace, mi.X, mi.TxtBounds.Dy()+Padding, mi.TxtColor)

	// draw the Item dropdown
	if mi.IsDropped {
		x, y := float64(mi.dropArea.Min.X), float64(mi.dropArea.Min.Y)
		border := ebiten.NewImage(mi.dropArea.Dx(), mi.dropArea.Dy())

		// draw border (it's a black image behind the dropdown)
		// pls refer to dropArea for Rect{} properties
		opb := &ebiten.DrawImageOptions{}
		opb.GeoM.Translate(x, y)
		border.Fill(res.Black)
		screen.DrawImage(border, opb)

		for i, opt := range mi.Options {
			opt.Draw(i, screen)
		}
	}
}

func (mi *Item) SwitchToGame(gameId int) {
	// set current option to active true, and others to false
	for _, v := range mi.Options {
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

func (mi *Item) HitBox() image.Rectangle {
	return image.Rect(mi.X, mi.Y, mi.X+mi.W, mi.Y+mi.H)
}

func (mi *Item) Hovered(cx, cy int) bool {
	return image.Pt(cx, cy).In(mi.HitBox())
}
