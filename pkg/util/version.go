package util

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
)

// Checks to see if a package's version is semver standards.
func ValidVersion(version string) (bool, error) {

    _, err := semver.NewVersion(version)
    if err != nil {
        return false, fmt.Errorf("Invalid version '%s' with error: %s\n", version, err)
    }

    return true, nil
}
