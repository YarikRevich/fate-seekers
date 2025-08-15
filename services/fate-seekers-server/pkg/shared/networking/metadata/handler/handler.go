package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/entity"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking/cache"
	metadatav1 "github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking/metadata/api"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/repository"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/repository/converter"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc"
)

var (
	ErrUserDoesNotExist             = errors.New("err happened user does not exist")
	ErrLobbySetDoesNotExist         = errors.New("err happened lobby set does not exist")
	ErrLobbyDoesNotExist            = errors.New("err happened lobby does not exist")
	ErrLobbyAlreadyExists           = errors.New("err happened lobby already exists")
	ErrUserDoesNotOwnSession        = errors.New("err happened user does not own session")
	ErrSessionHasMaxAmountOfLobbies = errors.New("err happened session has max amount of lobbies")
	ErrSessionHasLobbies            = errors.New("err happened session has lobbies")
)

// Handler represents handler implementation of metadatav1.MetadataServer.
type Handler struct {
	metadatav1.UnimplementedMetadataServiceServer
}

func (h *Handler) PingConnection(ctx context.Context, request *metadatav1.PingConnectionRequest) (*metadatav1.PingConnectionResponse, error) {
	// Leave empty. Used to simulation external call to check if client configuration is correct and
	// perform scheduled ping requests.

	return nil, nil
}

func (h *Handler) UpdateSessionActivity(request grpc.ClientStreamingServer[metadatav1.UpdateSessionActivityRequest, metadatav1.UpdateSessionActivityResponse]) error {
	// Leave empty. Used to simulation external call to check if client configuration is correct.

	// switch request.GetMode() {
	// 	case
	// }

	return nil
}

