package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/raphaelreyna/go-recon"
	"github.com/raphaelreyna/go-recon/sources"
	yaml "gopkg.in/yaml.v2"
)

type reconFile struct {
	Files      []*recon.File `json:"files" bson:"files" yaml:"files"`
	SourceDirs []string      `json:"source_dirs" bson:"source_dirs" yaml:"source_dirs"`
}

func main() {
	// Get the directory we're in
	here, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// Look for a recon file and read its contents if found.
	rf := reconFile{}
	rfData, err := ioutil.ReadFile(filepath.Join(here, "recon.yaml"))
	if err != nil {
		panic(err)
	}
	// Parse the recon file
	if err := yaml.Unmarshal(rfData, &rf); err != nil {
		panic(err)
	}

	// Create a new recon Dir struct which will add the files we requested into this directory
	d := &recon.Dir{
		Root:  here,
		Files: rf.Files,
		// Create a source chain from the directories listed in the recon file.
		// The first parameter, sources.SoftLink, signals that we want to soft-link our files into this directory when possible.
		SourceChain: sources.NewDirSourceChain(sources.SoftLink, rf.SourceDirs...),
		FilesPerm:   0644, // We want the files to have their permissions set to rw-r--r-- by default.
	}

	// Add the HTTPSource and ShellSource sources to the SourceChain.
	// If a file cant be found in one of the directories given in the recon file, try the internet and finally a shell script.
	d.SourceChain = append(d.SourceChain,
		&sources.HTTPSource{},
		&sources.ShellSource{
			WorkingDir: here,
			Shell:      "/bin/bash",
		},
	)

	// Reconcile this directories contents with the files listed in the recon file.
	if errs := d.Reconcile(); len(errs) != 0 {
		panic(errs)
	}
}
