package handler

import (
	"context"
	"errors"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/metadata/api"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/metadata/connector"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/store"
	"google.golang.org/grpc/status"
)

// PerformPingConnection performs ping connection request.
func PerformPingConnection(callback func(err error)) {
	go func() {
		_, err := connector.GetInstance().GetClient().PingConnection(
			context.Background(), &api.PingConnectionRequest{
				Issuer: store.GetRepositoryUUID(),
			})

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

// PerformGetSessions performs sessions retrieval request.
func PerformGetSessions(callback func(response *api.GetSessionsResponse, err error)) {
	go func() {
		response, err := connector.GetInstance().GetClient().GetSessions(
			context.Background(), &api.GetSessionsRequest{
				Issuer: store.GetRepositoryUUID(),
			})

		if err != nil {
			errRaw, ok := status.FromError(err)
			if !ok {
				callback(nil, err)

				return
			}

			callback(nil, errors.New(errRaw.Message()))

			return
		}

		callback(response, nil)
	}()
}

// PerformCreateSession performs session creation request.
func PerformCreateSession(name string, seed int64, callback func(err error)) {
	go func() {
		_, err := connector.GetInstance().GetClient().CreateSession(
			context.Background(), &api.CreateSessionRequest{
				Name:   name,
				Issuer: store.GetRepositoryUUID(),
				Seed:   &seed,
			})

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
