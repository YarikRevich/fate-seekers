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

// Run starts the repository sync worker, which takes latest updates
// from certain cache instances.
func Run() {

	// TODO: create some mapping of hashes, which would help to avoid not necessary updates.

	go func() {
		ticker := time.NewTicker(metadataTickerDuration)

		for range ticker.C {
			ticker.Stop()

			// cache.
			// 	GetInstance().
			// 	BeginMetadataTransaction()

			// cache.
			// 	GetInstance().
			// 	BeginLobbySetTransaction()

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
							CommitLobbySetTransaction()

						cache.
							GetInstance().
							CommitMetadataTransaction()

						logging.GetInstance().Fatal(err.Error())
					}

					if !exists {
						// cache.
						// 	GetInstance().
						// 	CommitLobbySetTransaction()

						// cache.
						// 	GetInstance().
						// 	CommitMetadataTransaction()

						logging.GetInstance().Fatal(ErrUserDoesNotExist.Error())
					}

					userID = user.ID

					cache.
						GetInstance().
						AddUser(key, userID)
				}

				for _, metadata := range value {
					err := repository.
						GetLobbiesRepository().
						InsertOrUpdate(
							dto.LobbiesRepositoryInsertOrUpdateRequest{
								UserID:     userID,
								SessionID:  metadata.SessionID,
								Skin:       metadata.Skin,
								Health:     metadata.Health,
								Active:     metadata.Active,
								Eliminated: metadata.Eliminated,
								Host:       metadata.Host,
								PositionX:  metadata.PositionX,
								PositionY:  metadata.PositionY,
							})
					if err != nil {
						// cache.
						// 	GetInstance().
						// 	CommitLobbySetTransaction()

						// cache.
						// 	GetInstance().
						// 	CommitMetadataTransaction()

						logging.GetInstance().Fatal(err.Error())
					}
				}
			}

			// cache.
			// 	GetInstance().
			// 	CommitLobbySetTransaction()

			// cache.
			// 	GetInstance().
			// 	CommitMetadataTransaction()

			ticker.Reset(metadataTickerDuration)
		}
	}()
}
