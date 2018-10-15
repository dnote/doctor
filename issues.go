package main

import (
	"github.com/dnote/doctor/semver"
)

type issue struct {
	name       string
	minVersion *semver.Version
	maxVersion *semver.Version
	desc       string
	fix        func() (bool, error)
}

var (
	v0_2_0 = semver.Version{Major: 0, Minor: 2, Patch: 0}
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

func (i issue) relevant(version semver.Version) bool {
	return (i.minVersion == nil || version.Gte(*i.minVersion)) &&
		(i.maxVersion == nil || version.Lte(*i.maxVersion))
}

var issues = []issue{
	i1,
}
