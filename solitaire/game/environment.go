package game

import (
	u "games-ebiten/resources/utils"
	"github.com/hajimehoshi/ebiten/v2"
)

type (
	Environment struct {
		BgImg        *ebiten.Image
		EmptySlotImg *ebiten.Image
		Columns      []CardColumn
		CardStores   []CardStore
		DrawnCardsSlot
		DrawCardSlot
		CardsVSpacer float64
		SpacerV      float64
		SpacerH      float64
	}

	DrawCardSlot struct {
		X, Y, W, H   float64
		GreenSlotImg *ebiten.Image
		RedSlotImg   *ebiten.Image
		Cards        []*Card
		IsGreen      bool
	}

	CardStore struct {
		X, Y, W, H float64
		Cards      []*Card
	}

	DrawnCardsSlot struct {
		Cards []*Card
	}

	CardColumn struct {
		X, Y, W, H float64
		Cards      []*Card
	}
)

// IsDrawCardHovered - returns true if the DrawCardSlot is hovered. Applies only when there are no cards in stack
func (e *Environment) IsDrawCardHovered(cx, cy int, th *Theme) bool {
	x, y, w, h := e.DrawCardSlot.X, e.DrawCardSlot.Y, e.DrawCardSlot.W, e.DrawCardSlot.H
	return int(x) <= cx && cx < int(x+w) && int(y) <= cy && cy < int(y+h)
}

func (cs *CardStore) GetStoreGeomData() (float64, float64, float64, float64) {
	return cs.X, cs.Y, cs.W, cs.H
}

func (cl *CardColumn) GetColumnGeoMData() (float64, float64, float64, float64) {
	return cl.X, cl.Y, cl.W, cl.H
}

func (cl *CardColumn) GetCountOfHidden() int {
	count := 0
	for _, v := range cl.Cards {
		if !v.IsRevealed {
			count++
		}
	}
	return count
}

func (e *Environment) DrawEnvironment(screen *ebiten.Image, th *Theme) {
	// Draw the BG Image
	opBg := &ebiten.DrawImageOptions{}
	opBg.GeoM.Scale(50, 50)
	screen.DrawImage(e.BgImg, opBg)

	// Draw the Draw Stack and the Drawn Stack
	for i := 1; i <= 2; i++ {
		var img *ebiten.Image
		opCardStack := &ebiten.DrawImageOptions{}
		opCardStack.GeoM.Scale(th.ScaleValue[th.Active][u.X], th.ScaleValue[th.Active][u.Y])

		x := (float64(th.FrontFaceFrameData[th.Active][u.FrW]) + e.SpacerH) * float64(i)
		y := e.SpacerV
		opCardStack.GeoM.Translate(x, y)

		if i == 2 {
			img = e.EmptySlotImg
		} else {
			img = e.GreenSlotImg
		}
		screen.DrawImage(img, opCardStack)
	}

	// Draw the Card Stores
	for i := 0; i < 4; i++ {
		opStackSlot := &ebiten.DrawImageOptions{}
		opStackSlot.GeoM.Scale(th.ScaleValue[th.Active][u.X], th.ScaleValue[th.Active][u.Y])

		// align Stacked Cards with the left columns of the Column slots
		x := (float64(th.FrontFaceFrameData[th.Active][u.FrW]) + e.SpacerH) * (float64(i) + 4)

		opStackSlot.GeoM.Translate(x, e.SpacerV)
		screen.DrawImage(e.EmptySlotImg, opStackSlot)
	}

	// Draw the colum slots
	for i := 1; i <= 7; i++ {
		opColumnSlot := &ebiten.DrawImageOptions{}
		opColumnSlot.GeoM.Scale(th.ScaleValue[th.Active][u.X], th.ScaleValue[th.Active][u.Y])

		x := (float64(th.FrontFaceFrameData[th.Active][u.FrW]) + e.SpacerH) * float64(i)
		y := float64(u.ScreenHeight / 3)

		opColumnSlot.GeoM.Translate(x, y)
		screen.DrawImage(e.EmptySlotImg, opColumnSlot)
	}
}
