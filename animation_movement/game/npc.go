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
	NPCScaleX      = 1.8
	NPCScaleY      = 1.8
)

// NPC - game character which implements InteractiveSprite interface
type NPC struct {
	Img          *ebiten.Image
	Name         string
	FrameNum     int
	Direction    int
	LX, LY       float64
	DX, DY       float64
	W, H         float64
	Speed        float64
	HitBox       map[string]float64
	IsMoving     bool
	IsNearMargin bool
	FrameCount   int // used to time an action (movement or idle)
	FrameLimit   int // used as limit to count frames for an action (the time for an action to complete)
}

func (npc *NPC) GetLocations() (float64, float64) {
	return npc.LX, npc.LY
}

func (npc *NPC) GetSize() (float64, float64) {
	return npc.W, npc.H
}

func (npc *NPC) GetFrameInfo() (int, int, int, int) {
	return NPCFrameOX, NPCFrameOY, NPCFrameWidth, NPCFrameHeight
}

func (npc *NPC) GetScaleVal() (float64, float64) {
	return NPCScaleX, NPCScaleY
}

func (npc *NPC) GetFrameNum() int {
	return npc.FrameNum
}

func (npc *NPC) GetDirection() int {
	return npc.Direction
}

func (npc *NPC) GetImg() *ebiten.Image {
	return npc.Img
}

func (npc *NPC) SetLocation(axis string, val float64) {
	if axis == u.X {
		npc.LX = val
	} else if axis == u.Y {
		npc.LY = val
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
	opNPC.GeoM.Scale(NPCScaleX, NPCScaleY)
	opNPC.GeoM.Translate(npc.LX, npc.LY)

	x, y := NPCFrameOX+npc.FrameNum*NPCFrameWidth, NPCFrameOY+npc.Direction*NPCFrameHeight
	screen.DrawImage(npc.Img.SubImage(image.Rect(x, y, x+NPCFrameWidth, y+NPCFrameHeight)).(*ebiten.Image), opNPC)
}

// GetHitBox - returns 4 values: minX, maxX, minY, maxY
func (npc *NPC) GetHitBox() (float64, float64, float64, float64) {
	return npc.HitBox[u.MinX], npc.HitBox[u.MaxX], npc.HitBox[u.MinY], npc.HitBox[u.MaxY]
}

func (npc *NPC) ValidateBoundaries(minX, maxX, minY, maxY float64) {
	if npc.LX <= minX || npc.LX >= maxX-npc.W || npc.LY <= minY || npc.LY >= maxY-npc.H {
		npc.IsNearMargin = true
		npc.IsMoving = false
		if npc.LX <= minX {
			npc.LX = minX
			npc.DX = 0
		}
		if npc.LX >= maxX-npc.W {
			npc.LX = maxX - npc.W
			npc.DX = 0
		}
		if npc.LY <= minY {
			npc.LY = minY
			npc.DY = 0
		}
		if npc.LY >= maxY-npc.H {
			npc.LY = maxY - npc.H
			npc.DY = 0
		}
	} else {
		npc.IsNearMargin = false
	}
}

func (npc *NPC) Move(minX, maxX, minY, maxY float64) {
	if npc.FrameCount == npc.FrameLimit-1 {
		npc.FrameCount = 0
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
			npc.FrameCount = 0
			npc.IsMoving = true
			npc.IsNearMargin = false
		} else {
			rand.Seed(time.Now().UnixNano())
			npc.Direction = rand.Intn(4)
		}
	}

	// update LX and LY based on Delta
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

		npc.LX += npc.DX
		npc.LY += npc.DY

		if npc.LX <= minX || npc.LX >= maxX-npc.W || npc.LY <= minY || npc.LY >= maxY-npc.H {
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
	npc.FrameCount++
}
