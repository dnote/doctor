package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/dnote/doctor/semver"
	"github.com/dnote/doctor/testutils"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"path/filepath"
)

func initCtx(t *testing.T, version semver.Version) Ctx {
	homePath, err := filepath.Abs("./tmp")
	if err != nil {
		t.Fatal(errors.Wrap(err, "pasrsing path").Error())
	}

	dnoteDir := fmt.Sprintf("%s/.dnote-tmp", homePath)
	if err := os.MkdirAll(dnoteDir, 0755); err != nil {
		t.Fatal(errors.Wrap(err, "setting up dir").Error())
	}

	return Ctx{
		version:      version,
		homeDirPath:  homePath,
		dnoteDirPath: dnoteDir,
	}
}

// teardownEnv cleans up the test env represented by the given context
func teardownEnv(ctx Ctx) {
	if err := os.RemoveAll(ctx.dnoteDirPath); err != nil {
		panic(err)
	}
}

func writeFile(ctx Ctx, content []byte, filename string) {
	dp, err := filepath.Abs(filepath.Join(ctx.dnoteDirPath, filename))
	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(dp, content, 0644); err != nil {
		panic(err)
	}
}

func readFile(ctx Ctx, filename string) []byte {
	path := filepath.Join(ctx.dnoteDirPath, filename)

	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return b
}

func TestIssue1(t *testing.T) {
	// set up
	ctx := initCtx(t, semver.Version{Major: 0, Minor: 4, Patch: 0})
	defer teardownEnv(ctx)

	dnote := dnoteV0_4_0V0_4_4{
		"js": bookV0_4_0V0_4_4{
			Name: "js",
			Notes: []noteV0_4_0V0_4_4{
				{
					UUID:     "uuid1",
					Content:  "content 1",
					AddedOn:  1,
					EditedOn: 2,
					Public:   false,
				},
				{
					UUID:     "uuid1",
					Content:  "content 1-edited-v1",
					AddedOn:  1,
					EditedOn: 3,
					Public:   false,
				},
				{
					UUID:     "uuid1",
					Content:  "content 1-edited-v2",
					AddedOn:  1,
					EditedOn: 4,
					Public:   false,
				},
				{
					UUID:     "uuid2",
					Content:  "content 2",
					AddedOn:  1,
					EditedOn: 2,
					Public:   false,
				},
				{
					UUID:     "uuid3",
					Content:  "content 3",
					AddedOn:  1,
					EditedOn: 2,
					Public:   false,
				},
				{
					UUID:     "uuid3",
					Content:  "content 3-edited",
					AddedOn:  1,
					EditedOn: 3,
					Public:   false,
				},
			},
		},
		"css": bookV0_4_0V0_4_4{
			Name: "css",
			Notes: []noteV0_4_0V0_4_4{
				{
					UUID:     "uuid4",
					Content:  "content 4",
					AddedOn:  1,
					EditedOn: 0,
					Public:   false,
				},
				{
					UUID:     "uuid5",
					Content:  "content 5",
					AddedOn:  1,
					EditedOn: 0,
					Public:   false,
				},
				{
					UUID:     "uuid5",
					Content:  "content 5-edited",
					AddedOn:  1,
					EditedOn: 3,
					Public:   false,
				},
			},
		},
	}

	b := testutils.MustMarshalJSON(t, dnote)
	writeFile(ctx, b, "dnote")

	// execute
	ok, err := i1.fix(ctx)
	if err != nil {
		t.Fatalf(errors.Wrap(err, "failing to fix").Error())
	}

	// test
	testutils.AssertEqual(t, ok, true, "diagnosed mismatch")

	expected := dnoteV0_4_0V0_4_4{
		"js": bookV0_4_0V0_4_4{
			Name: "js",
			Notes: []noteV0_4_0V0_4_4{
				{
					UUID:     "uuid1",
					Content:  "content 1-edited-v2",
					AddedOn:  1,
					EditedOn: 4,
					Public:   false,
				},
				{
					UUID:     "uuid2",
					Content:  "content 2",
					AddedOn:  1,
					EditedOn: 2,
					Public:   false,
				},
				{
					UUID:     "uuid3",
					Content:  "content 3-edited",
					AddedOn:  1,
					EditedOn: 3,
					Public:   false,
				},
			},
		},
		"css": bookV0_4_0V0_4_4{
			Name: "css",
			Notes: []noteV0_4_0V0_4_4{
				{
					UUID:     "uuid4",
					Content:  "content 4",
					AddedOn:  1,
					EditedOn: 0,
					Public:   false,
				},
				{
					UUID:     "uuid5",
					Content:  "content 5-edited",
					AddedOn:  1,
					EditedOn: 3,
					Public:   false,
				},
			},
		},
	}

	b = readFile(ctx, "dnote")
	var got dnoteV0_4_0V0_4_4
	testutils.MustUnmarshalJSON(t, b, &got)

	if ok := cmp.Equal(expected, got); !ok {
		t.Errorf("dnote content mismatch. diff: %s", cmp.Diff(expected, got))
	}
}
