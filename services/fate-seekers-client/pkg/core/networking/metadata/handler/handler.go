package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/metadata/api"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/metadata/connector"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/store"
	"google.golang.org/grpc/status"
)

// PerformPingConnection performs ping connection request.
func PerformPingConnection(callback func(err error)) {
	go func() {
		fmt.Println("BEFORE REQUEST")

		_, err := connector.GetInstance().GetClient().PingConnection(
			context.Background(), &api.PingConnectionRequest{
				Issuer: store.GetRepositoryUUID(),
			})

		fmt.Println("AFTER REQUEST")

		if err != nil {
			errRaw, ok := status.FromError(err)
			if !ok {
				callback(err)

				return
			}

			callback(errors.New(errRaw.Message()))

			return
		}

		callback(nil)
	}()
}
