package recon

import (
	"errors"
	"fmt"
	"os"
)

// Error is used for wrapping errors.
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

// SourceName is the name of a Source
type SourceName string

// Source represents anything that can take a string and use that to create a file at a given path.
type Source interface {
	// AddFileAs makes the file named name available at destination.
	// The argument destination is NOT a directory name, it is the full path to a file.
	// The returned boolean is true if the file with the name 'name' was found and could be written to destination
	AddFileAs(name, destination string, perm os.FileMode) bool
}

// ErrNoSource is raised when no source in a source chain contains a requested file.
var ErrNoSource error = errors.New("no source found")

// SourceChain is a list of Sources.
// A SourceChain searches for files starting with the first Source and moving down the list if the file is not found.
type SourceChain []Source

// AddFileAs searches the SourceChain for the first Source that can create a file with the name name at the path destination.
func (sc SourceChain) AddFileAs(name, destination string, perm os.FileMode) error {
	for _, s := range sc {
		if s.AddFileAs(name, destination, perm) {
			return nil
		}
	}

	return &Error{"error searching for" + name, ErrNoSource}
}
