package main

import (
	"regexp"
	"strconv"

	"github.com/pkg/errors"
)

type semver struct {
	Major int
	Minor int
	Patch int
}

type issue struct {
	name       string
	minVersion *semver
	maxVersion *semver
	desc       string
	fix        func() (bool, error)
}

func (s1 semver) lte(s2 semver) bool {
	return s1.Major < s2.Major || s1.Minor < s2.Minor || s1.Patch < s2.Patch
}
func (s1 semver) gte(s2 semver) bool {
	return s1.Major >= s2.Major || s1.Minor >= s2.Minor || s1.Patch >= s2.Patch
}

func (i issue) relevant(version semver) bool {
	return (i.minVersion == nil || version.gte(*i.minVersion)) &&
		(i.maxVersion == nil || version.lte(*i.maxVersion))
}

var (
	v0_2_0 = semver{Major: 0, Minor: 2, Patch: 0}
)

var i1 = issue{
	name:       "duplicate-json-note-uuid",
	minVersion: &v0_2_0,
	maxVersion: nil,
	desc: `Under 0.4.4, some notes have duplicate uuids if they were edited.
Duplicates have the same added_on but successively incrementing edited_on values.
Therefore the fix is to simply take the note with the latest edited_on value and discard the outdated ones.
The cause may be related to sync but is unknown. But will no longer possible in v0.4.5 and above because SQLite imposes uniqueness constraint on note uuid.`,
	fix: func() (bool, error) {
		return false, nil
	},
}

var issues = []issue{
	i1,
}

func parseSemver(version string) (semver, error) {
	re := regexp.MustCompile(`(\d*)\.(\d*)\.(\d*)`)
	match := re.FindStringSubmatch(version)

	if len(match) != 4 {
		return semver{}, errors.Errorf("invalid semver %s", version)
	}

	major, err := strconv.Atoi(match[1])
	if err != nil {
		return semver{}, errors.Wrap(err, "converting major version to int")
	}
	minor, err := strconv.Atoi(match[2])
	if err != nil {
		return semver{}, errors.Wrap(err, "converting minor version to int")
	}
	patch, err := strconv.Atoi(match[3])
	if err != nil {
		return semver{}, errors.Wrap(err, "converting patch version to int")
	}

	ret := semver{
		Major: major,
		Minor: minor,
		Patch: patch,
	}

	return ret, nil
}

func scanIssues(version semver) ([]issue, error) {
	var ret []issue

	for _, i := range issues {
		if i.relevant(version) {
			ret = append(ret, i)
		}
	}

	return ret, nil
}
