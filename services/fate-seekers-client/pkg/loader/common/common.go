package common

import (
	"errors"
	"io/fs"
	"os"

	"github.com/YarikRevich/fate-seekers/assets"
)

// ReadFile performs file read operation using both shared and client assets.
func ReadFile(sharedPath, clientPath string) ([]byte, error) {
	file, err := fs.ReadFile(assets.AssetsShared, sharedPath)
	if errors.Is(err, os.ErrNotExist) {
		file, err = fs.ReadFile(assets.AssetsClient, clientPath)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return file, nil
}
