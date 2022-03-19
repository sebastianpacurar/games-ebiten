package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"log"
)

const (
	screenWidth  = 1280
	screenHeight = 960
	MinX         = "minX"
	MinY         = "minY"
	MaxX         = "maxX"
	MaxY         = "maxY"
)

type Game struct {
	bounds map[string]float64

	Player
}

type Player struct {
	incrementX float64
	incrementY float64
	posx       float64
	posy       float64
	size       float64
}

func NewGame() *Game {
	return &Game{
		bounds: map[string]float64{MinX: 0, MaxX: screenWidth, MinY: 0, MaxY: screenHeight},
		Player: Player{
			size: 100,
			posx: float64(screenWidth / 2),
			posy: float64(screenHeight / 2),
		},
	}
}

func (g *Game) handleMovement() {
	switch {
	case ebiten.IsKeyPressed(ebiten.KeyD):
		if g.posx <= g.bounds[MaxX] {
			g.incrementX += 0.2
		}
	case ebiten.IsKeyPressed(ebiten.KeyW):
		if g.posy >= g.bounds[MinY] {
			g.incrementY -= 0.2
		}
	case ebiten.IsKeyPressed(ebiten.KeyS):
		if g.posy <= g.bounds[MaxY] {
			g.incrementY += 0.2
		}
	case ebiten.IsKeyPressed(ebiten.KeyA):
		if g.posx >= g.bounds[MinX] {
			g.incrementX -= 0.2
		}
	default:
		if g.incrementX > 0 {
			g.incrementX -= 0.1
			if g.incrementX < 0.00001 {
				g.incrementX = 0
			}
		} else if g.incrementX < 0 {
			g.incrementX += 0.1
			if g.incrementX > 0 {
				g.incrementX = 0
			}
		}

		if g.incrementY > 0 {
			g.incrementY -= 0.1
			if g.incrementY < 0.00001 {
				g.incrementY = 0
			}
		} else if g.incrementY < 0 {
			g.incrementY += 0.1
			if g.incrementY > 0 {
				g.incrementY = 0
			}
		}
	}
}

func (g *Game) Update() error {
	g.handleMovement()
	g.posx += g.incrementX
	g.posy += g.incrementY

	if g.posx <= g.bounds[MinX] {
		g.posx = g.bounds[MinX]
		g.incrementX = 0
	}
	if g.posx >= g.bounds[MaxX]-g.size {
		g.posx = g.bounds[MaxX] - g.size
		g.incrementX = 0
	}
	if g.posy <= g.bounds[MinY] {
		g.posy = g.bounds[MinY]
		g.incrementY = 0
	}
	if g.posy >= g.bounds[MaxY]-g.size {
		g.posy = g.bounds[MaxY] - g.size
		g.incrementY = 0
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)

	rect := ebiten.NewImage(int(g.size), int(g.size))
	rect.Fill(color.NRGBA{G: 255, A: 255})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.posx, g.posy)
	screen.DrawImage(rect, op)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("game")

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
