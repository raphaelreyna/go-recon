package recon

import (
	"os"
)

// Dir manages the directory at Root by reconciling its contents with Files.
// Missing files will be searched for in the sources listed in the SourceChain
type Dir struct {
	Root        string
	Files       []*File
	SourceChain SourceChain
	FilesPerm   os.FileMode
}

// NewDir is a helper function for creating a partially configured Dir struct.
func NewDir(root string, sc SourceChain, files ...string) *Dir {
	ff := []*File{}
	for _, f := range files {
		ff = append(ff, &File{
			Name: f,
		})
	}

	return &Dir{
		Root:        root,
		Files:       ff,
		SourceChain: sc,
		FilesPerm:   0644,
	}
}

// MissingFiles returns the Files that are missing from this Dir structs Root directory.
func (d *Dir) MissingFiles() ([]*File, error) {
	// Grab the names of the files currently in the directory
	root, err := os.Open(d.Root)
	if err != nil {
		return nil, err
	}

	currentFileNames, err := root.Readdirnames(0)
	root.Close()
	if err != nil {
		return nil, err
	}

	// for each file we want, loop through the current file names we have and look for matches
	missing := []*File{}
	for _, rf := range d.Files {
		var found bool

		for _, cfn := range currentFileNames {
			if rf.Name == cfn {
				found = true
				break
			}
		}

		if !found {
			missing = append(missing, rf)
		}
	}

	return missing, nil
}

// Reconcile adds the missing files to this Dir structs Root directory.
func (d *Dir) Reconcile() error {
	mf, err := d.MissingFiles()
	if err != nil {
		return err
	}

	for _, file := range mf {
		if err := file.AddTo(d.Root, int(d.FilesPerm), d.SourceChain); err != nil {
			return err
		}
	}

	return nil
}
