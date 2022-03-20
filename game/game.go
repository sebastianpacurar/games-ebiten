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
	ImgSize      = float64(30)
	ScreenWidth  = 1280
	ScreenHeight = 960
	MinX         = "minX"
	MinY         = "minY"
	MaxX         = "maxX"
	MaxY         = "maxY"
)

var (
	Bounds = map[string]float64{MinX: 0, MaxX: ScreenWidth, MinY: 0, MaxY: ScreenHeight}
)

type Game struct {
	data.Player
	data.Food
}

func NewGame() *Game {
	playerLocX := float64(ScreenWidth)/2 - ImgSize/2
	playerLocY := float64(ScreenHeight)/2 - ImgSize/2

	foodLocX, foodLocY := GenerateRandomLocation(Bounds[MinX], Bounds[MaxX], Bounds[MinY], Bounds[MaxY])

	return &Game{
		Player: data.Player{
			Size:   ImgSize,
			Speed:  1,
			Breaks: 1,
			LocX:   playerLocX,
			LocY:   playerLocY,
		},
		Food: data.Food{
			Size:        10,
			LocX:        foodLocX,
			LocY:        foodLocY,
			IsDisplayed: true,
		},
	}
}

func (g *Game) Update() error {
	g.Player.HitBox = HitBox(g.Player.LocX, g.Player.LocY, g.Player.Size, g.Player.Size)
	g.Food.HitBox = HitBox(g.Food.LocX, g.Food.LocY, g.Food.Size, g.Food.Size)

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

	// verify if player and food shape areas overlap
	if isCollision(g.Player.HitBox[MinX], g.Player.HitBox[MinY], g.Player.Size, g.Player.Size,
		g.Food.HitBox[MinX], g.Food.HitBox[MinY], g.Food.Size, g.Food.Size) {
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
	opPlayer.GeoM.Translate(g.Player.LocX, g.Player.LocY)
	screen.DrawImage(playerRect, opPlayer)

	// draw food
	if g.Food.IsDisplayed {
		foodRect := ebiten.NewImage(int(g.Food.Size), int(g.Food.Size))
		foodRect.Fill(color.NRGBA{G: 255, A: 255})
		opFood := &ebiten.DrawImageOptions{}
		opFood.GeoM.Translate(g.Food.LocX, g.Food.LocY)
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

// isCollision - returns true if images overlap in any way
func isCollision(x1, y1, w1, h1, x2, y2, w2, h2 float64) bool {
	return x1 < x2+w2 && x1+w1 > x2 && y1 < y2+h2 && h1+y1 > y2
}

// HitBox - generate the shape's hitbox (minX, maxX, minY, maxY)
func HitBox(x, y, w, h float64) map[string]float64 {
	return map[string]float64{
		"minX": x,
		"maxX": x + w,
		"minY": y,
		"maxY": y + h,
	}
}
