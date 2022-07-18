package main

import (
	res "games-ebiten/resources"
	r "games-ebiten/router"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"math/rand"
	"time"
)

// init - applies MPlus1pRegular_ttf font
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

	if err := ebiten.RunGame(r.NewRouter()); err != nil {
		log.Fatal(err)
	}
}
