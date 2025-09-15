package handler

import (
	"context"
	"errors"
	"fmt"

	metadatav1 "github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/metadata/api"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/metadata/common"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/metadata/connector"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/dispatcher"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/store"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrLobbyAlreadyExists          = errors.New("err happened lobby already exists")
	ErrFilteredSessionDoesNotExist = errors.New("err happened filtered session does not exist")
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
				dispatcher.
					GetInstance().
					Dispatch(
						action.NewSetStateResetApplicationAction(
							value.STATE_RESET_APPLICATION_TRUE_VALUE))

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

// PerformGetUserSessions performs user sessions retrieval request.
func PerformGetUserSessions(callback func(response *metadatav1.GetUserSessionsResponse, err error)) {
	go func() {
		response, err := connector.
			GetInstance().
			GetClient().
			GetUserSessions(
				context.Background(),
				&metadatav1.GetUserSessionsRequest{
					Issuer: store.GetRepositoryUUID(),
				})

		if err != nil {
			if status.Code(err) == codes.Unavailable {
				dispatcher.
					GetInstance().
					Dispatch(
						action.NewSetStateResetApplicationAction(
							value.STATE_RESET_APPLICATION_TRUE_VALUE))

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

// PerformGetFilteredSessions performs filtered sessions retrieval request.
func PerformGetFilteredSessions(request dto.GetFilteredSessionsRequest, callback func(response *metadatav1.GetFilteredSessionResponse, err error)) {
	go func() {
		response, err := connector.
			GetInstance().
			GetClient().
			GetFilteredSession(
				context.Background(),
				&metadatav1.GetFilteredSessionRequest{
					Name: request.Name,
				})

		if err != nil {
			if status.Code(err) == codes.Unavailable {
				dispatcher.
					GetInstance().
					Dispatch(
						action.NewSetStateResetApplicationAction(
							value.STATE_RESET_APPLICATION_TRUE_VALUE))

				dispatcher.GetInstance().Dispatch(
					action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_MENU_VALUE))

				callback(nil, common.ErrConnectionLost)

				return
			}

			if status.Code(err) == codes.NotFound {
				callback(nil, ErrFilteredSessionDoesNotExist)

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
				dispatcher.
					GetInstance().
					Dispatch(
						action.NewSetStateResetApplicationAction(
							value.STATE_RESET_APPLICATION_TRUE_VALUE))

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
				dispatcher.
					GetInstance().
					Dispatch(
						action.NewSetStateResetApplicationAction(
							value.STATE_RESET_APPLICATION_TRUE_VALUE))

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
			fmt.Println("EXTERNAL ERROR", err)

			if status.Code(err) == codes.Unavailable {
				dispatcher.
					GetInstance().
					Dispatch(
						action.NewSetStateResetApplicationAction(
							value.STATE_RESET_APPLICATION_TRUE_VALUE))

				dispatcher.GetInstance().Dispatch(
					action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_MENU_VALUE))

				callback(common.ErrConnectionLost)

				return
			}

			if status.Code(err) == codes.AlreadyExists {
				callback(ErrLobbyAlreadyExists)

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
func PerformRemoveLobby(sessionID int64, callback func(err error)) {
	go func() {
		_, err := connector.
			GetInstance().
			GetClient().
			RemoveLobby(
				context.Background(),
				&metadatav1.RemoveLobbyRequest{
					SessionId: sessionID,
					Issuer:    store.GetRepositoryUUID(),
				})

		if err != nil {
			if status.Code(err) == codes.Unavailable {
				dispatcher.
					GetInstance().
					Dispatch(
						action.NewSetStateResetApplicationAction(
							value.STATE_RESET_APPLICATION_TRUE_VALUE))

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
