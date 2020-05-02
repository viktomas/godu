package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIgnoreBasedOnIgnoreFile(t *testing.T) {
	ignored := []string{"node_modules"}
	ignoreFunction := ignoreBasedOnIgnoreFile(ignored)
	assert.True(t, ignoreFunction("something/node_modules"))
	assert.False(t, ignoreFunction("something/notIgnored"))
}
