# Go-Recon

A Go package for robustly reconciling directory contents with a list of files.
File data is pulled from a chain of plugin-defined sources so you can grab your data wherever it may be!


## Example
When run, the following snippet of Go code will search for files named `a.txt` and `b.txt` in directorires `/foo` and `/bar/baz` and soft-link them into `/home/recon`.

```
package main

import (
	"github.com/raphaelreyna/recon/sources"
	"github.com/raphaelreyna/recon"
)

func main() {
	files := []*recon.File{
		&recon.File{Name: "a.txt"},
		&recon.File{Name: "b.txt"},
	}

	d := &recon.Dir{
		Root: "/home/recon",
		Files: files,
		SourceChain: sources.NewDirSourceChain(sources.SoftLink, "/foo", "/bar/baz"),
		FilesPerm: 0644,
	}

	if err := d.Reconcile(); err != nil {
		panic(err)
	}
}
```
