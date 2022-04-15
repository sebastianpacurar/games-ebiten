package game

import (
	"fmt"
	u "games-ebiten/resources/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	_ "image/png"
)

// Bounds - represents the main screen edges
var (
	Bounds = map[string]float64{u.MinX: 0, u.MaxX: u.ScreenWidth, u.MinY: 0, u.MaxY: u.ScreenHeight}
	offset = float64(10)
)

type Game struct {
	Player
	Food
	NPCs []*NPC
}

func NewGame() *Game {
	playerLocX := float64(u.ScreenWidth)/2 - PlayerFrameWidth/2
	playerLocY := float64(u.ScreenHeight)/2 - PlayerFrameHeight/2
	npcLocations := make(map[string][]float64)

	// generate the random location (x, y) for every NPC (there is a 10 pixel inset for safety)
	for i := 1; i <= 5; i++ {
		npcTag := fmt.Sprintf("npc%d", i)
		if _, ok := npcLocations[npcTag]; !ok {
			x, y := u.GenerateRandomLocation(Bounds[u.MinX]+offset, Bounds[u.MaxX]-(NPCFrameWidth+offset), Bounds[u.MinY]+offset, Bounds[u.MaxY]-(NPCFrameHeight+offset))
			npcLocations[npcTag] = []float64{x, y}
		}
	}
	foodLocX, foodLocY := u.GenerateRandomLocation(Bounds[u.MinX], Bounds[u.MaxX]-FoodFrameWidth, Bounds[u.MinY], Bounds[u.MaxY]-FoodFrameHeight)

	return &Game{
		Player: Player{
			Img:       ebiten.NewImageFromImage(u.LoadSpriteImage("resources/images/character1.png")),
			W:         PlayerFrameWidth * PlayerScale,
			H:         PlayerFrameHeight * PlayerScale,
			Speed:     3,
			Direction: 0,
			LocX:      playerLocX,
			LocY:      playerLocY,
		},
		Food: Food{
			Img:  ebiten.NewImageFromImage(u.LoadSpriteImage("resources/images/food.png")),
			W:    FoodFrameWidth * FoodScale,
			H:    FoodFrameWidth * FoodScale,
			LocX: foodLocX,
			LocY: foodLocY,
		},
		NPCs: []*NPC{
			{
				Img:        ebiten.NewImageFromImage(u.LoadSpriteImage("resources/images/character2.png")),
				Name:       NPC1,
				W:          NPCFrameWidth * NPCScale,
				H:          NPCFrameHeight * NPCScale,
				Speed:      3,
				Direction:  2,
				LocX:       npcLocations[NPC1][0],
				LocY:       npcLocations[NPC1][1],
				FrameLimit: 45,
			},
			{
				Img:        ebiten.NewImageFromImage(u.LoadSpriteImage("resources/images/character3.png")),
				Name:       NPC2,
				W:          NPCFrameWidth * NPCScale,
				H:          NPCFrameHeight * NPCScale,
				Speed:      3,
				Direction:  2,
				LocX:       npcLocations[NPC2][0],
				LocY:       npcLocations[NPC2][1],
				FrameLimit: 25,
			},
			{
				Img:        ebiten.NewImageFromImage(u.LoadSpriteImage("resources/images/character4.png")),
				Name:       NPC3,
				W:          NPCFrameWidth * NPCScale,
				H:          NPCFrameHeight * NPCScale,
				Speed:      3,
				Direction:  2,
				LocX:       npcLocations[NPC3][0],
				LocY:       npcLocations[NPC3][1],
				FrameLimit: 45,
			},
			{
				Img:        ebiten.NewImageFromImage(u.LoadSpriteImage("resources/images/character5.png")),
				Name:       NPC4,
				W:          NPCFrameWidth * NPCScale,
				H:          NPCFrameHeight * NPCScale,
				Speed:      3,
				Direction:  2,
				LocX:       npcLocations[NPC4][0],
				LocY:       npcLocations[NPC4][1],
				FrameLimit: 30,
			},
			{
				Img:        ebiten.NewImageFromImage(u.LoadSpriteImage("resources/images/character6.png")),
				Name:       NPC5,
				W:          NPCFrameWidth * NPCScale,
				H:          NPCFrameHeight * NPCScale,
				Speed:      3,
				Direction:  2,
				LocX:       npcLocations[NPC5][0],
				LocY:       npcLocations[NPC5][1],
				FrameLimit: 30,
			},
		},
	}
}

func (g *Game) Update() error {
	g.Player.HitBox = u.HitBox(g.Player.LocX, g.Player.LocY, g.Player.W, g.Player.H)
	g.Food.HitBox = u.HitBox(g.Food.LocX, g.Food.LocY, g.Food.W, g.Food.H)

	for i := range g.NPCs {
		g.NPCs[i].HitBox = u.HitBox(g.NPCs[i].LocX, g.NPCs[i].LocY, g.NPCs[i].W, g.NPCs[i].H)
		g.NPCs[i].Move(Bounds[u.MinX], Bounds[u.MaxX], Bounds[u.MinY], Bounds[u.MaxY])
	}

	if ebiten.IsKeyPressed(ebiten.KeyShiftLeft) {
		g.Speed = 5
	} else if !ebiten.IsKeyPressed(ebiten.KeyShiftLeft) {
		g.Speed = 3
	}

	// player movement
	g.Player.HandleMovement(Bounds[u.MinX], Bounds[u.MaxX], Bounds[u.MinY], Bounds[u.MaxY])

	// update the food state if any NPC collides with food shape
	for i := range g.NPCs {
		if u.IsCollision(g.NPCs[i].HitBox[u.MinX], g.NPCs[i].HitBox[u.MinY], g.NPCs[i].W, g.NPCs[i].H, g.Food.HitBox[u.MinX], g.Food.HitBox[u.MinY], g.Food.W, g.Food.H) {
			g.Food.UpdateFoodState()
		}
	}

	// update the food state if the player and food shape areas overlap
	if u.IsCollision(g.Player.HitBox[u.MinX], g.Player.HitBox[u.MinY], g.Player.W, g.Player.H, g.Food.HitBox[u.MinX], g.Food.HitBox[u.MinY], g.Food.W, g.Food.H) {
		g.Food.UpdateFoodState()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.NRGBA{R: 255, G: 255, B: 255, A: 175})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(5, 4)
	screen.DrawImage(ebiten.NewImageFromImage(u.LoadSpriteImage("resources/images/bg-grass.png")), op)
	g.Player.DrawImage(screen)

	for i := range g.NPCs {
		g.NPCs[i].DrawImage(screen)
	}

	g.Food.DrawImage(screen)

	ebitenutil.DebugPrintAt(screen, "W A S D to move", 0, 0)
	ebitenutil.DebugPrintAt(screen, "LEFT SHIFT to speed up", 0, 25)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return u.ScreenWidth, u.ScreenHeight
}
