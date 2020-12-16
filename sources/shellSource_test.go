package sources

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestShellDirSource_AddFileAs(t *testing.T) {
	root, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(root); err != nil {
			t.Log(err)
		}
	}()

	ds := &ShellSource{
		WorkingDir: root,
		Shell:      "/bin/bash",
	}

	if err := ds.AddFileAs("echo pass", filepath.Join(root, "pass.txt"), 0644); err != nil {
		t.Fatal(err)
	}

	// Make sure the file was placed into the root folder
	data, err := ioutil.ReadFile(filepath.Join(root, "pass.txt"))
	if err != nil {
		t.Fatal(err)
	}

	if ds := string(data); ds != "pass\n" {
		t.Fatalf("wrong file contents, expected: pass\treceived: %s", ds)
	}
}
