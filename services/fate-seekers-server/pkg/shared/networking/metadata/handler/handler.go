package handler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/entity"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking/cache"
	metadatav1 "github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking/metadata/api"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/repository"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/repository/converter"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrUserDoesNotExist                   = errors.New("err happened user does not exist")
	ErrLobbySetDoesNotExist               = errors.New("err happened lobby set does not exist")
	ErrLobbyDoesNotExist                  = errors.New("err happened lobby does not exist")
	ErrLobbyAlreadyExists                 = errors.New("err happened lobby already exists")
	ErrSessionDoesNotExists               = errors.New("err happened session does not exist")
	ErrSessionAlreadyExists               = errors.New("err happened session already exists")
	ErrSessionAlreadyStarted              = errors.New("err happened session already started")
	ErrFilteredSessionDoesNotExists       = errors.New("err happened filtered session does not exist")
	ErrUserIsNotLobbyHost                 = errors.New("err happened user is not a host of a lobby")
	ErrUserDoesNotOwnSession              = errors.New("err happened user does not own session")
	ErrSessionHasMaxAmountOfLobbies       = errors.New("err happened session has max amount of lobbies")
	ErrSessionHasLobbies                  = errors.New("err happened session has lobbies")
	ErrSessionMetadataRetrievalNotAllowed = errors.New("err happened session metadata retrieval not allowed")
)

// Describes constant values used for handler management.
const (
	getSessionMetadataFrequency = time.Second * 2
	getLobbySetFrequency        = time.Second * 1
	getUserMetadataFrequency    = time.Second * 2
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

func (h *Handler) UpdateSessionActivity(stream grpc.ClientStreamingServer[metadatav1.UpdateSessionActivityRequest, metadatav1.UpdateSessionActivityResponse]) error {
	for {
		_, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&metadatav1.UpdateSessionActivityResponse{
				// fill response fields here
			})
		}
		if err != nil {
			return err
		}
	}
}

func (h *Handler) CreateUserIfNotExists(ctx context.Context, request *metadatav1.CreateUserIfNotExistsRequest) (*metadatav1.CreateUserIfNotExistsResponse, error) {
	exists, err := repository.
		GetUsersRepository().
		ExistsByName(request.GetIssuer())

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

	return new(metadatav1.CreateUserIfNotExistsResponse), nil
}

func (h *Handler) GetUserSessions(ctx context.Context, request *metadatav1.GetUserSessionsRequest) (*metadatav1.GetUserSessionsResponse, error) {
	response := new(metadatav1.GetUserSessionsResponse)

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

			sessions = append(
				sessions,
				converter.ConvertSessionEntityToCacheSessionEntity(rawSession))
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

func (h *Handler) GetFilteredSession(ctx context.Context, request *metadatav1.GetFilteredSessionRequest) (*metadatav1.GetFilteredSessionResponse, error) {
	response := new(metadatav1.GetFilteredSessionResponse)

	cache.
		GetInstance().
		BeginSessionsTransaction()

	var found bool

	for _, value := range cache.
		GetInstance().
		GetSessionsMappings() {
		if value.Name == request.GetName() {
			response.Session = &metadatav1.Session{
				SessionId: value.ID,
				Seed:      value.Seed,
				Name:      value.Name,
			}

			found = true

			break
		}
	}

	if !found {
		session, exists, err := repository.
			GetSessionsRepository().
			GetByName(request.GetName())

		if err != nil {
			cache.
				GetInstance().
				CommitSessionsTransaction()

			return nil, err
		}

		if !exists {
			cache.
				GetInstance().
				CommitSessionsTransaction()

			return nil, status.Errorf(codes.NotFound, ErrFilteredSessionDoesNotExists.Error())
		}

		response.Session = &metadatav1.Session{
			SessionId: session.ID,
			Seed:      uint64(session.Seed),
			Name:      session.Name,
		}

		cache.GetInstance().EvictSessions(session.ID)

		cache.
			GetInstance().
			AddSessions(
				session.ID,
				converter.ConvertSessionEntityToCacheSessionEntity(session))
	}

	cache.
		GetInstance().
		CommitSessionsTransaction()

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

	exists, err := repository.
		GetSessionsRepository().
		ExistsByName(request.GetName())
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, ErrSessionAlreadyExists
	}

	err = repository.
		GetSessionsRepository().
		InsertOrUpdate(dto.SessionsRepositoryInsertOrUpdateRequest{
			Name:   request.GetName(),
			Seed:   int64(request.GetSeed()),
			Issuer: userID,
		})
	if err != nil {
		return nil, err
	}

	return new(metadatav1.CreateSessionResponse), nil
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
		var lobbySet []dto.CacheLobbySetEntity

		for _, lobby := range lobbies {
			lobbySet = append(lobbySet, dto.CacheLobbySetEntity{
				Issuer: lobby.UserEntity.Name,
				Skin:   uint64(lobby.Skin),
				Host:   lobby.Host,
			})
		}

		cache.
			GetInstance().
			EvictLobbySet(request.GetSessionId())

		cache.
			GetInstance().
			AddLobbySet(request.GetSessionId(), lobbySet)

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

	return new(metadatav1.RemoveSessionResponse), nil
}

