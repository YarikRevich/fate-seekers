package ping

import (
	"sync"
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/metadata/handler"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/logging"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/dispatcher"
)

var (
	// GetInstance retrieves instance of the networking metadata ping, performing initilization if needed.
	GetInstance = sync.OnceValue[*NetworkingMetadataPing](newNetworkingMetadataPing)
)

const (
	// Describes metadata ping worker ticker duration.
	metadataPingTimer = time.Second
)

// NetworkingMetadataPing represents networking metadata ping worker.
type NetworkingMetadataPing struct {
	// Represents configured ticker, which serves to stop the worker.
	close chan bool
}

// Run starts networking metadata ping worker.
func (nmp *NetworkingMetadataPing) Run() {
	go func() {
		var wg sync.WaitGroup

		nmp.ticker.Reset(metadataPingTimer)

		for {
			select {
				case <- 
			}
		}

		for range nmp.ticker.C {
			nmp.ticker.Stop()

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

			nmp.ticker.Reset(metadataPingTimer)
		}
	}()
}

// Clean performs ping stop operation.
func (nmp *NetworkingMetadataPing) Clean() {
	close(nmp.ticker)
}

func newNetworkingMetadataPing() *NetworkingMetadataPing {
	ticker := time.NewTicker(metadataPingTimer)

	return &NetworkingMetadataPing{
		ticker: ticker,
	}
}
