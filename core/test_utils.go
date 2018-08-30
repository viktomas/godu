package core

// NewTestFolder is providing easy interface to create foders for automated tests
// Never use in production code!
func NewTestFolder(name string, files ...*File) *File {
	folder := &File{name, nil, 0, true, []*File{}}
	if files == nil {
		return folder
	}
	for _, file := range files {
		file.Parent = folder
	}
	folder.Files = files
	folder.UpdateSize()
	return folder
}

// NewTestFile provides easy interface to craete files for automated tests
// Never use in production code!
func NewTestFile(name string, size int64) *File {
	return &File{name, nil, size, false, []*File{}}
}

// FindTestFile helps testing by returning first occurance of file with given name.
// Never use in produciton code!
func FindTestFile(folder *File, name string) *File {
	if folder.Name == name {
		return folder
	}
	for _, file := range folder.Files {
		result := FindTestFile(file, name)
		if result != nil {
			return result
		}
	}
	return nil
}
