package recon

import (
	"os"
	"path/filepath"
)

// File represents a file that can may be obtained from a Source in the SourceChain.
type File struct {
	// Name is the name of the file in the managed directory.
	Name string `json:"name" bson:"name" yaml:"name"`
	// Location is the location of the file, whatever that may mean to each Source.
	Location string `json:"location" bson:"location" yaml:"location"` // optional
	// SourceChain is a list of Sources which will be queried for the file.
	SourceChain SourceChain `json:"-"`
	// Perm is the permissions this file should have in the managed directory.
	Perm os.FileMode `json:"perm" bson:"perm" yaml:"perm"`
}

// AddTo searches through sc SourceChain for the file and adds it to the directory with the given permissions perm.
func (f *File) AddTo(dir string, perm int, sc SourceChain) (SourceErrs, error) {
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
