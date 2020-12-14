package recon

import (
	"os"
	"path/filepath"
)

type File struct {
	Name        string      `json:"name" bson:"name" yaml:"name"`
	Location    string      `json:"location" bson:"location" yaml:"location"` // optional
	SourceChain SourceChain `json:"-"`
	Perm        os.FileMode `json:"perm" bson:"perm" yaml:"perm"`
}

func (f *File) AddTo(dir string, perm int, sc SourceChain) error {
	ssc := sc
	if f.SourceChain != nil {
		ssc = f.SourceChain
	}
	// Figure out the FileMode.
	// The one passed in by the user should get priority over the one given by File.
	// If neither the caller nor file provide a FileMode, default to 644.
	p := os.FileMode(perm)
	if p == 0 {
		p = f.Perm
	}
	if p == 0 {
		p = 644
	}

	loc := f.Location
	if loc == "" {
		loc = f.Name
	}

	return ssc.AddFileAs(loc, filepath.Join(dir, f.Name), p)
}

func (f *File) String() string {
	if f.Name == f.Location || f.Location == "" {
		return f.Name
	}

	return f.Name + "@" + f.Location
}
