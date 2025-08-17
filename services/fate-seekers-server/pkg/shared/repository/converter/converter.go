package converter

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/entity"
)

// ConvertLobbyEntityToCacheMetadataEntity converts provided entity.LobbyEntity
// array to an array of dto.CacheMetadataEntity instances.
func ConvertLobbyEntityToCacheMetadataEntity(
	input []*entity.LobbyEntity) []dto.CacheMetadataEntity {
	var output []dto.CacheMetadataEntity

	for _, lobby := range input {
		output = append(output, dto.CacheMetadataEntity{
			SessionID:  lobby.SessionID,
			PositionX:  lobby.PositionX,
			PositionY:  lobby.PositionY,
			Skin:       uint64(lobby.Skin),
			Health:     uint64(lobby.Health),
			Active:     lobby.Active,
			Eliminated: lobby.Eliminated,
			Host:       lobby.Host,
		})
	}

	return output
}
