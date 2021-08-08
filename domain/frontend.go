package domain

import (
	"fmt"
	"path"

	"github.com/hypha-dao/envctl/service"
)

const AppProcessId = "runBy=envctl"

type Frontend struct {
	InitSettings map[string]interface{}
}

func NewFrontend(initSettings map[string]interface{}) *Frontend {
	return &Frontend{
		InitSettings: initSettings,
	}
}

func (m *Frontend) Start() error {
	err := m.Stop()
	if err != nil {
		return err
	}
	err = m.runApps(m.InitSettings["run-apps"].(map[string]interface{}))
	if err != nil {
		return err
	}
	return nil
}

func (m *Frontend) Stop() error {
	err := m.killApps()
	if err != nil {
		return fmt.Errorf("failed to kill app processes, error: %v", err)
	}
	return nil
}

func (m *Frontend) killApps() error {
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

func (m *Frontend) runApps(run map[string]interface{}) error {
	apps := run["app"].([]interface{})
	for _, appI := range apps {
		app := appI.(map[interface{}]interface{})
		appPath := app["path"].(string)
		envSrcFile := app["env-file"].(string)
		port := app["port"].(int)
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
	}
	return nil
}
