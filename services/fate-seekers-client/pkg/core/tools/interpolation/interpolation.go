package interpolation

import (
	"fmt"
	"math"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
)

// GetDelayedPositions calculates delayed positions between the first position
// and the second provided one.
func GetDelayedPositions(previousPosition, position dto.Position) []dto.Position {
	var (
		result      []dto.Position
		newPosition dto.Position = previousPosition
	)

	for {
		dx, dy := position.X-newPosition.X, position.Y-newPosition.Y

		dist := math.Hypot(dx, dy)
		if dist == 0 || dist <= 1 {
			fmt.Println("IN A LOOP")

			break
		}

		k := 1 / dist

		newPosition = dto.Position{X: math.Round(newPosition.X + dx*k), Y: math.Round(newPosition.Y + dy*k)}

		result = append(result, newPosition)
	}

	return result
}
