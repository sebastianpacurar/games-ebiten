package data

type Food struct {
	PosX, PosY  float64
	Size        float64
	HitBox      map[string]float64
	IsDisplayed bool
}

// GetHitBox - generate the player's boundaries for minX, maxX and minY, maxY
func (f *Food) GetHitBox() map[string]float64 {
	return map[string]float64{
		"minX": f.PosX,
		"maxX": f.PosX + f.Size,
		"minY": f.PosY,
		"maxY": f.PosY + f.Size,
	}
}
