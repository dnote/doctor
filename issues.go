package main

type issue struct {
	name       string
	minVersion string
	maxVersion string
	desc       string
	fix        func() error
}

var i1 = issue{
	name:       "duplicate-note-uuid-0.4.4",
	minVersion: "0.2.0",
	maxVersion: "",
	desc: `Under 0.4.4, some notes have duplicate uuids if they were edited.
Duplicates have the same added_on but successively incrementing edited_on values.
Therefore the fix is to simply take the note with the latest edited_on value and discard the outdated ones.
The cause may be related to sync but is unknown. But will no longer possible in v0.4.5 and above because SQLite imposes uniqueness constraint on note uuid.`,
	fix: func() error {
		return nil
	},
}
