package data

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	LocX, LocY     float64
	DeltaX, DeltaY float64
	Size           float64
	Speed          float64
	Breaks         float64
	HitBox         map[string]float64 // X Y min and max values
}

func (p *Player) HandleMovement(minX, maxX, minY, maxY float64) {
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		p.DeltaX += 0.05 * p.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		p.DeltaY -= 0.05 * p.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		p.DeltaY += 0.05 * p.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		p.DeltaX -= 0.05 * p.Speed
	}

	if !ebiten.IsKeyPressed(ebiten.KeyD) {
		if p.DeltaX > 0 {
			p.DeltaX -= 0.1 * p.Breaks
			if p.DeltaX < 0.00001 {
				p.DeltaX = 0
			}
		}
	}
	if !ebiten.IsKeyPressed(ebiten.KeyA) {
		if p.DeltaX < 0 {
			p.DeltaX += 0.1 * p.Breaks
			if p.DeltaX > 0 {
				p.DeltaX = 0
			}
		}
	}
	if !ebiten.IsKeyPressed(ebiten.KeyS) {
		if p.DeltaY > 0 {
			p.DeltaY -= 0.1 * p.Breaks
			if p.DeltaY < 0.00001 {
				p.DeltaY = 0
			}
		}
	}

	if !ebiten.IsKeyPressed(ebiten.KeyW) {
		if p.DeltaY < 0 {
			p.DeltaY += 0.1 * p.Breaks
			if p.DeltaY > 0 {
				p.DeltaY = 0
			}
		}
	}

	// update the position of the player
	p.LocX += p.DeltaX
	p.LocY += p.DeltaY

	// prevent player to go over the screen boundaries
	if p.LocX <= minX {
		p.LocX = minX
		p.DeltaX = 0
	}
	if p.LocX >= maxX-p.Size {
		p.LocX = maxX - p.Size
		p.DeltaX = 0
	}
	if p.LocY <= minY {
		p.LocY = minY
		p.DeltaY = 0
	}
	if p.LocY >= maxY-p.Size {
		p.LocY = maxY - p.Size
		p.DeltaY = 0
	}
}
