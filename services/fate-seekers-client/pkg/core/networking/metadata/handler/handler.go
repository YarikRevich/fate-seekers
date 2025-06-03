package handler

type Handler struct {
}

func (h *Handler) PerformPingConnection(callback func(err error)) {

}

func newHandler() *Handler {
	return new(Handler)
}
