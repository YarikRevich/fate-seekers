package utils

import (
	"math/rand"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/dto"
)

const (
	// Represents value of max amount of point regeneration requests
	maxPointsGenerationRetries = 1000
)

/// TODO: take into account exlusion points on the map.

// generateRandomPosition represents random positions generator implementation.
func generateRandomPosition(width, height, amount, radius int, seed int64) []dto.Position {
	randInstance := rand.New(rand.NewSource(seed))

	var (
		result []dto.Position
		tries  int
		found  bool
	)

	radius2 := radius * radius

	for len(result) < amount && tries < maxPointsGenerationRetries {
		tries++

		position := dto.Position{
			X: randInstance.Intn(width),
			Y: randInstance.Intn(height),
		}

		found = true

		for _, value := range result {
			dx := position.X - value.X
			dy := position.Y - value.Y
			if dx*dx+dy*dy < radius2 {
				found = false

				break
			}
		}

		if found {
			result = append(result, position)
		}
	}
	return result
}

// GenerateChestPositions performs chests positions generation.
func GenerateChestPositions(seed int64) []dto.Position {
	return generateRandomPosition(
		config.GetOperationGenerationAreaWidth(),
		config.GetOperationGenerationAreaHeight(),
		config.GetOperationMaxChestsAmount(),
		config.GetOperationGenerationMaxRadius(),
		seed,
	)
}

// GenerateHealthPackPositions performs health pack positions generation.
func GenerateHealthPackPositions(seed int64) []dto.Position {
	return generateRandomPosition(
		config.GetOperationGenerationAreaWidth(),
		config.GetOperationGenerationAreaHeight(),
		config.GetOperationMaxChestsAmount(),
		config.GetOperationGenerationMaxRadius(),
		seed,
	)
}
