package cache

import (
	"fmt"
	"os"
	"path/filepath"
    "encoding/json"
)

var Cache = filepath.Join(os.Getenv("HOME"), ".grog", "cache")

type LockFile struct {
	isLatest     bool              `json:"isLatest"`
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

// func ReadLockFile(name string) (LockFile, error) {
//     var lockFile LockFile
//     homeDir, err := os.UserHomeDir()
//     if err != nil {
//         return lockFile, fmt.Errorf("unable to get user home directory: %w", err)
//     }
//
//     lock := filepath.Join(homeDir, ".grog", "cache", name, "grog-lock.json")
//     if _, err := os.Stat(lock); err != nil {
//         if os.IsNotExist(err) {
//                     }
//
//         return lockFile, err
//     }
//
//     return lockFile, nil
// }

func CreateLockFile(name, version string, isLatest bool, dependencies map[string]string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("unable to get user home directory: %w", err)
	}

	lockFileDir := filepath.Join(homeDir, ".grog", "cache", name, version, "grog-lock.json")
	if _, err := os.Stat(lockFileDir); err != nil {
		if os.IsNotExist(err) {
			return err
		}
		return err
	}

	lockFile := LockFile{
		isLatest:     isLatest,
		Dependencies: dependencies,
	}

	json, err := json.Marshal(lockFile)
	if err != nil {
		return fmt.Errorf("failed to marshal lock file: %w", err)
	}

	if err := os.WriteFile(lockFileDir, json, 0644); err != nil {
		return fmt.Errorf("failed to write lock file: %w", err)
	}

	return nil
}
