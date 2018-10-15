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
	backupName := fmt.Sprintf(".dnote-backup-%d", time.Now().Unix())
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

func checkVersion() (semver, error) {
	var ret semver

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
	ret, err = parseSemver(v)
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

func main() {
	os.Setenv("DNOTE_DOCTOR_DEBUG", "1")

	if err := parseFlag(); err != nil {
		panic(errors.Wrap(err, "parsing flag"))
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

}
