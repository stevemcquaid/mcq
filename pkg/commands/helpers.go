package commands

import (
	"fmt"
	"io/ioutil"
	"os/user"
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
	if len(modulePaths) == 0 {
		return "", "", fmt.Errorf("module not found")
	} else if len(modulePaths) == 1 {
		userName, err := GetUserName()
		if err != nil {
			return "", "", err
		}
		gitOrg = userName
		gitRepo = modulePaths[0]
		return gitOrg, gitRepo, nil
	}

	gitOrg = modulePaths[len(modulePaths)-2]
	gitRepo = modulePaths[len(modulePaths)-1]

	return gitOrg, gitRepo, nil
}

func GetUserName() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("unable to get current user")
	}
	return user.Name, nil
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
