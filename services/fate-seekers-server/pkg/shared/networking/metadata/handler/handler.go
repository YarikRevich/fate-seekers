package handler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/db"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/entity"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking/cache"
	metadatav1 "github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking/metadata/api"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking/metadata/events"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/repository"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/repository/converter"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

var (
	ErrUserDoesNotExist                     = errors.New("err happened user does not exist")
	ErrLobbySetDoesNotExist                 = errors.New("err happened lobby set does not exist")
	ErrLobbyDoesNotExist                    = errors.New("err happened lobby does not exist")
	ErrLobbiesAmountExceedsSpawnablesAmount = errors.New("err happened lobbies exceed spawnables amount")
	ErrLobbyAlreadyStarted                  = errors.New("err happened lobby already started")
	ErrLobbyAlreadyExists                   = errors.New("err happened lobby already exists")
	ErrSessionDoesNotExists                 = errors.New("err happened session does not exist")
	ErrSessionAlreadyExists                 = errors.New("err happened session already exists")
	ErrSessionAlreadyStarted                = errors.New("err happened session already started")
	ErrSessionNotStarted                    = errors.New("err happened session has not been started yet")
	ErrFilteredSessionDoesNotExists         = errors.New("err happened filtered session does not exist")
	ErrUserIsNotLobbyHost                   = errors.New("err happened user is not a host of a lobby")
	ErrUserDoesNotOwnSession                = errors.New("err happened user does not own session")
	ErrSessionHasMaxAmountOfLobbies         = errors.New("err happened session has max amount of lobbies")
	ErrSessionHasLobbies                    = errors.New("err happened session has lobbies")
	ErrSessionMetadataRetrievalNotAllowed   = errors.New("err happened session metadata retrieval not allowed")
)

// Describes constant values used for handler management.
const (
	getSessionMetadataFrequency = time.Second * 2
	getLobbySetFrequency        = time.Second * 1
	getUserMetadataFrequency    = time.Millisecond * 25
	getChestsFrequency          = time.Second
	getHealthPacksFrequency     = time.Second
	getEventsFrequency          = time.Millisecond * 100
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
				Seed:      uint64(cachedSession.Seed),
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
				Seed:      uint64(value.Seed),
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

	cache.
		GetInstance().
		BeginSessionsTransaction()

	fmt.Println("BEFORE 3")

	cache.
		GetInstance().
		BeginLobbySetTransaction()

	cache.
		GetInstance().
		BeginUserSessionsTransaction()

	err = repository.
		GetSessionsRepository().
		InsertOrUpdate(dto.SessionsRepositoryInsertOrUpdateRequest{
			Name:   request.GetName(),
			Seed:   int64(request.GetSeed()),
			Issuer: userID,
		})
	if err != nil {
		cache.
			GetInstance().
			CommitSessionsTransaction()

		cache.
			GetInstance().
			CommitLobbySetTransaction()

		cache.
			GetInstance().
			CommitUserSessionsTransaction()

		return nil, err
	}

	cache.
		GetInstance().
		CommitSessionsTransaction()

	cache.
		GetInstance().
		CommitLobbySetTransaction()

	cache.
		GetInstance().
		CommitUserSessionsTransaction()

	return new(metadatav1.CreateSessionResponse), nil
}

