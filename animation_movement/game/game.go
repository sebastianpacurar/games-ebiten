package game

import (
	"fmt"
	u "games-ebiten/resources/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	_ "image/png"
)

// Bounds - represents the main screen edges
var (
	Bounds = map[string]float64{u.MinX: 0, u.MaxX: u.ScreenWidth, u.MinY: 0, u.MaxY: u.ScreenHeight}
	inset  = float64(15)
)

type Game struct {
	Food
	Players []*Player
	NPCs    []*NPC
}

func NewGame() *Game {
	playerLocX := float64(u.ScreenWidth)/2 - PlayerFrameWidth/2
	playerLocY := float64(u.ScreenHeight)/2 - PlayerFrameHeight/2
	npcLocations := make(map[string][]float64)

	// generate the random location (x, y) for every NPC (there is a 15 pixel inset for safety)
	for i := 1; i <= 5; i++ {
		npcTag := fmt.Sprintf("npc%d", i)
		if _, ok := npcLocations[npcTag]; !ok {
			x, y := u.GenerateRandomLocation(Bounds[u.MinX]+NPCFrameWidth, Bounds[u.MaxX]-(NPCFrameWidth+inset), Bounds[u.MinY]+NPCFrameHeight, Bounds[u.MaxY]-(NPCFrameHeight+inset))
			npcLocations[npcTag] = []float64{x, y}
		}
	}
	foodLocX, foodLocY := u.GenerateRandomLocation(Bounds[u.MinX], Bounds[u.MaxX]-FoodFrameWidth, Bounds[u.MinY], Bounds[u.MaxY]-FoodFrameHeight)

	return &Game{
		Players: []*Player{
			{
				Img:       ebiten.NewImageFromImage(u.LoadSpriteImage("resources/images/misc/character1.png")),
				Tag:       Player1,
				W:         PlayerFrameWidth * PlayerScaleX,
				H:         PlayerFrameHeight * PlayerScaleY,
				Speed:     3,
				Direction: 0,
				LX:        playerLocX,
				LY:        playerLocY,
			},
		},
		Food: Food{
			Img: ebiten.NewImageFromImage(u.LoadSpriteImage("resources/images/misc/food.png")),
			W:   FoodFrameWidth * FoodScale,
			H:   FoodFrameWidth * FoodScale,
			LX:  foodLocX,
			LY:  foodLocY,
		},
		NPCs: []*NPC{
			{
				Img:        ebiten.NewImageFromImage(u.LoadSpriteImage("resources/images/misc/character2.png")),
				Name:       NPC1,
				W:          NPCFrameWidth * NPCScaleX,
				H:          NPCFrameHeight * NPCScaleY,
				Speed:      3,
				Direction:  2,
				LX:         npcLocations[NPC1][0],
				LY:         npcLocations[NPC1][1],
				FrameLimit: 45,
			},
			{
				Img:        ebiten.NewImageFromImage(u.LoadSpriteImage("resources/images/misc/character3.png")),
				Name:       NPC2,
				W:          NPCFrameWidth * NPCScaleX,
				H:          NPCFrameHeight * NPCScaleY,
				Speed:      3,
				Direction:  2,
				LX:         npcLocations[NPC2][0],
				LY:         npcLocations[NPC2][1],
				FrameLimit: 25,
			},
			{
				Img:        ebiten.NewImageFromImage(u.LoadSpriteImage("resources/images/misc/character4.png")),
				Name:       NPC3,
				W:          NPCFrameWidth * NPCScaleX,
				H:          NPCFrameHeight * NPCScaleY,
				Speed:      3,
				Direction:  2,
				LX:         npcLocations[NPC3][0],
				LY:         npcLocations[NPC3][1],
				FrameLimit: 45,
			},
			{
				Img:        ebiten.NewImageFromImage(u.LoadSpriteImage("resources/images/misc/character5.png")),
				Name:       NPC4,
				W:          NPCFrameWidth * NPCScaleX,
				H:          NPCFrameHeight * NPCScaleY,
				Speed:      3,
				Direction:  2,
				LX:         npcLocations[NPC4][0],
				LY:         npcLocations[NPC4][1],
				FrameLimit: 30,
			},
			{
				Img:        ebiten.NewImageFromImage(u.LoadSpriteImage("resources/images/misc/character6.png")),
				Name:       NPC5,
				W:          NPCFrameWidth * NPCScaleX,
				H:          NPCFrameHeight * NPCScaleY,
				Speed:      3,
				Direction:  2,
				LX:         npcLocations[NPC5][0],
				LY:         npcLocations[NPC5][1],
				FrameLimit: 30,
			},
		},
	}
}

func (g *Game) Update() error {
	g.Food.HitBox = u.HitBox(g.Food.LX, g.Food.LY, g.Food.W, g.Food.H)

	for i := range g.Players {
		g.Players[i].HitBox = u.HitBox(g.Players[i].LX, g.Players[i].LY, g.Players[i].W, g.Players[i].H)
		g.Players[i].HandleMovement(Bounds[u.MinX], Bounds[u.MaxX], Bounds[u.MinY], Bounds[u.MaxY])
	}

	for i := range g.NPCs {
		g.NPCs[i].HitBox = u.HitBox(g.NPCs[i].LX, g.NPCs[i].LY, g.NPCs[i].W, g.NPCs[i].H)
		g.NPCs[i].Move(Bounds[u.MinX], Bounds[u.MaxX], Bounds[u.MinY], Bounds[u.MaxY])
	}

	// player1 speed up
	if ebiten.IsKeyPressed(ebiten.KeyShiftLeft) {
		g.Players[0].Speed = 5
	} else if !ebiten.IsKeyPressed(ebiten.KeyShiftLeft) {
		g.Players[0].Speed = 3
	}

	// update the food state if any NPC collides with food shape
	for i := range g.NPCs {
		if u.IsCollision(g.NPCs[i].HitBox[u.MinX], g.NPCs[i].HitBox[u.MinY], g.NPCs[i].W, g.NPCs[i].H, g.Food.HitBox[u.MinX], g.Food.HitBox[u.MinY], g.Food.W, g.Food.H) {
			g.Food.UpdateFoodState()
		}
	}

	// update the food state if the player and food shape areas overlap
	for i := range g.Players {
		if u.IsCollision(g.Players[i].HitBox[u.MinX], g.Players[i].HitBox[u.MinY], g.Players[i].W, g.Players[i].H, g.Food.HitBox[u.MinX], g.Food.HitBox[u.MinY], g.Food.W, g.Food.H) {
			g.Food.UpdateFoodState()
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(5, 4)
	screen.DrawImage(ebiten.NewImageFromImage(u.LoadSpriteImage("resources/images/misc/bg-grass.png")), op)

	for i := range g.Players {
		g.Players[i].DrawImg(screen)
	}

	for i := range g.NPCs {
		g.NPCs[i].DrawImg(screen)
	}

	g.Food.DrawImage(screen)

	ebitenutil.DebugPrintAt(screen, "W A S D to move", 0, 0)
	ebitenutil.DebugPrintAt(screen, "LEFT SHIFT to speed up", 0, 25)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return u.ScreenWidth, u.ScreenHeight
}
