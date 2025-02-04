package lib

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"fmt"
	"io"
)

// Command : type string
type Command struct {
	name string
}

// NewCommand : get command
func NewCommand(name string) *Command {
	return &Command{name: name}
}

// PathList : get bin path list
func (cmd *Command) PathList() []string {
	path := os.Getenv("PATH")
	return strings.Split(path, string(os.PathListSeparator))
}

func isDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil || os.IsNotExist(err) {
		return false
	}
	return fileInfo.IsDir()
}

func isExecutable(path string) bool {
	if isDir(path) {
		return false
	}

	fileInfo, err := os.Stat(path)
	if err != nil || os.IsNotExist(err) {
		return false
	}

	if runtime.GOOS == "windows" {
		return true
	}

	if fileInfo.Mode()&0111 != 0 {
		return true
	}

	return false
}

// Find : find all bin path
func (cmd *Command) Find() func() string {
	pathChan := make(chan string)
	go func() {
		for _, p := range cmd.PathList() {
			if !isDir(p) {
				continue
			}
			fileList, err := ioutil.ReadDir(p)
			if err != nil {
				continue
			}

			for _, f := range fileList {
				path := filepath.Join(p, f.Name())
				if isExecutable(path) && f.Name() == cmd.name {
					pathChan <- path
				}
			}
		}
		pathChan <- ""
	}()

	return func() string {
		return <-pathChan
	}
}

func copyFile(sourceFile string, destinationFile string) bool {
	// Open the source file
	src, err := os.Open(sourceFile)
	if err != nil {
		fmt.Println("Error opening source file:", err)
		return false
	}
	defer src.Close()

	// Create the destination file
	dst, err := os.Create(destinationFile)
	if err != nil {
		fmt.Println("Error creating destination file:", err)
		return false
	}
	defer dst.Close()

	// Copy the contents
	_, err = io.Copy(dst, src)
	if err != nil {
		fmt.Println("Error copying file:", err)
		return false
	}

	fmt.Println("File copied successfully: ", sourceFile, "-", destinationFile)
	return true
}