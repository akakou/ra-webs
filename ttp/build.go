package ttp

import (
	"fmt"
	"os/exec"

	"github.com/akakou/ra_webs/ttp/ent"
)

var REPOSITORIES = "./data/repositories"

func compile(tainfo *ent.TAInfo) (string, uint16) {
	folderName := fmt.Sprintf("%v/%v", REPOSITORIES, tainfo.ID)
	exec.Command("mkdir", "-p", REPOSITORIES)

	exec.Command("git", "clone", tainfo.GitRepository, folderName)

	commitId := ""
	uniqueId := uint16(0)

	return commitId, uniqueId
}
