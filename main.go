package main

import (
	"games-ebiten/animation_movement"
	"games-ebiten/card_games/free_cell"
	"games-ebiten/card_games/klondike"
	"games-ebiten/match_pairs"
	"games-ebiten/menu"
	res "games-ebiten/resources"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"log"
	"math/rand"
	"time"
)

// apply MPlus1pRegular font
func init() {
	res.InitFonts()
}

// generate random seed
func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	ebiten.SetWindowSize(res.ScreenWidth, res.ScreenHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Games")

	if err := ebiten.RunGame(NewRouter()); err != nil {
		log.Fatal(err)
	}
}

// Router - responsible with the game states. used to navigate between games
// IsSelected - refers to the current game
type Router struct {
	*menu.Menu
}

func NewRouter() *Router {
	return &Router{
		Menu: menu.NewMenu(),
	}
}

func (r *Router) Draw(screen *ebiten.Image) {
	switch res.ActiveGame.(type) {
	case free_cell.Game:
		g := res.ActiveGame.(free_cell.Game)
		g.Draw(screen)
	case klondike.Game:
		g := res.ActiveGame.(klondike.Game)
		g.Draw(screen)
	case match_pairs.Game:
		g := res.ActiveGame.(match_pairs.Game)
		g.Draw(screen)
	case animation_movement.Game:
		g := res.ActiveGame.(animation_movement.Game)
		g.Draw(screen)
	}

	r.Menu.Draw(screen)
}

func (r *Router) Update() error {
	var err error

	switch res.ActiveGame.(type) {
	case free_cell.Game:
		g := res.ActiveGame.(free_cell.Game)
		err = g.Update()
	case klondike.Game:
		g := res.ActiveGame.(klondike.Game)
		err = g.Update()
	case match_pairs.Game:
		g := res.ActiveGame.(match_pairs.Game)
		err = g.Update()
	case animation_movement.Game:
		g := res.ActiveGame.(animation_movement.Game)
		err = g.Update()
	}

	for _, mi := range r.MenuItems {

		// handle click on menu item
		if res.IsAreaHovered(mi) {
			mi.TxtColor.A = 125
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				mi.IsDropped = !mi.IsDropped
			}
		} else {
			mi.TxtColor.A = 255
		}

		// handle change game
		if mi.IsDropped {
			for _, opt := range mi.Options {
				if opt.IsSelected {
					opt.Color = res.Green
					opt.TxtColor = res.Black
					opt.TxtColor.A = 255
				} else {
					opt.TxtColor = res.Black
					opt.Color = res.White
					if res.IsAreaHovered(opt) {
						opt.Color.A = 200
						opt.TxtColor.A = 200
						if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && !opt.IsSelected {
							mi.SwitchToGame(opt.Id)
							mi.IsDropped = false
						}
					} else {
						opt.Color.A = 255
					}
				}
			}
		}
	}

	return err
}

func (r *Router) Layout(outsideWidth, outsideHeight int) (int, int) {
	return res.ScreenWidth, res.ScreenHeight
}