func (h *Handler) RemoveSession(ctx context.Context, request *metadatav1.RemoveSessionRequest) (*metadatav1.RemoveSessionResponse, error) {
	var isCacheSessionsPresent bool

	cache.
		GetInstance().
		BeginSessionsTransaction()

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
				cache.
					GetInstance().
					CommitSessionsTransaction()

				return nil, err
			}

			if !exists {
				cache.
					GetInstance().
					CommitSessionsTransaction()

				return nil, ErrUserDoesNotExist
			}

			userID = user.ID
		}

		sessions, err := repository.
			GetSessionsRepository().
			GetByIssuer(userID)
		if err != nil {
			cache.
				GetInstance().
				CommitSessionsTransaction()

			return nil, err
		}

		if !slices.ContainsFunc(
			sessions,
			func(value *entity.SessionEntity) bool {
				return value.ID == request.GetSessionId()
			}) {
			cache.
				GetInstance().
				CommitSessionsTransaction()

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
			cache.
				GetInstance().
				CommitSessionsTransaction()

			return nil, ErrUserDoesNotExist
		}

		userID = user.ID
	}

	fmt.Println("BEFORE 4")

	cache.
		GetInstance().
		BeginLobbySetTransaction()

	cachedLobbySet, ok := cache.
		GetInstance().
		GetLobbySet(cachedUserID)
	if ok && len(cachedLobbySet) != 0 {
		cache.
			GetInstance().
			CommitSessionsTransaction()

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
			CommitSessionsTransaction()

		cache.
			GetInstance().
			CommitLobbySetTransaction()

		return nil, err
	}

	if exists {
		var lobbySet []dto.CacheLobbySetEntity

		for _, lobby := range lobbies {
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

		cache.
			GetInstance().
			CommitSessionsTransaction()

		cache.
			GetInstance().
			CommitLobbySetTransaction()

		return nil, ErrSessionHasLobbies
	}

	cache.
		GetInstance().
		BeginUserSessionsTransaction()

	err = repository.
		GetSessionsRepository().
		DeleteByID(request.GetSessionId())
	if err != nil {
		cache.
			GetInstance().
			CommitSessionsTransaction()

		cache.
			GetInstance().
			CommitLobbySetTransaction()

		cache.
			GetInstance().
			CommitUserSessionsTransaction()

		return nil, err
	}

	cache.
		GetInstance().
		CommitLobbySetTransaction()

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

	lobbies, exists, err := repository.
		GetLobbiesRepository().
		GetBySessionID(request.GetSessionId())
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

	if len(request.GetSpawnables()) < len(lobbies) {
		cache.
			GetInstance().
			CommitMetadataTransaction()

		return nil, ErrLobbiesAmountExceedsSpawnablesAmount
	}

	randomSpawnables := rand.Perm(len(request.GetSpawnables()))

	err = db.GetInstance().Transaction(func(tx *gorm.DB) error {
		for i, lobby := range lobbies {
			spawnable := request.GetSpawnables()[randomSpawnables[i]]

			err = repository.
				GetLobbiesRepository().
				InsertOrUpdateWithTransaction(
					tx,
					dto.LobbiesRepositoryInsertOrUpdateRequest{
						UserID:         userID,
						SessionID:      lobby.SessionID,
						Skin:           uint64(lobby.Skin),
						Health:         uint64(lobby.Health),
						Active:         lobby.Active,
						Eliminated:     lobby.Eliminated,
						Host:           lobby.Host,
						PositionX:      spawnable.GetX(),
						PositionY:      spawnable.GetY(),
						PositionStatic: lobby.PositionStatic,
					})
			if err != nil {
				return ErrLobbyDoesNotExist
			}

			cache.GetInstance().EvictMetadata(lobby.UserEntity.Name)
		}

		return nil
	})
	if err != nil {
		cache.
			GetInstance().
			CommitMetadataTransaction()

		return nil, err
	}

	cache.
		GetInstance().
		BeginSessionsTransaction()

	var (
		sessionName string
		sessionSeed int64
	)

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
				CommitMetadataTransaction()

			cache.
				GetInstance().
				CommitSessionsTransaction()

			return nil, err
		}

		if session.Started {
			cache.
				GetInstance().
				CommitMetadataTransaction()

			cache.
				GetInstance().
				CommitSessionsTransaction()

			return nil, ErrSessionAlreadyStarted
		}

		sessionName = session.Name
		sessionSeed = session.Seed

		cache.
			GetInstance().
			AddSessions(
				request.GetSessionId(),
				converter.ConvertSessionEntityToCacheSessionEntity(session))
	} else {
		if cachedSession.Started {
			cache.
				GetInstance().
				CommitMetadataTransaction()

			cache.
				GetInstance().
				CommitSessionsTransaction()

			return nil, ErrSessionAlreadyStarted
		}

		sessionName = cachedSession.Name
		sessionSeed = cachedSession.Seed
	}

	fmt.Println("BEFORE 5")

	cache.
		GetInstance().
		BeginLobbySetTransaction()

	cache.
		GetInstance().
		BeginUserSessionsTransaction()

	err = repository.
		GetSessionsRepository().
		InsertOrUpdate(
			dto.SessionsRepositoryInsertOrUpdateRequest{
				ID:      request.GetSessionId(),
				Name:    sessionName,
				Issuer:  userID,
				Started: true,
			})
	if err != nil {
		cache.
			GetInstance().
			CommitUserSessionsTransaction()

		cache.
			GetInstance().
			CommitLobbySetTransaction()

		cache.
			GetInstance().
			CommitMetadataTransaction()

		cache.
			GetInstance().
			CommitSessionsTransaction()

		return nil, err
	}

	cache.
		GetInstance().
		CommitUserSessionsTransaction()

	cache.
		GetInstance().
		CommitLobbySetTransaction()

	fmt.Println(sessionSeed)

	// TODO: replace with the generation based on the provided available positions.

	// chests := utils.GenerateChestPositions(sessionSeed)

	// healthPacks := utils.GenerateHealthPackPositions(sessionSeed)

	// fmt.Println(chests, len(chests), healthPacks, len(healthPacks))

	session, _, err := repository.
		GetSessionsRepository().
		GetByID(request.GetSessionId())
	if err != nil {
		cache.
			GetInstance().
			CommitMetadataTransaction()

		cache.
			GetInstance().
			CommitSessionsTransaction()

		return nil, err
	}

	cache.GetInstance().
		EvictSessions(
			request.GetSessionId())

	cache.
		GetInstance().
		AddSessions(
			request.GetSessionId(),
			converter.ConvertSessionEntityToCacheSessionEntity(session))

	cache.
		GetInstance().
		CommitMetadataTransaction()

	cache.
		GetInstance().
		CommitSessionsTransaction()

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

			fmt.Println("BEFORE LOBBY SET LOCK")

			fmt.Println("BEFORE 6")

			cache.
				GetInstance().
				BeginLobbySetTransaction()

			fmt.Println("AFTER LOBBY SET LOCK")

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
		cache.
			GetInstance().
			BeginSessionsTransaction()

		var cachedSession dto.CacheSessionEntity

		cachedSession, ok = cache.
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

			if session.Started {
				cache.
					GetInstance().
					CommitSessionsTransaction()

				return nil, status.Errorf(codes.Aborted, ErrLobbyAlreadyStarted.Error())
			}
		} else {
			if cachedSession.Started {
				cache.
					GetInstance().
					CommitSessionsTransaction()

				return nil, status.Errorf(codes.Aborted, ErrLobbyAlreadyStarted.Error())
			}
		}

		cache.
			GetInstance().
			CommitSessionsTransaction()

		return nil, status.Errorf(codes.AlreadyExists, ErrLobbyAlreadyExists.Error())
	}

	cache.
		GetInstance().
		BeginSessionsTransaction()

	var cachedSession dto.CacheSessionEntity

	cachedSession, ok = cache.
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

		if session.Started {
			cache.
				GetInstance().
				CommitSessionsTransaction()

			return nil, status.Errorf(codes.InvalidArgument, ErrSessionAlreadyStarted.Error())
		}
	} else {
		if cachedSession.Started {
			cache.
				GetInstance().
				CommitSessionsTransaction()

			return nil, status.Errorf(codes.InvalidArgument, ErrSessionAlreadyStarted.Error())
		}
	}

	cache.
		GetInstance().
		CommitSessionsTransaction()

	fmt.Println("BEFORE 7")

	cache.
		GetInstance().
		BeginLobbySetTransaction()

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
		var (
			lobbySet      []dto.CacheLobbySetEntity
			reservedSkins = make(map[int64]bool)
		)

		for _, lobby := range lobbies {
			lobbySet = append(lobbySet, dto.CacheLobbySetEntity{
				ID:     lobby.ID,
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
		cache.
			GetInstance().
			CommitLobbySetTransaction()

		return nil, err
	}

	cache.
		GetInstance().
		CommitLobbySetTransaction()

	cache.
		GetInstance().
		BeginMetadataTransaction()

	cache.
		GetInstance().
		EvictMetadata(request.GetIssuer())

	cache.
		GetInstance().
		CommitMetadataTransaction()

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
	}

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

			return nil, err
		}

		if session.Started {
			cache.
				GetInstance().
				CommitSessionsTransaction()

			return nil, ErrSessionAlreadyStarted
		}

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
	}

	cache.
		GetInstance().
		CommitSessionsTransaction()

	fmt.Println("BEFORE 8")

	cache.
		GetInstance().
		BeginLobbySetTransaction()

	lobbies, exists, err := repository.
		GetLobbiesRepository().
		GetByUserID(userID)
	if err != nil {
		cache.
			GetInstance().
			CommitLobbySetTransaction()

		return nil, err
	}

	if exists {
		for _, lobby := range lobbies {
			if lobby.SessionID == request.GetSessionId() {
				cache.
					GetInstance().
					EvictLobbySet(lobby.SessionID)

				break
			}
		}
	}

	repository.
		GetLobbiesRepository().
		Lock()

	err = db.BeginTransaction(func(transaction *gorm.DB) error {
		lobbies, _, err = repository.
			GetLobbiesRepository().
			GetBySessionID(request.GetSessionId())
		if err != nil {
			return err
		}

		if slices.ContainsFunc(lobbies, func(value *entity.LobbyEntity) bool {
			return value.Host && value.UserID == userID
		}) {
			var availableLobbies []*entity.LobbyEntity

			for _, lobby := range lobbies {
				if lobby.UserID != userID {
					availableLobbies = append(availableLobbies, lobby)
				}
			}

			if len(availableLobbies) > 0 {
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
					return err
				}
			}
		}

		err = repository.
			GetLobbiesRepository().
			DeleteByUserIDAndSessionID(userID, request.GetSessionId())
		if err != nil {
			return err
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

		return nil
	})
	if err != nil {
		repository.
			GetLobbiesRepository().
			Unlock()

		cache.
			GetInstance().
			CommitLobbySetTransaction()

		return nil, err
	}

	repository.
		GetLobbiesRepository().
		Unlock()

	cache.
		GetInstance().
		CommitLobbySetTransaction()

	return new(metadatav1.RemoveLobbyResponse), nil
}

