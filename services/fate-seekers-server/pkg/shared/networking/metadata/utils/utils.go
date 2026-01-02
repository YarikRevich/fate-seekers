package utils

import (
	"math/rand"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/dto"
	metadatav1 "github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking/metadata/api"
	"github.com/google/uuid"
)

// Describes all the available chest configurations.
const (
	CHEST_STANDARD_TYPE = "standard"

	CHEST_ITEM_HEALTH_PACK_TYPE = "standard_health_pack"

	CHEST_ITEM_LETTER_TYPE = "letter"
)

// Describes all the available health pack configurations.
const (
	HEALTH_PACK_FROG_TYPE = "frog"
)

// Describes sets of available chest, chest item and health pack types.
var (
	CHEST_TYPES = []string{
		CHEST_STANDARD_TYPE,
	}

	CHEST_ITEM_TYPES = []string{
		CHEST_ITEM_HEALTH_PACK_TYPE,
		CHEST_ITEM_LETTER_TYPE,
	}

	HEALTH_PACK_TYPES = []string{
		HEALTH_PACK_FROG_TYPE,
	}
)

// GenerateChests performs chests generation according to the provided available
// positions and seed used for deterministic random selection.
func GenerateChests(positions []*metadatav1.Position, seed int64) []dto.GeneratedChest {
	selector := rand.New(rand.NewSource(seed))

	availablePositions := make([]*metadatav1.Position, len(positions))
	copy(availablePositions, positions)

	selector.Shuffle(len(availablePositions), func(i, j int) {
		availablePositions[i], availablePositions[j] = availablePositions[j], availablePositions[i]
	})

	chestsAmount := config.GetOperationMaxChestsAmount()
	if len(availablePositions) < chestsAmount {
		chestsAmount = len(availablePositions)
	}

	result := make([]dto.GeneratedChest, 0, chestsAmount)

	for i := 0; i < chestsAmount; i++ {
		chestItemsAmount := selector.Intn(config.MAX_CHEST_ITEMS_PER_CHEST)

		chestItems := make([]dto.ChestItem, 0, chestItemsAmount)
		for k := 0; k < chestItemsAmount; k++ {
			chestItems = append(chestItems, dto.ChestItem{
				Instance: uuid.NewString(),
				Name:     CHEST_ITEM_TYPES[selector.Intn(len(CHEST_ITEM_TYPES))],
			})
		}

		chest := dto.GeneratedChest{
			Instance: uuid.NewString(),
			Position: dto.Position{
				X: int(availablePositions[i].X),
				Y: int(availablePositions[i].Y),
			},
			Name:       CHEST_TYPES[selector.Intn(len(CHEST_TYPES))],
			ChestItems: chestItems,
		}

		result = append(result, chest)
	}

	return result
}

// GenerateHealthPacks performs health packs generation according to the provided available
// positions and seed used for deterministic random selection.
func GenerateHealthPacks(positions []*metadatav1.Position, seed int64) []dto.GeneratedHealthPack {
	selector := rand.New(rand.NewSource(seed))

	availablePositions := make([]*metadatav1.Position, len(positions))
	copy(availablePositions, positions)

	selector.Shuffle(len(availablePositions), func(i, j int) {
		availablePositions[i], availablePositions[j] = availablePositions[j], availablePositions[i]
	})

	healthPacksAmount := config.GetOperationMaxHealthPacksAmount()
	if len(availablePositions) < healthPacksAmount {
		healthPacksAmount = len(availablePositions)
	}

	result := make([]dto.GeneratedHealthPack, 0, healthPacksAmount)

	for i := 0; i < healthPacksAmount; i++ {
		healthPack := dto.GeneratedHealthPack{
			Instance: uuid.NewString(),
			Position: dto.Position{
				X: int(availablePositions[i].X),
				Y: int(availablePositions[i].Y),
			},
			Name: HEALTH_PACK_TYPES[selector.Intn(len(HEALTH_PACK_TYPES))],
		}

		result = append(result, healthPack)
	}

	return result
}
