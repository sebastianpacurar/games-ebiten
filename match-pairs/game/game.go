package game

import (
	u "games-ebiten/resources/utils"
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	_ "image/jpeg"
	"math/rand"
	"strings"
	"time"
)

var (
	iconsSprite = ebiten.NewImageFromImage(u.LoadSpriteImage("resources/images/misc/icons.jpeg"))
	FrOX        = 19
	FrOY        = 41
	FrW         = 34
	FrH         = 34
	SpacingV    = 25
	SpacingH    = 50
)

type Game struct {
	Icons      []*Icon
	IconPairs  []*Icon
	Rows, Cols int
}

func NewGame() *Game {
	g := &Game{
		Rows: 4,
		Cols: 4,
	}
	g.Icons = GenerateAvailableIcons()
	g.GeneratePairs(2)
	return g
}

func (g *Game) Update() error {
	cx, cy := ebiten.CursorPosition()
	for i := range g.IconPairs {
		if u.IsAreaHovered(g.IconPairs[i], cx, cy) {
			if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
				g.IconPairs[i].IsRevealed = !g.IconPairs[i].IsRevealed
			}
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	count := 0
	for i := 0; i < g.Rows; i++ {
		for j := 0; j < g.Cols; j++ {
			icon := g.IconPairs[count]

			icon.DrawIcon(screen)
			count++
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return u.ScreenWidth, u.ScreenHeight
}

//TODO: broken
func (g *Game) GeneratePairs(p int) {
	var id string
	maxNum := g.Rows * g.Cols

	rand.Seed(time.Now().UnixNano())

	// shuffle g.Icons array
	rand.Shuffle(len(g.Icons), func(i, j int) {
		g.Icons[i], g.Icons[j] = g.Icons[j], g.Icons[i]
	})

	// iterate p number of times (pair elements in a pair)
	for j := 0; j < p; j++ {
		// append only the first half of all the icons with the correct ID
		for i := 0; i < maxNum/2; i++ {
			if i%p == 0 {
				id = strings.Split(uuid.New().String(), "-")[1]
			}
			for dup := 0; dup < p; dup++ {
				g.IconPairs = append(g.IconPairs, g.Icons[i])
				g.IconPairs[i].ID = id
			}
		}
	}
	// shuffle g.IconPairs array
	rand.Shuffle(len(g.IconPairs), func(i, j int) {
		g.IconPairs[i], g.IconPairs[j] = g.IconPairs[j], g.IconPairs[i]
	})

}

//TODO: broken
func GenerateAvailableIcons() []*Icon {
	icons := make([]*Icon, 0, 300)

	for row := 0; row < 15; row++ {
		for col := 0; col < 15; col++ {
			x, y := FrOX+row*FrW, FrOY+col*FrH

			icon := &Icon{
				Img: iconsSprite.SubImage(image.Rect(x, y, x+FrW, y+FrH)).(*ebiten.Image),
				W:   float64(FrW),
				H:   float64(FrH),
			}
			icons = append(icons, icon)
		}
	}
	return icons
}
