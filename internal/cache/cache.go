package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	ver "github.com/LOTaher/grog/internal/version"
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

func RemovePackageDependenciesGlobally(name, version string) error {

	if len(version) == 1 {
		v := ver.GetMostRecentVersion(name)
		version = v
	}

	if strings.ContainsAny(version, "<>~^=") {
		foundVersion, err := ver.BestMatchingVersion(name, version)
		if err != nil {
			return fmt.Errorf("failed to resolve version constraint '%s' : %w", version, err)
		}
		version = foundVersion
	}

	lockfile, err := ReadLockFile(name, version)
	if err != nil {
		return fmt.Errorf("failed to read lock file: %w", err)
	}

	if lockfile.Dependencies != nil {
		for dep, ver := range lockfile.Dependencies {
			os.RemoveAll("./node_modules/" + dep)
			os.RemoveAll(Cache + "/" + dep)
			if err := RemovePackageDependenciesGlobally(dep, ver); err != nil {
				return fmt.Errorf("failed to remove dependencies: %w", err)
			}
		}
	}

	return nil
}

func RemovePackageDependenciesLocally(name, version string) error {

	if len(version) == 1 {
		v := ver.GetMostRecentVersion(name)
		version = v
	}

	if strings.ContainsAny(version, "<>~^=") {
		foundVersion, err := ver.BestMatchingVersion(name, version)
		if err != nil {
			return fmt.Errorf("failed to resolve version constraint '%s' : %w", version, err)
		}
		version = foundVersion
	}

	lockfile, err := ReadLockFile(name, version)
	if err != nil {
		return fmt.Errorf("failed to read lock file: %w", err)
	}

	if lockfile.Dependencies != nil {
		for dep, ver := range lockfile.Dependencies {
			os.RemoveAll("./node_modules/" + dep)
			if err := RemovePackageDependenciesLocally(dep, ver); err != nil {
				return fmt.Errorf("failed to remove dependencies: %w", err)
			}
		}
	}

	return nil
}
