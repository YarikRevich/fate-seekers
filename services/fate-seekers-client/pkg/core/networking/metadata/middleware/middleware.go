package middleware

import (
	"context"
	"errors"

	"buf.build/go/protovalidate"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// Represents all the headers used for middlewares management.
const (
	AuthenticationHeader = "authentication"
)

var (
	ErrMessageValidationFailed = errors.New("err happened message validation failed")
)

// AuthenticationMiddleware represents authentication middleware.
type AuthenticationMiddleware struct{}

func (am *AuthenticationMiddleware) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		AuthenticationHeader: config.GetSettingsNetworkingEncryptionKey(),
	}, nil
}

func (am *AuthenticationMiddleware) RequireTransportSecurity() bool {
	return false
}

// CheckValidationMiddleware represents protobuf API validation middleware.
func CheckValidationMiddleware(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	err := protovalidate.Validate(req.(proto.Message))
	if err != nil {
		return status.Errorf(codes.InvalidArgument, ErrMessageValidationFailed.Error())
	}

	return invoker(ctx, method, req, reply, cc, opts...)
}
