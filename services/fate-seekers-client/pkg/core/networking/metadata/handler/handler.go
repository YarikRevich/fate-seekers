package handler

import (
	"context"
	"errors"
	"fmt"

	metadatav1 "github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/metadata/api"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/metadata/connector"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/store"
	"google.golang.org/grpc/status"
)

// PerformPingConnection performs ping connection request.
func PerformPingConnection(callback func(err error)) {
	go func() {
		fmt.Println(store.GetRepositoryUUID(), "PERFORMING PING CONNECTION")

		_, err := connector.
			GetInstance().
			GetClient().
			PingConnection(
				context.Background(),
				&metadatav1.PingConnectionRequest{
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

// PerformCreateUserIfNotExists performs user creation attempt request.
func PerformCreateUserIfNotExists(callback func(err error)) {
	go func() {
		_, err := connector.
			GetInstance().
			GetClient().
			CreateUserIfNotExists(
				context.Background(),
				&metadatav1.CreateUserIfNotExistsRequest{
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
func PerformGetSessions(callback func(response *metadatav1.GetSessionsResponse, err error)) {
	go func() {
		response, err := connector.
			GetInstance().
			GetClient().
			GetSessions(
				context.Background(),
				&metadatav1.GetSessionsRequest{
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
		_, err := connector.
			GetInstance().
			GetClient().
			CreateSession(
				context.Background(),
				&metadatav1.CreateSessionRequest{
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

// PerformGetLobbySet performs lobby set retrieval request.
func PerformGetLobbySet(sessionId int64, callback func(response *metadatav1.GetLobbySetResponse, err error)) {
	go func() {
		response, err := connector.
			GetInstance().
			GetClient().
			GetLobbySet(
				context.Background(),
				&metadatav1.GetLobbySetRequest{
					Issuer:    store.GetRepositoryUUID(),
					SessionId: sessionId,
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

// PerformCreateLobby performs lobby creation request.
func PerformCreateLobby(sessionId int64, callback func(err error)) {
	go func() {
		_, err := connector.
			GetInstance().
			GetClient().
			CreateLobby(
				context.Background(),
				&metadatav1.CreateLobbyRequest{
					Issuer:    store.GetRepositoryUUID(),
					SessionId: sessionId,
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

// PerformRemoveLobby performs lobby removal request.
func PerformRemoveLobby(callback func(err error)) {
	go func() {
		_, err := connector.
			GetInstance().
			GetClient().
			RemoveLobby(
				context.Background(),
				&metadatav1.RemoveLobbyRequest{
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
