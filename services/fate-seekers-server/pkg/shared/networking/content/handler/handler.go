package handler

import "github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking/content/api"

// Handler performs content connector state management.
type Handler struct {
}

func (h *Handler) Process(key string, value []byte) error {
	switch key {
	case api.GET_USER_METADATA_POSITIONS:

	}

	return nil
}

func NewHandler() *Handler {
	return new(Handler)
}
