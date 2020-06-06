package worker

import "os"

type Path struct {
	path string
	info os.FileInfo
}

func NewPath(path string, info os.FileInfo) Path {
	return Path{
		path: path,
		info: info,
	}
}

func (p *Path) Info() os.FileInfo {
	return p.info
}

func (p *Path) Path() string {
	return p.path
}
