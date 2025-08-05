package middleware

import (
	"sync"
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
	middlewares := []func(callback func() error) error{}

	return &NetworkingContentMiddlewarePipeline{
		middlewares: middlewares,
	}
}
