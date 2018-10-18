// Package semver parses and compares semantic version strings
package semver

import (
	"regexp"
	"strconv"

	"github.com/pkg/errors"
)

// Version holds parsed semver detail
type Version struct {
	Major      int
	Minor      int
	Patch      int
	PreRelease string
}

// Lte checks if the given version is less than or equal to another
func (s1 Version) Lte(s2 Version) bool {
	return s1.Major < s2.Major || s1.Minor < s2.Minor || s1.Patch < s2.Patch
}

// Gte checks if the given version is greater than or equal to another
func (s1 Version) Gte(s2 Version) bool {
	return s1.Major >= s2.Major || s1.Minor >= s2.Minor || s1.Patch >= s2.Patch
}

// Parse parses the given semver string
func Parse(version string) (Version, error) {
	re := regexp.MustCompile(`(\d+)\.(\d+)\.(\d+)-?(.*)`)
	match := re.FindStringSubmatch(version)

	if len(match) != 5 {
		return Version{}, errors.Errorf("invalid semver %s", version)
	}

	major, err := strconv.Atoi(match[1])
	if err != nil {
		return Version{}, errors.Wrap(err, "converting major version to int")
	}
	minor, err := strconv.Atoi(match[2])
	if err != nil {
		return Version{}, errors.Wrap(err, "converting minor version to int")
	}
	patch, err := strconv.Atoi(match[3])
	if err != nil {
		return Version{}, errors.Wrap(err, "converting patch version to int")
	}

	ret := Version{
		Major:      major,
		Minor:      minor,
		Patch:      patch,
		PreRelease: match[4],
	}

	return ret, nil
}
