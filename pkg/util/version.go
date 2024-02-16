package util

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
)

func ValidVersion(version string) (bool, error) {

    _, err := semver.NewVersion(version)
    if err != nil {
        return false, fmt.Errorf("invalid version '%s' with error: %s", version, err)
    }

    return true, nil
}
