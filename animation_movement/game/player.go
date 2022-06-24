package game

import (
	u "games-ebiten/resources/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

const (
	PlayerFrameOX     = 0
	PlayerFrameOY     = 0
	PlayerFrameWidth  = 24
	PlayerFrameHeight = 32
	PlayerScX         = 3
	PlayerScY         = 3

	Player1 = "player1"
)

// Player - game character which implements InteractiveSprite interface
// Tag - highlights which player is related to the context (player1 or player2)
// FrameNum - animate movement - all actions are from left to right in the sprite sheet
// Direction - specifies which row to pick (0 = down, 1 = up, 2 = left, 3 = right)
// X, Y - the location of the image on the screen (starts from top left corner)
// W, H - represent the size (width, height) which is calculated by multiplying each with PlayerScale
type Player struct {
	Img        *ebiten.Image
	Tag        string
	FrameNum   int
	Direction  int
	X, Y, W, H int
	VX, VY     float64
	Speed      float64
}

func (p *Player) HitBox() image.Rectangle {
	return image.Rect(p.X, p.Y, p.X+p.W, p.Y+p.H)
}

func (p *Player) SetLocation(axis string, val int) {
	if axis == u.X {
		p.X = val
	} else if axis == u.Y {
		p.Y = val
	}
}

func (p *Player) SetDelta(axis string, val float64) {
	if axis == u.X {
		p.VX = val
	} else if axis == u.Y {
		p.VY = val
	}
}

func (p *Player) DrawInteractiveSprite(screen *ebiten.Image) {
	opPlayer := &ebiten.DrawImageOptions{}
	opPlayer.GeoM.Scale(PlayerScX, PlayerScY)
	opPlayer.GeoM.Translate(float64(p.X), float64(p.Y))

	x, y := PlayerFrameOX+p.FrameNum*PlayerFrameWidth, PlayerFrameOY+p.Direction*PlayerFrameHeight
	screen.DrawImage(p.Img.SubImage(image.Rect(x, y, x+PlayerFrameWidth, y+PlayerFrameHeight)).(*ebiten.Image), opPlayer)
}

// HandleMovement - takes the vertices as params for screen cross boundary prevention
func (p *Player) HandleMovement(minX, minY, maxX, maxY int) {
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		p.Direction = 0
		if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyD) {
			p.VX = 0
		}

		p.FrameNum++
		if p.FrameNum == 7 {
			p.FrameNum = 0
		}

		p.VY = 3 * p.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		p.Direction = 1
		if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyD) {
			p.VX = 0
		}

		p.FrameNum++
		if p.FrameNum == 7 {
			p.FrameNum = 0
		}

		p.VY = -3 * p.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		p.Direction = 2
		if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyS) {
			p.VY = 0
		}

		p.FrameNum++
		if p.FrameNum == 7 {
			p.FrameNum = 0
		}

		p.Direction = 2
		p.VX = -3 * p.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		p.Direction = 3
		if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyS) {
			p.VY = 0
		}

		p.FrameNum++
		if p.FrameNum == 7 {
			p.FrameNum = 0
		}

		p.VX = 3 * p.Speed
	}

	if !ebiten.IsKeyPressed(ebiten.KeyA) && !ebiten.IsKeyPressed(ebiten.KeyD) {
		p.VX = 0
	}

	if !ebiten.IsKeyPressed(ebiten.KeyS) && !ebiten.IsKeyPressed(ebiten.KeyW) {
		p.VY = 0
	}

	// when the player is not moving
	if p.VX == 0 && p.VY == 0 {
		p.FrameNum = 0
	}

	// update the position of the player
	p.X += int(p.VX)
	p.Y += int(p.VY)

	u.BoundaryValidation(p, minX, maxX, minY, maxY)
}
