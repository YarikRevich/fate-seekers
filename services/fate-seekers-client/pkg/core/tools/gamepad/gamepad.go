package gamepad

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const deadzone = 0.25

type Direction int

const (
	DirNone Direction = iota
	DirUp
	DirDown
	DirLeft
	DirRight
	DirUpLeft
	DirUpRight
	DirDownLeft
	DirDownRight
)

// getGamepadStickDirection retrieves gamepad stick direction.
func getGamepadStickDirection(id ebiten.GamepadID, x, y float64) Direction {
	result := DirNone

	if math.Abs(x) < deadzone && math.Abs(y) < deadzone {
		return result
	}

	angle := math.Atan2(y, x)
	pi := math.Pi

	sector := pi / 8.0

	if angle > -sector && angle <= sector {
		result = DirRight
	} else if angle > sector && angle <= 3*sector {
		result = DirUpRight
	} else if angle > 3*sector && angle <= 5*sector {
		result = DirUp
	} else if angle > 5*sector && angle <= 7*sector {
		result = DirUpLeft
	} else if angle > 7*sector || angle <= -7*sector {
		result = DirLeft
	} else if angle > -7*sector && angle <= -5*sector {
		result = DirDownLeft
	} else if angle > -5*sector && angle <= -3*sector {
		result = DirDown
	} else if angle > -3*sector && angle <= -sector {
		result = DirDownRight
	}

	return result
}

// GetGamepadLeftStickDirection retrieves left gamepad stick direction.
func GetGamepadLeftStickDirection(id ebiten.GamepadID) Direction {
	x := ebiten.GamepadAxisValue(id, 0)
	y := ebiten.GamepadAxisValue(id, 1)

	return getGamepadStickDirection(id, x, y)
}

// GetGamepadRightStickDirection retrieves right gamepad stick direction.
func GetGamepadRightStickDirection(id ebiten.GamepadID) Direction {
	x := ebiten.GamepadAxisValue(id, 2)
	y := ebiten.GamepadAxisValue(id, 5)

	return getGamepadStickDirection(id, x, y)
}
