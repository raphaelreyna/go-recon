package recon

import (
	"os"
	"sync"
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

// FileErrs maps a file to the errors encountered while trying to add it to a directory.
type FileErrs map[*File]SourceErrs

// Reconcile adds the missing files to this Dir structs Root directory.
// The returned FileErrs maps Files that could not be added to the directory to the errors raised while sourcing the file data.
func (d *Dir) Reconcile() FileErrs {
	mf, err := d.MissingFiles()
	if err != nil {
		return nil
	}

	fe := FileErrs{}
	for _, file := range mf {
		if se, err := file.AddTo(d.Root, int(d.FilesPerm), d.SourceChain); err != nil {
			fe[file] = se
		}
	}

	return fe
}

// ReconcileC adds the missing files to this Dir structs Root directory concurrently.
// The returned FileErrs maps Files that could not be added to the directory to the errors raised while sourcing the file data.
func (d *Dir) ReconcileC(workerCount uint) FileErrs {
	if workerCount == 0 {
		workerCount = 1
	}
	mf, err := d.MissingFiles()
	if err != nil {
		return nil
	}

	mu := sync.Mutex{}
	wg := sync.WaitGroup{}

	fe := FileErrs{}
	workChan := make(chan *File, workerCount)

	// Launch the workers
	for i := uint(0); i < workerCount; i++ {
		go func() {
			for file := range workChan {
				se, err := file.AddTo(d.Root, int(d.FilesPerm), d.SourceChain)
				if err != nil {
					mu.Lock()
					fe[file] = se
					mu.Unlock()
				}
				wg.Done()
			}
		}()
	}

	// Send files into chan for workers to pick up
	for _, file := range mf {
		wg.Add(1)
		workChan <- file
	}
	close(workChan)

	wg.Wait()

	return fe
}
