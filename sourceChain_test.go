package recon

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestSourceChain_AddFileAs(t *testing.T) {
	root, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(root); err != nil {
			t.Log(err)
		}
	}()

	sc := SourceChain{&ts1, &ts2}
	dst := filepath.Join(root, "a.txt")

	if err := sc.AddFileAs("a.txt", dst, 0644); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(dst); err != nil {
		t.Fatal(err)
	}
}
