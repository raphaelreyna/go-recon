package recon

import (
	"errors"
	"os"
	"fmt"
)

type Error struct {
	msg string
	err error
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
}

func (e *Error) Unwrap() error {
	return e.err
}




type Source interface {
	// AddFileAs makes the file named name available at destination.
	// The argument destination is NOT a directory name, it is the full path to a file.
	// The returned boolean is true if the file with the name 'name' was found and could be written to destination
	AddFileAs(name, destination string, perm os.FileMode) bool
}

var ErrNoSource error = errors.New("no source found")

type SourceChain []Source

func (sc SourceChain) AddFileAs(file, destination string, perm os.FileMode) error {
	for _, s := range sc {
		if s.AddFileAs(file, destination, perm) {
			return nil
		}
	}

	return &Error{"error searching for" + file, ErrNoSource}
}
