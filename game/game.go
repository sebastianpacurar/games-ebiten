package game

import (
	"fmt"
	"game_ebiten/game/data"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	"math/rand"
	"time"
)

const (
	ScreenWidth  = 1280
	ScreenHeight = 960
	MinX         = "minX"
	MinY         = "minY"
	MaxX         = "maxX"
	MaxY         = "maxY"
)

var Bounds = map[string]float64{MinX: 0, MaxX: ScreenWidth, MinY: 0, MaxY: ScreenHeight}

type Game struct {
	data.Player
	data.Food
}

func NewGame() *Game {
	playerSize := float64(25)
	playerOrigin := playerSize / 2
	playerPosX := float64(ScreenWidth)/2 - playerOrigin
	playerPosY := float64(ScreenHeight)/2 - playerOrigin

	foodSize := float64(10)
	foodX, foodY := GenerateRandomLocation(Bounds[MinX], Bounds[MaxX], Bounds[MinY], Bounds[MaxY])

	return &Game{
		Player: data.Player{
			Size:   playerSize,
			Origin: playerOrigin,
			Speed:  1,
			Breaks: 1,
			PosX:   playerPosX,
			PosY:   playerPosY,
		},
		Food: data.Food{
			Size:        foodSize,
			PosX:        foodX,
			PosY:        foodY,
			IsDisplayed: true,
		},
	}
}

func (g *Game) Update() error {
	g.Player.HitBox = g.Player.GetHitBox()
	g.Food.HitBox = g.Food.GetHitBox()

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.Speed = 5
	} else if !ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.Speed = 1
	}

	if ebiten.IsKeyPressed(ebiten.KeyShiftLeft) {
		g.Breaks = 2.5
	} else if !ebiten.IsKeyPressed(ebiten.KeyShiftLeft) {
		g.Breaks = 1
	}

	g.HandleMovement(Bounds[MinX], Bounds[MaxX], Bounds[MinY], Bounds[MaxY])

	// verify if player and food hitboxes overlap
	// TODO: not currently working
	if (g.Player.HitBox[MinX] >= g.Food.HitBox[MinX] && g.Player.HitBox[MaxX] <= g.Food.HitBox[MaxX]) ||
		(g.Player.HitBox[MinY] >= g.Food.HitBox[MinY] && g.Player.HitBox[MaxY] <= g.Food.HitBox[MaxY]) {
		fmt.Println("test")
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	// draw player
	playerRect := ebiten.NewImage(int(g.Player.Size), int(g.Player.Size))
	playerRect.Fill(color.NRGBA{R: 255, G: 255, B: 255, A: 175})
	opPlayer := &ebiten.DrawImageOptions{}
	opPlayer.GeoM.Translate(g.Player.PosX, g.Player.PosY)
	screen.DrawImage(playerRect, opPlayer)

	// draw food
	if g.Food.IsDisplayed {
		foodRect := ebiten.NewImage(int(g.Food.Size), int(g.Food.Size))
		foodRect.Fill(color.NRGBA{G: 255, A: 255})
		opFood := &ebiten.DrawImageOptions{}
		opFood.GeoM.Translate(g.Food.PosX, g.Food.PosY)
		screen.DrawImage(foodRect, opFood)
	}

	ebitenutil.DebugPrintAt(screen, "W A S D to move", 0, 0)
	ebitenutil.DebugPrintAt(screen, "SPACE BAR to speed up", 0, 25)
	ebitenutil.DebugPrintAt(screen, "LEFT SHIFT to hit the breaks", 0, 50)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

// GenerateRandomLocation - generate a random location for X and y
// (x is between 0 and ScreenWidth) and (y is between 0 and ScreenHeight)
func GenerateRandomLocation(minX, maxX, minY, maxY float64) (float64, float64) {
	rand.Seed(time.Now().UnixNano())
	return (rand.Float64() * (maxX - minX)) + minX, (rand.Float64() * (maxY - minY)) + minY
}
