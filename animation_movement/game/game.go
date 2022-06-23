package game

import (
	"fmt"
	u "games-ebiten/resources/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	_ "image/png"
)

type Game struct {
	Items   []*Item
	Players []*Player
	NPCs    []*NPC
}

func NewGame() *Game {
	playerLocX := u.ScreenWidth/2 - PlayerFrameWidth/2
	playerLocY := u.ScreenHeight/2 - PlayerFrameHeight/2
	npcLocations := make(map[string][]int)

	// generate the random location (x, y) for every NP
	for i := 1; i <= 5; i++ {
		npcTag := fmt.Sprintf("npc%d", i)
		if _, ok := npcLocations[npcTag]; !ok {
			x, y := u.GenerateRandomPosition(0, 0, u.ScreenWidth-NPCFrameWidth, u.ScreenHeight-NPCFrameHeight)
			npcLocations[npcTag] = []int{x, y}
		}
	}
	ItemLocX, ItemLocY := u.GenerateRandomPosition(0, 0, u.ScreenWidth-ItemFrameWidth, u.ScreenHeight-ItemFrameHeight)

	return &Game{
		Players: []*Player{
			{
				Img:       ebiten.NewImageFromImage(u.LoadSpriteImage("resources/images/misc/character1.png")),
				Tag:       Player1,
				W:         PlayerFrameWidth * PlayerScX,
				H:         PlayerFrameHeight * PlayerScY,
				Speed:     3,
				Direction: 0,
				X:         playerLocX,
				Y:         playerLocY,
			},
		},
		Items: []*Item{
			{
				Img: ebiten.NewImageFromImage(u.LoadSpriteImage("resources/images/misc/food.png")),
				W:   ItemFrameWidth * ItemScale,
				H:   ItemFrameWidth * ItemScale,
				X:   ItemLocX,
				Y:   ItemLocY,
			},
		},
		NPCs: []*NPC{
			{
				Img:           ebiten.NewImageFromImage(u.LoadSpriteImage("resources/images/misc/character2.png")),
				Name:          u.NPC1,
				W:             NPCFrameWidth * NPCScX,
				H:             NPCFrameHeight * NPCScY,
				Speed:         3,
				Direction:     2,
				X:             npcLocations[u.NPC1][0],
				Y:             npcLocations[u.NPC1][1],
				FrTimingLimit: 45,
			},
			{
				Img:           ebiten.NewImageFromImage(u.LoadSpriteImage("resources/images/misc/character3.png")),
				Name:          u.NPC2,
				W:             NPCFrameWidth * NPCScX,
				H:             NPCFrameHeight * NPCScY,
				Speed:         3,
				Direction:     2,
				X:             npcLocations[u.NPC2][0],
				Y:             npcLocations[u.NPC2][1],
				FrTimingLimit: 25,
			},
			{
				Img:           ebiten.NewImageFromImage(u.LoadSpriteImage("resources/images/misc/character4.png")),
				Name:          u.NPC3,
				W:             NPCFrameWidth * NPCScX,
				H:             NPCFrameHeight * NPCScY,
				Speed:         3,
				Direction:     2,
				X:             npcLocations[u.NPC3][0],
				Y:             npcLocations[u.NPC3][1],
				FrTimingLimit: 45,
			},
			{
				Img:           ebiten.NewImageFromImage(u.LoadSpriteImage("resources/images/misc/character5.png")),
				Name:          u.NPC4,
				W:             NPCFrameWidth * NPCScX,
				H:             NPCFrameHeight * NPCScY,
				Speed:         3,
				Direction:     2,
				X:             npcLocations[u.NPC4][0],
				Y:             npcLocations[u.NPC4][1],
				FrTimingLimit: 30,
			},
			{
				Img:           ebiten.NewImageFromImage(u.LoadSpriteImage("resources/images/misc/character6.png")),
				Name:          u.NPC5,
				W:             NPCFrameWidth * NPCScX,
				H:             NPCFrameHeight * NPCScY,
				Speed:         3,
				Direction:     2,
				X:             npcLocations[u.NPC5][0],
				Y:             npcLocations[u.NPC5][1],
				FrTimingLimit: 30,
			},
		},
	}
}

func (g *Game) Update() error {
	for _, p := range g.Players {
		p.HandleMovement(0, 0, u.ScreenWidth, u.ScreenWidth)
	}

	for _, npc := range g.NPCs {
		npc.Move(0, 0, u.ScreenWidth, u.ScreenWidth)
	}

	// player1 speed up
	if ebiten.IsKeyPressed(ebiten.KeyShiftLeft) {
		g.Players[0].Speed = 5
	} else if !ebiten.IsKeyPressed(ebiten.KeyShiftLeft) {
		g.Players[0].Speed = 3
	}

	// update the Item state if any NPC collides with Item shape
	for _, npc := range g.NPCs {
		if u.IsCollision(npc.GetGeomData(), g.Items[0].GetGeomData()) {
			g.Items[0].UpdateItemState()
		}
	}

	// update the Item state if the player and Item shape areas overlap
	for _, p := range g.Players {
		if u.IsCollision(p.GetGeomData(), g.Items[0].GetGeomData()) {
			g.Items[0].UpdateItemState()
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(5, 4)
	screen.DrawImage(ebiten.NewImageFromImage(u.LoadSpriteImage("resources/images/misc/bg-grass.png")), op)

	for i := range g.Players {
		g.Players[i].DrawInteractiveSprite(screen)
	}

	for i := range g.NPCs {
		g.NPCs[i].DrawInteractiveSprite(screen)
	}

	for i := range g.Items {
		g.Items[i].DrawStaticSprite(screen)
	}

	ebitenutil.DebugPrintAt(screen, "W Ace S D to move", 0, 0)
	ebitenutil.DebugPrintAt(screen, "LEFT SHIFT to speed up", 0, 25)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return u.ScreenWidth, u.ScreenHeight
}
