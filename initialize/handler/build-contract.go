package handler

import (
	"fmt"
	"os"
	"path"

	"github.com/hypha-dao/envctl/service"
)

type BuildContract struct {
}

func NewBuildContract() *BuildContract {
	return &BuildContract{}
}

func (m *BuildContract) Handle(data map[interface{}]interface{}, config map[interface{}]interface{}, initOp InitializeOp) error {
	basePath := config["base-path"].(string)
	name := data["name"].(string)
	repoPath := path.Join(basePath, name)
	fmt.Printf("Building repo: %v\n", repoPath)
	buildPath := path.Join(repoPath, "build")
	err := os.Mkdir(buildPath, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}
	isDirEmpty, err := service.IsDirEmpty(buildPath)
	if err != nil {
		return fmt.Errorf("failed to check if build path: %v is empty, error: %v", buildPath, err)
	}
	if initOp == InitializeOp_Restart || isDirEmpty {
		err = service.ExecCmd(&service.CmdArgs{
			Name: "cmake",
			Args: []string{".."},
			Dir:  buildPath,
		})
		if err != nil {
			return fmt.Errorf("failed running cmake for repo: %v, build path: %v, error %v", repoPath, buildPath, err)
		}

		err = service.ExecCmd(&service.CmdArgs{
			Name: "make",
			Dir:  buildPath,
		})
		if err != nil {
			return fmt.Errorf("failed running make for repo: %v, build path: %v, error %v", repoPath, buildPath, err)
		}
	}
	return nil
}
