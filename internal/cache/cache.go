package cache

import (
	"fmt"
	"os"
	"path/filepath"
)

var Cache = filepath.Join(os.Getenv("HOME"), ".grog", "cache")

func PackageCached(name, version string) (bool, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false, fmt.Errorf("unable to get user home directory: %w", err)
	}

	cacheDir := filepath.Join(homeDir, ".grog", "cache", name, version)
	if _, err := os.Stat(cacheDir); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
