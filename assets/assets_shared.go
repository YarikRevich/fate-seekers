//go:build shared
// +build shared

package assets

import "embed"

var (
	// Represents shared assets to be embedded to the application executable.
	//go:embed shared
	AssetsShared embed.FS
)
