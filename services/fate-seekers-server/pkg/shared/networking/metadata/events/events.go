package events

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking/cache"
)

var (
	ErrUserDoesNotExist = errors.New("err happened user does not exist")
)

const (
	// Represents ticker duration used for events processing worker.
	eventsTickerDuration = time.Second * 1

	// Represents events processing time, which is used as a pause between events.
	eventsProcessingDuration = time.Second * 90
)

// GetSessionEvents retrieves instance of the session events map, performing initilization if needed.
var GetSessionEvents = sync.OnceValue[map[string]*dto.SessionEvent](func() map[string]*dto.SessionEvent {
	return make(map[string]*dto.SessionEvent)
})

// Run starts the repository sync worker, which takes latest updates
// from certain cache instances.
func Run() {
	go func() {
		ticker := time.NewTicker(eventsTickerDuration)

		for range ticker.C {
			ticker.Stop()

			fmt.Println("BEFORE ITERATION BLOCK")

			fmt.Println("BEFORE ITERATION", cache.
				GetInstance().
				GetLobbySetMappings())

			for key, value := range cache.
				GetInstance().
				GetLobbySetMappings() {

				fmt.Println("ITERATION BEGINNING")

				cachedSession, ok := cache.GetInstance().GetSessions(key)
				if !ok {
					ticker.Reset(eventsTickerDuration)

					continue
				}

				fmt.Println("ITERATION START CHECK")

				if !cachedSession.Started {
					ticker.Reset(eventsTickerDuration)

					continue
				}

				var sessionEvent *dto.SessionEvent

				sessionEvent, ok = GetSessionEvents()[cachedSession.Name]
				if !ok {
					sessionEvent = new(dto.SessionEvent)

					GetSessionEvents()[cachedSession.Name] = sessionEvent
				}

				fmt.Println(sessionEvent, "SESSION EVENT")

				if sessionEvent.EndRate.Before(time.Now()) {
					fmt.Println(sessionEvent.Name, "SESSION EVENT HAS ENDED")

					if !sessionEvent.PauseRate.Before(time.Now()) {
						sessionEvent.Name = dto.EVENT_NAME_EMPTY

						ticker.Reset(eventsTickerDuration)

						continue
					}

					if rand.Intn(2) == 0 {
						fmt.Println("SESSION EVENT HAS BEEN MISSED")

						sessionEvent.PauseRate.Add(eventsProcessingDuration)

						ticker.Reset(eventsTickerDuration)

						continue
					}

					fmt.Println("SESSION EVENT HAS BEEN SELECTED")

					selectedEvent := dto.EVENTS_NAME_MAP[rand.Intn(len(dto.EVENTS_NAME_MAP))]

					sessionEvent.Name = selectedEvent

					fmt.Println(selectedEvent, "SESSION EVENT HAS BEEN CHOSEN")

					switch selectedEvent {
					case dto.EVENT_NAME_TOXIC_RAIN:
						sessionEvent.FrequencyRate = time.Now().Add(dto.EVENT_FREQUENCY_RATE_TOXIC_RAIN)
						sessionEvent.EndRate = time.Now().Add(dto.EVENT_DURATION_TIME_TOXIC_RAIN)
						sessionEvent.PauseRate = sessionEvent.EndRate.Add(eventsProcessingDuration)
					}
				} else if sessionEvent.FrequencyRate.Before(time.Now()) {
					fmt.Println(sessionEvent.Name, "SESSION EVENT FREQ RATE")

					for _, lobby := range value {
						metadataSet, ok := cache.GetInstance().GetMetadata(lobby.Issuer)
						if ok {
							for _, metadata := range metadataSet {
								switch sessionEvent.Name {
								case dto.EVENT_NAME_TOXIC_RAIN:
									metadata.Health -= dto.EVENT_HIT_RATE_TOXIC_RAIN
								}
							}
						}
					}

					switch sessionEvent.Name {
					case dto.EVENT_NAME_TOXIC_RAIN:
						sessionEvent.FrequencyRate = time.Now().Add(dto.EVENT_FREQUENCY_RATE_TOXIC_RAIN)
					}
				}
			}

			ticker.Reset(eventsTickerDuration)
		}
	}()
}
