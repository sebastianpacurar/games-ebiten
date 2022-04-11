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
	ScreenWidth  = 1400
	ScreenHeight = 860
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
	Friends []data.NPC
	Enemies []data.NPC
}

func NewGame() *Game {
	playerLocX := float64(ScreenWidth)/2 - data.PlayerFrameWidth/2
	playerLocY := float64(ScreenHeight)/2 - data.PlayerFrameHeight/2

	friend1LocX, friend1LocY := GenerateRandomLocation(Bounds[MinX], Bounds[MaxX]-data.NPCFrameWidth, Bounds[MinY], Bounds[MaxY]-data.NPCFrameHeight)
	friend2LocX, friend2LocY := GenerateRandomLocation(Bounds[MinX], Bounds[MaxX]-data.NPCFrameWidth, Bounds[MinY], Bounds[MaxY]-data.NPCFrameHeight)
	friend3LocX, friend3LocY := GenerateRandomLocation(Bounds[MinX], Bounds[MaxX]-data.NPCFrameWidth, Bounds[MinY], Bounds[MaxY]-data.NPCFrameHeight)

	enemy1LocX, enemy1LocY := GenerateRandomLocation(Bounds[MinX], Bounds[MaxX]-data.NPCFrameWidth, Bounds[MinY], Bounds[MaxY]-data.NPCFrameHeight)
	enemy2LocX, enemy2LocY := GenerateRandomLocation(Bounds[MinX], Bounds[MaxX]-data.NPCFrameWidth, Bounds[MinY], Bounds[MaxY]-data.NPCFrameHeight)

	foodLocX, foodLocY := GenerateRandomLocation(Bounds[MinX], Bounds[MaxX]-data.FoodFrameWidth, Bounds[MinY], Bounds[MaxY]-data.FoodFrameHeight)

	return &Game{
		Player: data.Player{
			Img:       ebiten.NewImageFromImage(LoadSpriteImage("spritesheets/character1.png")),
			W:         float64(data.PlayerFrameWidth * data.PlayerScale),
			H:         float64(data.PlayerFrameHeight * data.PlayerScale),
			Speed:     3,
			Direction: 0,
			LocX:      playerLocX,
			LocY:      playerLocY,
		},
		Food: data.Food{
			Img:  ebiten.NewImageFromImage(LoadSpriteImage("spritesheets/food-drink.png")),
			W:    data.FoodFrameWidth * data.FoodScale,
			H:    data.FoodFrameWidth * data.FoodScale,
			LocX: foodLocX,
			LocY: foodLocY,
		},
		Friends: []data.NPC{
			{
				Img:        ebiten.NewImageFromImage(LoadSpriteImage("spritesheets/character2.png")),
				W:          data.NPCFrameOX * data.NPCScale,
				H:          data.NPCFrameHeight * data.NPCScale,
				Speed:      3,
				Direction:  2,
				LocX:       friend1LocX,
				LocY:       friend1LocY,
				FrameLimit: 45,
			},
			{
				Img:        ebiten.NewImageFromImage(LoadSpriteImage("spritesheets/character3.png")),
				W:          data.NPCFrameOX * data.NPCScale,
				H:          data.NPCFrameHeight * data.NPCScale,
				Speed:      3,
				Direction:  2,
				LocX:       friend2LocX,
				LocY:       friend2LocY,
				FrameLimit: 25,
			},
			{
				Img:        ebiten.NewImageFromImage(LoadSpriteImage("spritesheets/character4.png")),
				W:          data.NPCFrameOX * data.NPCScale,
				H:          data.NPCFrameHeight * data.NPCScale,
				Speed:      3,
				Direction:  2,
				LocX:       friend3LocX,
				LocY:       friend3LocY,
				FrameLimit: 45,
			},
		},

		Enemies: []data.NPC{
			{
				Img:        ebiten.NewImageFromImage(LoadSpriteImage("spritesheets/enemy1.png")),
				W:          data.NPCFrameOX * data.NPCScale,
				H:          data.NPCFrameHeight * data.NPCScale,
				Speed:      3,
				Direction:  2,
				LocX:       enemy1LocX,
				LocY:       enemy1LocY,
				FrameLimit: 30,
			},
			{
				Img:        ebiten.NewImageFromImage(LoadSpriteImage("spritesheets/enemy2.png")),
				W:          data.NPCFrameOX * data.NPCScale,
				H:          data.NPCFrameHeight * data.NPCScale,
				Speed:      3,
				Direction:  2,
				LocX:       enemy2LocX,
				LocY:       enemy2LocY,
				FrameLimit: 30,
			},
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

	// friendly NPCs
	g.Friends[0].HitBox = HitBox(g.Friends[0].LocX, g.Friends[0].LocY, g.Friends[0].W, g.Friends[0].H)
	g.Friends[1].HitBox = HitBox(g.Friends[1].LocX, g.Friends[1].LocY, g.Friends[1].W, g.Friends[1].H)
	g.Friends[2].HitBox = HitBox(g.Friends[2].LocX, g.Friends[2].LocY, g.Friends[2].W, g.Friends[2].H)

	// enemy NPCs
	g.Enemies[0].HitBox = HitBox(g.Enemies[0].LocX, g.Enemies[0].LocY, g.Enemies[0].W, g.Enemies[0].H)
	g.Enemies[1].HitBox = HitBox(g.Enemies[1].LocX, g.Enemies[1].LocY, g.Enemies[1].W, g.Enemies[1].H)

	if ebiten.IsKeyPressed(ebiten.KeyShiftLeft) {
		g.Speed = 5
	} else if !ebiten.IsKeyPressed(ebiten.KeyShiftLeft) {
		g.Speed = 3
	}

	// NPCs random movement (all have same random pattern for now)
	g.Friends[0].Move(Bounds[MinX], Bounds[MaxX], Bounds[MinY], Bounds[MaxY])
	g.Friends[1].Move(Bounds[MinX], Bounds[MaxX], Bounds[MinY], Bounds[MaxY])
	g.Friends[2].Move(Bounds[MinX], Bounds[MaxX], Bounds[MinY], Bounds[MaxY])
	g.Enemies[0].Move(Bounds[MinX], Bounds[MaxX], Bounds[MinY], Bounds[MaxY])
	g.Enemies[1].Move(Bounds[MinX], Bounds[MaxX], Bounds[MinY], Bounds[MaxY])

	// player movement
	g.HandleMovement(Bounds[MinX], Bounds[MaxX], Bounds[MinY], Bounds[MaxY])

	// verify if player and food shape areas overlap
	if isCollision(g.Player.HitBox[MinX], g.Player.HitBox[MinY], g.Player.W, g.Player.H, g.Food.HitBox[MinX], g.Food.HitBox[MinY], g.Food.W, g.Food.H) {
		g.Food.IsDisplayed = false
		g.Food.ImgNo++
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.NRGBA{R: 255, G: 255, B: 255, A: 175})

	g.Player.DrawImage(screen)
	g.Friends[0].DrawImage(screen)
	g.Friends[1].DrawImage(screen)
	g.Friends[2].DrawImage(screen)

	g.Enemies[0].DrawImage(screen)
	g.Enemies[1].DrawImage(screen)

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
