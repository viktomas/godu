package files

// NewTestFolder is providing easy interface to create folders for automated tests
// Never use in production code!
func NewTestFolder(name string, files ...*File) *File {
	folder := &File{Name: name, IsDir: true, Files: []*File{}}
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

// NewTestFile provides easy interface to create files for automated tests
// Never use in production code!
func NewTestFile(name string, size int64) *File {
	return &File{Name: name, Size: size, Files: []*File{}}
}

// FindTestFile helps testing by returning first occurrence of file with given name.
// Never use in production code!
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
