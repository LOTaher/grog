package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

var Cache = filepath.Join(os.Getenv("HOME"), ".grog", "cache")

type LockFile struct {
	IsLatest     bool              `json:"isLatest"`
	Dependencies map[string]string `json:"dependencies"`
}

func IsVersionCached(name, version string) (bool, error) {
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

func CreateLockFile(name, version string, isLatest bool, dependencies map[string]string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("unable to get user home directory: %w", err)
	}

	versionDir := filepath.Join(homeDir, ".grog", "cache", name, version, "package")
	if _, err := os.Stat(versionDir); os.IsNotExist(err) {
		if err := os.MkdirAll(versionDir, 0755); err != nil { 
			return fmt.Errorf("unable to create directory %s: %w", versionDir, err)
		}
	} else if err != nil {
		return fmt.Errorf("error checking directory %s: %w", versionDir, err)
	}

    lockFilePath := filepath.Join(versionDir, "grog-lock.json")

	lockFile := LockFile{
		IsLatest:     isLatest,
		Dependencies: dependencies,
	}

	json, err := json.Marshal(lockFile)
	if err != nil {
		return fmt.Errorf("failed to marshal lock file: %w", err)
	}

	if err := os.WriteFile(lockFilePath, json, 0644); err != nil {
		return fmt.Errorf("failed to write lock file: %w", err)
	}

	return nil
}
