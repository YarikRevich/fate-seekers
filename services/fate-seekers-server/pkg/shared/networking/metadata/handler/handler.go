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
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/monitoring/services"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking/cache"
	metadatav1 "github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking/metadata/api"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking/metadata/events"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking/metadata/utils"
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
	ErrSessionChestLocationsNotEnough       = errors.New("err happened session chest locations don't fulfil min chests per session amount")
	ErrSessionHealthPacksLocationsNotEnough = errors.New("err happened session health packs locations don't fulfil min health packs per session amount")
	ErrInventoryCapacityExceeded            = errors.New("err happened inventory capacity has been exceeded")
	ErrHealthPackDoesNotExist               = errors.New("err happened health pack does not exist")
	ErrFilteredSessionDoesNotExists         = errors.New("err happened filtered session does not exist")
	ErrUserIsNotLobbyHost                   = errors.New("err happened user is not a host of a lobby")
	ErrUserIsNotInLobby                     = errors.New("err happened user is not in a lobby")
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

	fmt.Println("BEFORE SESSION LOCK")

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

	services.IncAvailableSession()

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

	fmt.Println("BEFORE 4")

	cache.
		GetInstance().
		BeginLobbySetTransaction()

	cachedLobbySet, ok := cache.
		GetInstance().
		GetLobbySet(request.GetSessionId())
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

	cache.GetInstance().EvictUserSessions(request.GetIssuer())

	cache.
		GetInstance().
		CommitLobbySetTransaction()

	cache.
		GetInstance().
		CommitUserSessionsTransaction()

	cache.
		GetInstance().
		CommitSessionsTransaction()

	services.DecAvailableSession()

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

		inventory, _, err := repository.
			GetInventoryRepository().
			GetBySessionIDAndUserID(request.GetSessionId(), userID)
		if err != nil {
			cache.
				GetInstance().
				CommitMetadataTransaction()

			return nil, err
		}

		cache.
			GetInstance().
			AddMetadata(
				request.GetIssuer(),
				converter.ConvertLobbyEntityToCacheMetadataEntity(
					lobbies, inventory))

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

	if len(request.GetChestLocations()) < config.GetOperationMinChestsAmount() {
		cache.
			GetInstance().
			CommitMetadataTransaction()

		return nil, ErrSessionChestLocationsNotEnough
	}

	if len(request.GetHealthPackLocations()) < config.GetOperationMinHealthPacksAmount() {
		cache.
			GetInstance().
			CommitMetadataTransaction()

		return nil, ErrSessionHealthPacksLocationsNotEnough
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

	randomSpawnables := rand.Perm(len(request.GetSpawnables()))

	err = db.GetInstance().Transaction(func(tx *gorm.DB) error {
		chests := utils.GenerateChests(request.GetChestLocations(), sessionSeed)

		for _, chest := range chests {
			err = repository.
				GetGenerationRepository().
				InsertOrUpdateWithTransaction(tx, dto.GenerationsRepositoryInsertOrUpdateRequest{
					SessionID: request.GetSessionId(),
					Instance:  chest.Instance,
					Name:      chest.Name,
					Type:      dto.CHEST_GENERATION_TYPE,
					Active:    true,
					PositionX: float64(chest.Position.X),
					PositionY: float64(chest.Position.Y),
				})
			if err != nil {
				return err
			}

			var generation *entity.GenerationsEntity

			generation, _, err = repository.
				GetGenerationRepository().
				GetChestTypeByInstanceAndSessionIDWithTransaction(tx, chest.Instance, request.GetSessionId())
			if err != nil {
				return err
			}

			fmt.Println(generation.ID, chest.ChestItems, "INGENERATION")

			for _, chestItem := range chest.ChestItems {
				err = repository.
					GetAssociationsRepository().
					InsertOrUpdateWithTransaction(tx, dto.AssociationsRepositoryInsertOrUpdateRequest{
						SessionID:    request.GetSessionId(),
						GenerationID: generation.ID,
						Instance:     chestItem.Instance,
						Name:         chestItem.Name,
						Active:       true,
					})
				if err != nil {
					return err
				}
			}
		}

		healthPacks := utils.GenerateHealthPacks(request.GetHealthPackLocations(), sessionSeed)

		for _, healthPack := range healthPacks {
			err = repository.
				GetGenerationRepository().
				InsertOrUpdateWithTransaction(tx, dto.GenerationsRepositoryInsertOrUpdateRequest{
					SessionID: request.GetSessionId(),
					Instance:  healthPack.Instance,
					Name:      healthPack.Name,
					Type:      dto.HEALTH_PACK_GENERATION_TYPE,
					Active:    true,
					PositionX: float64(healthPack.Position.X),
					PositionY: float64(healthPack.Position.Y),
				})
			if err != nil {
				return err
			}
		}

		for i, lobby := range lobbies {
			spawnable := request.GetSpawnables()[randomSpawnables[i]]

			err = repository.
				GetLobbiesRepository().
				InsertOrUpdateWithTransaction(
					tx,
					dto.LobbiesRepositoryInsertOrUpdateRequest{
						UserID:         lobby.UserID,
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

		err = repository.
			GetSessionsRepository().
			InsertOrUpdateWithTransaction(
				tx,
				dto.SessionsRepositoryInsertOrUpdateRequest{
					ID:      request.GetSessionId(),
					Name:    sessionName,
					Issuer:  userID,
					Started: true,
				})
		if err != nil {
			return err
		}

		return nil
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

		inventory, _, err := repository.
			GetInventoryRepository().
			GetBySessionIDAndUserID(request.GetSessionId(), userID)
		if err != nil {
			cache.
				GetInstance().
				CommitMetadataTransaction()

			return err
		}

		cache.
			GetInstance().
			AddMetadata(
				request.GetIssuer(),
				converter.ConvertLobbyEntityToCacheMetadataEntity(
					lobbies, inventory))

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
				cache.
					GetInstance().
					CommitSessionsTransaction()

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
			cache.
				GetInstance().
				CommitSessionsTransaction()

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

	services.IncAvailableLobby()

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

	services.DecAvailableLobby()

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

		inventory, _, err := repository.
			GetInventoryRepository().
			GetBySessionIDAndUserID(request.GetSessionId(), userID)
		if err != nil {
			cache.
				GetInstance().
				CommitMetadataTransaction()

			return nil, err
		}

		metadata = converter.ConvertLobbyEntityToCacheMetadataEntity(lobbies, inventory)

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
				BeginLobbySetTransaction()

			cache.
				GetInstance().
				BeginMetadataTransaction()

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
						CommitMetadataTransaction()

					cache.
						GetInstance().
						CommitLobbySetTransaction()

					return err
				}

				if !exists {
					cache.
						GetInstance().
						CommitMetadataTransaction()

					cache.
						GetInstance().
						CommitLobbySetTransaction()

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
								CommitMetadataTransaction()

							cache.
								GetInstance().
								CommitLobbySetTransaction()

							return err
						}

						if !exists {
							cache.
								GetInstance().
								CommitMetadataTransaction()

							cache.
								GetInstance().
								CommitLobbySetTransaction()

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

						cache.
							GetInstance().
							CommitLobbySetTransaction()

						return err
					}

					if !exists {
						cache.
							GetInstance().
							CommitMetadataTransaction()

						cache.
							GetInstance().
							CommitLobbySetTransaction()

						return ErrLobbyDoesNotExist
					}

					inventory, _, err := repository.
						GetInventoryRepository().
						GetBySessionIDAndUserID(request.GetSessionId(), userID)
					if err != nil {
						cache.
							GetInstance().
							CommitMetadataTransaction()

						cache.
							GetInstance().
							CommitLobbySetTransaction()

						return err
					}

					metadata := converter.ConvertLobbyEntityToCacheMetadataEntity(lobbies, inventory)

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
							Inventory: converter.ConvertCacheInventoryEntityToInventory(
								metadata.Inventory),
						})
					}
				}
			}

			cache.
				GetInstance().
				CommitMetadataTransaction()

			cache.
				GetInstance().
				CommitLobbySetTransaction()

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

