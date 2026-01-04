package converter

import (
	metadatav1 "github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/metadata/api"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
)

// ConvertGetUserSessionsResponseToRetrievedSessionsMetadata converts provided metadatav1.GetSessionsResponse
// instance to an array of dto.RetrievedSessionMetadata instances.
func ConvertGetUserSessionsResponseToRetrievedSessionsMetadata(
	input *metadatav1.GetUserSessionsResponse) []dto.RetrievedSessionMetadata {
	var output []dto.RetrievedSessionMetadata

	for _, session := range input.GetSessions() {
		output = append(output, dto.RetrievedSessionMetadata{
			SessionID: session.GetSessionId(),
			Name:      session.GetName(),
			Seed:      session.GetSeed(),
		})
	}

	return output
}

// ConvertGetUserSessionsResponseToRetrievedSessionsMetadata converts provided metadatav1.GetChestsResponse
// instance to an array of dto.RetrievedChestsMetadata instances.
func ConvertGetChestsResponseToRetrievedChests(
	input *metadatav1.GetChestsResponse) []dto.RetrievedChests {
	var output []dto.RetrievedChests

	for _, chest := range input.GetChests() {
		var chestItems []dto.RetrievedChestItems

		for _, chestItem := range chest.GetChestItems() {
			chestItems = append(chestItems, dto.RetrievedChestItems{
				ID:     chestItem.GetChestItemId(),
				Name:   chestItem.GetName(),
				Active: chestItem.GetActive(),
			})
		}

		output = append(output, dto.RetrievedChests{
			SessionID: chest.GetSessionId(),
			ID:        chest.GetChestId(),
			Position: dto.Position{
				X: chest.GetPosition().GetX(),
				Y: chest.GetPosition().GetY(),
			},
			ChestItems: chestItems,
		})
	}

	return output
}

// ConvertGetHealthPacksResponseToRetrievedHealthPacks converts provided metadatav1.GetHealthPacksResponse
// instance to an array of dto.RetrievedHealthPacksMetadata instances.
func ConvertGetHealthPacksResponseToRetrievedHealthPacks(
	input *metadatav1.GetHealthPacksResponse) []dto.RetrievedHealthPack {
	var output []dto.RetrievedHealthPack

	for _, healthPacks := range input.GetHealthPacks() {
		output = append(output, dto.RetrievedHealthPack{
			SessionID: healthPacks.GetSessionId(),
			ID:        healthPacks.GetHealthPackId(),
			Name:      healthPacks.GetName(),
			Position: dto.Position{
				X: healthPacks.GetPosition().GetX(),
				Y: healthPacks.GetPosition().GetY(),
			},
		})
	}

	return output
}

// ConvertGetUserSessionsResponseToListEntries converts provided metadatav1.GetSessionsResponse instance
// to an array of list entries used by UI component.
func ConvertGetUserSessionsResponseToListEntries(
	input *metadatav1.GetUserSessionsResponse) []interface{} {
	var output []interface{}

	for _, session := range input.GetSessions() {
		output = append(output, session.Name)
	}

	return output
}

// ConvertGetLobbySetResponseToRetrievedLobbySetMetadata converts provided metadatav1.GetLobbySetResponse
// instance to an array of dto.RetrievedLobbySetMetadata instances.
func ConvertGetLobbySetResponseToRetrievedLobbySetMetadata(
	input *metadatav1.GetLobbySetResponse) []dto.RetrievedLobbySetMetadata {
	var output []dto.RetrievedLobbySetMetadata

	for _, lobby := range input.GetLobbySet() {
		output = append(output, dto.RetrievedLobbySetMetadata{
			Issuer: lobby.GetIssuer(),
			Skin:   lobby.GetSkin(),
			Host:   lobby.GetHost(),
		})
	}

	return output
}

// ConvertGetLobbySetResponseToListEntries converts provided metadatav1.GetLobbySetResponse instance
// to an array of list entries used by UI component.
func ConvertGetLobbySetResponseToListEntries(
	input []*metadatav1.LobbySetUnit) []interface{} {
	var output []interface{}

	for _, lobby := range input {
		output = append(output, lobby.Skin)
	}

	return output
}

// ConvertPositionsToStartSessionSpawnables converts provided dto.Position array
// to an array of metadatav1.Positions used as spawnables.
func ConvertPositionsToStartSessionSpawnables(input []dto.Position) []*metadatav1.Position {
	var output []*metadatav1.Position

	for _, lobby := range input {
		output = append(output, &metadatav1.Position{
			X: lobby.X,
			Y: lobby.Y,
		})
	}

	return output
}

// ConvertPositionsToStartSessionChestLocations converts provided dto.Position array
// to an array of metadatav1.Positions used as chest locations.
func ConvertPositionsToStartSessionChestLocations(input []dto.Position) []*metadatav1.Position {
	var output []*metadatav1.Position

	for _, lobby := range input {
		output = append(output, &metadatav1.Position{
			X: lobby.X,
			Y: lobby.Y,
		})
	}

	return output
}

// ConvertPositionsToStartSessionHealthPackLocations converts provided dto.Position array
// to an array of metadatav1.Positions used as health pack locations.
func ConvertPositionsToStartSessionHealthPackLocations(input []dto.Position) []*metadatav1.Position {
	var output []*metadatav1.Position

	for _, lobby := range input {
		output = append(output, &metadatav1.Position{
			X: lobby.X,
			Y: lobby.Y,
		})
	}

	return output
}
