package builder

import (
	"bytes"
	"embed"
	"fmt"
	"os/exec"
	"strings"

	extractembed "github.com/akakou/extract-embed"
)

const (
	ERROR_RUNNING_CMD       = "Error running command"
	ERROR_OUTPUT_SIZE_WRONG = "Error outpu size is wrong"
)

//go:embed build.sh
var embedFiles embed.FS

var Build = build

const BASE_REPO_PATH = "./repo/"
const BASE_PROGRAM_PATH = "."

const BUILD_SCRIPT = BASE_REPO_PATH + "build.sh"

const COMMIT_ID_INDEX = 0
const UNIQUE_ID_INDEX = 1

const BRANCH = "main"

func build(name, repo string) (string, string, error) {
	extractembed.Extract(BASE_REPO_PATH, &embedFiles)

	var outBuf, errBuf bytes.Buffer
	cmd := exec.Command("sh", BUILD_SCRIPT, repo, BASE_REPO_PATH+"/"+name, BRANCH, BASE_PROGRAM_PATH)

	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()
	fmt.Print(errBuf.String())

	if err != nil {
		return "", "", fmt.Errorf("%v: %v", ERROR_RUNNING_CMD, err)
	}

	lines := strings.Split(outBuf.String(), "\n")
	if len(lines) != 2 {
		return "", "", fmt.Errorf("%v: expected 2 lines, but got %v", ERROR_OUTPUT_SIZE_WRONG, len(lines))
	}

	return lines[COMMIT_ID_INDEX], lines[UNIQUE_ID_INDEX], nil
}