func (h *Handler) DropInventoryItem(context context.Context, request *metadatav1.DropInventoryItemRequest) (*metadatav1.DropInventoryItemResponse, error) {
	response := new(metadatav1.DropInventoryItemResponse)

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

	cache.GetInstance().BeginMetadataTransaction()

	err := repository.
		GetInventoryRepository().
		DeleteByUserIDAndID(request.GetInventoryId(), userID)
	if err != nil {
		cache.GetInstance().CommitMetadataTransaction()

		return nil, err
	}

	cache.GetInstance().EvictMetadata(request.GetIssuer())

	cache.GetInstance().CommitMetadataTransaction()

	return response, nil
}

func (h *Handler) TakeChestItem(context context.Context, request *metadatav1.TakeChestItemRequest) (*metadatav1.TakeChestItemResponse, error) {
	response := new(metadatav1.TakeChestItemResponse)

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

		if !session.Started {
			cache.
				GetInstance().
				CommitSessionsTransaction()

			return nil, ErrSessionNotStarted
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

			return nil, ErrSessionNotStarted
		}

		sessionName = cachedSession.Name
	}

	cache.
		GetInstance().
		CommitSessionsTransaction()

	inventoryCount, err := repository.
		GetInventoryRepository().
		CountByLobbyIDAndUserID(request.GetLobbyId(), userID)
	if err != nil {
		return nil, err
	}

	if inventoryCount >= dto.MAX_INVENTORY_CAPACITY {
		return nil, ErrInventoryCapacityExceeded
	}

	association, _, err := repository.
		GetAssociationsRepository().
		GetByID(request.GetAssociationId())
	if err != nil {
		return nil, err
	}

	err = repository.
		GetInventoryRepository().
		InsertOrUpdate(dto.InventoryRepositoryInsertOrUpdateRequest{
			UserID:    userID,
			LobbyID:   request.GetLobbyId(),
			SessionID: request.GetSessionId(),
			Name:      association.Name,
		})
	if err != nil {
		return nil, err
	}

	cache.
		GetInstance().
		BeginMetadataTransaction()

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

	inventory, _, err := repository.
		GetInventoryRepository().
		GetBySessionIDAndUserID(request.GetSessionId(), userID)
	if err != nil {
		cache.
			GetInstance().
			CommitMetadataTransaction()

		return nil, err
	}

	cache.GetInstance().EvictMetadata(request.GetIssuer())

	cache.
		GetInstance().
		AddMetadata(
			request.GetIssuer(),
			converter.ConvertLobbyEntityToCacheMetadataEntity(lobbies, inventory))

	err = repository.
		GetAssociationsRepository().
		InsertOrUpdate(dto.AssociationsRepositoryInsertOrUpdateRequest{
			ID:     request.GetAssociationId(),
			Active: false,
		})
	if err != nil {
		cache.
			GetInstance().
			CommitMetadataTransaction()

		return nil, err
	}

	cache.
		GetInstance().
		CommitMetadataTransaction()

	cache.
		GetInstance().
		EvictGeneratedChests(sessionName)

	return response, nil
}

