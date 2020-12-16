package sources

import "github.com/raphaelreyna/go-recon"

var DefaultSources []recon.SourceName = []recon.SourceName{
	DirSrc,
	FlatDirSrc,
	HTTPSrc,
}

func IsDefaultName(name recon.SourceName) bool {
	for _, sn := range DefaultSources {
		if sn == name {
			return true
		}
	}

	return false
}
