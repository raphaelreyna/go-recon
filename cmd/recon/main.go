package main

import (
	"github.com/raphaelreyna/recon"
	"github.com/raphaelreyna/recon/sources"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

type reconFile struct {
	Files      []*recon.File `json:"files" bson:"files" yaml:"files"`
	SourceDirs []string      `json:"dirs" bson:"dirs" yaml:"dirs"`
}

func main() {
	here, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	rf := reconFile{}
	rfData, err := ioutil.ReadFile(filepath.Join(here, ".recon.yaml"))
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(rfData, &rf); err != nil {
		panic(err)
	}

	d := &recon.Dir{
		Root:        here,
		Files:       rf.Files,
		SourceChain: sources.NewDirSourceChain(sources.SoftLink, rf.SourceDirs...),
		FilesPerm:   0644,
	}

	d.SourceChain = append(d.SourceChain,
		&sources.HTTPSource{},
		&sources.ShellSource{
			WorkingDir: here,
			Shell:      "/bin/bash",
		},
	)

	if err := d.Reconcile(); err != nil {
		panic(err)
	}
}
