package handler

// Handler performs content connector state management.
type Handler struct {
}

func (h *Handler) Process(key string, value []byte) error {
	return nil
}

func NewHandler() *Handler {
	return new(Handler)
}
