package version

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
    "os"
    "path/filepath"

	"github.com/Masterminds/semver/v3"
)

type Version struct {
	Versions map[string]interface{} `json:"versions"`
}

func (v *Version) reqRegistry(packageName string) error {
	url := fmt.Sprintf("https://registry.npmjs.org/%s", packageName)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &v); err != nil {
		fmt.Printf("Error unmarshaling JSON: %v\n", err)
		return err
	}

	return nil
}

func BestMatchingVersion(packageName, constraintStr string) (string, error) {
	versions := Version{}
	if err := versions.reqRegistry(packageName); err != nil {
		return "", err
	}

	var parsedVersions []*semver.Version
	for vStr := range versions.Versions {
		v, err := semver.NewVersion(vStr)
		if err != nil {
			continue
		}
		parsedVersions = append(parsedVersions, v)
	}

	sort.Slice(parsedVersions, func(i, j int) bool {
		return parsedVersions[i].LessThan(parsedVersions[j])
	})

	constraint, err := semver.NewConstraint(constraintStr)
	if err != nil {
		return "", err
	}

	for _, v := range parsedVersions {
		if constraint.Check(v) {
			return v.String(), nil
		}
	}

	return "", fmt.Errorf("no version found that satisfies the constraint '%s'", constraintStr)
}

func ValidVersion(version string) (bool, error) {
	_, err := semver.NewVersion(version)
	if err != nil {
		return false, fmt.Errorf("invalid version '%s' with error: %s", version, err)
	}

	return true, nil
}

func GetMostRecentVersion(pkg string) string {
	versions := Version{}
	if err := versions.reqRegistry(pkg); err != nil {
		fmt.Printf("Error requesting registry: %v\n", err)
		return ""
	}

	var highestVersion *semver.Version
	for version := range versions.Versions {
		parsedVersion, err := semver.NewVersion(version)
		if err != nil {
			fmt.Printf("Error parsing version '%s': %v\n", version, err)
			continue
		}

		if highestVersion == nil || parsedVersion.GreaterThan(highestVersion) {
			highestVersion = parsedVersion
		}
	}

	if highestVersion != nil {
		return highestVersion.String()
	}

	return ""
}

func IsLatestVersion(pkg, version string) (bool, error) {
    mostRecentVersion := GetMostRecentVersion(pkg)
    if mostRecentVersion == "" {
        return false, fmt.Errorf("failed to get most recent version for %s", pkg)
    }

    return version == mostRecentVersion, nil
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

func FindCorrectVersion(name, versionConstraint string) (string, error) {
	versions, err := GetVersions(name)
	if err != nil {
		return "", fmt.Errorf("failed to get versions: %w", err)
	}

	constraint, err := semver.NewConstraint(versionConstraint)
	if err != nil {
		return "", fmt.Errorf("failed to parse version constraint '%s': %w", versionConstraint, err)
	}

	var satisfyingVersions []*semver.Version
	for _, v := range versions {
		version, err := semver.NewVersion(v)
		if err != nil {
            return "", fmt.Errorf("failed to parse version '%s': %w", v, err)
		}

		if constraint.Check(version) {
			satisfyingVersions = append(satisfyingVersions, version)
		}
	}

	if len(satisfyingVersions) == 0 {
		return "", fmt.Errorf("no version found that satisfies the constraint '%s' for package %s", versionConstraint, name)
	}

	sort.Slice(satisfyingVersions, func(i, j int) bool {
		return satisfyingVersions[i].LessThan(satisfyingVersions[j])
	})

	return satisfyingVersions[len(satisfyingVersions)-1].String(), nil
}

