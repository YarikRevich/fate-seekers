package direction

import (
	"math"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
)

// GetAnimationDirection retrieves animation direction based on a provided position delta.
func GetAnimationDirection(previousX, previousY, x, y float64) string {
	dx := x - previousX
	dy := y - previousY

	ang := math.Atan2(dy, dx)

	step := math.Pi / 4
	idx := int(math.Round(ang / step))

	idx = ((idx % 8) + 8) % 8

	switch idx {
	case 0:
		return dto.RightMovableRotation
	case 1:
		return dto.UpRightMovableRotation
	case 2:
		return dto.UpMovableRotation
	case 3:
		return dto.UpLeftMovableRotation
	case 4:
		return dto.LeftMovableRotation
	case 5:
		return dto.DownLeftMovableRotation
	case 6:
		return dto.DownMovableRotation
	case 7:
		return dto.DownRightMovableRotation
	default:
		return ""
	}
}
