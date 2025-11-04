package events

import (
	"errors"
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
	eventsProcessingDuration = time.Minute * 1
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

			cache.
				GetInstance().
				BeginMetadataTransaction()

			cache.
				GetInstance().
				BeginLobbySetTransaction()

			cache.
				GetInstance().
				BeginSessionsTransaction()

			for key, value := range cache.
				GetInstance().
				GetLobbySetMappings() {

				cachedSession, ok := cache.GetInstance().GetSessions(key)
				if !ok {
					ticker.Reset(eventsTickerDuration)

					continue
				}

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

				if sessionEvent.EndRate.Before(time.Now()) {
					if sessionEvent.Name != dto.EVENT_NAME_EMPTY {
						sessionEvent.Name = dto.EVENT_NAME_EMPTY
					}

					if sessionEvent.PauseRate.IsZero() {
						sessionEvent.PauseRate = time.Now().Add(eventsProcessingDuration)
					}

					if sessionEvent.PauseRate.After(time.Now()) {
						ticker.Reset(eventsTickerDuration)

						continue
					}

					if rand.Intn(2) == 0 {
						sessionEvent.PauseRate = time.Now().Add(eventsProcessingDuration)

						ticker.Reset(eventsTickerDuration)

						continue
					}

					selectedEvent := dto.EVENTS_NAME_MAP[rand.Intn(len(dto.EVENTS_NAME_MAP))]

					sessionEvent.Name = selectedEvent

					switch selectedEvent {
					case dto.EVENT_NAME_TOXIC_RAIN:
						sessionEvent.FrequencyRate = time.Now().Add(dto.EVENT_FREQUENCY_RATE_TOXIC_RAIN)
						sessionEvent.EndRate = time.Now().Add(dto.EVENT_DURATION_TIME_TOXIC_RAIN)
						sessionEvent.PauseRate = sessionEvent.EndRate.Add(eventsProcessingDuration)
					}
				} else if sessionEvent.FrequencyRate.Before(time.Now()) {
					for _, lobby := range value {
						metadataSet, ok := cache.
							GetInstance().
							GetMetadata(lobby.Issuer)
						if ok {
							for _, metadata := range metadataSet {
								if !metadata.Eliminated {
									switch sessionEvent.Name {
									case dto.EVENT_NAME_TOXIC_RAIN:
										if metadata.Health-dto.EVENT_HIT_RATE_TOXIC_RAIN >= 0 {
											metadata.Health -= dto.EVENT_HIT_RATE_TOXIC_RAIN
										} else if metadata.Health != 0 {
											metadata.Health = 0
											metadata.Eliminated = true
										}
									}
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

			cache.
				GetInstance().
				CommitSessionsTransaction()

			cache.
				GetInstance().
				CommitLobbySetTransaction()

			cache.
				GetInstance().
				CommitMetadataTransaction()

			ticker.Reset(eventsTickerDuration)
		}
	}()
}
