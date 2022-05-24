package game

import (
	u "games-ebiten/resources/utils"
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
	NPCScX         = 1.8
	NPCScY         = 1.8
)

// NPC - game character which implements InteractiveSprite interface
// FrTiming - used to time an action (movement or idle)
// FrTimingLimit - used as limit to count frames for an action (the time for an action to complete)
type NPC struct {
	Img           *ebiten.Image
	Name          string
	FrameNum      int
	Direction     int
	X, Y          float64
	DX, DY        float64
	W, H          float64
	Speed         float64
	HitBox        map[string]float64
	IsMoving      bool
	IsNearMargin  bool
	FrTiming      int
	FrTimingLimit int
}

func (npc *NPC) GetLocations() (float64, float64) {
	return npc.X, npc.Y
}

func (npc *NPC) GetSize() (float64, float64) {
	return npc.W, npc.H
}

func (npc *NPC) SetLocation(axis string, val float64) {
	if axis == u.X {
		npc.X = val
	} else if axis == u.Y {
		npc.Y = val
	}
}

func (npc *NPC) SetDelta(axis string, val float64) {
	if axis == u.X {
		npc.DX = val
	} else if axis == u.Y {
		npc.DY = val
	}
}

func (npc *NPC) DrawInteractiveSprite(screen *ebiten.Image) {
	opNPC := &ebiten.DrawImageOptions{}
	opNPC.GeoM.Scale(NPCScX, NPCScY)
	opNPC.GeoM.Translate(npc.X, npc.Y)

	x, y := NPCFrameOX+npc.FrameNum*NPCFrameWidth, NPCFrameOY+npc.Direction*NPCFrameHeight
	screen.DrawImage(npc.Img.SubImage(image.Rect(x, y, x+NPCFrameWidth, y+NPCFrameHeight)).(*ebiten.Image), opNPC)
}

func (npc *NPC) ValidateBoundaries(minX, maxX, minY, maxY float64) {
	if npc.X <= minX || npc.X >= maxX-npc.W || npc.Y <= minY || npc.Y >= maxY-npc.H {
		npc.IsNearMargin = true
		npc.IsMoving = false
		if npc.X <= minX {
			npc.X = minX
			npc.DX = 0
		}
		if npc.X >= maxX-npc.W {
			npc.X = maxX - npc.W
			npc.DX = 0
		}
		if npc.Y <= minY {
			npc.Y = minY
			npc.DY = 0
		}
		if npc.Y >= maxY-npc.H {
			npc.Y = maxY - npc.H
			npc.DY = 0
		}
	} else {
		npc.IsNearMargin = false
	}
}

func (npc *NPC) Move(minX, maxX, minY, maxY float64) {
	if npc.FrTiming == npc.FrTimingLimit-1 {
		npc.FrTiming = 0
		npc.IsMoving = !npc.IsMoving

		// force the NPC to walk the opposite way, no matter the IsMoving state
		if npc.IsNearMargin {
			switch npc.Direction {
			case 0:
				npc.Direction = 2
			case 1:
				npc.Direction = 3
			case 2:
				npc.Direction = 0
			case 3:
				npc.Direction = 1
			}
			npc.FrTiming = 0
			npc.IsMoving = true
			npc.IsNearMargin = false
		} else {
			rand.Seed(time.Now().UnixNano())
			npc.Direction = rand.Intn(4)
		}
	}

	// update X and Y based on Delta
	if npc.IsMoving {
		switch npc.Direction {

		// north
		case 0:
			npc.FrameNum++
			if npc.FrameNum == 8 {
				npc.FrameNum = 0
			}
			npc.DY = -3 * npc.Speed

		// west
		case 1:
			npc.FrameNum++
			if npc.FrameNum == 7 {
				npc.FrameNum = 0
			}
			npc.DX = -3 * npc.Speed

		// south
		case 2:
			npc.FrameNum++
			if npc.FrameNum == 7 {
				npc.FrameNum = 0
			}

			npc.DY = 3 * npc.Speed

		// east
		case 3:
			npc.FrameNum++
			if npc.FrameNum == 7 {
				npc.FrameNum = 0
			}

			npc.DX = 3 * npc.Speed
		}

		npc.X += npc.DX
		npc.Y += npc.DY

		if npc.X <= minX || npc.X >= maxX-npc.W || npc.Y <= minY || npc.Y >= maxY-npc.H {
			npc.IsNearMargin = true
			npc.IsMoving = false

			u.BoundaryValidation(npc, minX, maxX, minY, maxY)
		} else {
			npc.IsNearMargin = false
		}
	} else {
		npc.FrameNum = 0
		npc.DX = 0
		npc.DY = 0
	}
	npc.FrTiming++
}
