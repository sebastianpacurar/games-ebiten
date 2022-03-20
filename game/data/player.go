package data

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	PosX, PosY     float64 // the centre of the image
	DeltaX, DeltaY float64
	Size           float64 // the image's size
	Origin         float64 // the centre of the image
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
	p.PosX += p.DeltaX
	p.PosY += p.DeltaY

	// prevent player to go over the screen boundaries
	if p.PosX <= minX {
		p.PosX = minX
		p.DeltaX = 0
	}
	if p.PosX >= maxX-p.Size {
		p.PosX = maxX - p.Size
		p.DeltaX = 0
	}
	if p.PosY <= minY {
		p.PosY = minY
		p.DeltaY = 0
	}
	if p.PosY >= maxY-p.Size {
		p.PosY = maxY - p.Size
		p.DeltaY = 0
	}
}

// GetHitBox - generate the player's boundaries for minX, maxX and minY, maxY
func (p *Player) GetHitBox() map[string]float64 {
	return map[string]float64{
		"minX": p.PosX + p.Origin,
		"maxX": p.PosX + p.Origin + p.Size,
		"minY": p.PosY + p.Origin,
		"maxY": p.PosY + p.Origin + p.Size,
	}
}