func (h *Handler) OpenChest(context context.Context, request *metadatav1.OpenChestRequest) (*metadatav1.OpenChestResponse, error) {
	response := new(metadatav1.OpenChestResponse)

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

			return nil, err
		}

		if !exists {
			cache.
				GetInstance().
				CommitSessionsTransaction()

			return nil, ErrSessionDoesNotExists
		}

		if !session.Started {
			cache.
				GetInstance().
				CommitSessionsTransaction()

			return nil, ErrSessionNotStarted
		}

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

			return nil, ErrSessionNotStarted
		}
	}

	cache.
		GetInstance().
		CommitSessionsTransaction()

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

			return nil, err
		}

		if !exists {
			cache.
				GetInstance().
				CommitLobbySetTransaction()

			return nil, ErrLobbySetDoesNotExist
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

	var found bool

	for _, value := range cachedLobbySet {
		if value.Issuer == request.GetIssuer() {
			found = true
		}
	}

	if !found {
		cache.
			GetInstance().
			CommitLobbySetTransaction()

		return nil, ErrUserIsNotInLobby
	}

	cache.
		GetInstance().
		CommitLobbySetTransaction()

	cache.
		GetInstance().
		BeginGeneratedChestsTransaction()

	err := repository.
		GetGenerationRepository().
		InsertOrUpdate(dto.GenerationsRepositoryInsertOrUpdateRequest{
			ID:     request.GetGenerationId(),
			Active: false,
		})
	if err != nil {
		cache.
			GetInstance().
			CommitGeneratedChestsTransaction()

		return nil, err
	}

	cache.
		GetInstance().
		CommitGeneratedChestsTransaction()

	return response, nil
}

func (h *Handler) GetChests(request *metadatav1.GetChestsRequest, stream grpc.ServerStreamingServer[metadatav1.GetChestsResponse]) error {
	response := new(metadatav1.GetChestsResponse)

	ticker := time.NewTicker(getEventsFrequency)

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
			cache.
				GetInstance().
				CommitSessionsTransaction()

			return ErrSessionDoesNotExists
		}

		if !session.Started {
			cache.
				GetInstance().
				CommitSessionsTransaction()

			return ErrSessionNotStarted
		}

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
	}

	cache.
		GetInstance().
		CommitSessionsTransaction()

	for {
		select {
		case <-ticker.C:
			ticker.Stop()

			response.Chests = response.Chests[:0]

			chests, err := repository.
				GetGenerationRepository().
				GetChestTypeBySessionID(request.GetSessionId())
			if err != nil {
				return err
			}

			for _, chest := range chests {
				var chestItems []*metadatav1.ChestItem

				associations, _, err := repository.
					GetAssociationsRepository().
					GetByGenerationID(chest.ID)
				if err != nil {
					return err
				}

				for _, association := range associations {
					chestItems = append(chestItems, &metadatav1.ChestItem{
						ChestItemId: association.ID,
						Name:        association.Name,
						Active:      association.Active,
					})
				}

				response.Chests = append(response.Chests, &metadatav1.Chest{
					SessionId: chest.SessionID,
					ChestId:   chest.ID,
					Active:    chest.Active,
					Position: &metadatav1.Position{
						X: chest.PositionX,
						Y: chest.PositionY,
					},
					Instance:   chest.Instance,
					ChestItems: chestItems,
				})
			}

			err = stream.Send(response)
			if err != nil {
				return err
			}

			ticker.Reset(getEventsFrequency)
		case <-stream.Context().Done():
			return nil
		}
	}
}

