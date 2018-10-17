package main

import (
	"testing"

	"github.com/dnote/doctor/semver"
	"github.com/pkg/errors"
	"path/filepath"
)

func initCtx(t *testing.T, version semver.Version) Ctx {
	homePath, err := filepath.Abs("./tmp")
	if err != nil {
		t.Fatal(errors.Wrap(err, "pasrsing path").Error())
	}

	return Ctx{
		version:     version,
		homeDirPath: homePath,
	}
}

func TestIssue1(t *testing.T) {
}