func (h *Handler) StartSession(ctx context.Context, request *metadatav1.StartSessionRequest) (*metadatav1.StartSessionResponse, error) {
	cache.
		GetInstance().
		BeginMetadataTransaction()

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

			return nil, err
		}

		if !exists {
			cache.
				GetInstance().
				CommitMetadataTransaction()

			return nil, ErrUserDoesNotExist
		}

		userID = user.ID

		cache.
			GetInstance().
			AddUser(request.GetIssuer(), userID)
	}

	metadata, ok := cache.
		GetInstance().
		GetMetadata(request.GetIssuer())
	if !ok {
		lobbies, exists, err := repository.
			GetLobbiesRepository().
			GetByUserID(userID)
		if err != nil {
			cache.
				GetInstance().
				CommitMetadataTransaction()

			return nil, err
		}

		if !exists {
			cache.
				GetInstance().
				CommitMetadataTransaction()

			return nil, ErrLobbyDoesNotExist
		}

		cache.
			GetInstance().
			AddMetadata(
				request.GetIssuer(),
				converter.ConvertLobbyEntityToCacheMetadataEntity(
					lobbies))

		var selectedLobby *entity.LobbyEntity

		for _, lobby := range lobbies {
			if lobby.ID == request.GetLobbyId() &&
				lobby.SessionID == request.GetSessionId() {
				selectedLobby = lobby
			}
		}

		if selectedLobby == nil {
			cache.
				GetInstance().
				CommitMetadataTransaction()

			return nil, ErrLobbyDoesNotExist
		}

		if !selectedLobby.Host {
			cache.
				GetInstance().
				CommitMetadataTransaction()

			return nil, ErrUserIsNotLobbyHost
		}
	} else {
		var selectedLobby *dto.CacheMetadataEntity

		for _, value := range metadata {
			if value.LobbyID == request.GetLobbyId() &&
				value.SessionID == request.GetSessionId() {
				selectedLobby = value
			}
		}

		if selectedLobby == nil {
			cache.
				GetInstance().
				CommitMetadataTransaction()

			return nil, ErrLobbyDoesNotExist
		}

		if !selectedLobby.Host {
			cache.
				GetInstance().
				CommitMetadataTransaction()

			return nil, ErrUserIsNotLobbyHost
		}
	}

	cache.
		GetInstance().
		CommitMetadataTransaction()

	cache.
		GetInstance().
		BeginSessionsTransaction()

	var sessionName string

	cachedSession, ok := cache.
		GetInstance().
		GetSessions(request.GetSessionId())
	if !ok {
		session, _, err := repository.
			GetSessionsRepository().
			GetByID(request.GetSessionId())
		if err != nil {
			cache.
				GetInstance().
				CommitSessionsTransaction()

			return nil, err
		}

		if session.Started {
			cache.
				GetInstance().
				CommitSessionsTransaction()

			return nil, ErrSessionAlreadyStarted
		}

		sessionName = session.Name

		cache.
			GetInstance().
			AddSessions(
				request.GetSessionId(),
				converter.ConvertSessionEntityToCacheSessionEntity(session))
	} else {
		if cachedSession.Started {
			cache.
				GetInstance().
				CommitSessionsTransaction()

			return nil, ErrSessionAlreadyStarted
		}

		sessionName = cachedSession.Name
	}

	cache.
		GetInstance().
		CommitSessionsTransaction()

	err := repository.
		GetSessionsRepository().
		InsertOrUpdate(
			dto.SessionsRepositoryInsertOrUpdateRequest{
				ID:      request.GetSessionId(),
				Name:    sessionName,
				Issuer:  userID,
				Started: true,
			})

	return new(metadatav1.StartSessionResponse), err
}

