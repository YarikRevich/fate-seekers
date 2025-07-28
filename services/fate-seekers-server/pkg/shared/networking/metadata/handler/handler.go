package handler

import (
	"context"
	"errors"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking/cache"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking/metadata/api"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/repository"
	"google.golang.org/grpc"
)

var (
	ErrLobbySetDoesNotExist = errors.New("err happened lobby set does not exist")
	ErrLobbyDoesNotExist    = errors.New("err happened lobby does not exist")
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
	response := new(api.GetSessionsResponse)

	cachedSessions, ok := cache.
		GetInstance().
		GetSessions(request.GetIssuer())
	if ok {
		for _, cachedSession := range cachedSessions {
			response.Sessions = append(response.Sessions, &api.Session{
				SessionId: cachedSession.ID,
				Seed:      cachedSession.Seed,
				Name:      cachedSession.Name,
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
				SessionId: rawSession.ID,
				Seed:      rawSession.Seed,
				Name:      rawSession.Name,
			})

			sessions = append(sessions, dto.CacheSessionEntity{
				ID:   rawSession.ID,
				Seed: rawSession.Seed,
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

	err := repository.
		GetSessionsRepository().
		Insert(request.GetName(), userID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *Handler) RemoveSession(ctx context.Context, request *api.RemoveSessionRequest) (*api.RemoveSessionResponse, error) {
	err := repository.
		GetSessionsRepository().
		DeleteByID(request.GetSessionId())
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *Handler) GetLobbySet(ctx context.Context, request *api.GetLobbySetRequest) (*api.GetLobbySetResponse, error) {
	response := new(api.GetLobbySetResponse)

	issuers, ok := cache.
		GetInstance().
		GetLobbySet(request.GetSessionId())
	if !ok {
		return nil, ErrLobbySetDoesNotExist
	}

	response.Issuers = issuers

	return response, nil
}

func (h *Handler) CreateLobby(ctx context.Context, request *api.CreateLobbyRequest) (*api.CreateLobbyResponse, error) {
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

	err := repository.
		GetLobbiesRepository().
		InsertOrUpdate(
			dto.LobbiesRepositoryInsertOrUpdateRequest{
				UserID:    userID,
				SessionID: request.GetSessionId(),
			})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *Handler) RemoveLobby(context context.Context, request *api.RemoveLobbyRequest) (*api.RemoveLobbyResponse, error) {
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

	err := repository.
		GetLobbiesRepository().
		DeleteByUserID(userID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *Handler) GetUserMetadata(ctx context.Context, request *api.GetUserMetadataRequest) (*api.GetUserMetadataResponse, error) {
	response := new(api.GetUserMetadataResponse)

	metadata, ok := cache.
		GetInstance().
		GetMetadata(request.GetIssuer())
	if !ok {
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

		lobby, exists, err := repository.
			GetLobbiesRepository().
			GetByUserID(userID)
		if err != nil {
			return nil, err
		}

		if !exists {
			return nil, ErrLobbyDoesNotExist
		}

		cache.
			GetInstance().
			AddMetadata(
				request.GetIssuer(), dto.CacheMetadataEntity{
					SessionID:  lobby.SessionID,
					PositionX:  lobby.PositionX,
					PositionY:  lobby.PositionY,
					Skin:       lobby.Skin,
					Health:     lobby.Health,
					Eliminated: lobby.Eliminated,
				})

		response.UserMetadata = &api.UserMetadata{
			Health:     lobby.Health,
			Skin:       lobby.Skin,
			Eliminated: lobby.Eliminated,
			Position: &api.Position{
				X: lobby.PositionX,
				Y: lobby.PositionY,
			},
		}
	} else {
		response.UserMetadata = &api.UserMetadata{
			Health:     metadata.Health,
			Skin:       metadata.Skin,
			Eliminated: metadata.Eliminated,
			Position: &api.Position{
				X: metadata.PositionX,
				Y: metadata.PositionY,
			},
		}
	}

	return response, nil
}

func (h *Handler) GetChests(context.Context, *api.GetChestsRequest) (*api.GetChestsResponse, error) {
	return nil, nil
}

func (h *Handler) GetMap(request *api.GetMapRequest, stream grpc.ServerStreamingServer[api.GetMapResponse]) error {
	return nil
}

func (h *Handler) GetChatMessages(request *api.GetChatMessagesRequest, stream grpc.ServerStreamingServer[api.GetChatMessagesResponse]) error {
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
