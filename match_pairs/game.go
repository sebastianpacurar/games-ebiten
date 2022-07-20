package match_pairs

import (
	res "games-ebiten/resources"
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	_ "image/jpeg"
	"math/rand"
	"strings"
)

var (
	iconsSprite = ebiten.NewImageFromImage(res.LoadSpriteImage("resources/assets/misc/icons.jpeg"))
	FrOX        = 19
	FrOY        = 41
	FrW         = 34
	FrH         = 34
	ScX         = 3.0
	ScY         = 3.0
)

// Game
// Icons - represents all the generated icons (around 300)
// IconPairs - all the generated pairs for the current game
// RevealedIDs - stores uuid of the current revealed cards
// PairNum - the number of icons in a pair
type Game struct {
	Icons       []*Icon
	IconPairs   []*Icon
	RevealedIDs map[string]int
	PairNum     int
	Rows, Cols  int
}

func NewGame() *Game {
	g := &Game{
		Rows:    4,
		Cols:    4,
		PairNum: 2,
		Icons:   GenerateAvailableIcons(),
	}
	g.RevealedIDs = make(map[string]int, (g.Rows*g.Cols)/2)
	g.GeneratePairs(g.PairNum)
	return g
}

func (g *Game) Update() error {
	g.HandleRevealLogic()

	for _, v := range g.IconPairs {
		if res.IsAreaHovered(v) {
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				if !v.Revealed() {
					v.SetRevealed(true)
					g.RevealedIDs[v.ID]++
					if g.RevealedIDs[v.ID] == g.PairNum {
						g.RemovePairs(v.ID)
					}
				} else {
					v.SetRevealed(false)
					g.RevealedIDs[v.ID]--
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
	return res.ScreenWidth, res.ScreenHeight
}

// GeneratePairs - generates the given pairs. if smaller than 2 or bigger than 3, then default to 2
func (g *Game) GeneratePairs(p int) {
	count := 0

	// append only the first half of all the icons with the correct ID
	for i := 0; i < (g.Rows*g.Cols)/p; i++ {
		id := strings.Split(uuid.New().String(), "-")[1]
		g.RevealedIDs[id] = 0
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
	quads := res.GridQuadrants(g.Rows, g.Cols)
	for i := 0; i < g.Rows; i++ {
		for j := 0; j < g.Cols; j++ {
			icon := g.IconPairs[count]

			icon.X = res.CenterOnX(icon.W, quads[i][j])
			icon.Y = res.CenterOnY(icon.W, quads[i][j])
			count++
		}
	}
}

// HandleRevealLogic - handles basic reveal, hide and match (paired icons removal) functionality
func (g *Game) HandleRevealLogic() {
	if count, ids := g.CountOfUnmatched(); count == 2 {
		for _, id := range ids {
			for _, ic := range g.IconPairs {
				if id == ic.ID {
					ic.SetRevealed(false)
				}
			}
			g.RevealedIDs[id] = 0
		}
	}
}

// CountOfUnmatched - returns the count to be used as condition, and the list of matching ids
func (g *Game) CountOfUnmatched() (int, []string) {
	count := 0
	ids := make([]string, 0)
	for k, v := range g.RevealedIDs {
		if v == 1 {
			count++
			ids = append(ids, k)
		}
	}
	return count, ids
}

// RemovePairs - set the pair icon's removedState to true
func (g *Game) RemovePairs(byId string) {
	for _, i := range g.IconPairs {
		if i.ID == byId {
			i.SetRemoved(true)
		}
	}
	delete(g.RevealedIDs, byId)
}
