package handler

import (
	"errors"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking/cache"
	contentv1 "github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking/content/api"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/repository"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/repository/converter"
	"google.golang.org/protobuf/proto"
)

var (
	ErrUserDoesNotExist  = errors.New("err happened user does not exist")
	ErrLobbyDoesNotExist = errors.New("err happened lobby does not exist")
	ErrUserIsEliminated  = errors.New("err happened user has been eliminated")
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

			newLobbies := converter.ConvertLobbyEntityToCacheMetadataEntity(lobbies)

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

			newLobbies := converter.ConvertLobbyEntityToCacheMetadataEntity(lobbies)

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
	case contentv1.OPEN_GENERATED_CHEST:
	case contentv1.OPEN_GENERATED_HEALTH_PACK:
	case contentv1.SEND_CHAT_MESSAGE:
	}

	return nil
}

// NewHandler initializes Handler.
func NewHandler() *Handler {
	return new(Handler)
}
