package handler

import (
	"context"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking/metadata/api"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking/metadata/cache"
	"google.golang.org/grpc"
)

// Handler represents handler implementation of api.MetadataServer.
type Handler struct {
	api.UnimplementedMetadataServer
}

func (h *Handler) PingConnection(ctx context.Context, request *api.PingConnectionRequest) (*api.PingConnectionResponse, error) {
	// Leave empty. Used to simulation external call to check if client configuration is correct.

	return nil, nil
}

func (h *Handler) GetSessions(ctx context.Context, request *api.GetSessionsRequest) (*api.GetSessionsResponse, error) {
	// Should use LRU cache over here for safety

	value, ok := cache.GetInstance().
		GetSessions().
		Get(request.GetIssuer())

	if !ok {

	}

	return nil, nil
}

func (h *Handler) CreateSession(ctx context.Context, request *api.CreateSessionRequest) (*api.CreateSessionResponse, error) {

	return nil, nil
}

func (h *Handler) RemoveSession(context.Context, *api.RemoveSessionRequest) (*api.RemoveSessionResponse, error) {
	return nil, nil
}

func (h *Handler) JoinToSession(context.Context, *api.JoinToSessionRequest) (*api.JoinToSessionResponse, error) {
	return nil, nil
}

func (h *Handler) GetChests(context.Context, *api.GetChestsRequest) (*api.GetChestsResponse, error) {
	return nil, nil
}

func (h *Handler) GetMap(*api.GetMapRequest, grpc.ServerStreamingServer[api.GetMapResponse]) error {
	return nil
}

func (h *Handler) GetChat(*api.GetChatRequest, grpc.ServerStreamingServer[api.GetChatResponse]) error {
	// TODO: messages would be retrieved from memory(not lru cache??????)

	return nil
}

func (h *Handler) CreateChatMessage(context.Context, *api.CreateChatMessageRequest) (*api.CreateChatMessageResponse, error) {
	// TODO: add to a delayed batch, not to overload the database.

	return nil, nil
}

// NewHandler initializes implementation of api.MetadataServer.
func NewHandler() api.MetadataServer {
	return new(Handler)
}
