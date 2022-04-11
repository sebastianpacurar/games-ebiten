package data

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"math/rand"
	"time"
)

const (
	NPCFrameOX     = 0
	NPCFrameOY     = 0
	NPCFrameWidth  = 64
	NPCFrameHeight = 64
	NPCScale       = 1.5
)

type NPC struct {
	Img            *ebiten.Image
	FrameNum       int
	Direction      int
	LocX, LocY     float64
	DeltaX, DeltaY float64
	W, H           float64
	Speed          float64
	HitBox         map[string]float64
	isMoving       bool
	FrameCount     int // used to time an action (movement or idle)
	FrameLimit     int // used as limit to count frames for an action (the time for an action to complete)
}

func (npc *NPC) Move(minX, maxX, minY, maxY float64) {
	rand.Seed(time.Now().UnixNano())

	if npc.FrameCount == npc.FrameLimit-1 {
		npc.FrameCount = 0
		npc.isMoving = !npc.isMoving
		npc.Direction = rand.Intn(4)
	}
	if npc.isMoving {
		switch npc.Direction {

		// north
		case 0:
			npc.FrameNum++
			if npc.FrameNum == 8 {
				npc.FrameNum = 0
			}
			npc.DeltaY = -3 * npc.Speed

		// west
		case 1:
			npc.FrameNum++
			if npc.FrameNum == 7 {
				npc.FrameNum = 0
			}
			npc.DeltaX = -3 * npc.Speed

		// south
		case 2:
			npc.FrameNum++
			if npc.FrameNum == 7 {
				npc.FrameNum = 0
			}

			npc.DeltaY = 3 * npc.Speed

		// east
		case 3:
			npc.FrameNum++
			if npc.FrameNum == 7 {
				npc.FrameNum = 0
			}

			npc.DeltaX = 3 * npc.Speed
		}

		npc.LocX += npc.DeltaX
		npc.LocY += npc.DeltaY

		// prevent player to go over the screen boundaries
		if npc.LocX <= minX {
			npc.LocX = minX
			npc.DeltaX = 0
		}
		if npc.LocX >= maxX-npc.W {
			npc.LocX = maxX - npc.W
			npc.DeltaX = 0
		}
		if npc.LocY <= minY {
			npc.LocY = minY
			npc.DeltaY = 0
		}
		if npc.LocY >= maxY-npc.H {
			npc.LocY = maxY - npc.H
			npc.DeltaY = 0
		}

	} else {
		npc.FrameNum = 0
		npc.DeltaX = 0
		npc.DeltaY = 0
	}
	npc.FrameCount++
}

func (npc *NPC) DrawImage(screen *ebiten.Image) {
	opNPC := &ebiten.DrawImageOptions{}
	opNPC.GeoM.Scale(NPCScale, NPCScale)
	opNPC.GeoM.Translate(npc.LocX, npc.LocY)

	// load every sub image based on the received key input
	x, y := NPCFrameOX+npc.FrameNum*NPCFrameWidth, NPCFrameOY+npc.Direction*NPCFrameHeight
	screen.DrawImage(npc.Img.SubImage(image.Rect(x, y, x+NPCFrameWidth, y+NPCFrameHeight)).(*ebiten.Image), opNPC)
}
