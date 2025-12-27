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
var GetSessionEvents = sync.OnceValue[*sync.Map](func() *sync.Map {
	return new(sync.Map)
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
				BeginSessionsTransaction()

			cache.
				GetInstance().
				BeginMetadataTransaction()

			cache.
				GetInstance().
				BeginLobbySetTransaction()

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

				raw, ok := GetSessionEvents().Load(cachedSession.Name)
				if ok {
					sessionEvent = raw.(*dto.SessionEvent)
				} else {
					sessionEvent = new(dto.SessionEvent)

					GetSessionEvents().Store(cachedSession.Name, sessionEvent)
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
								if metadata.SessionID == key {
									if !metadata.Eliminated {
										switch sessionEvent.Name {
										case dto.EVENT_NAME_TOXIC_RAIN:
											if metadata.Health-dto.EVENT_HIT_RATE_TOXIC_RAIN >= 0 {
												metadata.Health -= dto.EVENT_HIT_RATE_TOXIC_RAIN

												if metadata.Health == 0 {
													metadata.Eliminated = true
												}
											} else {
												metadata.Health = 0
												metadata.Eliminated = true
											}
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
				CommitLobbySetTransaction()

			cache.
				GetInstance().
				CommitMetadataTransaction()

			cache.
				GetInstance().
				CommitSessionsTransaction()

			ticker.Reset(eventsTickerDuration)
		}
	}()
}
