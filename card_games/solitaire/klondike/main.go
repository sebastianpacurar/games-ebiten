package main

import (
	"games-ebiten/card_games/solitaire/klondike/k_game"
	u "games-ebiten/resources/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func main() {
	ebiten.SetWindowSize(u.ScreenWidth, u.ScreenHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Klondike Solitaire")

	if err := ebiten.RunGame(k_game.NewGame()); err != nil {
		log.Fatal(err)
	}
}
