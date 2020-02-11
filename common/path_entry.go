package common

import "fmt"

type PathEntry struct {
	Path             string `json:"path"`
	MD5              string `json:"md5"`
	LastModifiedTime int64  `json:"lastModifiedTime"`
}

func (p *PathEntry) String() string {
	return fmt.Sprintf("PathEntry{path: %s, MD5: %s, LastModifiedTime: %d}",
		p.Path, p.MD5, p.LastModifiedTime)
}
