//go:build server
// +build server

package assets

import "embed"

var (
	// Represents server assets to be embedded to the application executable.
	//go:embed server
	AssetsServer embed.FS
)
