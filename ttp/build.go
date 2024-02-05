package ttp

import (
	"fmt"
	"os/exec"

	"github.com/akakou/ra_webs/ttp/ent"
)

var REPOSITORIES = "./data/repositories"

func compile(tainfo *ent.TAInfo) (string, []byte) {
	folderName := fmt.Sprintf("%v/%v", REPOSITORIES, tainfo.ID)
	exec.Command("mkdir", "-p", REPOSITORIES)

	exec.Command("git", "clone", tainfo.GitRepository, folderName)

	commitId := ""
	uniqueId := []byte{}

	return commitId, uniqueId
}
