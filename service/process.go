package service

import (
	"fmt"
	"strings"

	"github.com/shirou/gopsutil/process"
)

func FindProcessByCmdLineSuffix(cmdLineSuffix string) ([]*process.Process, error) {
	procs, err := process.Processes()

	if err != nil {
		return nil, fmt.Errorf("failed to get processes, error: %v", err)
	}

	found := make([]*process.Process, 0)
	for _, proc := range procs {
		cmdline, err := proc.Cmdline()
		if err != nil {
			return nil, fmt.Errorf("failed getting cmdLine, error: %v", err)
		}
		if strings.HasSuffix(cmdline, cmdLineSuffix) {
			found = append(found, proc)
		}
	}
	return found, nil
}
