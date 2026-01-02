package converter

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/entity"
	metadatav1 "github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking/metadata/api"
)

// ConvertSessionEntityToCacheSessionEntity converts provided entity.SessionEntity
// to dto.CacheSessionEntity instance.
func ConvertSessionEntityToCacheSessionEntity(
	input *entity.SessionEntity) dto.CacheSessionEntity {

	return dto.CacheSessionEntity{
		ID:      input.ID,
		Seed:    input.Seed,
		Name:    input.Name,
		Started: input.Started,
	}
}

// ConvertLobbyEntityToCacheMetadataEntity converts provided entity.LobbyEntity
// array to an array of dto.CacheMetadataEntity instances.
func ConvertLobbyEntityToCacheMetadataEntity(
	input1 []*entity.LobbyEntity, input2 []*entity.InventoryEntity) []*dto.CacheMetadataEntity {
	var output []*dto.CacheMetadataEntity

	for _, lobby := range input1 {
		var inventory []dto.CacheInventoryEntity

		for _, item := range input2 {
			if item.LobbyID == lobby.ID {
				inventory = append(inventory, dto.CacheInventoryEntity{
					ID:   item.ID,
					Name: item.Name,
				})
			}
		}

		output = append(output, &dto.CacheMetadataEntity{
			LobbyID:        lobby.ID,
			SessionID:      lobby.SessionID,
			PositionX:      lobby.PositionX,
			PositionY:      lobby.PositionY,
			PositionStatic: lobby.PositionStatic,
			Skin:           uint64(lobby.Skin),
			Health:         uint64(lobby.Health),
			Active:         lobby.Active,
			Eliminated:     lobby.Eliminated,
			Host:           lobby.Host,
			Inventory:      inventory,
		})
	}

	return output
}

// ConvertCacheInventoryEntityToInventory converts provided dto.CacheInventoryEntity
// array to an array of metadatav1.Inventory instances.
func ConvertCacheInventoryEntityToInventory(
	input []dto.CacheInventoryEntity) []*metadatav1.Inventory {
	var output []*metadatav1.Inventory

	for _, inventory := range input {
		output = append(output, &metadatav1.Inventory{
			InventoryId: inventory.ID,
			Name:        inventory.Name,
		})
	}

	return output
}
