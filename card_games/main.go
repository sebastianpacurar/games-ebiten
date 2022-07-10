package main

import (
	c "games-ebiten/card_games/router"
	u "games-ebiten/resources/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func main() {
	ebiten.SetWindowSize(u.ScreenWidth, u.ScreenHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("FreeCell Solitaire")

	if err := ebiten.RunGame(c.NewRouter()); err != nil {
		log.Fatal(err)
	}
}
