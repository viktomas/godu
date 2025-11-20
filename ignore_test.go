package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIgnoreBasedOnIgnoreFile(t *testing.T) {
	//these tests might fail on windows
	tests := []struct {
		ignored []string
		file string
		match   bool
	}{
		{[]string{"**/what/test"}, "what/what/test", true},
		{[]string{"**node_modules"}, "something/node_modules", true},
		{[]string{"**node_modules"}, "something/notIgnored", false},
		{[]string{"**/what/test"}, "/home/user/test", false},
		{[]string{"**user/test"}, "user/test", true},
		{[]string{"**test"}, "test", true},
		{[]string{"**png"}, "path/to/test.png", true},
		{[]string{"**folder**png"}, "path/to/folder/test.png", true},
		{[]string{"**folder**png"}, "path/folder/to/test.png", true},
		{[]string{"**/*folder*/**png"}, "path/hellofolder/to/test.png", true},
		{[]string{"**/*folder*/**png"}, "path/folder2/to/test.png", true},
		{[]string{"**/*folder*/**png"}, "path/folder/to/test.png", true},
		{[]string{"**/*folder*/**png"}, "path/foler/to/test.png", false},
		{[]string{"**/test.png"}, "path/to/test.png", true},
		{[]string{"*png"}, "test.png", true},
		{[]string{"dir/*png"}, "dir/test.png", true},
		{[]string{"*dir*/*png"}, "testdir2/test.png", true},
		{[]string{"*dir*/*png"}, "testdir/test.png", true},
		{[]string{"*dir*dir*/*png"}, "dirtestdir/test.png", true},
		{[]string{"*dir*dir*/*png"}, "dir/testdir/test.png", false},
		{[]string{"wrong/pattern"}, "test", false},
	}

	for _, tt := range tests {
		ignoreFunction := ignoreBasedOnIgnoreFile(tt.ignored)
		assert.True(t, tt.match == ignoreFunction(tt.file))
	}
}
