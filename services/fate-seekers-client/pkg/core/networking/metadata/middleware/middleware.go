package middleware

import (
	"context"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
)

// AuthenticationMiddleware represents authentication middleware.
type AuthenticationMiddleware struct{}

func (am *AuthenticationMiddleware) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": config.GetSettingsNetworkingEncryptionKey(),
	}, nil
}

func (am *AuthenticationMiddleware) RequireTransportSecurity() bool {
	return false
}
