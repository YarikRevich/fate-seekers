package activity

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
	metadataTickerDuration = time.Second * 10
)

// Run starts the repository sync worker, which takes latest updates
// from certain cache instances.
func Run() {
	go func() {
		ticker := time.NewTicker(metadataTickerDuration)

		for range ticker.C {
			ticker.Stop()

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

				for _, lobby := range value {
					err := repository.
						GetLobbiesRepository().
						InsertOrUpdate(
							dto.LobbiesRepositoryInsertOrUpdateRequest{
								UserID:     userID,
								SessionID:  lobby.SessionID,
								Skin:       lobby.Skin,
								Health:     lobby.Health,
								Eliminated: lobby.Eliminated,
								PositionX:  lobby.PositionX,
								PositionY:  lobby.PositionY,
							})
					if err != nil {
						cache.
							GetInstance().
							CommitLobbySetTransaction()

						cache.
							GetInstance().
							CommitMetadataTransaction()

						logging.GetInstance().Fatal(err.Error())
					}
				}
			}

			cache.
				GetInstance().
				CommitLobbySetTransaction()

			cache.
				GetInstance().
				CommitMetadataTransaction()

			ticker.Reset(metadataTickerDuration)
		}
	}()
}
