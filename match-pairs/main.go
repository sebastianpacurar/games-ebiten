package main

import (
	"games-ebiten/match-pairs/game"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func main() {
	ebiten.SetWindowSize(1800, 860)
	ebiten.SetWindowTitle("match pairs")

	ebiten.SetMaxTPS(30)
	if err := ebiten.RunGame(game.NewGame()); err != nil {
		log.Fatal(err)
	}
}
