package syntax

type File struct {
	Types map[string]Type
    TypeList []Type

}

func NewFile() *File {
	file := &File{Types: make(map[string]Type)}
	for k, t := range types {
		file.Types[k] = t
	}
	return file
}

func (f *File) AddType(name string, t Type) {
    f.Types[name] = t
    f.TypeList = append(f.TypeList, t)
}
