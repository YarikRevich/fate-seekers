package handler

import (
	"context"
	"fmt"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/metadata/api"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/metadata/connector"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/store"
)

// PerformPingConnection performs ping connection request.
func PerformPingConnection(callback func(err error)) {
	go func() {
		_, err := connector.GetInstance().GetClient().PingConnection(
			context.Background(), &api.PingConnectionRequest{
				Issuer: store.GetRepositoryUUID(),
			})
		fmt.Println(err, "ERROR")
		if err != nil {
			callback(err)

			return
		}

		callback(nil)
	}()
}
