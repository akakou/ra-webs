package builder

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	extractembed "github.com/akakou/extract-embed"
)

const (
	ERROR_BUILD_CONTAINER = "failed to execute build bulding container"
	ERROR_RUN_CONTAINER   = "failed to execute run building container"
)

const BASE_PATH = "./repo"
const IMAGE_NAME = "ra-webs-builder"

const COMMIT_ID_INDEX = 0
const UNIQUE_ID_INDEX = 3

//go:embed *.sh Dockerfile
var embedFiles embed.FS

func docker_run(name string, arguments ...string) ([]byte, error) {
	directory := fmt.Sprintf("%s/%s", BASE_PATH, name)
	os.Mkdir(directory, 0600)

	params := []string{
		"run",
		"-v",
		fmt.Sprintf("%s:/repo/%s:z", directory, name),
		"-w",
		fmt.Sprintf("/repo/%s", name),
		"--rm",
		IMAGE_NAME,
		"sh",
	}

	args := append(params, arguments...)

	fmt.Printf("docker args: %s\n", args)

	return exec.Command("docker", args...).CombinedOutput()
}

var Build = build

func build(name, repo string) (string, string, error) {
	extractembed.Extract(BASE_PATH, &embedFiles)

	current, err := os.Getwd()
	if err != nil {
		return "", "", err
	}

	err = os.Chdir(BASE_PATH)
	if err != nil {
		return "", "", err
	}

	commitId, uniqueId, err := execBuild(name, repo)

	if err != nil {
		return "", "", err
	}

	err = os.Chdir(current)
	if err != nil {
		return "", "", err
	}

	return commitId, uniqueId, nil
}

func execBuild(name, repo string) (string, string, error) {
	cmd := exec.Command("docker", "build", "-t", IMAGE_NAME, ".")
	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Fprintf(os.Stderr, "stdout: %s\nerr: %s", output, err)
		return "", "", errors.New(ERROR_BUILD_CONTAINER)
	}

	output, err = docker_run(name, "/build/build.sh", repo)

	fmt.Fprintf(os.Stderr, "stdout: %s\nerr: %s", output, err)
	if err != nil {
		return "", "", errors.New(ERROR_RUN_CONTAINER + "1")
	}

	output, err = docker_run(name, "/build/show.sh", name)

	fmt.Fprintf(os.Stderr, "cmd: %s\nstdout: %s\nerr: %s", cmd.Args, output, err)
	if err != nil {
		return "", "", errors.New(ERROR_RUN_CONTAINER + "2")
	}

	lines := strings.Split(string(output), "\n")
	return lines[COMMIT_ID_INDEX], lines[UNIQUE_ID_INDEX], nil
}
