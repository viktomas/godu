package core

// TstFolder is providing easy interface to create foders for automated tests
// Never use in production code!
func TstFolder(name string, files ...*File) *File {
	folder := &File{name, nil, 0, true, []*File{}}
	if files == nil {
		return folder
	}
	var size int64
	for _, file := range files {
		size += file.Size
		file.Parent = folder
	}
	folder.Size = size
	folder.Files = files
	return folder
}

// TstFile provides easy interface to craete files for automated tests
// Never use in production code!
func TstFile(name string, size int64) *File {
	return &File{name, nil, size, false, []*File{}}
}
