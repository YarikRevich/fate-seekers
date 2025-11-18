package gamepad

import (
	"math"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
	"github.com/hajimehoshi/ebiten/v2"
)

// Represents configured stick deadzone.
const DEADZONE = 0.25

// getInvertedGamepadStickDirection retrieves inverted gamepad stick direction.
func getInvertedGamepadStickDirection(id ebiten.GamepadID, x, y float64) dto.Direction {
	result := dto.DirNone

	if math.Abs(x) < DEADZONE && math.Abs(y) < DEADZONE {
		return result
	}

	angle := math.Atan2(y, x)
	pi := math.Pi

	sector := pi / 8.0

	if angle > -sector && angle <= sector {
		result = dto.DirRight
	} else if angle > sector && angle <= 3*sector {
		result = dto.DirUpRight
	} else if angle > 3*sector && angle <= 5*sector {
		result = dto.DirUp
	} else if angle > 5*sector && angle <= 7*sector {
		result = dto.DirUpLeft
	} else if angle > 7*sector || angle <= -7*sector {
		result = dto.DirLeft
	} else if angle > -7*sector && angle <= -5*sector {
		result = dto.DirDownLeft
	} else if angle > -5*sector && angle <= -3*sector {
		result = dto.DirDown
	} else if angle > -3*sector && angle <= -sector {
		result = dto.DirDownRight
	}

	return result
}

// getNormalGamepadStickDirection retrieves normal gamepad stick direction.
func getNormalGamepadStickDirection(id ebiten.GamepadID, x, y float64) dto.Direction {
	result := dto.DirNone

	if math.Abs(x) < DEADZONE && math.Abs(y) < DEADZONE {
		return result
	}

	angle := math.Atan2(y, x)
	pi := math.Pi

	sector := pi / 8.0

	if angle > -sector && angle <= sector {
		result = dto.DirRight
	} else if angle > sector && angle <= 3*sector {
		result = dto.DirDownRight
	} else if angle > 3*sector && angle <= 5*sector {
		result = dto.DirDown
	} else if angle > 5*sector && angle <= 7*sector {
		result = dto.DirDownLeft
	} else if angle > 7*sector || angle <= -7*sector {
		result = dto.DirLeft
	} else if angle > -7*sector && angle <= -5*sector {
		result = dto.DirUpLeft
	} else if angle > -5*sector && angle <= -3*sector {
		result = dto.DirUp
	} else if angle > -3*sector && angle <= -sector {
		result = dto.DirUpRight
	}

	return result
}

// GetGamepadLeftStickDirection retrieves left gamepad stick direction.
func GetGamepadLeftStickDirection(id ebiten.GamepadID) dto.Direction {
	x := ebiten.GamepadAxisValue(id, 0)
	y := ebiten.GamepadAxisValue(id, 1)

	return getNormalGamepadStickDirection(id, x, y)
}

// GetGamepadRightStickDirection retrieves right gamepad stick direction.
func GetGamepadRightStickDirection(id ebiten.GamepadID) dto.Direction {
	x := ebiten.GamepadAxisValue(id, 2)
	y := ebiten.GamepadAxisValue(id, 5)

	return getInvertedGamepadStickDirection(id, x, y)
}
