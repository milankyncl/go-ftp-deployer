package ftp

type File struct {
	content []byte
}

func (f *File) Write(p []byte) (n int, err error) {
	f.content = p
	return 0, nil
}

func (f *File) Content() []byte {
	return f.content
}
