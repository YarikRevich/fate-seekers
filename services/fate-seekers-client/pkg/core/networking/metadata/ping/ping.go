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
	cancel chan bool
}

// Run starts networking metadata ping worker.
func (nmp *NetworkingMetadataPing) Run() {
	nmp.cancel = make(chan bool)

	go func() {
		var wg sync.WaitGroup

		ticker := time.NewTicker(metadataPingTimer)

		for {
			select {
			case <-ticker.C:
				ticker.Stop()

				wg.Add(1)

				start := time.Now()

				handler.PerformPingConnection(func(err error) {
					if err != nil {
						logging.GetInstance().Error(err.Error())

						wg.Done()

						return
					}

					dispatcher.GetInstance().Dispatch(
						action.NewSetStatisticsMetadataPing(
							time.Since(start).Milliseconds()))

					wg.Done()
				})

				wg.Wait()

				ticker.Reset(metadataPingTimer)
			case <-nmp.cancel:
				ticker.Stop()

				close(nmp.cancel)

				nmp.cancel = nil

				return
			}
		}
	}()
}

// Clean performs ping stop operation.
func (nmp *NetworkingMetadataPing) Clean(callback func()) {
	if nmp.cancel == nil {
		callback()
	}

	go func() {
		nmp.cancel <- true

		callback()
	}()
}

func newNetworkingMetadataPing() *NetworkingMetadataPing {
	return &NetworkingMetadataPing{
		cancel: make(chan bool),
	}
}