func (h *Handler) LeaveLobby(context context.Context, request *metadatav1.LeaveLobbyRequest) (*metadatav1.LeaveLobbyResponse, error) {
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

				return nil, err
			}

			if !exists {
				cache.
					GetInstance().
					CommitMetadataTransaction()

				return nil, ErrUserDoesNotExist
			}

			userID = user.ID
		}

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

		metadata = converter.ConvertLobbyEntityToCacheMetadataEntity(lobbies)

		cache.
			GetInstance().
			AddMetadata(request.GetIssuer(), metadata)

		var found bool

		for _, value := range metadata {
			if value.SessionID == request.GetSessionId() {
				found = true

				value.Active = false

				break
			}
		}

		if !found {
			cache.
				GetInstance().
				CommitMetadataTransaction()

			return nil, ErrSessionMetadataRetrievalNotAllowed
		}
	} else {
		var found bool

		for _, value := range metadata {
			if value.SessionID == request.GetSessionId() {
				found = true

				value.Active = false

				break
			}
		}

		if !found {
			cache.
				GetInstance().
				CommitMetadataTransaction()

			return nil, ErrSessionMetadataRetrievalNotAllowed
		}
	}

	cache.
		GetInstance().
		CommitMetadataTransaction()

	return new(metadatav1.LeaveLobbyResponse), nil
}

