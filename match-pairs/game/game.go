package game

import (
	u "games-ebiten/resources/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	_ "image/jpeg"
)

var (
	iconsSprite = ebiten.NewImageFromImage(u.LoadSpriteImage("resources/images/misc/icons.jpeg"))
	FrOX        = 19
	FrOY        = 41
	FrW         = 34
	FrH         = 34
)

type Game struct {
	AllIcons []*Icon
}

func NewGame() *Game {
	return &Game{
		AllIcons: GenerateAvailableCards(),
	}
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	count := 0
	for i := 0; i < 10; i++ {
		for j := 0; j < 18; j++ {
			icon := g.AllIcons[count]
			icon.Y = (icon.H+50)*float64(i) + 25
			icon.X = (icon.W+50)*float64(j) + 50

			icon.DrawIcon(screen)
			count++
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return u.ScreenWidth, u.ScreenHeight
}

func GenerateAvailableCards() []*Icon {
	cards := make([]*Icon, 0, 300)

	for row := 0; row < 15; row++ {
		for col := 0; col < 15; col++ {
			x, y := FrOX+row*FrW, FrOY+col*FrH

			c := &Icon{
				Img: iconsSprite.SubImage(image.Rect(x, y, x+FrW, y+FrH)).(*ebiten.Image),
				W:   float64(FrW),
				H:   float64(FrH),
			}
			cards = append(cards, c)
		}
	}
	return cards
}
