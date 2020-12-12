package commands

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	modfile "golang.org/x/mod/modfile"
)

func ReadModFile() (string, error) {
	goModBytes, err := ioutil.ReadFile("go.mod")
	if err != nil {
		return "", err
	}

	modName := modfile.ModulePath(goModBytes)
	return modName, nil
}

func GetModules() (gitOrg string, gitRepo string, err error) {
	mod, err := ReadModFile()
	if err != nil {
		return "", "", err
	}

	modulePaths := strings.Split(mod, "/")
	if len(modulePaths) <= 1 {
		return "", "", errors.New("unable to split module parts")
	}

	gitOrg = modulePaths[len(modulePaths)-2]
	gitRepo = modulePaths[len(modulePaths)-1]

	return gitOrg, gitRepo, nil
}

func GetDockerImage() (string, error) {
	gitOrg, gitRepo, err := GetModules()
	if err != nil {
		return "", err
	}
	dockerBase := path.Join(gitOrg, gitRepo)
	dockerImage := fmt.Sprintf("%s:%s", dockerBase, "latest")

	return dockerImage, nil
}
