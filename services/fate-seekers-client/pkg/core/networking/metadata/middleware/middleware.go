package middleware

import (
	"context"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
)

// Represents all the headers used for middlewares management.
const (
	AuthenticationHeader = "authentication"
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
