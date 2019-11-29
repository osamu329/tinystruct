package syntax

type File struct {
	Types map[string]Type
}

func NewFile() *File {
	file := &File{Types: make(map[string]Type)}
	for k, t := range types {
		file.Types[k] = t
	}
	return file
}