func (h *Handler) GetSessionMetadata(request *metadatav1.GetSessionMetadataRequest, stream grpc.ServerStreamingServer[metadatav1.GetSessionMetadataResponse]) error {
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

		var found bool

		for _, lobby := range lobbies {
			if lobby.SessionID == request.GetSessionId() {
				found = true

				break
			}
		}

		if !found {
			cache.
				GetInstance().
				CommitMetadataTransaction()

			return ErrSessionMetadataRetrievalNotAllowed
		}
	} else {
		var found bool

		for _, value := range metadata {
			if value.SessionID == request.GetSessionId() {
				found = true

				break
			}
		}

		if !found {
			cache.
				GetInstance().
				CommitMetadataTransaction()

			return ErrSessionMetadataRetrievalNotAllowed
		}
	}

	cache.
		GetInstance().
		CommitMetadataTransaction()

	ticker := time.NewTicker(getSessionMetadataFrequency)

	for {
		select {
		case <-ticker.C:
			ticker.Stop()

			var started bool

			cache.
				GetInstance().
				BeginSessionsTransaction()

			cachedSession, ok := cache.
				GetInstance().
				GetSessions(request.GetSessionId())
			if !ok {
				session, _, err := repository.
					GetSessionsRepository().
					GetByID(request.GetSessionId())
				if err != nil {
					cache.
						GetInstance().
						CommitSessionsTransaction()

					return err
				}

				started = session.Started

				cache.
					GetInstance().
					AddSessions(
						request.GetSessionId(),
						converter.ConvertSessionEntityToCacheSessionEntity(session))
			} else {
				started = cachedSession.Started
			}

			cache.
				GetInstance().
				CommitSessionsTransaction()

			err := stream.Send(&metadatav1.GetSessionMetadataResponse{
				Started: started,
			})
			if err != nil {
				return err
			}

			ticker.Reset(getSessionMetadataFrequency)
		case <-stream.Context().Done():
			return nil
		}
	}
}

func (h *Handler) GetLobbySet(request *metadatav1.GetLobbySetRequest, stream grpc.ServerStreamingServer[metadatav1.GetLobbySetResponse]) error {
	response := new(metadatav1.GetLobbySetResponse)

	ticker := time.NewTicker(getLobbySetFrequency)

	for {
		select {
		case <-ticker.C:
			ticker.Stop()

			response.LobbySet = response.LobbySet[:0]

			cache.
				GetInstance().
				BeginLobbySetTransaction()

			cachedLobbySet, ok := cache.
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

					return err
				}

				if !exists {
					cache.
						GetInstance().
						CommitLobbySetTransaction()

					return ErrLobbySetDoesNotExist
				}

				var lobbySet []dto.CacheLobbySetEntity

				for _, lobby := range lobbies {
					response.LobbySet = append(response.LobbySet, &metadatav1.LobbySetUnit{
						LobbyId: lobby.ID,
						Issuer:  lobby.UserEntity.Name,
						Skin:    uint64(lobby.Skin),
						Host:    lobby.Host,
					})

					lobbySet = append(lobbySet, dto.CacheLobbySetEntity{
						ID:     lobby.ID,
						Issuer: lobby.UserEntity.Name,
						Skin:   uint64(lobby.Skin),
						Host:   lobby.Host,
					})
				}

				cache.
					GetInstance().
					EvictLobbySet(request.GetSessionId())

				cache.
					GetInstance().
					AddLobbySet(request.GetSessionId(), lobbySet)
			} else {
				for _, cachedLobby := range cachedLobbySet {
					response.LobbySet = append(response.LobbySet, &metadatav1.LobbySetUnit{
						LobbyId: cachedLobby.ID,
						Issuer:  cachedLobby.Issuer,
						Skin:    cachedLobby.Skin,
						Host:    cachedLobby.Host,
					})
				}
			}

			cache.
				GetInstance().
				CommitLobbySetTransaction()

			err := stream.Send(response)
			if err != nil {
				return err
			}

			ticker.Reset(getLobbySetFrequency)
		case <-stream.Context().Done():
			return nil
		}
	}
}

