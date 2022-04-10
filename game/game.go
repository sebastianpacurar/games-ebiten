package game

import (
	"game_ebiten/game/data"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	_ "image/png"
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

// Bounds - represents the main screen edges
var Bounds = map[string]float64{MinX: 0, MaxX: ScreenWidth, MinY: 0, MaxY: ScreenHeight}

type Game struct {
	data.Player
	data.Food
}

func NewGame() *Game {
	playerLocX := float64(ScreenWidth)/2 - data.PlayerFrameWidth/2
	playerLocY := float64(ScreenHeight)/2 - data.PlayerFrameHeight/2

	foodLocX, foodLocY := GenerateRandomLocation(Bounds[MinX], Bounds[MaxX]-data.FoodFrameWidth, Bounds[MinY], Bounds[MaxY]-data.FoodFrameHeight)
	foodImg := LoadSpriteImage("spritesheets/food-drink.png")
	playerImg := LoadSpriteImage("spritesheets/character.png")

	return &Game{
		Player: data.Player{
			Img:       ebiten.NewImageFromImage(playerImg),
			W:         float64(data.PlayerFrameWidth * data.PlayerScale),
			H:         float64(data.PlayerFrameHeight * data.PlayerScale),
			Speed:     3,
			Direction: 0,
			LocX:      playerLocX,
			LocY:      playerLocY,
		},
		Food: data.Food{
			Img:  ebiten.NewImageFromImage(foodImg),
			W:    data.FoodFrameWidth * data.FoodScale,
			H:    data.FoodFrameWidth * data.FoodScale,
			LocX: foodLocX,
			LocY: foodLocY,
		},
	}
}

func (g *Game) Update() error {
	if !g.Food.IsDisplayed {
		g.Food.LocX, g.Food.LocY = GenerateRandomLocation(Bounds[MinX], Bounds[MaxX]-data.FoodFrameWidth, Bounds[MinY], Bounds[MaxY]-data.FoodFrameHeight)
		g.Food.IsDisplayed = true
	}

	g.Player.HitBox = HitBox(g.Player.LocX, g.Player.LocY, g.Player.W, g.Player.H)
	g.Food.HitBox = HitBox(g.Food.LocX, g.Food.LocY, g.Food.W, g.Food.H)

	if ebiten.IsKeyPressed(ebiten.KeyShiftLeft) {
		g.Speed = 5
	} else if !ebiten.IsKeyPressed(ebiten.KeyShiftLeft) {
		g.Speed = 3
	}

	g.HandleMovement(Bounds[MinX], Bounds[MaxX], Bounds[MinY], Bounds[MaxY])

	// verify if player and food shape areas overlap
	if isCollision(g.Player.HitBox[MinX], g.Player.HitBox[MinY], g.Player.W, g.Player.H, g.Food.HitBox[MinX], g.Food.HitBox[MinY], g.Food.W, g.Food.H) {
		g.Food.IsDisplayed = false
		g.Food.ImgNo++
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)

	g.Player.DrawImage(screen)

	// draw food
	if g.Food.IsDisplayed {
		g.Food.DrawImage(screen)
	}

	ebitenutil.DebugPrintAt(screen, "W A S D to move", 0, 0)
	ebitenutil.DebugPrintAt(screen, "LEFT SHIFT to speed up", 0, 25)
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
