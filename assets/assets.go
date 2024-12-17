package assets

import "embed"

var (
	// Represents assets to be embedded to the application executable.
	//go:embed dist
	Assets embed.FS
)
