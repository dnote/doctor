package main

import (
	"encoding/json"
	"io/ioutil"
	"sort"

	"github.com/dnote/doctor/semver"
	"github.com/dnote/fileutils"
	"github.com/pkg/errors"
)

type issue struct {
	name       string
	minVersion *semver.Version
	maxVersion *semver.Version
	desc       string
	fix        func(Ctx) (bool, error)
}

var (
	v0_2_0 = semver.Version{Major: 0, Minor: 4, Patch: 0}
)

var i1 = issue{
	name:       "duplicate-json-note-uuid",
	minVersion: &v0_2_0,
	maxVersion: nil,
	desc: `Under 0.4.4, some notes have duplicate uuids if they were edited.
Duplicates have the same added_on but successively incrementing edited_on values.
Therefore the fix is to simply take the note with the latest edited_on value and discard the outdated ones.
The cause may be related to sync but is unknown. But will no longer possible in v0.4.5 and above because SQLite imposes uniqueness constraint on note uuid.`,
	fix: func(ctx Ctx) (bool, error) {
		// if a legacy json dnote is not found, do not proceed
		notePath := getJSONDnotePath(ctx)
		if !fileutils.Exists(notePath) {
			return false, nil
		}

		diagnosed := false

		rawDnote, err := readJSONDnote(ctx)
		if err != nil {
			return false, errors.Wrap(err, "getting dnote")
		}

		var dnote dnoteV0_4_0V0_4_4
		if err := json.Unmarshal(rawDnote, &dnote); err != nil {
			return false, errors.Wrap(err, "unmarshalling notes")
		}

		type tmpNote struct {
			content  noteV0_4_0V0_4_4
			bookName string
		}

		// build a flat slice of notes and mark notes by bookname
		tmpNotes := []tmpNote{}
		for _, book := range dnote {
			for _, note := range book.Notes {
				tmp := tmpNote{
					content:  note,
					bookName: book.Name,
				}

				tmpNotes = append(tmpNotes, tmp)
			}
		}

		// sort by uuid
		sort.Slice(tmpNotes, func(i, j int) bool {
			return tmpNotes[i].content.UUID <= tmpNotes[j].content.UUID
		})

		// remove duplicates
		deduped := []tmpNote{}
		for i := 0; i < len(tmpNotes); {
			current := tmpNotes[i]

			if i == len(tmpNotes)-1 {
				deduped = append(deduped, current)
				break
			}

			for j := i + 1; j < len(tmpNotes); j++ {
				next := tmpNotes[j]

				if current.content.UUID != next.content.UUID {
					i++
					break
				}

				diagnosed = true
				i = j + 1

				if next.content.EditedOn > current.content.EditedOn {
					current = next
				}
			}

			deduped = append(deduped, current)
		}

		// put notes back to dnote structure
		dnote = dnoteV0_4_0V0_4_4{}
		for _, tmpNote := range deduped {
			bookName := tmpNote.bookName

			_, ok := dnote[bookName]

			var notes []noteV0_4_0V0_4_4
			if ok {
				notes = append(dnote[bookName].Notes, tmpNote.content)
			} else {
				notes = []noteV0_4_0V0_4_4{tmpNote.content}
			}

			dnote[bookName] = bookV0_4_0V0_4_4{
				Name:  bookName,
				Notes: notes,
			}
		}

		d, err := json.MarshalIndent(dnote, "", "  ")
		if err != nil {
			return diagnosed, errors.Wrap(err, "marhsalling deduplicated dnote")
		}

		err = ioutil.WriteFile(notePath, d, 0644)
		if err != nil {
			return diagnosed, errors.Wrap(err, "writing dnote file")
		}

		return diagnosed, nil
	},
}

func (i issue) relevant(version semver.Version) bool {
	return (i.minVersion == nil || version.Gte(*i.minVersion)) &&
		(i.maxVersion == nil || version.Lte(*i.maxVersion))
}

var issues = []issue{
	i1,
}
