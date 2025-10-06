package handler

import (
	"fmt"

	contentv1 "github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking/content/api"
)

// Handler performs content connector state management.
type Handler struct {
}

func (h *Handler) Process(key string, value []byte) error {
	switch key {
	case contentv1.GET_USER_METADATA_POSITIONS:
	case contentv1.UPDATE_USER_METADATA_POSITIONS:
		fmt.Println(key, string(value))
	}

	return nil
}

// NewHandler initializes Handler.
func NewHandler() *Handler {
	return new(Handler)
}
