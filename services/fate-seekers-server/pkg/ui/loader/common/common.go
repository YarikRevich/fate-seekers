package common

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/Frabjous-Studios/asebiten"
	"github.com/YarikRevich/fate-seekers/assets"
)

// Represents all the available base paths.
const (
	ServerBasePath = "server"
	SharedBasePath = "shared"
)

// ReadFile performs file read operation using both shared and client assets.
func ReadFile(path string) ([]byte, error) {
	file, err := fs.ReadFile(assets.AssetsShared, filepath.Join(SharedBasePath, path))
	if errors.Is(err, os.ErrNotExist) {
		file, err = fs.ReadFile(assets.AssetsServer, filepath.Join(ServerBasePath, path))
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return file, nil
}

// LoadAnimation performs asebiten load animation operation.
func LoadAnimation(path string) (*asebiten.Animation, error) {
	animation, err := asebiten.LoadAnimation(
		assets.AssetsShared, filepath.Join(SharedBasePath, path))
	if errors.Is(err, os.ErrNotExist) {
		animation, err = asebiten.LoadAnimation(
			assets.AssetsServer, filepath.Join(ServerBasePath, path))
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return animation, nil
}
