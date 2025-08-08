package middleware

import (
	"context"

	"buf.build/go/protovalidate"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/config"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// Represents all the headers used for middlewares management.
const (
	AuthenticationHeader = "authentication"
)

var (
	ErrMissingMetadata            = errors.New("err happened missing metadata")
	ErrAuthorizationHeaderInvalid = errors.New("err happened during authorization header validation")
	ErrMessageValidationFailed    = errors.New("err happened message validation failed")
)

// CheckValidationMiddleware represents protobuf API validation middleware.
func CheckValidationMiddleware(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	err := protovalidate.Validate(req.(proto.Message))
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			errors.Wrap(err, ErrMessageValidationFailed.Error()).Error())
	}

	return handler(ctx, req)
}

// CheckAuthenticationMiddleware performs authentication middleware validation.
func CheckAuthenticationMiddleware(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, ErrMissingMetadata.Error())
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
