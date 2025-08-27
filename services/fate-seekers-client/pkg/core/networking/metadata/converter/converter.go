package converter

import (
	metadatav1 "github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/metadata/api"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
)

// ConvertGetSessionsResponseToRetrievedSessionsMetadata converts provided metadatav1.GetSessionsResponse
// instance to an array of dto.RetrievedSessionMetadata instances.
func ConvertGetSessionsResponseToRetrievedSessionsMetadata(
	input *metadatav1.GetSessionsResponse) []dto.RetrievedSessionMetadata {
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

// ConvertGetSessionsResponseToListEntries converts provided metadatav1.GetSessionsResponse instance
// to an array of list entries used by UI component.
func ConvertGetSessionsResponseToListEntries(
	input *metadatav1.GetSessionsResponse) []interface{} {
	var output []interface{}

	for _, session := range input.GetSessions() {
		output = append(output, session.Name)
	}

	return output
}

// ConvertGetLobbySetResponseToListEntries converts provided metadatav1.GetLobbySetResponse instance
// to an array of list entries used by UI component.
func ConvertGetLobbySetResponseToListEntries(
	input *metadatav1.GetLobbySetResponse) []interface{} {
	var output []interface{}

	for _, lobby := range input.GetIssuers() {
		output = append(output, lobby)
	}

	return output
}
