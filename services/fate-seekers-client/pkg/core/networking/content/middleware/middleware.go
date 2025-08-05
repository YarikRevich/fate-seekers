package middleware

import (
	"sync"
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/dispatcher"
)

var (
	// GetInstance retrieves instance of the networking content middleware pipeline,
	// performing initilization if needed.
	GetInstance = sync.OnceValue[*NetworkingContentMiddlewarePipeline](newNetworkingContentMiddlewarePipeline)
)

// NetworkingContentMiddlewarePipeline represents networking content middleware pipeline.
type NetworkingContentMiddlewarePipeline struct {
	middlewares []func(callback func() error) error
}

// Run starts middlewares pipeline execution.
func (ncmp *NetworkingContentMiddlewarePipeline) Run(callback func() error) error {
	var err error

	for _, middleware := range ncmp.middlewares {
		err = middleware(callback)
		if err != nil {
			return err
		}
	}

	return nil
}

// newNetworkingContentMiddlewarePipeline initializes NetworkingContentMiddlewarePipeline.
func newNetworkingContentMiddlewarePipeline() *NetworkingContentMiddlewarePipeline {
	middlewares := []func(callback func() error) error{
		CheckLatencyMiddleware,
	}

	return &NetworkingContentMiddlewarePipeline{
		middlewares: middlewares,
	}
}

// CheckLatencyMiddleware performs latency check.
func CheckLatencyMiddleware(callback func() error) error {
	start := time.Now()

	err := callback()

	dispatcher.GetInstance().Dispatch(
		action.NewSetStatisticsContentPing(
			time.Since(start).Milliseconds()))

	return err
}
