package data

import (
	u "games-ebiten/utils"
)

// Translation - used for acquiring the right card index while getting the card SubImage from the Image
// DraggedCard - used to force the hovered card to overlap other images, while dragged
// CardRanks - smallest is "Ace"(0), while highest is "King"(12)
var (
	DraggedCard interface{}
	Translation = map[string]map[int]string{
		u.PixelatedTheme: {
			1: "2", 2: "3", 3: "4", 4: "5", 5: "6", 6: "7",
			7: "8", 8: "9", 9: "10", 10: "J", 11: "Q", 12: "K", 13: "A",
		},
		u.ClassicTheme: {
			0: "A", 1: "2", 2: "3", 3: "4", 4: "5", 5: "6", 6: "7",
			7: "8", 8: "9", 9: "10", 10: "J", 11: "Q", 12: "K",
		},
		//u.AbstractTheme: {
		//	0: "2", 1: "3", 2: "4", 3: "5", 4: "6", 5: "7",
		//	6: "8", 7: "9", 8: "10", 9: "Jack", 10: "Queen", 11: "King", 12: "Ace",
		//},
	}
	CardRanks = map[string]int{
		u.Ace:   0,
		"2":     1,
		"3":     2,
		"4":     3,
		"5":     4,
		"6":     5,
		"7":     6,
		"8":     7,
		"9":     8,
		"10":    9,
		u.Jack:  10,
		u.Queen: 11,
		u.King:  12,
	}
)
