package sources

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
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

	if err := ds.AddFileAs("https://golang.org/doc/gopher/project.png", filepath.Join(root, "photo.png"), 0644); err != nil {
		t.Fatal(err)
	}

	// Make sure the file was placed into the root folder
	if _, err := os.Stat(filepath.Join(root, "photo.png")); err != nil {
		t.Fatal(err)
	}
}
