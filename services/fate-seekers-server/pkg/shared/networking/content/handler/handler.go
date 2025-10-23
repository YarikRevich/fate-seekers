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
	case contentv1.UPDATE_USER_METADATA_POSITIONS:
		var message contentv1.UpdateUserMetadataPositionsRequest
		if err := proto.Unmarshal(value, &message); err != nil {
			return err
		}

		fmt.Println(message)

		// TODO: safe user position in the lobby metadata.
	case contentv1.OPEN_GENERATED_CHEST:
	case contentv1.OPEN_GENERATED_HEALTH_PACK:
	case contentv1.SEND_CHAT_MESSAGE:
	}

	return nil
}

// NewHandler initializes Handler.
func NewHandler() *Handler {
	return new(Handler)
}
