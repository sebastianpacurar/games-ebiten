package utils

import "github.com/hajimehoshi/ebiten/v2"

// InteractiveSprite - represents any type of image which can be animated and holds several states.
type (
	InteractiveSprite interface {

		// GetLocations - returns (LX, LY)
		GetLocations() (float64, float64)

		// GetSize - returns (width, height)
		GetSize() (float64, float64)

		// GetScaleVal - returns (ScaleX, ScaleY)
		GetScaleVal() (float64, float64)

		// GetFramePosition - returns (FrameOX, FrameWidth, FrameOY, FrameHeight)
		GetFramePosition() (int, int, int, int)

		// GetFrameNum - returns the current column from the current sprite sheet image row
		GetFrameNum() int

		// GetDirection - returns the direction the sprite is currently facing
		GetDirection() int

		// GetImg - returns the Img of the sprite
		GetImg() *ebiten.Image

		// SetLocation - updates LX or LY, based on X or Y axis
		SetLocation(string, float64)

		// SetDelta - updates DX or DY, based on X or Y axis
		SetDelta(string, float64)

		// DrawImg - Draws the image (shorthand for writing it in Draw() directly)
		DrawImg(image *ebiten.Image)
	}

	StaticSprite interface {
	}
)