func (h *Handler) GetUsersMetadata(request *metadatav1.GetUsersMetadataRequest, stream grpc.ServerStreamingServer[metadatav1.GetUsersMetadataResponse]) error {
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

			fmt.Println("BEFORE 9")

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

					cache.
						GetInstance().
						CommitMetadataTransaction()

					return err
				}

				if !exists {
					cache.
						GetInstance().
						CommitLobbySetTransaction()

					cache.
						GetInstance().
						CommitMetadataTransaction()

					return ErrLobbySetDoesNotExist
				}

				var lobbySet []dto.CacheLobbySetEntity

				for _, lobby := range lobbies {
					lobbySet = append(lobbySet, dto.CacheLobbySetEntity{
						ID:     lobby.ID,
						Issuer: lobby.UserEntity.Name,
						Skin:   uint64(lobby.Skin),
						Host:   lobby.Host,
					})
				}

				cachedLobbySet = lobbySet

				cache.
					GetInstance().
					AddLobbySet(request.GetSessionId(), lobbySet)
			}

			for _, lobbySet := range cachedLobbySet {
				cachedMetadata, ok := cache.
					GetInstance().
					GetMetadata(lobbySet.Issuer)
				if !ok {
					var userID int64

					cachedUserID, ok := cache.
						GetInstance().
						GetUsers(lobbySet.Issuer)
					if ok {
						userID = cachedUserID
					} else {
						user, exists, err := repository.
							GetUsersRepository().
							GetByName(lobbySet.Issuer)
						if err != nil {
							cache.
								GetInstance().
								CommitLobbySetTransaction()

							cache.
								GetInstance().
								CommitMetadataTransaction()

							return err
						}

						if !exists {
							cache.
								GetInstance().
								CommitLobbySetTransaction()

							cache.
								GetInstance().
								CommitMetadataTransaction()

							return ErrUserDoesNotExist
						}

						userID = user.ID
					}

					lobbies, exists, err := repository.
						GetLobbiesRepository().
						GetByUserID(userID)
					if err != nil {
						cache.
							GetInstance().
							CommitLobbySetTransaction()

						cache.
							GetInstance().
							CommitMetadataTransaction()

						return err
					}

					if !exists {
						cache.
							GetInstance().
							CommitLobbySetTransaction()

						cache.
							GetInstance().
							CommitMetadataTransaction()

						return ErrLobbyDoesNotExist
					}

					metadata := converter.ConvertLobbyEntityToCacheMetadataEntity(lobbies)

					cachedMetadata = metadata

					cache.
						GetInstance().
						AddMetadata(lobbySet.Issuer, metadata)
				}

				for _, metadata := range cachedMetadata {
					if metadata.SessionID == request.GetSessionId() {
						response.UserMetadata = append(response.UserMetadata, &metadatav1.UserMetadata{
							Issuer:     lobbySet.Issuer,
							Health:     metadata.Health,
							Skin:       metadata.Skin,
							Active:     metadata.Active,
							Eliminated: metadata.Eliminated,
							Position: &metadatav1.Position{
								X: metadata.PositionX,
								Y: metadata.PositionY,
							},
							Static: metadata.PositionStatic,
						})
					}
				}
			}

			cache.
				GetInstance().
				CommitLobbySetTransaction()

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

