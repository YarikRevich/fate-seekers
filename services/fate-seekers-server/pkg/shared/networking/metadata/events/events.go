package events

import (
	"errors"
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking/cache"
)

var (
	ErrUserDoesNotExist = errors.New("err happened user does not exist")
)

const (
	// Represents ticker duration used for events processing worker.
	eventsTickerDuration = time.Second * 2
)

var (
	// Represents session events holder.
	sessionEvents map[string]dto.SessionEvent = make(map[string]dto.SessionEvent)
)

// TODO: create logic to randomly select events for a session

// Run starts the repository sync worker, which takes latest updates
// from certain cache instances.
func Run() {
	go func() {
		ticker := time.NewTicker(eventsTickerDuration)

		for range ticker.C {
			ticker.Stop()

			for key, value := range cache.
				GetInstance().
				GetLobbySetMappings() {

				session, ok := cache.GetInstance().GetSessions(key)
				if !ok {

				}

				if session.Started {
					for _, lobby := range value {
						metadataSet, ok := cache.GetInstance().GetMetadata(lobby.Issuer)
						if ok {
							for _, metadata := range metadataSet {
								metadata.Health -= 20
							}
						}
					}
				}

			}

			// 	for key, value := range cache.
			// 		GetInstance().
			// 		GetMetadataMappings() {
			// 		var userID int64

			// 		cachedUserID, ok := cache.
			// 			GetInstance().
			// 			GetUsers(key)
			// 		if ok {
			// 			userID = cachedUserID
			// 		} else {
			// 			user, exists, err := repository.
			// 				GetUsersRepository().
			// 				GetByName(key)
			// 			if err != nil {
			// 				logging.GetInstance().Fatal(err.Error())
			// 			}

			// 			if !exists {
			// 				logging.GetInstance().Fatal(ErrUserDoesNotExist.Error())
			// 			}

			// 			userID = user.ID
			// 		}

			// 		for _, lobby := range value {
			// 			err := repository.
			// 				GetLobbiesRepository().
			// 				InsertOrUpdate(
			// 					dto.LobbiesRepositoryInsertOrUpdateRequest{
			// 						UserID:     userID,
			// 						SessionID:  lobby.SessionID,
			// 						Skin:       lobby.Skin,
			// 						Health:     lobby.Health,
			// 						Eliminated: lobby.Eliminated,
			// 						PositionX:  lobby.PositionX,
			// 						PositionY:  lobby.PositionY,
			// 					})
			// 			if err != nil {
			// 				logging.GetInstance().Fatal(err.Error())
			// 			}
			// 		}
			// 	}

			ticker.Reset(eventsTickerDuration)
		}
	}()
}
