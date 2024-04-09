package builder

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	ERROR_BUILD_CONTAINER = "failed to execute build bulding container"
	ERROR_RUN_CONTAINER   = "failed to execute run building container"
)

const BASE_PATH = "./repo"

const COMMIT_ID_INDEX = 0
const UNIQUE_ID_INDEX = 3

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
		"ra-webs-builder",
		"sh",
	}

	args := append(params, arguments...)

	fmt.Printf("docker args: %s\n", args)

	return exec.Command("docker", args...).CombinedOutput()
}

func Build(name string, repo string) (string, string, error) {
	cmd := exec.Command("docker", "build", "-t", "ra-webs-builder", ".")
	output, err := cmd.Output()

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
