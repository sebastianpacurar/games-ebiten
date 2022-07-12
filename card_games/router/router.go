package router

import (
	"games-ebiten/card_games/solitaire/free_cell"
	"games-ebiten/card_games/solitaire/klondike"
	u "games-ebiten/resources/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
)

// Router - responsible with the game states. used to navigate between games
// Active - refers to the current game
type Router struct {
	Active interface{}
	fcGame free_cell.Game
	kGame  klondike.Game
	*Menu
}

func NewRouter() *Router {
	r := &Router{
		Menu:   NewMenu(),
		fcGame: *free_cell.NewGame(),
		kGame:  *klondike.NewGame(),
	}
	r.Active = r.fcGame
	return r
}

func (r *Router) Draw(screen *ebiten.Image) {
	switch r.Active.(type) {
	case free_cell.Game:
		g := r.Active.(free_cell.Game)
		g.Draw(screen)
	case klondike.Game:
		g := r.Active.(klondike.Game)
		g.Draw(screen)
	}

	r.DrawMenu(screen)

	ebitenutil.DebugPrintAt(screen, "Press K for classic version and F for Free Cell", 10, u.ScreenHeight-135)
}

func (r *Router) Update() error {
	var err error
	cx, cy := ebiten.CursorPosition()

	switch {
	case inpututil.IsKeyJustReleased(ebiten.KeyF):
		r.Active = r.fcGame
	case inpututil.IsKeyJustReleased(ebiten.KeyK):
		r.Active = r.kGame
	}

	switch r.Active.(type) {
	case free_cell.Game:
		g := r.Active.(free_cell.Game)
		err = g.Update()
	case klondike.Game:
		g := r.Active.(klondike.Game)
		err = g.Update()
	}

	for _, mi := range r.MenuItems {
		if image.Pt(cx, cy).In(GetItemHitBox(mi)) {
			mi.TxtColor.A = 125
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				mi.IsDropped = !mi.IsDropped
			}
		} else {
			mi.TxtColor.A = 200
		}
	}

	return err
}

func (r *Router) Layout(outsideWidth, outsideHeight int) (int, int) {
	return u.ScreenWidth, u.ScreenHeight
}
