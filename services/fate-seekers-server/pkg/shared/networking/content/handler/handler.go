package handler

import (
	"fmt"

	contentv1 "github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking/content/api"
	"google.golang.org/protobuf/proto"
)

// Handler performs content connector state management.
type Handler struct {
}

func (h *Handler) Process(key string, value []byte) error {
	switch key {
	case contentv1.GET_USER_METADATA_POSITIONS:
	case contentv1.UPDATE_USER_METADATA_POSITIONS:
		var message contentv1.UpdateUserMetadataPositionsRequest
		if err := proto.Unmarshal(value, &message); err != nil {
			return err
		}

		fmt.Println(message)

		// TODO: safe user position in the lobby metadata.
	}

	return nil
}

// NewHandler initializes Handler.
func NewHandler() *Handler {
	return new(Handler)
}
