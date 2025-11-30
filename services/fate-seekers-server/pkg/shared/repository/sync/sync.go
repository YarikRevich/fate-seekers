package sync

import (
	"errors"
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/logging"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking/cache"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/repository"
)

var (
	ErrUserDoesNotExist = errors.New("err happened user does not exist")
)

const (
	// Represents ticker duration used for metadata synchronization worker.
	metadataTickerDuration = time.Minute
)

var (
	// Represents a map of affected sessions, which should be returned to cache.
	affectedSessions map[int64]bool = make(map[int64]bool)
)

// Run starts the repository sync worker, which takes latest updates
// from certain cache instances.
func Run() {
	// TODO: create some mapping of hashes, which would help to avoid not necessary updates.

	go func() {
		ticker := time.NewTicker(metadataTickerDuration)

		for range ticker.C {
			ticker.Stop()

			clear(affectedSessions)

			cache.
				GetInstance().
				BeginLobbySetTransaction()

			cache.
				GetInstance().
				BeginMetadataTransaction()

			for key, value := range cache.
				GetInstance().
				GetMetadataMappings() {
				var userID int64

				cachedUserID, ok := cache.
					GetInstance().
					GetUsers(key)
				if ok {
					userID = cachedUserID
				} else {
					user, exists, err := repository.
						GetUsersRepository().
						GetByName(key)
					if err != nil {
						cache.
							GetInstance().
							CommitMetadataTransaction()

						cache.
							GetInstance().
							CommitLobbySetTransaction()

						logging.GetInstance().Fatal(err.Error())
					}

					if !exists {
						cache.
							GetInstance().
							CommitMetadataTransaction()

						cache.
							GetInstance().
							CommitLobbySetTransaction()

						logging.GetInstance().Fatal(ErrUserDoesNotExist.Error())
					}

					userID = user.ID
				}

				for _, metadata := range value {
					err := repository.
						GetLobbiesRepository().
						InsertOrUpdate(
							dto.LobbiesRepositoryInsertOrUpdateRequest{
								UserID:         userID,
								SessionID:      metadata.SessionID,
								Skin:           metadata.Skin,
								Health:         metadata.Health,
								Active:         metadata.Active,
								Eliminated:     metadata.Eliminated,
								Host:           metadata.Host,
								PositionX:      metadata.PositionX,
								PositionY:      metadata.PositionY,
								PositionStatic: metadata.PositionStatic,
							})
					if err != nil {
						cache.
							GetInstance().
							CommitMetadataTransaction()

						cache.
							GetInstance().
							CommitLobbySetTransaction()

						logging.GetInstance().Fatal(err.Error())
					}

					affectedSessions[metadata.SessionID] = true
				}
			}

			for sessionID := range affectedSessions {
				lobbies, exists, err := repository.
					GetLobbiesRepository().
					GetBySessionID(sessionID)
				if err != nil {
					cache.
						GetInstance().
						CommitMetadataTransaction()

					cache.
						GetInstance().
						CommitLobbySetTransaction()

					logging.GetInstance().Fatal(err.Error())
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
						EvictLobbySet(sessionID)

					cache.
						GetInstance().
						AddLobbySet(sessionID, lobbySet)
				}
			}

			cache.
				GetInstance().
				CommitMetadataTransaction()

			cache.
				GetInstance().
				CommitLobbySetTransaction()

			ticker.Reset(metadataTickerDuration)
		}
	}()
}
