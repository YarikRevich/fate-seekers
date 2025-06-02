package middleware

import (
	"context"
	"errors"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Represents all the headers used for middlewares management.
const (
	AuthenticationHeader = "authentication"
)

var (
	ErrAuthorizationHeaderInvalid = errors.New("err happened during authorization header validation")
)

// CheckAuthenticationMiddleware performs authentication middleware validation.
func CheckAuthenticationMiddleware(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
	}

	if _, ok := md[AuthenticationHeader]; !ok {
		return nil, status.Errorf(codes.Unauthenticated, ErrAuthorizationHeaderInvalid.Error())
	}

	if len(md[AuthenticationHeader]) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, ErrAuthorizationHeaderInvalid.Error())
	}

	if md[AuthenticationHeader][0] != config.GetSettingsNetworkingEncryptionKey() {
		return nil, status.Errorf(codes.Unauthenticated, ErrAuthorizationHeaderInvalid.Error())
	}

	return handler(ctx, req)
}
