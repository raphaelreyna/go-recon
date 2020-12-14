package sources

import (
	"os/exec"
	"os"
	"github.com/raphaelreyna/recon"
)

const ShellSrc recon.SourceName = "shell_source"

type ShellSource struct {
	WorkingDir string
	Shell string
}

func (ss *ShellSource) AddFileAs(name, destination string, perm os.FileMode) bool {
	cmd := exec.Command(ss.Shell, "-c", name)
	file, err := os.OpenFile(destination, os.O_CREATE | os.O_RDWR, perm)
	if err != nil {
		return false
	}

	rollback := true
	defer func() {
		file.Close()
		if rollback {
			os.Remove(file.Name())
		}
	}()

	cmd.Stdout = file
	cmd.Stderr = file
	cmd.Dir = ss.WorkingDir

	if err := cmd.Run(); err != nil {
		return false
	}

	rollback = false
	return true
}