func (h *Handler) GetChests(request *metadatav1.GetChestsRequest, stream grpc.ServerStreamingServer[metadatav1.GetChestsResponse]) error {
	return nil

	// response := new(metadatav1.GetEventsResponse)

	// ticker := time.NewTicker(getEventsFrequency)

	// var sessionName string

	// cache.
	// 	GetInstance().
	// 	BeginSessionsTransaction()

	// cachedSession, ok := cache.
	// 	GetInstance().
	// 	GetSessions(request.GetSessionId())

	// if !ok {
	// 	session, exists, err := repository.
	// 		GetSessionsRepository().
	// 		GetByID(request.GetSessionId())
	// 	if err != nil {
	// 		cache.
	// 			GetInstance().
	// 			CommitSessionsTransaction()

	// 		return err
	// 	}

	// 	if !exists {
	// 		return ErrSessionDoesNotExists
	// 	}

	// 	if !session.Started {
	// 		cache.
	// 			GetInstance().
	// 			CommitSessionsTransaction()

	// 		return ErrSessionNotStarted
	// 	}

	// 	sessionName = session.Name

	// 	cache.
	// 		GetInstance().
	// 		AddSessions(
	// 			request.GetSessionId(),
	// 			converter.ConvertSessionEntityToCacheSessionEntity(session))
	// } else {
	// 	if !cachedSession.Started {
	// 		cache.
	// 			GetInstance().
	// 			CommitSessionsTransaction()

	// 		return ErrSessionNotStarted
	// 	}

	// 	sessionName = cachedSession.Name
	// }

	// for {
	// 	select {
	// 	case <-ticker.C:
	// 		ticker.Stop()

	// 		response.Name = ""

	// 		sessionEvent, ok := events.GetSessionEvents()[sessionName]
	// 		if ok {
	// 			response.Name = sessionEvent.Name
	// 		}

	// 		err := stream.Send(response)
	// 		if err != nil {
	// 			return err
	// 		}

	// 		ticker.Reset(getEventsFrequency)
	// 	case <-stream.Context().Done():
	// 		return nil
	// 	}
	// }
}

func (h *Handler) GetHealthPacks(request *metadatav1.GetHealthPacksRequest, stream grpc.ServerStreamingServer[metadatav1.GetHealthPacksResponse]) error {

	return nil
}

func (h *Handler) GetEvents(request *metadatav1.GetEventsRequest, stream grpc.ServerStreamingServer[metadatav1.GetEventsResponse]) error {
	response := new(metadatav1.GetEventsResponse)

	ticker := time.NewTicker(getEventsFrequency)

	var sessionName string

	cache.
		GetInstance().
		BeginSessionsTransaction()

	cachedSession, ok := cache.
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

			return err
		}

		if !exists {
			return ErrSessionDoesNotExists
		}

		if !session.Started {
			cache.
				GetInstance().
				CommitSessionsTransaction()

			return ErrSessionNotStarted
		}

		sessionName = session.Name

		cache.
			GetInstance().
			AddSessions(
				request.GetSessionId(),
				converter.ConvertSessionEntityToCacheSessionEntity(session))
	} else {
		if !cachedSession.Started {
			cache.
				GetInstance().
				CommitSessionsTransaction()

			return ErrSessionNotStarted
		}

		sessionName = cachedSession.Name
	}

	cache.
		GetInstance().
		CommitSessionsTransaction()

	for {
		select {
		case <-ticker.C:
			ticker.Stop()

			response.Name = ""

			sessionEvent, ok := events.GetSessionEvents()[sessionName]
			if ok {
				response.Name = sessionEvent.Name
			}

			err := stream.Send(response)
			if err != nil {
				return err
			}

			ticker.Reset(getEventsFrequency)
		case <-stream.Context().Done():
			return nil
		}
	}
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
