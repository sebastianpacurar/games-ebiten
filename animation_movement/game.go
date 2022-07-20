package animation_movement

import (
	"fmt"
	res "games-ebiten/resources"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	_ "image/png"
	"math"
)

type Game struct {
	Items   []*Item
	Players []*Player
	NPCs    []*NPC
}

func NewGame() *Game {
	ebiten.SetMaxTPS(40)
	playerLocX := res.ScreenWidth/2 - PlayerFrameWidth/2
	playerLocY := (res.ScreenHeight/2 + res.MainMenuH) - PlayerFrameHeight/2
	npcLocations := make(map[string][]int)

	// generate the random location (x, y) for every NP
	for i := 1; i <= 5; i++ {
		npcTag := fmt.Sprintf("npc%d", i)
		if _, ok := npcLocations[npcTag]; !ok {
			x, y := res.GenerateRandomPosition(0, res.MainMenuH, res.ScreenWidth-NPCFrameWidth, res.ScreenHeight-NPCFrameHeight)
			npcLocations[npcTag] = []int{x, y}
		}
	}
	ItemLocX, ItemLocY := res.GenerateRandomPosition(0, res.MainMenuH, res.ScreenWidth-ItemFrameWidth, res.ScreenHeight-ItemFrameHeight)

	return &Game{
		Players: []*Player{
			{
				Img:       ebiten.NewImageFromImage(res.LoadSpriteImage("resources/assets/misc/character1.png")),
				Tag:       Player1,
				W:         PlayerFrameWidth * PlayerScX,
				H:         PlayerFrameHeight * PlayerScY,
				Speed:     1,
				Direction: 0,
				X:         playerLocX,
				Y:         playerLocY,
			},
		},
		Items: []*Item{
			{
				Img: ebiten.NewImageFromImage(res.LoadSpriteImage("resources/assets/misc/food.png")),
				W:   int(math.Floor(ItemFrameWidth) * ItemScale),
				H:   int(math.Floor(ItemFrameHeight) * ItemScale),
				X:   ItemLocX,
				Y:   ItemLocY,
			},
		},
		NPCs: []*NPC{
			{
				Img:           ebiten.NewImageFromImage(res.LoadSpriteImage("resources/assets/misc/character2.png")),
				Name:          "npc1",
				W:             int(math.Floor(NPCFrameWidth) * NPCScX),
				H:             int(math.Floor(NPCFrameHeight) * NPCScY),
				Speed:         4,
				Direction:     2,
				X:             npcLocations["npc1"][0],
				Y:             npcLocations["npc1"][1],
				FrTimingLimit: 45,
			},
			{
				Img:           ebiten.NewImageFromImage(res.LoadSpriteImage("resources/assets/misc/character3.png")),
				Name:          "npc2",
				W:             int(math.Floor(NPCFrameWidth) * NPCScX),
				H:             int(math.Floor(NPCFrameHeight) * NPCScY),
				Speed:         4,
				Direction:     2,
				X:             npcLocations["npc2"][0],
				Y:             npcLocations["npc2"][1],
				FrTimingLimit: 45,
			},
			{
				Img:           ebiten.NewImageFromImage(res.LoadSpriteImage("resources/assets/misc/character4.png")),
				Name:          "npc3",
				W:             int(math.Floor(NPCFrameWidth) * NPCScX),
				H:             int(math.Floor(NPCFrameHeight) * NPCScY),
				Speed:         4,
				Direction:     2,
				X:             npcLocations["npc3"][0],
				Y:             npcLocations["npc3"][1],
				FrTimingLimit: 45,
			},
		},
	}
}

func (g *Game) Update() error {
	for _, p := range g.Players {
		p.HandleMovement(0, 0, res.ScreenWidth, res.ScreenHeight)
	}

	for _, npc := range g.NPCs {
		npc.Move(0, res.MainMenuH, res.ScreenWidth, res.ScreenHeight)
	}

	// player1 speed up
	if ebiten.IsKeyPressed(ebiten.KeyShiftLeft) {
		g.Players[0].Speed = 5
	} else if !ebiten.IsKeyPressed(ebiten.KeyShiftLeft) {
		g.Players[0].Speed = 3
	}

	// update the Item state if any NPC collides with Item shape
	for _, npc := range g.NPCs {
		if res.IsCollision(npc.HitBox(), g.Items[0].HitBox()) {
			g.Items[0].UpdateItemState()
		}
	}

	// update the Item state if the player and Item shape areas overlap
	for _, p := range g.Players {
		if res.IsCollision(p.HitBox(), g.Items[0].HitBox()) {
			g.Items[0].UpdateItemState()
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(5, 4)
	screen.DrawImage(ebiten.NewImageFromImage(res.LoadSpriteImage("resources/assets/misc/bg-grass.png")), op)

	for _, p := range g.Players {
		p.DrawSprite(screen)
	}

	for _, npc := range g.NPCs {
		npc.DrawSprite(screen)
	}

	for _, i := range g.Items {
		i.DrawSprite(screen)
	}

	ebitenutil.DebugPrintAt(screen, "W A S D to move", 10, res.ScreenHeight-65)
	ebitenutil.DebugPrintAt(screen, "LEFT SHIFT to speed up", 10, res.ScreenHeight-35)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return res.ScreenWidth, res.ScreenHeight
}
