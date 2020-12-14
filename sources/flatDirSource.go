package sources

import (
	"os"
	"path/filepath"
	"io"
	"github.com/raphaelreyna/recon"
)

type LinkingType uint

const (
	NoLink LinkingType = iota
	HardLink
	SoftLink
)

type FlatDirSource struct {
	Root string
	Linking LinkingType
}

func (ds *FlatDirSource) AddFileAs(name, destination string, perm os.FileMode) bool {
	srcFile := filepath.Join(ds.Root, filepath.Base(name))

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
		return err == nil
	case HardLink:
		linkFunc = os.Link
	case SoftLink:
		linkFunc = os.Symlink
	}

	return linkFunc(srcFile, destination) == nil
}

func NewDirSourceChain(linking LinkingType, dirs ...string) recon.SourceChain {
	sc := recon.SourceChain{}
	for _, dir := range dirs {
		sc = append(sc, &FlatDirSource{
			Root: dir,
			Linking: linking,
		})
	}
	return sc
}
