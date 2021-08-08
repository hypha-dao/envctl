package service_test

import (
	"fmt"
	"testing"

	"github.com/hypha-dao/envctl/service"
	"gotest.tools/assert"
)

func TestFindProcess(t *testing.T) {
	procs, err := service.FindProcessByCmdLineSuffix("--port=8085 runBy=envctl")
	assert.NilError(t, err)
	for _, proc := range procs {
		err = proc.Kill()
		fmt.Println("Killing process: ", proc.Pid)
		assert.NilError(t, err)
	}
}
