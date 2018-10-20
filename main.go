package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"time"

	"github.com/dnote/doctor/semver"
	"github.com/dnote/fileutils"
	"github.com/pkg/errors"
)

var (
	homeDirPath = flag.String("homeDir", "", "the full path to the home directory")
)

const (
	backupModeCopy = iota
	backupModeRename
)

var versionTag = "master"
var helpText = `dnote-doctor

Automatically diagnose and fix any issues with local dnote copy.

Usage:
$ dnote-doctor
`

// Ctx holds runtime configuration of dnote doctor
type Ctx struct {
	version      semver.Version
	homeDirPath  string
	dnoteDirPath string
}

func debug(msg string, v ...interface{}) {
	if os.Getenv("DNOTE_DOCTOR_DEBUG") == "1" {
		fmt.Printf("DEBUG: %s\n", fmt.Sprintf(msg, v...))
	}
}

func getDnoteDirPath() string {
	return fmt.Sprintf("%s/.dnote", *homeDirPath)
}

// backupDnoteDir backs up the dnote directory to a temporary backup directory
func backupDnoteDir(mode int) (string, error) {
	dnoteDirPath := getDnoteDirPath()
	backupName := fmt.Sprintf(".dnote-backup-%d", time.Now().UnixNano())
	backupPath := fmt.Sprintf("%s/%s", *homeDirPath, backupName)

	debug("backing up %s to %s", dnoteDirPath, backupPath)

	var err error
	switch mode {
	case backupModeCopy:
		err = fileutils.CopyDir(dnoteDirPath, backupPath)
	case backupModeRename:
		err = os.Rename(dnoteDirPath, backupPath)
	}

	if err != nil {
		return backupPath, errors.Wrapf(err, "backing up %s using %d mode", dnoteDirPath, mode)
	}

	return backupPath, nil

}

func restoreBackup(backupPath string) error {
	var err error

	defer func() {
		if err != nil {
			fmt.Printf(`Failed to restore backup from dnote doctor.
	Don't worry. Your data is still intact in the backup.
	Please reach out on https://github.com/dnote/cli/issues so that we can help you.`)
		}
	}()

	srcPath := getDnoteDirPath()
	debug("restoring %s to %s", backupPath, srcPath)

	if err = os.RemoveAll(srcPath); err != nil {
		return errors.Wrapf(err, "Failed to clear current dnote data at %s", backupPath)
	}

	if err = os.Rename(backupPath, srcPath); err != nil {
		return errors.Wrap(err, `Failed to copy backup data to the original directory.`)
	}

	return nil
}

func fixIssue(i issue, ctx Ctx) (bool, error) {
	_, err := backupDnoteDir(backupModeCopy)
	if err != nil {
		return false, errors.Wrap(err, "backing up dnote")
	}

	ok, err := i.fix(ctx)
	if err != nil {
		return false, errors.Wrap(err, "diagnosing")
	}

	return ok, nil
}

func scanIssues(version semver.Version) ([]issue, error) {
	var ret []issue

	for _, i := range issues {
		if i.relevant(version) {
			ret = append(ret, i)
		}
	}

	return ret, nil
}

func checkVersion() (semver.Version, error) {
	var ret semver.Version

	backupPath, err := backupDnoteDir(backupModeRename)
	if err != nil {
		return ret, errors.Wrap(err, "backing up dnote")
	}

	cmd := exec.Command("dnote", "version")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return ret, errors.Wrap(err, "running dnote version")
	}

	out := stdout.String()
	r := regexp.MustCompile(`dnote (\d+\.\d+\.\d+)`)
	matches := r.FindStringSubmatch(out)
	if len(matches) == 0 {
		return ret, errors.Errorf("unrecognized version output: %s", stdout.String())
	}

	v := matches[1]
	ret, err = semver.Parse(v)
	if err != nil {
		return ret, errors.Wrap(err, "parsing semver")
	}

	err = restoreBackup(backupPath)
	if err != nil {
		return ret, errors.Wrap(err, "restoring backup")
	}

	return ret, nil
}

func parseFlag() error {
	flag.Parse()

	if *homeDirPath == "" {
		usr, err := user.Current()
		if err != nil {
			return errors.Wrap(err, "getting the current user")
		}

		// set home dir
		homeDirPath = &usr.HomeDir
	}

	return nil
}

func newCtx(version semver.Version) Ctx {
	return Ctx{
		version:      version,
		homeDirPath:  *homeDirPath,
		dnoteDirPath: fmt.Sprintf("%s/.dnote", *homeDirPath),
	}
}

func main() {
	if err := parseFlag(); err != nil {
		panic(errors.Wrap(err, "parsing flag"))
	}

	args := os.Args
	if len(args) > 1 {
		cmd := args[1]

		if cmd == "version" {
			fmt.Printf("dnote-doctor %s\n", versionTag)
			return
		} else if cmd == "help" {
			fmt.Println(helpText)
			return
		}

		fmt.Printf("unknwon command %s\n", cmd)
		return
	}

	version, err := checkVersion()
	if err != nil {
		panic(errors.Wrap(err, "checking version"))
	}

	debug("using version %d.%d.%d", version.Major, version.Minor, version.Patch)

	issues, err := scanIssues(version)
	if err != nil {
		panic(errors.Wrap(err, "scanning issues"))
	}

	debug("%d issues apply to this version", len(issues))

	ctx := newCtx(version)

	for _, i := range issues {
		fmt.Printf("diagnosing: %s...\n", i.name)

		ok, err := fixIssue(i, ctx)
		if err != nil {
			fmt.Println(errors.Wrapf(err, "⨯ Failed to diagnose %s", i.name))
			continue
		}

		if ok {
			fmt.Println("✔ fixed")
		} else {
			fmt.Println("✔ no issue found")
		}
	}

	fmt.Println("✔ done")
}
