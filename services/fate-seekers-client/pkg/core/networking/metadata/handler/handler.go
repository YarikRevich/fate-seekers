package handler

import (
	"context"
	"errors"

	metadatav1 "github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/metadata/api"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/metadata/common"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/metadata/connector"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/dispatcher"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/store"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PerformPingConnection performs ping connection request.
func PerformPingConnection(callback func(err error)) {
	go func() {
		_, err := connector.
			GetInstance().
			GetClient().
			PingConnection(
				context.Background(),
				&metadatav1.PingConnectionRequest{
					Issuer: store.GetRepositoryUUID(),
				})

		if err != nil {
			if status.Code(err) == codes.Unavailable {
				callback(common.ErrConnectionLost)

				return
			}

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
			if status.Code(err) == codes.Unavailable {
				dispatcher.GetInstance().Dispatch(
					action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_MENU_VALUE))

				callback(common.ErrConnectionLost)

				return
			}

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
			if status.Code(err) == codes.Unavailable {
				dispatcher.GetInstance().Dispatch(
					action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_MENU_VALUE))

				callback(nil, common.ErrConnectionLost)

				return
			}

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
func PerformCreateSession(name string, seed uint64, callback func(err error)) {
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
			if status.Code(err) == codes.Unavailable {
				dispatcher.GetInstance().Dispatch(
					action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_MENU_VALUE))

				callback(common.ErrConnectionLost)

				return
			}

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

// PerformRemoveSession performs session removal request.
func PerformRemoveSession(sessionID int64, callback func(err error)) {
	go func() {
		_, err := connector.
			GetInstance().
			GetClient().
			RemoveSession(
				context.Background(),
				&metadatav1.RemoveSessionRequest{
					SessionId: sessionID,
					Issuer:    store.GetRepositoryUUID(),
				})

		if err != nil {
			if status.Code(err) == codes.Unavailable {
				dispatcher.GetInstance().Dispatch(
					action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_MENU_VALUE))

				callback(common.ErrConnectionLost)

				return
			}

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
func PerformGetLobbySet(sessionID int64, callback func(response *metadatav1.GetLobbySetResponse, err error)) {
	go func() {
		response, err := connector.
			GetInstance().
			GetClient().
			GetLobbySet(
				context.Background(),
				&metadatav1.GetLobbySetRequest{
					Issuer:    store.GetRepositoryUUID(),
					SessionId: sessionID,
				})

		if err != nil {
			if status.Code(err) == codes.Unavailable {
				dispatcher.GetInstance().Dispatch(
					action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_MENU_VALUE))

				callback(nil, common.ErrConnectionLost)

				return
			}

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
func PerformCreateLobby(sessionID int64, callback func(err error)) {
	go func() {
		_, err := connector.
			GetInstance().
			GetClient().
			CreateLobby(
				context.Background(),
				&metadatav1.CreateLobbyRequest{
					Issuer:    store.GetRepositoryUUID(),
					SessionId: sessionID,
				})

		if err != nil {
			if status.Code(err) == codes.Unavailable {
				dispatcher.GetInstance().Dispatch(
					action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_MENU_VALUE))

				callback(common.ErrConnectionLost)

				return
			}

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
			if status.Code(err) == codes.Unavailable {
				dispatcher.GetInstance().Dispatch(
					action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_MENU_VALUE))

				callback(common.ErrConnectionLost)

				return
			}

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
