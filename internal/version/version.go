package version

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"

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
