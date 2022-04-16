package main

import (
	"games-ebiten/animation_movement/game"
	u "games-ebiten/resources/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func main() {
	ebiten.SetWindowSize(u.ScreenWidth, u.ScreenHeight)
	ebiten.SetWindowTitle("animation_movement")

	// slow down ticks per second by half
	ebiten.SetMaxTPS(30)
	if err := ebiten.RunGame(game.NewGame()); err != nil {
		log.Fatal(err)
	}
}
