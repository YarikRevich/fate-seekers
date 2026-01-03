package handler

import (
	"errors"
	"math"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking/cache"
	contentv1 "github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking/content/api"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/repository"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/repository/converter"
	"google.golang.org/protobuf/proto"
)

var (
	ErrUserDoesNotExist     = errors.New("err happened user does not exist")
	ErrLobbyDoesNotExist    = errors.New("err happened lobby does not exist")
	ErrUserIsEliminated     = errors.New("err happened user has been eliminated")
	ErrLobbySetDoesNotExist = errors.New("err happened lobby set does not exist")
	ErrUserIsNotInLobby     = errors.New("err happened user is not in a lobby")
)

// Handler performs content connector state management.
type Handler struct {
}

func (h *Handler) Process(key string, value []byte) error {
	switch key {
	case contentv1.UPDATE_USER_METADATA_POSITIONS:
		var message contentv1.UpdateUserMetadataPositionsRequest
		if err := proto.Unmarshal(value, &message); err != nil {
			return err
		}

		cache.
			GetInstance().
			BeginMetadataTransaction()

		metadata, ok := cache.
			GetInstance().
			GetMetadata(message.GetIssuer())
		if !ok {
			var userID int64

			cachedUserID, ok := cache.
				GetInstance().
				GetUsers(message.GetIssuer())
			if ok {
				userID = cachedUserID
			} else {
				user, exists, err := repository.
					GetUsersRepository().
					GetByName(message.GetIssuer())
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
				GetBySessionIDAndUserID(message.GetSessionId(), userID)
			if err != nil {
				cache.
					GetInstance().
					CommitMetadataTransaction()

				return err
			}

			newLobbies := converter.ConvertLobbyEntityToCacheMetadataEntity(lobbies, inventory)

			for _, newLobby := range newLobbies {
				if newLobby.LobbyID == message.GetLobbyId() {
					if newLobby.Eliminated {
						cache.
							GetInstance().
							CommitMetadataTransaction()

						return ErrUserIsEliminated
					}

					newLobby.Active = true

					newLobby.PositionX = message.GetPosition().X
					newLobby.PositionY = message.GetPosition().Y
				}
			}

			cache.
				GetInstance().
				AddMetadata(message.GetIssuer(), newLobbies)
		} else {
			for _, lobby := range metadata {
				if lobby.LobbyID == message.GetLobbyId() {
					if lobby.Eliminated {
						return ErrUserIsEliminated
					}

					lobby.Active = true

					lobby.PositionX = message.GetPosition().X
					lobby.PositionY = message.GetPosition().Y
				}
			}
		}

		cache.
			GetInstance().
			CommitMetadataTransaction()
	case contentv1.UPDATE_USER_METADATA_STATIC:
		var message contentv1.UpdateUserMetadataStaticRequest
		if err := proto.Unmarshal(value, &message); err != nil {
			return err
		}

		cache.
			GetInstance().
			BeginMetadataTransaction()

		metadata, ok := cache.
			GetInstance().
			GetMetadata(message.GetIssuer())
		if !ok {
			var userID int64

			cachedUserID, ok := cache.
				GetInstance().
				GetUsers(message.GetIssuer())
			if ok {
				userID = cachedUserID
			} else {
				user, exists, err := repository.
					GetUsersRepository().
					GetByName(message.GetIssuer())
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
				GetBySessionIDAndUserID(message.GetSessionId(), userID)
			if err != nil {
				cache.
					GetInstance().
					CommitMetadataTransaction()

				return err
			}

			newLobbies := converter.ConvertLobbyEntityToCacheMetadataEntity(lobbies, inventory)

			for _, newLobby := range newLobbies {
				if newLobby.LobbyID == message.GetLobbyId() {
					if newLobby.Eliminated {
						cache.
							GetInstance().
							CommitMetadataTransaction()

						return ErrUserIsEliminated
					}

					newLobby.PositionStatic = message.GetStatic()
				}
			}

			cache.
				GetInstance().
				AddMetadata(message.GetIssuer(), newLobbies)
		} else {
			for _, lobby := range metadata {
				if lobby.LobbyID == message.GetLobbyId() {
					if lobby.Eliminated {
						return ErrUserIsEliminated
					}

					lobby.PositionStatic = message.GetStatic()
				}
			}
		}

		cache.
			GetInstance().
			CommitMetadataTransaction()
	case contentv1.HIT_PLAYER_WITH_FIST_REQUEST:
		var message contentv1.HitPlayerWithFistRequest
		if err := proto.Unmarshal(value, &message); err != nil {
			return err
		}

		cache.
			GetInstance().
			BeginLobbySetTransaction()

		cachedLobbySet, ok := cache.
			GetInstance().
			GetLobbySet(message.GetSessionId())
		if !ok {
			lobbies, exists, err := repository.
				GetLobbiesRepository().
				GetBySessionID(message.GetSessionId())
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
				AddLobbySet(message.GetSessionId(), lobbySet)
		}

		cache.
			GetInstance().
			BeginMetadataTransaction()

		cachedMetadata, ok := cache.
			GetInstance().
			GetMetadata(message.GetIssuer())
		if !ok {
			var userID int64

			cachedUserID, ok := cache.
				GetInstance().
				GetUsers(message.GetIssuer())
			if ok {
				userID = cachedUserID
			} else {
				user, exists, err := repository.
					GetUsersRepository().
					GetByName(message.GetIssuer())
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
				GetBySessionIDAndUserID(message.GetSessionId(), userID)
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
				AddMetadata(message.GetIssuer(), metadata)
		}

		var (
			mainPositionX float64
			mainPositionY float64

			mainFound bool
		)

		for _, metadata := range cachedMetadata {
			if metadata.SessionID == message.GetSessionId() {
				mainPositionX = metadata.PositionX
				mainPositionY = metadata.PositionY

				mainFound = true
			}
		}

		if !mainFound {
			return ErrUserIsNotInLobby
		}

		for _, lobbySet := range cachedLobbySet {
			if lobbySet.Issuer != message.GetIssuer() {
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
							GetByName(message.GetIssuer())
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
						GetBySessionIDAndUserID(message.GetSessionId(), userID)
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
					if metadata.SessionID == message.GetSessionId() {
						dist := math.Hypot(metadata.PositionX-mainPositionX, metadata.PositionY-mainPositionY)

						if dist <= dto.HIT_PLAYER_WITH_FIST_DISTANCE {
							if metadata.Health-dto.HIT_PLAYER_WITH_FIST_RATE <= 0 {
								metadata.Eliminated = true
							} else {
								metadata.Health -= dto.HIT_PLAYER_WITH_FIST_RATE
							}
						}
					}
				}
			}
		}

		cache.
			GetInstance().
			CommitMetadataTransaction()

		cache.
			GetInstance().
			CommitLobbySetTransaction()
	}

	return nil
}

// NewHandler initializes Handler.
func NewHandler() *Handler {
	return new(Handler)
}
