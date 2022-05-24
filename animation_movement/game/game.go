package game

import (
	"fmt"
	u "games-ebiten/resources/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	_ "image/png"
)

// u.ScreenDims - represents the main screen edges
var inset = float64(15)

type Game struct {
	Items   []*Item
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
			x, y := u.GenerateRandomLocation(u.ScreenDims[u.MinX]+NPCFrameWidth, u.ScreenDims[u.MaxX]-(NPCFrameWidth+inset), u.ScreenDims[u.MinY]+NPCFrameHeight, u.ScreenDims[u.MaxY]-(NPCFrameHeight+inset))
			npcLocations[npcTag] = []float64{x, y}
		}
	}
	ItemLocX, ItemLocY := u.GenerateRandomLocation(u.ScreenDims[u.MinX], u.ScreenDims[u.MaxX]-ItemFrameWidth, u.ScreenDims[u.MinY], u.ScreenDims[u.MaxY]-ItemFrameHeight)

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
	for i := range g.Items {
		item := g.Items[i]
		item.HitBox = u.HitBox(item.X, item.Y, item.W, item.H)
	}

	for i := range g.Players {
		p := g.Players[i]
		p.HitBox = u.HitBox(p.X, p.Y, p.W, p.H)
		p.HandleMovement(u.ScreenDims[u.MinX], u.ScreenDims[u.MaxX], u.ScreenDims[u.MinY], u.ScreenDims[u.MaxY])
	}

	for i := range g.NPCs {
		npc := g.NPCs[i]
		npc.HitBox = u.HitBox(npc.X, npc.Y, npc.W, npc.H)
		npc.Move(u.ScreenDims[u.MinX], u.ScreenDims[u.MaxX], u.ScreenDims[u.MinY], u.ScreenDims[u.MaxY])
	}

	// player1 speed up
	if ebiten.IsKeyPressed(ebiten.KeyShiftLeft) {
		g.Players[0].Speed = 5
	} else if !ebiten.IsKeyPressed(ebiten.KeyShiftLeft) {
		g.Players[0].Speed = 3
	}

	// update the Item state if any NPC collides with Item shape
	for i := range g.NPCs {
		npc := g.NPCs[i]
		item := g.Items[0]
		if u.IsCollision(npc.HitBox[u.MinX], npc.HitBox[u.MinY], npc.W, npc.H, item.HitBox[u.MinX], item.HitBox[u.MinY], item.W, item.H) {
			g.Items[0].UpdateItemState()
		}
	}

	// update the Item state if the player and Item shape areas overlap
	for i := range g.Players {
		p := g.Players[i]
		item := g.Items[0]
		if u.IsCollision(p.HitBox[u.MinX], p.HitBox[u.MinY], p.W, p.H, item.HitBox[u.MinX], item.HitBox[u.MinY], item.W, item.H) {
			item.UpdateItemState()
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
