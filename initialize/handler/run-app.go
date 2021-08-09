package handler

import (
	"fmt"
	"path"

	"github.com/hypha-dao/envctl/service"
)

const AppProcessId = "runBy=envctl"

type RunApp struct {
}

func NewRunApp(appProcessId string) *RunApp {
	return &RunApp{}
}

func (m *RunApp) Handle(data map[interface{}]interface{}, config map[interface{}]interface{}, initOp InitializeOp) error {
	if initOp == InitializeOp_Start {
		return m.start(data)
	} else if initOp == InitializeOp_Stop {
		return m.stop()
	}
	return nil
}

func (m *RunApp) start(data map[interface{}]interface{}) error {
	err := m.stop()
	if err != nil {
		return err
	}
	err = m.runApp(data)
	if err != nil {
		return err
	}
	return nil
}

func (m *RunApp) stop() error {
	err := m.killApps()
	if err != nil {
		return fmt.Errorf("failed to kill app processes, error: %v", err)
	}
	return nil
}

func (m *RunApp) runApp(data map[interface{}]interface{}) error {
	appPath := data["path"].(string)
	envSrcFile := data["env-file"].(string)
	port := data["port"].(int)
	fmt.Printf("Running app: %v, with env file: %v, on port: %v\n", appPath, envSrcFile, port)
	envDstFile := path.Join(appPath, ".env")
	err := service.CopyFile(envSrcFile, envDstFile)
	if err != nil {
		return fmt.Errorf("failed copying env file: %v, to: %v, error: %v", envSrcFile, envDstFile, err)
	}
	err = service.ExecCmd(&service.CmdArgs{
		Name: "yarn",
		Args: []string{"install"},
		Dir:  appPath,
	})
	if err != nil {
		return fmt.Errorf("failed installing app dependencies for app: %v, error %v", appPath, err)
	}

	err = service.ExecCmd(&service.CmdArgs{
		Name:       "quasar",
		Args:       []string{"dev", fmt.Sprintf("--port=%v", port), AppProcessId},
		Dir:        appPath,
		Background: true,
	})
	if err != nil {
		return fmt.Errorf("failed running app: %v on port: %v, error %v", appPath, port, err)
	}
	return nil
}

func (m *RunApp) killApps() error {
	procs, err := service.FindProcessByCmdLineSuffix(AppProcessId)
	if err != nil {
		return fmt.Errorf("failed to find app processes, error: %v", err)
	}
	for _, proc := range procs {
		err := proc.Kill()
		if err != nil {
			return fmt.Errorf("failed to kill app process: %v, error: %v", proc.Pid, err)
		}
	}
	return nil
}
