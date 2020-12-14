package sources

import (
	"os"
	"path/filepath"
	"io"
	"github.com/raphaelreyna/recon"
	"errors"
	"sync"
)

const DirSrc recon.SourceName = "dir_source"

type DirSource struct {
	Root string `json:"root" bson:"root" yaml:"root"`
	Linking LinkingType `json:"linking" bson:"linking" yaml;"linking"`

	cache map[string]string `json:"-" bson:"-" yaml:"-"`
	sync.Mutex `json:"-" bson:"-" yaml:"-"`
}

var DoneWalking error = errors.New("done walking")

func (ds *DirSource) findFile(name string) (bool, error) {
	if ds.cache == nil {
		ds.cache = map[string]string{}
	}

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
			ds.cache[name] = n
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
			ds.Unlock()
			return false
		}
		if found {
			srcFile = ds.cache[name]
		} else {
			ds.Unlock()
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
