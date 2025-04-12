package builder

import (
	"bytes"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"fmt"
	"os/exec"
	"strings"
	"time"

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
const BASE_PROGRAM_PATH = "./ta/example"

const BUILD_SCRIPT = BASE_REPO_PATH + "build.sh"

const COMMIT_ID_INDEX = 0
const UNIQUE_ID_INDEX = 1

const EXECUTABLE = "example"

func build(repo, commitId string) ([]byte, error) {
	sha256 := sha256.Sum256([]byte(repo))
	folderName := fmt.Sprintf("%v-%x", time.Now().Unix(), sha256)

	uniqueIdString, err := buildCode(folderName, repo, commitId)
	if err != nil {
		return nil, err
	}

	uniqueId, _ := hex.DecodeString(uniqueIdString)

	return uniqueId, nil
}

func buildCode(name, repo, commitId string) (string, error) {
	extractembed.Extract(BASE_REPO_PATH, &embedFiles)

	var outBuf, errBuf bytes.Buffer
	cmd := exec.Command("bash", BUILD_SCRIPT, name, repo, commitId, BASE_REPO_PATH, BASE_PROGRAM_PATH, EXECUTABLE)

	fmt.Printf("Running command: %v\n", cmd.String())

	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()
	fmt.Print(errBuf.String())

	if err != nil {
		return "", fmt.Errorf("%v: %v", ERROR_RUNNING_CMD, err)
	}

	lines := strings.Split(outBuf.String(), "\n")
	if len(lines) != 3 {
		return "", fmt.Errorf("%v: expected 2 lines, but got %v", ERROR_OUTPUT_SIZE_WRONG, len(lines))
	}

	actualCommitId := lines[COMMIT_ID_INDEX]
	uniqueId := lines[UNIQUE_ID_INDEX]

	if actualCommitId != commitId {
		return "", fmt.Errorf("expected commit id %v, but got %v", commitId, actualCommitId)
	}

	return uniqueId, nil
}
