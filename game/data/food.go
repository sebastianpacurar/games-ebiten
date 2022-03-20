package data

type Food struct {
	LocX, LocY  float64
	Size        float64
	HitBox      map[string]float64
	IsDisplayed bool
}
