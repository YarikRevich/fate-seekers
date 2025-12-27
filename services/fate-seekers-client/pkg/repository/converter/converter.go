package converter

import "github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/entity"

// ConvertRetrievedCollectionsToListEntries converts provided retrieved collections instance
// to an array of list entries used by UI component.
func ConvertRetrievedCollectionsToListEntries(
	input []entity.CollectionEntity) []interface{} {
	var output []interface{}

	for _, lobby := range input {
		output = append(output, lobby.Name)
	}

	return output
}
