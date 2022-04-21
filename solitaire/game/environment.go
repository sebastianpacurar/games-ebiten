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
		DrawnSlot
		DrawCardSlot
		SpacerV int
		SpacerH int
	}

	DrawCardSlot struct {
		GreenSlotImg *ebiten.Image
		RedSlotImg   *ebiten.Image
		Cards        []*Card
		IsEmpty      bool
		IsGreen      bool
		TriesCount   bool
	}

	DrawnSlot struct {
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

func (e *Environment) DrawEnvironment(screen *ebiten.Image, th *Theme) {
	// Draw the BG Image
	opBg := &ebiten.DrawImageOptions{}
	opBg.GeoM.Scale(50, 50)
	screen.DrawImage(e.BgImg, opBg)

	// Draw the Draw Stack and Placeholder
	for i := 1; i <= 2; i++ {
		var img *ebiten.Image
		opCardStack := &ebiten.DrawImageOptions{}
		opCardStack.GeoM.Scale(th.ScaleValue[th.Active][u.X], th.ScaleValue[th.Active][u.Y])

		x := float64(u.ScreenWidth/7 + (th.FrontFaceFrameData[th.Active][u.FrW]+e.SpacerH)*i)
		y := float64(e.SpacerV)
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
		x := float64(u.ScreenWidth)/2.05 + float64((th.FrontFaceFrameData[th.Active][u.FrW]+e.SpacerH)*i)

		opStackSlot.GeoM.Translate(x, float64(e.SpacerV))
		screen.DrawImage(e.EmptySlotImg, opStackSlot)
	}

	// Draw the colum slots
	for i := 1; i <= 7; i++ {
		opColumnSlot := &ebiten.DrawImageOptions{}
		opColumnSlot.GeoM.Scale(th.ScaleValue[th.Active][u.X], th.ScaleValue[th.Active][u.Y])

		x := float64(u.ScreenWidth/7 + (th.FrontFaceFrameData[th.Active][u.FrW]+e.SpacerH)*i)
		y := float64(u.ScreenHeight / 3)

		opColumnSlot.GeoM.Translate(x, y)
		screen.DrawImage(e.EmptySlotImg, opColumnSlot)
	}
}
