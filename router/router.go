package router

import (
	"games-ebiten/animation_movement"
	"games-ebiten/card_games/free_cell"
	"games-ebiten/card_games/klondike"
	"games-ebiten/match-pairs"
	res "games-ebiten/resources"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
)

// Router - responsible with the game states. used to navigate between games
// Active - refers to the current game
type Router struct {
	Active       interface{}
	freeCell     free_cell.Game
	klondike     klondike.Game
	matchPairs   match_pairs.Game
	animMovement animation_movement.Game
	*Menu
}

func NewRouter() *Router {
	r := &Router{
		Menu:         NewMenu(),
		freeCell:     *free_cell.NewGame(),
		klondike:     *klondike.NewGame(),
		matchPairs:   *match_pairs.NewGame(),
		animMovement: *animation_movement.NewGame(),
	}
	r.Active = r.freeCell
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
	case match_pairs.Game:
		g := r.Active.(match_pairs.Game)
		g.Draw(screen)
	case animation_movement.Game:
		g := r.Active.(animation_movement.Game)
		g.Draw(screen)
	}

	r.DrawMenu(screen)
}

func (r *Router) Update() error {
	var err error
	cx, cy := ebiten.CursorPosition()

	switch {
	case inpututil.IsKeyJustReleased(ebiten.Key1):
		r.Active = r.freeCell
	case inpututil.IsKeyJustReleased(ebiten.Key2):
		r.Active = r.klondike
	case inpututil.IsKeyJustReleased(ebiten.Key3):
		r.Active = r.matchPairs
	case inpututil.IsKeyJustReleased(ebiten.Key4):
		r.Active = r.animMovement
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
