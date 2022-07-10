package router

import (
	"games-ebiten/card_games/solitaire/free_cell"
	"games-ebiten/card_games/solitaire/klondike"
	u "games-ebiten/resources/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"image"
	"image/color"
)

// Router - responsible with the game states. used to navigate between games
// active - refers to the current game
type Router struct {
	active interface{}
	fcGame free_cell.Game
	kGame  klondike.Game
	menu   *Menu
}

// Menu - global menu which overlaps all game screens.
type Menu struct {
	X, Y, W, H   int
	containerImg *ebiten.Image
	MenuItems    []*MenuItem
}

// MenuItem - consists of a rectangle formed from the MenuItem's name
type MenuItem struct {
	name    string
	bounds  image.Rectangle
	options []*Option
}

// Option - serves as the selected or available MenuItem option
type Option struct {
	name   string
	active bool
}

// NewMenuItems - construct the layout based on the text's size
func NewMenuItems() []*MenuItem {
	return []*MenuItem{
		{
			name: "Games",
			options: []*Option{
				{name: "Free Cell", active: true},
				{name: "Klondike", active: false},
			},
			bounds: text.BoundString(u.FontFace, "Game"),
		},
		{
			name: "Themes",
			options: []*Option{
				{name: "Classic", active: true},
				{name: "8 bit", active: false},
			},
			bounds: text.BoundString(u.FontFace, "Themes"),
		},
	}
}

func NewRouter() *Router {
	r := &Router{
		menu:   NewMenu(),
		fcGame: *free_cell.NewGame(),
		kGame:  *klondike.NewGame(),
	}
	r.active = r.fcGame
	return r
}

func NewMenu() *Menu {
	menuItems := NewMenuItems()
	m := &Menu{
		X:         menuItems[0].bounds.Min.X,
		Y:         menuItems[0].bounds.Min.Y,
		W:         u.ScreenWidth,
		H:         menuItems[0].bounds.Dy() * 3,
		MenuItems: make([]*MenuItem, 0),
	}
	m.containerImg = ebiten.NewImage(m.W, m.H)
	m.containerImg.Fill(color.NRGBA{R: 255, G: 255, B: 255, A: 255})
	m.MenuItems = menuItems

	return m
}

func (r *Router) Draw(screen *ebiten.Image) {
	switch r.active.(type) {
	case free_cell.Game:
		g := r.active.(free_cell.Game)
		g.Draw(screen)
	case klondike.Game:
		g := r.active.(klondike.Game)
		g.Draw(screen)
	}

	r.DrawMenu(screen)

	ebitenutil.DebugPrintAt(screen, "Press K for classic version and F for Free Cell", 10, u.ScreenHeight-135)
}

func (r *Router) Update() error {
	var err error

	switch {
	case inpututil.IsKeyJustReleased(ebiten.KeyF):
		r.active = r.fcGame
	case inpututil.IsKeyJustReleased(ebiten.KeyK):
		r.active = r.kGame
	}

	switch r.active.(type) {
	case free_cell.Game:
		g := r.active.(free_cell.Game)
		err = g.Update()
	case klondike.Game:
		g := r.active.(klondike.Game)
		err = g.Update()
	}
	return err
}

func (r *Router) Layout(outsideWidth, outsideHeight int) (int, int) {
	return u.ScreenWidth, u.ScreenHeight
}

// DrawMenu - draws the menu bar along with the MenuItem elements
func (r *Router) DrawMenu(screen *ebiten.Image) {
	opm := &ebiten.DrawImageOptions{}
	opm.GeoM.Translate(float64(r.menu.X), float64(r.menu.Y))
	screen.DrawImage(r.menu.containerImg, opm)

	r.DrawMenuItems(screen)
}

func (r *Router) DrawMenuItems(screen *ebiten.Image) {
	pos := 0
	for _, mi := range r.menu.MenuItems {
		text.Draw(screen, mi.name, u.FontFace, 10+pos, mi.bounds.Dy()+6, color.NRGBA{R: 0, G: 0, B: 0, A: 255})
		pos += mi.bounds.Dx() + 20
	}
}
