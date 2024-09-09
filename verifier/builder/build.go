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
const BASE_PROGRAM_PATH = "./ta/example/"

const BUILD_SCRIPT = BASE_REPO_PATH + "build.sh"

const COMMIT_ID_INDEX = 0
const UNIQUE_ID_INDEX = 1

const BRANCH = "feature/direct-building"
const EXECUTABLE = "example"

func build(name, repo string) (string, string, error) {
	extractembed.Extract(BASE_REPO_PATH, &embedFiles)

	var outBuf, errBuf bytes.Buffer
	cmd := exec.Command("bash", BUILD_SCRIPT, name, repo, BASE_REPO_PATH, BRANCH, BASE_PROGRAM_PATH, EXECUTABLE)

	fmt.Printf("Running command: %v\n", cmd.String())

	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()
	fmt.Print(errBuf.String())

	if err != nil {
		return "", "", fmt.Errorf("%v: %v", ERROR_RUNNING_CMD, err)
	}

	lines := strings.Split(outBuf.String(), "\n")
	if len(lines) != 3 {
		return "", "", fmt.Errorf("%v: expected 2 lines, but got %v", ERROR_OUTPUT_SIZE_WRONG, len(lines))
	}

	return lines[COMMIT_ID_INDEX], lines[UNIQUE_ID_INDEX], nil
}
