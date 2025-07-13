package handler

import (
	"context"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking/cache"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking/metadata/api"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/repository"
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
	var response *api.GetSessionsResponse

	cachedSessions, ok := cache.
		GetInstance().
		GetSessions(request.GetIssuer())
	if ok {
		for _, cachedSession := range cachedSessions {
			response.Sessions = append(response.Sessions, &api.Session{
				Id:   cachedSession.ID,
				Name: cachedSession.Name,
			})
		}
	} else {
		var userID int64

		cachedUserID, ok := cache.
			GetInstance().
			GetUsers(request.GetIssuer())
		if ok {
			userID = cachedUserID
		} else {
			user, _, err := repository.
				GetUsersRepository().
				GetByName(request.GetIssuer())
			if err != nil {
				return nil, err
			}

			userID = user.ID

			cache.
				GetInstance().
				AddUser(request.GetIssuer(), userID)
		}

		rawSessions, err := repository.
			GetSessionsRepository().
			GetByIssuer(userID)
		if err != nil {
			return nil, err
		}

		var sessions []dto.CacheSessionEntity

		for _, rawSession := range rawSessions {
			response.Sessions = append(response.Sessions, &api.Session{
				Id:   rawSession.ID,
				Name: rawSession.Name,
			})

			sessions = append(sessions, dto.CacheSessionEntity{
				ID:   rawSession.ID,
				Name: rawSession.Name,
			})
		}

		cache.
			GetInstance().
			AddSessions(request.GetIssuer(), sessions)
	}

	return response, nil
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

func (h *Handler) GetMap(request *api.GetMapRequest, stream grpc.ServerStreamingServer[api.GetMapResponse]) error {
	return nil
}

func (h *Handler) GetChat(request *api.GetChatRequest, stream grpc.ServerStreamingServer[api.GetChatResponse]) error {
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
