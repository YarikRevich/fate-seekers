package handler

import (
	"context"
	"fmt"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking/metadata/api"
	"google.golang.org/grpc"
)

// Handler represents handler implementation of api.MetadataServer.
type Handler struct {
	api.UnimplementedMetadataServer
}

func (h *Handler) PingConnection(ctx context.Context, request *api.PingConnectionRequest) (*api.PingConnectionResponse, error) {
	fmt.Println(request.Issuer, "REQUEST")

	return nil, nil
}

func (h *Handler) CreateSession(context.Context, *api.CreateSessionRequest) (*api.CreateSessionResponse, error) {
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
	return nil
}

// NewHandler initializes implementation of api.MetadataServer.
func NewHandler() api.MetadataServer {
	return new(Handler)
}