func (h *Handler) OpenHealthPack(context context.Context, request *metadatav1.OpenHealthPackRequest) (*metadatav1.OpenHealthPackResponse, error) {
	response := new(metadatav1.OpenHealthPackResponse)

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

			return nil, err
		}

		if !exists {
			cache.
				GetInstance().
				CommitSessionsTransaction()

			return nil, ErrSessionDoesNotExists
		}

		if !session.Started {
			cache.
				GetInstance().
				CommitSessionsTransaction()

			return nil, ErrSessionNotStarted
		}

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

			return nil, ErrSessionNotStarted
		}
	}

	cache.
		GetInstance().
		CommitSessionsTransaction()

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

			return nil, err
		}

		if !exists {
			cache.
				GetInstance().
				CommitLobbySetTransaction()

			return nil, ErrLobbySetDoesNotExist
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

	var found bool

	for _, value := range cachedLobbySet {
		if value.Issuer == request.GetIssuer() {
			found = true
		}
	}

	if !found {
		cache.
			GetInstance().
			CommitLobbySetTransaction()

		return nil, ErrUserIsNotInLobby
	}

	cache.
		GetInstance().
		CommitLobbySetTransaction()

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

	cache.
		GetInstance().
		BeginMetadataTransaction()

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

		inventory, exists, err := repository.
			GetInventoryRepository().
			GetBySessionIDAndUserID(request.GetSessionId(), userID)
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

			return nil, ErrHealthPackDoesNotExist
		}

		metadata = converter.ConvertLobbyEntityToCacheMetadataEntity(lobbies, inventory)

		cache.
			GetInstance().
			AddMetadata(request.GetIssuer(), metadata)
	}

	for _, value := range metadata {
		if value.SessionID == request.GetSessionId() {
			if !value.Eliminated {
				for index, item := range value.Inventory {
					if item.ID == request.GetInventoryId() {
						previous := value.Health

						if item.Name == utils.CHEST_ITEM_HEALTH_PACK_TYPE {
							if value.Health+dto.HEALTH_PACK_RATE >= 100 {
								value.Health = 100
							} else {
								value.Health += dto.HEALTH_PACK_RATE
							}
						}

						err := repository.
							GetInventoryRepository().
							DeleteByUserIDAndID(request.GetInventoryId(), userID)
						if err != nil {
							value.Health = previous

							cache.GetInstance().CommitMetadataTransaction()

							return nil, err
						}

						value.Inventory = append(value.Inventory[:index], value.Inventory[index+1:]...)

						break
					}
				}
			}
		}
	}

	cache.GetInstance().CommitMetadataTransaction()

	return response, nil
}

func (h *Handler) GetHealthPacks(request *metadatav1.GetHealthPacksRequest, stream grpc.ServerStreamingServer[metadatav1.GetHealthPacksResponse]) error {
	response := new(metadatav1.GetHealthPacksResponse)

	ticker := time.NewTicker(getEventsFrequency)

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
			cache.
				GetInstance().
				CommitSessionsTransaction()

			return ErrSessionDoesNotExists
		}

		if !session.Started {
			cache.
				GetInstance().
				CommitSessionsTransaction()

			return ErrSessionNotStarted
		}

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
	}

	cache.
		GetInstance().
		CommitSessionsTransaction()

	for {
		select {
		case <-ticker.C:
			ticker.Stop()

			response.HealthPacks = response.HealthPacks[:0]

			healthPacks, err := repository.
				GetGenerationRepository().
				GetHealthPackTypeBySessionID(request.GetSessionId())
			if err != nil {
				return err
			}

			for _, healthPack := range healthPacks {
				response.HealthPacks = append(response.HealthPacks, &metadatav1.HealthPack{
					HealthPackId: healthPack.ID,
					Active:       healthPack.Active,
					Position: &metadatav1.Position{
						X: healthPack.PositionX,
						Y: healthPack.PositionY,
					},
					Instance: healthPack.Instance,
				})
			}

			err = stream.Send(response)
			if err != nil {
				return err
			}

			ticker.Reset(getEventsFrequency)
		case <-stream.Context().Done():
			return nil
		}
	}
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
			cache.
				GetInstance().
				CommitSessionsTransaction()

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

			value, ok := events.GetSessionEvents().Load(sessionName)
			if ok {
				sessionEvent := value.(*dto.SessionEvent)

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

// NewHandler initializes implementation of metadatav1.MetadataServer.
func NewHandler() metadatav1.MetadataServiceServer {
	return new(Handler)
}