func (h *Handler) CreateUserIfNotExists(ctx context.Context, request *metadatav1.CreateUserIfNotExistsRequest) (*metadatav1.CreateUserIfNotExistsResponse, error) {
	exists, err := repository.
		GetUsersRepository().
		Exists(request.GetIssuer())

	if err != nil {
		return nil, err
	}

	if !exists {
		err = repository.
			GetUsersRepository().
			Insert(request.GetIssuer())

		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (h *Handler) GetSessions(ctx context.Context, request *metadatav1.GetSessionsRequest) (*metadatav1.GetSessionsResponse, error) {
	response := new(metadatav1.GetSessionsResponse)

	cache.
		GetInstance().
		BeginUserSessionsTransaction()

	cachedSessions, ok := cache.
		GetInstance().
		GetUserSessions(request.GetIssuer())
	if ok {
		for _, cachedSession := range cachedSessions {
			response.Sessions = append(response.Sessions, &metadatav1.Session{
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
			user, exists, err := repository.
				GetUsersRepository().
				GetByName(request.GetIssuer())
			if err != nil {
				cache.
					GetInstance().
					CommitUserSessionsTransaction()

				return nil, err
			}

			if !exists {
				cache.
					GetInstance().
					CommitUserSessionsTransaction()

				return nil, ErrUserDoesNotExist
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
			cache.
				GetInstance().
				CommitUserSessionsTransaction()

			return nil, err
		}

		var sessions []dto.CacheSessionEntity

		for _, rawSession := range rawSessions {
			response.Sessions = append(response.Sessions, &metadatav1.Session{
				SessionId: rawSession.ID,
				Seed:      uint64(rawSession.Seed),
				Name:      rawSession.Name,
			})

			sessions = append(sessions, dto.CacheSessionEntity{
				ID:   rawSession.ID,
				Seed: uint64(rawSession.Seed),
				Name: rawSession.Name,
			})
		}

		cache.
			GetInstance().
			AddUserSessions(request.GetIssuer(), sessions)
	}

	cache.
		GetInstance().
		CommitUserSessionsTransaction()

	return response, nil
}

func (h *Handler) CreateSession(ctx context.Context, request *metadatav1.CreateSessionRequest) (*metadatav1.CreateSessionResponse, error) {
	var userID int64

	cachedUserID, ok := cache.
		GetInstance().
		GetUsers(request.GetIssuer())
	if ok {
		userID = cachedUserID
	} else {
		user, exists, err := repository.
			GetUsersRepository().
			GetByName(request.GetIssuer())
		if err != nil {
			return nil, err
		}

		if !exists {
			return nil, ErrUserDoesNotExist
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

func (h *Handler) RemoveSession(ctx context.Context, request *metadatav1.RemoveSessionRequest) (*metadatav1.RemoveSessionResponse, error) {
	var isCacheSessionsPresent bool

	cachedSessions, ok := cache.
		GetInstance().
		GetUserSessions(request.GetIssuer())
	if ok {
		if slices.ContainsFunc(
			cachedSessions,
			func(value dto.CacheSessionEntity) bool {
				return value.ID == request.GetSessionId()
			}) {
			isCacheSessionsPresent = true
		}
	}

	var userID int64

	if !isCacheSessionsPresent {
		cachedUserID, ok := cache.
			GetInstance().
			GetUsers(request.GetIssuer())
		if ok {
			userID = cachedUserID
		} else {
			user, exists, err := repository.
				GetUsersRepository().
				GetByName(request.GetIssuer())
			if err != nil {
				return nil, err
			}

			if !exists {
				return nil, ErrUserDoesNotExist
			}

			userID = user.ID

			cache.
				GetInstance().
				AddUser(request.GetIssuer(), userID)
		}

		sessions, err := repository.
			GetSessionsRepository().
			GetByIssuer(userID)
		if err != nil {
			return nil, err
		}

		if !slices.ContainsFunc(
			sessions,
			func(value *entity.SessionEntity) bool {
				return value.ID == request.GetSessionId()
			}) {
			return nil, ErrUserDoesNotOwnSession
		}
	}

	cachedUserID, ok := cache.
		GetInstance().
		GetUsers(request.GetIssuer())
	if ok {
		userID = cachedUserID
	} else {
		user, exists, err := repository.
			GetUsersRepository().
			GetByName(request.GetIssuer())
		if err != nil {
			return nil, err
		}

		if !exists {
			return nil, ErrUserDoesNotExist
		}

		userID = user.ID

		cache.
			GetInstance().
			AddUser(request.GetIssuer(), userID)
	}

	cache.
		GetInstance().
		BeginLobbySetTransaction()

	cachedLobbySet, ok := cache.
		GetInstance().
		GetLobbySet(cachedUserID)
	if ok && len(cachedLobbySet) != 0 {
		cache.
			GetInstance().
			CommitLobbySetTransaction()

		return nil, ErrSessionHasLobbies
	}

	lobbies, exists, err := repository.
		GetLobbiesRepository().
		GetBySessionID(request.GetSessionId())
	if err != nil {
		cache.
			GetInstance().
			CommitLobbySetTransaction()

		return nil, err
	}

	if exists {
		var issuers []string

		for _, lobby := range lobbies {
			issuers = append(issuers, lobby.UserEntity.Name)
		}

		cache.
			GetInstance().
			EvictLobbySet(request.GetSessionId())

		cache.
			GetInstance().
			AddLobbySet(request.GetSessionId(), issuers)

		cache.
			GetInstance().
			CommitLobbySetTransaction()

		return nil, ErrSessionHasLobbies
	}

	cache.
		GetInstance().
		CommitLobbySetTransaction()

	err = repository.
		GetSessionsRepository().
		DeleteByID(request.GetSessionId())
	if err != nil {
		return nil, err
	}

	cache.
		GetInstance().
		BeginUserSessionsTransaction()

	cache.
		GetInstance().
		EvictUserSessions(request.GetIssuer())

	cache.
		GetInstance().
		CommitUserSessionsTransaction()

	return nil, nil
}

func (h *Handler) GetLobbySet(ctx context.Context, request *metadatav1.GetLobbySetRequest) (*metadatav1.GetLobbySetResponse, error) {
	response := new(metadatav1.GetLobbySetResponse)

	cache.
		GetInstance().
		BeginLobbySetTransaction()

	issuers, ok := cache.
		GetInstance().
		GetLobbySet(request.GetSessionId())
	if !ok {
		lobbies, exists, err := repository.
			GetLobbiesRepository().
			GetBySessionID(request.GetSessionId())
		if err != nil {
			cache.
				GetInstance().
				CommitLobbySetTransaction()

			return nil, err
		}

		if !exists {
			cache.
				GetInstance().
				CommitLobbySetTransaction()

			return nil, ErrLobbySetDoesNotExist
		}

		var issuers []string

		for _, lobby := range lobbies {
			issuers = append(issuers, lobby.UserEntity.Name)
		}

		cache.
			GetInstance().
			EvictLobbySet(request.GetSessionId())

		cache.
			GetInstance().
			AddLobbySet(request.GetSessionId(), issuers)
	}

	cache.
		GetInstance().
		CommitLobbySetTransaction()

	response.Issuers = issuers

	return response, nil
}

func (h *Handler) CreateLobby(ctx context.Context, request *metadatav1.CreateLobbyRequest) (*metadatav1.CreateLobbyResponse, error) {
	var userID int64

	fmt.Println("BEGINNING")

	cachedUserID, ok := cache.
		GetInstance().
		GetUsers(request.GetIssuer())
	if ok {
		userID = cachedUserID
	} else {
		user, exists, err := repository.
			GetUsersRepository().
			GetByName(request.GetIssuer())
		if err != nil {
			return nil, err
		}

		if !exists {
			return nil, ErrUserDoesNotExist
		}

		userID = user.ID

		cache.
			GetInstance().
			AddUser(request.GetIssuer(), userID)
	}

	var host bool

	cachedSessions, ok := cache.
		GetInstance().
		GetUserSessions(request.GetIssuer())
	if ok {
		if slices.ContainsFunc(
			cachedSessions,
			func(value dto.CacheSessionEntity) bool {
				return value.ID == request.GetSessionId()
			}) {
			host = true
		} else {
			cachedUserID, ok := cache.
				GetInstance().
				GetUsers(request.GetIssuer())
			if ok {
				userID = cachedUserID
			} else {
				user, exists, err := repository.
					GetUsersRepository().
					GetByName(request.GetIssuer())
				if err != nil {
					return nil, err
				}

				if !exists {
					return nil, ErrUserDoesNotExist
				}

				userID = user.ID

				cache.
					GetInstance().
					AddUser(request.GetIssuer(), userID)
			}

			sessions, err := repository.
				GetSessionsRepository().
				GetByIssuer(userID)
			if err != nil {
				return nil, err
			}

			if slices.ContainsFunc(
				sessions,
				func(value *entity.SessionEntity) bool {
					return value.ID == request.GetSessionId()
				}) {
				host = true
			}
		}
	}

	userLobbies, exists, err := repository.
		GetLobbiesRepository().
		GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	if slices.ContainsFunc(
		userLobbies,
		func(value *entity.LobbyEntity) bool {
			return value.SessionID == request.GetSessionId()
		}) {
		return nil, ErrLobbyAlreadyExists
	}

	cache.
		GetInstance().
		BeginLobbySetTransaction()

	cachedLobbySet, ok := cache.
		GetInstance().
		GetLobbySet(request.GetSessionId())
	if ok && len(cachedLobbySet) >= config.MAX_SESSION_USERS {
		cache.
			GetInstance().
			CommitLobbySetTransaction()

		return nil, ErrSessionHasMaxAmountOfLobbies
	}

	sessionLobbies, exists, err := repository.
		GetLobbiesRepository().
		GetBySessionID(request.GetSessionId())
	if err != nil {
		cache.
			GetInstance().
			CommitLobbySetTransaction()

		return nil, err
	}

	if exists {
		var issuers []string

		for _, lobby := range sessionLobbies {
			issuers = append(issuers, lobby.UserEntity.Name)
		}

		cache.
			GetInstance().
			EvictLobbySet(request.GetSessionId())

		cache.
			GetInstance().
			AddLobbySet(request.GetSessionId(), issuers)

		if len(cachedLobbySet) >= config.MAX_SESSION_USERS {
			cache.
				GetInstance().
				CommitLobbySetTransaction()

			return nil, ErrSessionHasMaxAmountOfLobbies
		}
	}

	cache.
		GetInstance().
		CommitLobbySetTransaction()

	fmt.Println("INSERTING OR UPDATING LOBBIES")

	err = repository.
		GetLobbiesRepository().
		InsertOrUpdate(
			dto.LobbiesRepositoryInsertOrUpdateRequest{
				UserID:    userID,
				SessionID: request.GetSessionId(),
				Host:      host,
			})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *Handler) RemoveLobby(context context.Context, request *metadatav1.RemoveLobbyRequest) (*metadatav1.RemoveLobbyResponse, error) {
	var userID int64

	cachedUserID, ok := cache.
		GetInstance().
		GetUsers(request.GetIssuer())
	if ok {
		userID = cachedUserID
	} else {
		user, exists, err := repository.
			GetUsersRepository().
			GetByName(request.GetIssuer())
		if err != nil {
			return nil, err
		}

		if !exists {
			return nil, ErrUserDoesNotExist
		}

		userID = user.ID

		cache.
			GetInstance().
			AddUser(request.GetIssuer(), userID)
	}

	lobbies, exists, err := repository.
		GetLobbiesRepository().
		GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	if exists {
		cache.
			GetInstance().
			BeginLobbySetTransaction()

		for _, lobby := range lobbies {
			if lobby.SessionID == request.GetSessionId() {
				cache.
					GetInstance().
					EvictLobbySet(lobby.SessionID)

				break
			}
		}

		cache.
			GetInstance().
			CommitLobbySetTransaction()
	}

	err = repository.
		GetLobbiesRepository().
		DeleteByUserID(userID)
	if err != nil {
		return nil, err
	}

	cache.
		GetInstance().
		BeginMetadataTransaction()

	cache.
		GetInstance().
		EvictMetadata(request.GetIssuer())

	cache.
		GetInstance().
		CommitMetadataTransaction()

	return nil, nil
}

func (h *Handler) GetUserMetadata(request *metadatav1.GetUserMetadataRequest, stream grpc.ServerStreamingServer[metadatav1.GetUserMetadataResponse]) error {
	response := new(metadatav1.GetUserMetadataResponse)

	cache.
		GetInstance().
		BeginMetadataTransaction()

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
			user, exists, err := repository.
				GetUsersRepository().
				GetByName(request.GetIssuer())
			if err != nil {
				cache.
					GetInstance().
					CommitMetadataTransaction()

				return err
			}

			if !exists {
				cache.
					GetInstance().
					CommitMetadataTransaction()

				return ErrUserDoesNotExist
			}

			userID = user.ID

			cache.
				GetInstance().
				AddUser(request.GetIssuer(), userID)
		}

		lobbies, exists, err := repository.
			GetLobbiesRepository().
			GetByUserID(userID)
		if err != nil {
			cache.
				GetInstance().
				CommitMetadataTransaction()

			return err
		}

		if !exists {
			cache.
				GetInstance().
				CommitMetadataTransaction()

			return ErrLobbyDoesNotExist
		}

		cache.
			GetInstance().
			AddMetadata(
				request.GetIssuer(),
				converter.ConvertLobbyEntityToCacheMetadataEntity(
					lobbies))

		for _, lobby := range lobbies {
			if lobby.SessionID == request.GetSessionId() {
				response.UserMetadata = &metadatav1.UserMetadata{
					Health:     uint64(lobby.Health),
					Skin:       uint64(lobby.Skin),
					Eliminated: lobby.Eliminated,
					Position: &metadatav1.Position{
						X: lobby.PositionX,
						Y: lobby.PositionY,
					},
				}

				break
			}
		}
	} else {
		for _, value := range metadata {
			if value.SessionID == request.GetSessionId() {
				response.UserMetadata = &metadatav1.UserMetadata{
					Health:     value.Health,
					Skin:       value.Skin,
					Eliminated: value.Eliminated,
					Position: &metadatav1.Position{
						X: value.PositionX,
						Y: value.PositionY,
					},
				}

				break
			}
		}
	}

	cache.
		GetInstance().
		CommitMetadataTransaction()

	stream.Send(response)

	return nil
}

func (h *Handler) GetChests(context.Context, *metadatav1.GetChestsRequest) (*metadatav1.GetChestsResponse, error) {
	return nil, nil
}

func (h *Handler) GetMap(request *metadatav1.GetMapRequest, stream grpc.ServerStreamingServer[metadatav1.GetMapResponse]) error {
	return nil
}

func (h *Handler) GetChatMessages(request *metadatav1.GetChatMessagesRequest, stream grpc.ServerStreamingServer[metadatav1.GetChatMessagesResponse]) error {
	// TODO: messages would be retrieved from memory(not lru cache??????)

	return nil
}

func (h *Handler) CreateChatMessage(context.Context, *metadatav1.CreateChatMessageRequest) (*metadatav1.CreateChatMessageResponse, error) {
	// TODO: add to a delayed batch, not to overload the database.

	return nil, nil
}

// NewHandler initializes implementation of metadatav1.MetadataServer.
func NewHandler() metadatav1.MetadataServiceServer {
	return new(Handler)
}
