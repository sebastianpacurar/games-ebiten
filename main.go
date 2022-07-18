package main

import (
	u "games-ebiten/resources"
	r "games-ebiten/router"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"math/rand"
	"time"
)

// init - applies MPlus1pRegular_ttf font
func init() {
	u.InitFonts()
}

// generate random seed
func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	ebiten.SetWindowSize(u.ScreenWidth, u.ScreenHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Card Games")

	if err := ebiten.RunGame(r.NewRouter()); err != nil {
		log.Fatal(err)
	}
}