func (h *Handler) CreateLobby(ctx context.Context, request *metadatav1.CreateLobbyRequest) (*metadatav1.CreateLobbyResponse, error) {
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

	cache.
		GetInstance().
		BeginSessionsTransaction()

	_, ok = cache.
		GetInstance().
		GetSessions(request.GetSessionId())
	if !ok {
		session, exists, err := repository.
			GetSessionsRepository().
			GetByID(request.GetSessionId())
		if err != nil {
			cache.
				GetInstance().
				CommitSessionsTransaction()

			return nil, err
		}

		if !exists {
			return nil, ErrSessionDoesNotExists
		}

		cache.
			GetInstance().
			AddSessions(
				request.GetSessionId(),
				converter.ConvertSessionEntityToCacheSessionEntity(session))
	}

	cache.
		GetInstance().
		CommitSessionsTransaction()

	fmt.Println("BEFORE 1")

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
		return nil, status.Errorf(codes.AlreadyExists, ErrLobbyAlreadyExists.Error())
	}

	cache.
		GetInstance().
		BeginLobbySetTransaction()

	fmt.Println("BEFORE 2")

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
		var lobbySet []dto.CacheLobbySetEntity

		for _, lobby := range sessionLobbies {
			lobbySet = append(lobbySet, dto.CacheLobbySetEntity{
				Issuer: lobby.UserEntity.Name,
				Skin:   uint64(lobby.Skin),
				Host:   lobby.Host,
			})
		}

		cache.
			GetInstance().
			EvictLobbySet(request.GetSessionId())

		cache.
			GetInstance().
			AddLobbySet(request.GetSessionId(), lobbySet)

		if len(sessionLobbies) >= config.MAX_SESSION_USERS {
			cache.
				GetInstance().
				CommitLobbySetTransaction()

			return nil, ErrSessionHasMaxAmountOfLobbies
		}
	}

	var (
		host bool
		skin uint64
	)

	cachedLobbySet, _ := cache.
		GetInstance().
		GetLobbySet(request.GetSessionId())
	if len(cachedLobbySet) >= config.MAX_SESSION_USERS {
		cache.
			GetInstance().
			CommitLobbySetTransaction()

		return nil, ErrSessionHasMaxAmountOfLobbies
	}
	fmt.Println("BEFORE 3")

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
		host = true

		skin = uint64(rand.Intn(config.MAX_SESSION_USERS))
	} else {
		fmt.Println("creating new skin")

		var (
			lobbySet      []dto.CacheLobbySetEntity
			reservedSkins = make(map[int64]bool)
		)

		for _, lobby := range lobbies {
			lobbySet = append(lobbySet, dto.CacheLobbySetEntity{
				Issuer: lobby.UserEntity.Name,
				Skin:   uint64(lobby.Skin),
				Host:   lobby.Host,
			})

			reservedSkins[lobby.Skin] = true
		}

		var availableSkins []int64

		for i := int64(0); i < config.MAX_SESSION_USERS; i++ {
			if _, ok := reservedSkins[i]; !ok {
				availableSkins = append(availableSkins, i)
			}
		}

		skin = uint64(availableSkins[uint64(rand.Intn(len(availableSkins)))])

		cache.
			GetInstance().
			EvictLobbySet(request.GetSessionId())

		cache.
			GetInstance().
			AddLobbySet(request.GetSessionId(), lobbySet)
	}

	fmt.Println("BEFORE 4")

	cache.
		GetInstance().
		CommitLobbySetTransaction()

	err = repository.
		GetLobbiesRepository().
		InsertOrUpdate(
			dto.LobbiesRepositoryInsertOrUpdateRequest{
				UserID:    userID,
				SessionID: request.GetSessionId(),
				Host:      host,
				Skin:      uint64(skin),
			})
	if err != nil {
		return nil, err
	}

	return new(metadatav1.CreateLobbyResponse), nil
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

	repository.
		GetLobbiesRepository().
		Lock()

	lobbies, _, err = repository.
		GetLobbiesRepository().
		GetBySessionID(request.GetSessionId())
	if err != nil {
		repository.
			GetLobbiesRepository().
			Unlock()

		return nil, err
	}

	if !slices.ContainsFunc(lobbies, func(value *entity.LobbyEntity) bool {
		return value.Host && value.UserID != userID
	}) {
		var availableLobbies []*entity.LobbyEntity

		for _, lobby := range lobbies {
			if lobby.UserID != userID {
				availableLobbies = append(availableLobbies, lobby)
			}
		}

		selectedLobby := availableLobbies[rand.Intn(len(availableLobbies))]

		err = repository.
			GetLobbiesRepository().
			InsertOrUpdate(
				dto.LobbiesRepositoryInsertOrUpdateRequest{
					UserID:    selectedLobby.UserID,
					SessionID: request.GetSessionId(),
					Host:      true,
					Skin:      uint64(selectedLobby.Skin),
				})
		if err != nil {
			repository.
				GetLobbiesRepository().
				Unlock()

			return nil, err
		}
	}

	err = repository.
		GetLobbiesRepository().
		DeleteByUserIDAndSessionID(userID, request.GetSessionId())
	if err != nil {
		repository.
			GetLobbiesRepository().
			Unlock()

		return nil, err
	}

	repository.
		GetLobbiesRepository().
		Unlock()

	cache.
		GetInstance().
		BeginMetadataTransaction()

	cache.
		GetInstance().
		EvictMetadata(request.GetIssuer())

	cache.
		GetInstance().
		CommitMetadataTransaction()

	return new(metadatav1.RemoveLobbyResponse), nil
}

