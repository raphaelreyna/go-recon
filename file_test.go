package recon

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

type testSource map[string][]byte

var ts1 testSource = map[string][]byte{
	"a.txt": []byte("file a"),
	"b.txt": []byte("file b"),
	"c.txt": []byte("file c"),
}

var ts2 testSource = map[string][]byte{
	"1.txt": []byte("file 1"),
	"2.txt": []byte("file 2"),
	"3.txt": []byte("file 3"),
}

func (ts *testSource) AddFileAs(name, destination string, perm os.FileMode) bool {
	file, err := os.OpenFile(destination, os.O_CREATE|os.O_WRONLY, perm)
	if err != nil {
		return false
	}
	defer file.Close()

	data, exists := (*ts)[name]
	if !exists {
		return false
	}

	_, err = file.Write(data)

	return err == nil
}

func TestFile_AddTo(t *testing.T) {
	root, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(root); err != nil {
			t.Log(err)
		}
	}()

	f := &File{
		Name: "a.txt",
		SourceChain: SourceChain{
			&ts1, &ts2,
		},
	}

	if err := f.AddTo(root, 0644, nil); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(root, f.Name)); err != nil {
		t.Fatal(err)
	}
}

func TestFile_String(t *testing.T) {
	f1 := &File{Name: "f1"}
	f2 := &File{Name: "f1", Location: "loc"}

	if expected := f1.Name; f1.String() != expected {
		t.Fatalf("got wrong file string:\n\texpected: %s\n\treceived: %s\n", expected, f1.String())
	}

	if expected := f2.Name + "@" + f2.Location; f2.String() != expected {
		t.Fatalf("got wrong file string:\n\texpected: %s\n\treceived: %s\n", expected, f2.String())
	}
}
