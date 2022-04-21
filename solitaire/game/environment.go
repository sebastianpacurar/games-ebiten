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
		Stores       []CardStore
		DrawnCardsSlot
		DrawCardSlot
		SpacerV float64
		SpacerH float64
	}

	DrawCardSlot struct {
		X            float64
		Y            float64
		W            float64
		H            float64
		GreenSlotImg *ebiten.Image
		RedSlotImg   *ebiten.Image
		Cards        []*Card
		IsEmpty      bool
		IsGreen      bool
		TriesCount   bool
	}

	DrawnCardsSlot struct {
		Cards   []*Card
		IsEmpty bool
	}

	CardStore struct {
		Cards   []*Card
		IsEmpty bool
	}

	CardColumn struct {
		Cards   []*Card
		IsEmpty bool
	}
)

// IsDrawCardHovered - returns true if the DrawCardSlot is hovered. Applies only when there are no cards in stack
func (e *Environment) IsDrawCardHovered(cx, cy int, th *Theme) bool {
	x, y, w, h := e.DrawCardSlot.X, e.DrawCardSlot.Y, e.DrawCardSlot.W, e.DrawCardSlot.H
	return int(x) <= cx && cx < int(x+w) && int(y) <= cy && cy < int(y+h)
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

	// Draw the to-fill stack slot area
	for i := 0; i < 4; i++ {
		opStackSlot := &ebiten.DrawImageOptions{}
		opStackSlot.GeoM.Scale(th.ScaleValue[th.Active][u.X], th.ScaleValue[th.Active][u.Y])

		// align StackedCards with the last column from the column Stack Slot
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
