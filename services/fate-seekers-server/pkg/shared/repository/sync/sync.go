package sync

import (
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/logging"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking/cache"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/repository"
)

// Run starts the repository sync worker, which takes latest updates
// from certain cache instances.
func Run() {
	go func() {
		ticker := time.NewTicker(time.Second * 15)

		for range ticker.C {
			ticker.Stop()

			for key, value := range cache.GetInstance().GetMetadataMappings() {
				var userID int64

				cachedUserID, ok := cache.
					GetInstance().
					GetUsers(key)
				if ok {
					userID = cachedUserID
				} else {
					user, _, err := repository.
						GetUsersRepository().
						GetByName(key)
					if err != nil {
						logging.GetInstance().Fatal(err.Error())
					}

					userID = user.ID

					cache.
						GetInstance().
						AddUser(key, userID)
				}

				err := repository.
					GetLobbiesRepository().
					InsertOrUpdate(
						dto.LobbiesRepositoryInsertOrUpdateRequest{
							UserID:     userID,
							SessionID:  value.SessionID,
							Skin:       value.Skin,
							Health:     value.Health,
							Eliminated: value.Eliminated,
							Position:   value.Position,
						})
				if err != nil {
					logging.GetInstance().Fatal(err.Error())
				}
			}

			ticker.Reset(time.Second * 10)
		}
	}()
}
