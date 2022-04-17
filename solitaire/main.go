package main

import (
	u "games-ebiten/resources/utils"
	"games-ebiten/solitaire/game"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func main() {
	ebiten.SetWindowSize(u.ScreenWidth, u.ScreenHeight)
	ebiten.SetWindowTitle("solitaire")

	if err := ebiten.RunGame(game.NewGame()); err != nil {
		log.Fatal(err)
	}
}
