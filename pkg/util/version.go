package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
    "sort"

	"github.com/Masterminds/semver/v3"
    version "github.com/hashicorp/go-version"
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

	var versions Version
	if err := json.Unmarshal(body, &versions); err != nil {
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

	var parsedVersions []*version.Version
	for vStr := range versions.Versions {
		v, err := version.NewVersion(vStr)
		if err != nil {
			continue 
		}
		parsedVersions = append(parsedVersions, v)
	}

	sort.Sort(sort.Reverse(version.Collection(parsedVersions)))

	constr, err := version.NewConstraint(constraintStr)
	if err != nil {
		return "", err
	}

	for _, v := range parsedVersions {
		if constr.Check(v) {
			return v.String(), nil 
		}
	}

	return "", fmt.Errorf("no version found that satisfies the constraint '%s'", constraintStr)
}

func HasConstraintSymbols(version string) bool {
    // TODO
} 

func ValidVersion(version string) (bool, error) {

	_, err := semver.NewVersion(version)
	if err != nil {
		return false, fmt.Errorf("invalid version '%s' with error: %s", version, err)
	}

	return true, nil
}
