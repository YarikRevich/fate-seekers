//go:build client
// +build client

package assets

import "embed"

var (
	// Represents client assets to be embedded to the application executable.
	//go:embed client
	AssetsClient embed.FS
)
