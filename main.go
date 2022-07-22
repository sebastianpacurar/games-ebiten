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
	"image"
	"log"
	"math/rand"
	"time"
)

// apply initial setup
func init() {
	res.InitFonts()
	rand.Seed(time.Now().UnixNano())
	res.ActiveCardsTheme = res.ClassicTheme
}

func main() {
	ebiten.SetWindowSize(res.ScreenWidth, res.ScreenHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Games")

	if err := ebiten.RunGame(NewRouter()); err != nil {
		log.Fatal(err)
	}
}

// Router - manages the game states - used to navigate between games.
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

	for _, s := range r.Menu.Sections {

		// handle click on menu top item.
		if res.IsAreaHovered(s.Header) {
			s.Header.TxtColor.A = 125
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				if s.Header.IsDroppable {
					// handle trigger dropdown functionality
					s.Header.IsDropped = !s.Header.IsDropped
				} else {
					// handle New Game functionality
					if s.Header.Name == "New Game" {
						HandleNewGameClick(res.ActiveGame)
					}
				}
			}
		} else {
			s.Header.TxtColor.A = 255
		}

		// handle change game
		if s.Header.IsDropped {
			firstRect := image.Rect(s.Header.X, s.Header.Y, s.Header.X+s.Header.W, s.DropArea.Max.Y)
			secondRect := s.DropArea
			if image.Pt(ebiten.CursorPosition()).In(firstRect) || image.Pt(ebiten.CursorPosition()).In(secondRect) {
				for _, item := range s.Items {
					if item.IsSelected {
						item.Color = res.Green
						item.TxtColor = res.Black
						item.TxtColor.A = 255
					} else {
						item.TxtColor = res.Black
						item.Color = res.White
						if res.IsAreaHovered(item) {
							item.Color.A = 200
							item.TxtColor.A = 200
							if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && !item.IsSelected {
								if s.Header.Name == "Games" {
									s.SwitchToGame(item.Id)
								} else if s.Header.Name == "Themes" {
									s.SwitchToCardTheme(item.Id)
								}
								s.Header.IsDropped = false
							}
						} else {
							item.Color.A = 255
						}
					}
				}
			} else {
				s.Header.IsDropped = false
			}
		}
	}

	return err
}

func (r *Router) Layout(outsideWidth, outsideHeight int) (int, int) {
	return res.ScreenWidth, res.ScreenHeight
}

func HandleNewGameClick(i interface{}) {
	switch i.(type) {
	case free_cell.Game:
		game := i.(free_cell.Game)
		game.BuildDeck()
	case klondike.Game:
		game := i.(klondike.Game)
		game.BuildDeck()
	case match_pairs.Game:
		game := i.(match_pairs.Game)
		game.Reset()
	}
}
