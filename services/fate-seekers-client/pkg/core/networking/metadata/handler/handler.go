package handler

import (
	"context"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/metadata/api"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/metadata/connector"
)

// PerformPingConnection performs ping connection request.
func PerformPingConnection(callback func(err error)) {
	go func() {
		connector.GetInstance().GetClient().PingConnection(
			context.Background(), &api.PingConnectionRequest{})
	}()
}
