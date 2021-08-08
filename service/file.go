package service

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func CopyFile(src, dst string) error {
	bytesRead, err := ioutil.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed reading file: %v, error: %v", src, err)
	}

	err = ioutil.WriteFile(dst, bytesRead, 0755)
	if err != nil {
		return fmt.Errorf("failed writing file: %v, error: %v", dst, err)
	}
	return nil
}

func IsDirEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}
