package main

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/viktomas/godu/commands"
	"github.com/viktomas/godu/files"
)

func captureStdout(t *testing.T, print func()) string {
	t.Helper()
	reader, writer, err := os.Pipe()
	assert.Nil(t, err)
	original := os.Stdout
	os.Stdout = writer
	print()
	writer.Close()
	os.Stdout = original
	output, err := io.ReadAll(reader)
	assert.Nil(t, err)
	return string(output)
}

func TestPrintMarkedFiles(t *testing.T) {
	root := files.NewTestFolder(".", files.NewTestFile("normal", 0), files.NewTestFile("multi\nline\nname", 0))
	state := &commands.State{
		MarkedFiles: map[*files.File]struct{}{
			root.Files[0]: {},
			root.Files[1]: {},
		},
	}
	output := captureStdout(t, func() { printMarkedFiles(state, false) })
	// FilesAsSlice prints longer paths first so that nested files are
	// deleted before the folders containing them
	assert.Equal(t, "\"multi\\nline\\nname\"\nnormal\n", output)
}

func TestPrintMarkedFilesNullTerminated(t *testing.T) {
	root := files.NewTestFolder(".", files.NewTestFile("multi\nline\nname", 0))
	state := &commands.State{
		MarkedFiles: map[*files.File]struct{}{
			root.Files[0]: {},
		},
	}
	output := captureStdout(t, func() { printMarkedFiles(state, true) })
	// paths are passed through verbatim so that xargs -0 can process them
	assert.Equal(t, "multi\nline\nname\x00", output)
}
