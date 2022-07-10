package router

import (
	"games-ebiten/card_games/solitaire/free_cell"
	"games-ebiten/card_games/solitaire/klondike"
	u "games-ebiten/resources/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
)

// Router - responsible with the game states. used to navigate between games
type Router struct {
	active interface{}
	fcGame free_cell.Game
	kGame  klondike.Game
	menu   *Menu
}

// Menu - global menu which overlaps all game screens.
type Menu struct {
	containerImg *ebiten.Image
}

type MenuItem struct{}

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
	m := &Menu{}
	m.containerImg = ebiten.NewImage(u.ScreenWidth, 25)
	m.containerImg.Fill(color.NRGBA{R: 255, G: 255, B: 255, A: 255})

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

func (r *Router) DrawMenu(screen *ebiten.Image) {
	opm := &ebiten.DrawImageOptions{}
	opm.GeoM.Translate(0, 0)
	screen.DrawImage(r.menu.containerImg, opm)

}
