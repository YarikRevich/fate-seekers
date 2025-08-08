package ping

import (
	"sync"
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/metadata/handler"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/logging"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/dispatcher"
)

const (
	// Describes metadata ping worker ticker duration.
	metadataPingTimer = time.Second
)

// Starts networking metadata ping worker.
func Run() {
	go func() {
		var wg sync.WaitGroup

		ticker := time.NewTicker(metadataPingTimer)

		for range ticker.C {
			ticker.Stop()

			wg.Add(1)

			start := time.Now()

			handler.PerformPingConnection(func(err error) {
				if err != nil {
					logging.GetInstance().Error(err.Error())

					return
				}

				dispatcher.GetInstance().Dispatch(
					action.NewSetStatisticsMetadataPing(
						time.Since(start).Milliseconds()))

				wg.Done()
			})

			wg.Wait()

			ticker.Reset(metadataPingTimer)
		}
	}()
}
