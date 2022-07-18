package router

import (
	"games-ebiten/animation_movement"
	"games-ebiten/card_games/free_cell"
	"games-ebiten/card_games/klondike"
	"games-ebiten/match-pairs"
	res "games-ebiten/resources"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
)

// Router - responsible with the game states. used to navigate between games
// Active - refers to the current game
type Router struct {
	Active interface{}
	*Menu
}

func NewRouter() *Router {
	return &Router{
		Menu:   NewMenu(),
		Active: *free_cell.NewGame(),
	}
}

func (r *Router) Draw(screen *ebiten.Image) {
	switch r.Active.(type) {
	case free_cell.Game:
		g := r.Active.(free_cell.Game)
		g.Draw(screen)
	case klondike.Game:
		g := r.Active.(klondike.Game)
		g.Draw(screen)
	case match_pairs.Game:
		ebiten.SetMaxTPS(60)
		g := r.Active.(match_pairs.Game)
		g.Draw(screen)
	case animation_movement.Game:
		g := r.Active.(animation_movement.Game)
		g.Draw(screen)
	}

	r.DrawMenu(screen)
	ebitenutil.DebugPrintAt(screen, "Press 1 - Free Cell", 10, res.ScreenHeight-125)
	ebitenutil.DebugPrintAt(screen, "Press 2 - Klondike Solitaire", 10, res.ScreenHeight-95)
	ebitenutil.DebugPrintAt(screen, "Press 3 - Match Pairs (currently broken)", 10, res.ScreenHeight-65)
	ebitenutil.DebugPrintAt(screen, "Press 4 - Animation Movement", 10, res.ScreenHeight-35)
}

func (r *Router) Update() error {
	var err error
	cx, cy := ebiten.CursorPosition()

	switch {
	case inpututil.IsKeyJustReleased(ebiten.Key1):
		ebiten.SetMaxTPS(60)
		r.Active = *free_cell.NewGame()
	case inpututil.IsKeyJustReleased(ebiten.Key2):
		ebiten.SetMaxTPS(60)
		r.Active = *klondike.NewGame()
	case inpututil.IsKeyJustReleased(ebiten.Key3):
		ebiten.SetMaxTPS(60)
		r.Active = *match_pairs.NewGame()
	case inpututil.IsKeyJustReleased(ebiten.Key4):
		ebiten.SetMaxTPS(40)
		r.Active = *animation_movement.NewGame()
	}

	switch r.Active.(type) {
	case free_cell.Game:
		g := r.Active.(free_cell.Game)
		err = g.Update()
	case klondike.Game:
		g := r.Active.(klondike.Game)
		err = g.Update()
	case match_pairs.Game:
		g := r.Active.(match_pairs.Game)
		err = g.Update()
	case animation_movement.Game:
		g := r.Active.(animation_movement.Game)
		err = g.Update()
	}

	for _, mi := range r.MenuItems {

		// handle click on menu item
		if image.Pt(cx, cy).In(GetItemHitBox(mi)) {
			mi.TxtColor.A = 125
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				mi.IsDropped = !mi.IsDropped
			}
		} else {
			mi.TxtColor.A = 200
		}

		//if mi.IsDropped {
		//	for _, opt := range mi.Options {
		//		if
		//	}
		//}
	}

	return err
}

func (r *Router) Layout(outsideWidth, outsideHeight int) (int, int) {
	return res.ScreenWidth, res.ScreenHeight
}
