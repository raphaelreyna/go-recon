package sources

import (
	"testing"
	"io/ioutil"
	"os"
	"path/filepath"
)

func TestHTTPSource_AddFileAs(t *testing.T) {
	root, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(root); err != nil {
			t.Log(err)
		}
	}()

	ds := &HTTPSource{}

	if found := ds.AddFileAs("https://golang.org/doc/gopher/project.png", filepath.Join(root, "photo.png"), 0644); !found {
		t.Fatal("could not find file in source dir")
	}

	// Make sure the file was placed into the root folder
	if _, err := os.Stat(filepath.Join(root, "photo.png")); err != nil {
		t.Fatal(err)
	}
}
