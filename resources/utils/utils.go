package utils

import (
	"bytes"
	"image"
	"io/ioutil"
	"log"
	"math/rand"
	"time"
)

const (
	ScreenWidth  = 1280
	ScreenHeight = 860
	X            = "x"
	Y            = "y"
	MinX         = "minX"
	MinY         = "minY"
	MaxX         = "maxX"
	MaxY         = "maxY"
)

func LoadSpriteImage(path string) image.Image {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}
	img, _, err := image.Decode(bytes.NewReader(file))
	if err != nil {
		log.Fatalln(err)
	}
	return img
}

// GenerateRandomLocation - generate a random location for X and y
// (x is between 0 and ScreenWidth) and (y is between 0 and ScreenHeight)
func GenerateRandomLocation(minX, maxX, minY, maxY float64) (float64, float64) {
	rand.Seed(time.Now().UnixNano())
	return (rand.Float64() * (maxX - minX)) + minX, (rand.Float64() * (maxY - minY)) + minY
}

// IsCollision - returns true if rectangle images overlap in any way
func IsCollision(x1, y1, w1, h1, x2, y2, w2, h2 float64) bool {
	return x1 < x2+w2 && x1+w1 > x2 && y1 < y2+h2 && h1+y1 > y2
}

// HitBox - generate the shape's hitbox (minX, maxX, minY, maxY)
func HitBox(x, y, w, h float64) map[string]float64 {
	return map[string]float64{
		"minX": x,
		"maxX": x + w,
		"minY": y,
		"maxY": y + h,
	}
}

// BoundaryValidation - prevents characters to move out of the view, in any of the 4 directions
func BoundaryValidation(i interface{}, minX, maxX, minY, maxY float64) {
	switch i.(type) {
	case InteractiveSprite:
		img := i.(InteractiveSprite)

		locX, locY := img.GetLocations()
		w, h := img.GetSize()

		if locX <= minX {
			img.SetLocation(X, minX)
			img.SetDelta(X, 0)
		}
		if locX >= maxX-w {
			img.SetLocation(X, maxX-w)
			img.SetDelta(X, 0)
		}
		if locY <= minY {
			img.SetLocation(Y, minY)
			img.SetDelta(Y, 0)
		}
		if locY >= maxY-h {
			img.SetLocation(Y, maxY-h)
			img.SetDelta(Y, 0)
		}
	}
}
