package data

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

const (
	PlayerFrameOX     = 0
	PlayerFrameOY     = 0
	PlayerFrameWidth  = 24
	PlayerFrameHeight = 32
	PlayerScale       = 3
)

// Player FrameNum - animate movement - all actions are from left to right in the sprite sheet
// Player Direction - specifies which row to pick (0 = down, 1 = up, 2 = left, 3 = right)
// Player LocX, LocY - the location of the image on the screen (starts from top left corner)
// Player W, H - represent the size (width, height) which is calculated by multiplying each with PlayerScale
type Player struct {
	Img            *ebiten.Image
	FrameNum       int
	Direction      int
	LocX, LocY     float64
	DeltaX, DeltaY float64
	W, H           float64
	Speed          float64
	HitBox         map[string]float64 // X Y min and max values
}

// HandleMovement - takes the vertices as params for screen cross boundary prevention
func (p *Player) HandleMovement(minX, maxX, minY, maxY float64) {
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		p.Direction = 0
		if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyD) {
			p.DeltaX = 0
		}

		p.FrameNum++
		if p.FrameNum == 7 {
			p.FrameNum = 0
		}

		p.DeltaY = 3 * p.Speed
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		p.Direction = 1
		if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyD) {
			p.DeltaX = 0
		}

		p.FrameNum++
		if p.FrameNum == 7 {
			p.FrameNum = 0
		}

		p.DeltaY = -3 * p.Speed
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyS) {
			p.DeltaY = 0
		}

		p.FrameNum++
		if p.FrameNum == 7 {
			p.FrameNum = 0
		}

		p.Direction = 2
		p.DeltaX = -3 * p.Speed
	}

	if ebiten.IsKeyPressed(ebiten.KeyD) {
		p.Direction = 3
		if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyS) {
			p.DeltaY = 0
		}

		p.FrameNum++
		if p.FrameNum == 7 {
			p.FrameNum = 0
		}

		p.DeltaX = 3 * p.Speed
	}

	if !ebiten.IsKeyPressed(ebiten.KeyA) && !ebiten.IsKeyPressed(ebiten.KeyD) {
		p.DeltaX = 0
	}

	if !ebiten.IsKeyPressed(ebiten.KeyS) && !ebiten.IsKeyPressed(ebiten.KeyW) {
		p.DeltaY = 0
	}

	if p.DeltaX == 0 && p.DeltaY == 0 {
		p.FrameNum = 0
	}

	// update the position of the player
	p.LocX += p.DeltaX
	p.LocY += p.DeltaY

	// prevent player to go over the screen boundaries
	if p.LocX <= minX {
		p.LocX = minX
		p.DeltaX = 0
	}
	if p.LocX >= maxX-p.W {
		p.LocX = maxX - p.W
		p.DeltaX = 0
	}
	if p.LocY <= minY {
		p.LocY = minY
		p.DeltaY = 0
	}
	if p.LocY >= maxY-p.H {
		p.LocY = maxY - p.H
		p.DeltaY = 0
	}
}

func (p *Player) DrawImage(screen *ebiten.Image) {
	opPlayer := &ebiten.DrawImageOptions{}
	opPlayer.GeoM.Scale(PlayerScale, PlayerScale)
	opPlayer.GeoM.Translate(p.LocX, p.LocY)

	// load every sub image based on the received key input
	px, py := PlayerFrameOX+p.FrameNum*PlayerFrameWidth, PlayerFrameOY+p.Direction*PlayerFrameHeight
	screen.DrawImage(p.Img.SubImage(image.Rect(px, py, px+PlayerFrameWidth, py+PlayerFrameHeight)).(*ebiten.Image), opPlayer)
}

func (p *Player) GenerateHitBox() {

}
