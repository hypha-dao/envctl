package service

import (
	"fmt"
	"os"
	"os/exec"
)

type CmdArgs struct {
	Name       string
	Args       []string
	Dir        string
	Background bool
}

func (m *CmdArgs) String() string {
	return fmt.Sprintf(
		`
			CmdArgs {
				Name: %v,
				Args: %v,
				Dir: %v,
				Background: %v,
			}	
		`,
		m.Name,
		m.Args,
		m.Dir,
		m.Background,
	)
}

func ExecCmd(args *CmdArgs) error {
	fmt.Println("Executing cmd: ", args)
	cmd := exec.Command(args.Name, args.Args...)
	cmd.Dir = args.Dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if args.Background {
		err := cmd.Start()
		if err != nil {
			return fmt.Errorf("error starting command: %v, error: %v", args, err)
		}
	} else {
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("error executing command: %v, error: %v", args, err)
		}
	}
	return nil
}
