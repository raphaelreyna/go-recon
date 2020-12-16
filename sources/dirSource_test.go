package sources

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestDirSource_AddFileAs(t *testing.T) {
	root, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	srcDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(root); err != nil {
			t.Log(err)
		}
		if err := os.RemoveAll(srcDir); err != nil {
			t.Log(err)
		}
	}()

	// Add the source file into the source dir
	srcFile, err := os.Create(filepath.Join(srcDir, "test.txt"))
	if err != nil {
		t.Fatal(err)
	}
	srcFile.Write([]byte("pass"))
	srcFile.Close()

	ds := &DirSource{
		Root:    srcDir,
		Linking: NoLink,
	}

	if err := ds.AddFileAs("test.txt", filepath.Join(root, "pass.txt"), 0644); err != nil {
		t.Fatal(err)
	}

	// Make sure the file was placed into the root folder
	if _, err := os.Stat(filepath.Join(root, "pass.txt")); err != nil {
		t.Fatal(err)
	}
}
