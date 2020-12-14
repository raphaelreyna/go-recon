package sources

import (
	"os"
	"path/filepath"
	"io"
	"github.com/raphaelreyna/recon"
	"errors"
	"sync"
)

type DirSource struct {
	Root string
	Linking LinkingType

	cache map[string]string
	sync.Mutex
}

var DoneWalking error = errors.New("done walking")

func (ds *DirSource) findFile(name string) (bool, error) {
	ds.Lock()
	if ds.cache == nil {
		ds.cache = map[string]string{}
	}
	ds.Unlock()

	var found bool
	err := filepath.Walk(ds.Root, filepath.WalkFunc(func(n string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		if filepath.Base(name) == filepath.Base(n) {
			found = true
			ds.Lock()
			ds.cache[name] = n
			ds.Unlock()
			return DoneWalking
		}

		return nil
	}))
	if err != nil && err != DoneWalking {
		return false, err
	}

	return found, nil
}

func (ds *DirSource) AddFileAs(name, destination string, perm os.FileMode) bool {
	ds.Lock()
	srcFile, exists := ds.cache[name]
	if !exists {
		found, err := ds.findFile(name)
		if err != nil {
			return false
		}
		if found {
			srcFile = ds.cache[name]
		} else {
			return false
		}
	}
	ds.Unlock()

	var linkFunc func(string, string) error
	switch ds.Linking {
	case NoLink:
		nf, err := os.OpenFile(destination, os.O_CREATE | os.O_WRONLY, perm)
		if err != nil {
			return false
		}
		defer nf.Close()


		sf, err := os.Open(srcFile)
		if err != nil {
			return false
		}
		defer sf.Close()

		_, err = io.Copy(nf, sf)
		return false
	case HardLink:
		linkFunc = os.Link
	case SoftLink:
		linkFunc = os.Symlink
	}

	return linkFunc(srcFile, destination) == nil
}

func NewRecDirSourceChain(linking LinkingType, dirs ...string) recon.SourceChain {
	sc := recon.SourceChain{}
	for _, dir := range dirs {
		sc = append(sc, &DirSource{
			Root: dir,
			Linking: linking,
		})
	}
	return sc
}
