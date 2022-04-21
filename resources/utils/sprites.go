package utils

import "github.com/hajimehoshi/ebiten/v2"

type (
	// InteractiveSprite - is implemented by structs which hold multiple states and which tend to be dynamic
	InteractiveSprite interface {

		// GetLocations - returns (LX, LY)
		GetLocations() (float64, float64)

		// GetSize - returns (width, height)
		GetSize() (float64, float64)

		// GetScaleVal - returns (ScaleX, ScaleY)
		GetScaleVal() (float64, float64)

		// GetFrameInfo - returns (FrOX, FrOY, FrameWidth, FrameHeight)
		GetFrameInfo() (int, int, int, int)

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

		// DrawInteractiveSprite - Draws the image
		DrawInteractiveSprite(*ebiten.Image)

		GetHitBox() (float64, float64, float64, float64)
	}

	// StaticSprite - is implemented by structs which don't contain many states to get updated frequently
	StaticSprite interface {

		// GetLocations - returns (LX, LY)
		GetLocations() (float64, float64)

		// GetSize - returns (width, height)
		GetSize() (float64, float64)

		// GetFrameInfo - returns (FrOX, FrameWidth, FrOY, FrameHeight)
		GetFrameInfo() (int, int, int, int)

		// GetImg - returns the Img of the sprite
		GetImg() *ebiten.Image

		// SetLocation - updates LX or LY, based on X or Y axis
		SetLocation(string, float64)

		// DrawStaticSprite - Draws the image
		DrawStaticSprite(*ebiten.Image)

		GetHitBox() (float64, float64, float64, float64)
	}
)