func (h *Handler) GetUserMetadata(request *metadatav1.GetUsersMetadataRequest, stream grpc.ServerStreamingServer[metadatav1.GetUsersMetadataResponse]) error {
	response := new(metadatav1.GetUsersMetadataResponse)

	ticker := time.NewTicker(getUserMetadataFrequency)

	for {
		select {
		case <-ticker.C:
			ticker.Stop()

			response.UserMetadata = response.UserMetadata[:0]

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
						response.UserMetadata = append(response.UserMetadata, &metadatav1.UserMetadata{
							Health:     uint64(lobby.Health),
							Skin:       uint64(lobby.Skin),
							Eliminated: lobby.Eliminated,
							Position: &metadatav1.Position{
								X: lobby.PositionX,
								Y: lobby.PositionY,
							},
						})
					}
				}
			} else {
				for _, value := range metadata {
					if value.SessionID == request.GetSessionId() {
						response.UserMetadata = append(response.UserMetadata, &metadatav1.UserMetadata{
							Health:     value.Health,
							Skin:       value.Skin,
							Eliminated: value.Eliminated,
							Position: &metadatav1.Position{
								X: value.PositionX,
								Y: value.PositionY,
							},
						})
					}
				}
			}

			cache.
				GetInstance().
				CommitMetadataTransaction()

			err := stream.Send(response)
			if err != nil {
				return err
			}

			ticker.Reset(getUserMetadataFrequency)
		case <-stream.Context().Done():
			return nil
		}
	}
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
