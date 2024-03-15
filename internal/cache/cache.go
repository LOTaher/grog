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

func IsPackageCached(name string) (bool, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return false, fmt.Errorf("unable to get user home directory: %w", err)
    }

    cacheDir := filepath.Join(homeDir, ".grog", "cache", name)
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

func ReadLockFile(name, version string) (LockFile, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return LockFile{}, fmt.Errorf("unable to get user home directory: %w", err)
    }

    lockFilePath := filepath.Join(homeDir, ".grog", "cache", name, version, "package", "grog-lock.json")

    file, err := os.ReadFile(lockFilePath)
    if err != nil {
        return LockFile{}, fmt.Errorf("failed to read lock file: %w", err)
    }

    var lockFile LockFile
    if err := json.Unmarshal(file, &lockFile); err != nil {
        return LockFile{}, fmt.Errorf("failed to unmarshal lock file: %w", err)
    }

    return lockFile, nil
}

func IsLockFileLatest(name, version string) (bool, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return false, fmt.Errorf("unable to get user home directory: %w", err)
    }

    lockFilePath := filepath.Join(homeDir, ".grog", "cache", name, version, "package", "grog-lock.json")

	file, err := os.ReadFile(lockFilePath)
	if err != nil {
		return false, fmt.Errorf("failed to read lock file: %w", err)
	}

	if string(file[13]) == "f" {
		return false, nil
    }

    return true, nil 
}

func GetVersions(name string) ([]string, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return nil, fmt.Errorf("unable to get user home directory: %w", err)
    }

    versionsDir := filepath.Join(homeDir, ".grog", "cache", name)
    versions, err := os.ReadDir(versionsDir)
    if err != nil {
        return nil, fmt.Errorf("failed to read directory: %w", err)
    }

    var versionStrings []string
    for _, version := range versions {
        if version.IsDir() {
            versionStrings = append(versionStrings, version.Name())
        }
    }

    return versionStrings, nil
}

func GetLatestVersion(name string) (string, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return "", fmt.Errorf("unable to get user home directory: %w", err)
    }

    versionsDir := filepath.Join(homeDir, ".grog", "cache", name)
    versions, err := os.ReadDir(versionsDir)
    if err != nil {
        return "", fmt.Errorf("failed to read directory: %w", err)
    }

    var latestVersion string
    for _, version := range versions {
        if version.IsDir() {
            latestVersion = version.Name()
        }
    }

    return latestVersion, nil
}
