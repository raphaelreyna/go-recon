package recon

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestDir_MissingFiles(t *testing.T) {
	root, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(root); err != nil {
			t.Log(err)
		}
	}()

	file, err := os.Create(filepath.Join(root, "a.txt"))
	if err != nil {
		t.Fatal(err)
	}
	file.Close()

	file, err = os.Create(filepath.Join(root, "b.txt"))
	if err != nil {
		t.Fatal(err)
	}
	file.Close()

	d := &Dir{Root: root}
	d.Files = []*File{
		{Name: "a.txt"},
		{Name: "b.txt"},
		{Name: "c.txt"},
		{Name: "d.txt"},
	}

	mf, err := d.MissingFiles()
	if err != nil {
		t.Fatal(err)
	}

	if len(mf) != 2 {
		t.Fatalf("expected 2 missing files, found: %d", len(mf))
	}

	if mf[0].Name != "c.txt" {
		t.Fatalf("expected first missing file to be c.txt, instead its: %s", mf[0].Name)
	}
	if mf[1].Name != "d.txt" {
		t.Fatalf("expected second missing file to be d.txt, instead its: %s", mf[1].Name)
	}
}

func TestDir_Reconcile(t *testing.T) {
	root, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(root); err != nil {
			t.Log(err)
		}
	}()

	file, err := os.Create(filepath.Join(root, "a.txt"))
	if err != nil {
		t.Fatal(err)
	}
	file.Close()

	file, err = os.Create(filepath.Join(root, "b.txt"))
	if err != nil {
		t.Fatal(err)
	}
	file.Close()

	d := &Dir{Root: root}
	d.Files = []*File{
		{Name: "a.txt"},
		{Name: "b.txt"},
		{Name: "c.txt"},
		{Name: "1.txt"},
	}
	d.SourceChain = SourceChain{
		&ts1, &ts2,
	}
	d.FilesPerm = 0644

	errs := d.Reconcile()
	if len(errs) != 0 {
		t.Fatal(errs)
	}

	data, err := ioutil.ReadFile(filepath.Join(root, "c.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(data, ts1["c.txt"]) != 0 {
		t.Fatalf("wrong file contents:\n\texpected: %s\n\treceived: %s",
			string(ts1["c.txt"]), string(data),
		)
	}

	data, err = ioutil.ReadFile(filepath.Join(root, "1.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(data, ts2["1.txt"]) != 0 {
		t.Fatalf("wrong file contents:\n\texpected: %s\n\treceived: %s",
			string(ts1["1.txt"]), string(data),
		)
	}
}

func TestDir_ReconcileC(t *testing.T) {
	root, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(root); err != nil {
			t.Log(err)
		}
	}()

	file, err := os.Create(filepath.Join(root, "a.txt"))
	if err != nil {
		t.Fatal(err)
	}
	file.Close()

	file, err = os.Create(filepath.Join(root, "b.txt"))
	if err != nil {
		t.Fatal(err)
	}
	file.Close()

	d := &Dir{Root: root}
	d.Files = []*File{
		{Name: "a.txt"},
		{Name: "b.txt"},
		{Name: "c.txt"},
		{Name: "1.txt"},
	}
	d.SourceChain = SourceChain{
		&ts1, &ts2,
	}
	d.FilesPerm = 0644

	errs := d.ReconcileC(2)
	if len(errs) != 0 {
		t.Fatal(errs)
	}

	data, err := ioutil.ReadFile(filepath.Join(root, "c.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(data, ts1["c.txt"]) != 0 {
		t.Fatalf("wrong file contents:\n\texpected: %s\n\treceived: %s",
			string(ts1["c.txt"]), string(data),
		)
	}

	data, err = ioutil.ReadFile(filepath.Join(root, "1.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(data, ts2["1.txt"]) != 0 {
		t.Fatalf("wrong file contents:\n\texpected: %s\n\treceived: %s",
			string(ts1["1.txt"]), string(data),
		)
	}
}
