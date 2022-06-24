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
	SpacingV    = 25
	SpacingH    = 50
	FrOX        = 19
	FrOY        = 41
	FrW         = 34
	FrH         = 34
	ScX         = float64(2)
	ScY         = float64(2)
)

// Game
// Icons - represents all the generated icons (around 300)
// IconPairs - all the generated pairs for the current game
// RevealedIDs - stores uuid of the current revealed cards
// PairNum - the number of icons in a pair
type Game struct {
	Icons       []*Icon
	IconPairs   []*Icon
	RevealedIDs []string
	PairNum     int
	Rows, Cols  int
}

func NewGame() *Game {
	g := &Game{
		RevealedIDs: make([]string, 0),
		PairNum:     2,
	}
	g.Icons = GenerateAvailableIcons()
	g.GeneratePairs(g.PairNum)
	return g
}

func (g *Game) Update() error {
	g.HandleRevealLogic()

	for _, v := range g.IconPairs {
		if u.IsAreaHovered(v) {
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				if !v.IsRevealed() {
					v.SetRevealedState(true)
					g.RevealedIDs = append(g.RevealedIDs, v.ID)
				} else {
					v.SetRevealedState(false)
					g.RevealedIDs = g.RemoveRevealedID(v.ID)
				}
			}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, ic := range g.IconPairs {
		ic.DrawIcon(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return u.ScreenWidth, u.ScreenHeight
}

// GeneratePairs - generates the given pairs. if smaller than 2 or bigger than 3, then default to 2
func (g *Game) GeneratePairs(p int) {
	if p < 2 || p > 3 {
		p = 2
	}

	// 16 icons for 2 pairs and 36 icons for 3 pairs
	if p%2 == 0 {
		g.Rows = 4
		g.Cols = 4
	} else if p%3 == 0 {
		g.Rows = 6
		g.Cols = 6
	}
	maxNum := g.Rows * g.Cols
	count := 0

	// append only the first half of all the icons with the correct ID
	for i := 0; i < maxNum/p; i++ {
		id := strings.Split(uuid.New().String(), "-")[1]
		for dup := 0; dup < p; dup++ {

			// create a unique icon, add it to the pairs array, then assign the id
			g.IconPairs = append(g.IconPairs, &Icon{
				Img: g.Icons[i].Img,
				W:   g.Icons[i].W,
				H:   g.Icons[i].H,
			})
			g.IconPairs[count].ID = id
			count++
		}
	}

	//shuffle g.IconPairs array
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(g.IconPairs), func(i, j int) {
		g.IconPairs[i], g.IconPairs[j] = g.IconPairs[j], g.IconPairs[i]
	})

	g.GeneratePositions()
}

func GenerateAvailableIcons() []*Icon {
	icons := make([]*Icon, 0, 300)

	for row := 0; row < 15; row++ {
		for col := 0; col < 15; col++ {
			x, y := FrOX+row*FrW, FrOY+col*FrH

			icon := &Icon{
				Img: iconsSprite.SubImage(image.Rect(x, y, x+FrW, y+FrH)).(*ebiten.Image),
				W:   int(float64(FrW) * ScX),
				H:   int(float64(FrH) * ScY),
			}
			icons = append(icons, icon)
		}
	}
	return icons
}

// GeneratePositions - add the x, y coords for every icon
func (g *Game) GeneratePositions() {
	count := 0
	for i := 0; i < g.Rows; i++ {
		for j := 0; j < g.Cols; j++ {
			icon := g.IconPairs[count]

			icon.X = (icon.W+50)*j + SpacingV
			icon.Y = (icon.H+50)*i + SpacingH
			count++
		}
	}
}

func (g *Game) HandleRevealLogic() {
	if len(g.RevealedIDs) > 1 {
		allTrue := true

		firstId := g.RevealedIDs[0]
		for _, id := range g.RevealedIDs[1:] {
			if firstId != id {
				allTrue = false
				break
			}
		}

		if allTrue {
			// if the pair is a match, set all icons to removed and clear array
			for _, ic := range g.IconPairs {
				if firstId == ic.ID {
					ic.SetRemovedState(true)
				}
			}
			g.RevealedIDs = nil
		} else {
			// if the pair isn't a match, set all icons to hidden and clear array
			for _, ic := range g.IconPairs {
				ic.SetRevealedState(false)
			}
			g.RevealedIDs = nil
		}
	}
}

// RemoveRevealedID - remove the uid from RevealedIDs, based on the uid string
func (g *Game) RemoveRevealedID(val string) []string {
	res := make([]string, 0)
	i := 0
	if len(g.RevealedIDs) > 1 {
		for _, v := range g.IconPairs {
			if v.ID == val {
				break
			}
			i++
		}
		res = append(g.RevealedIDs[:i], g.RevealedIDs[i+1:]...)
	}
	return res
}
