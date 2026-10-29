package main

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
	"github.com/viktomas/godu/commands"
	"github.com/viktomas/godu/files"
	"github.com/viktomas/godu/interactive"
)

func TestGetCell(t *testing.T) {
	vs := visualState{
		folders: []interactive.Line{
			{Text: []rune("Hello"), IsMarked: false},
			{Text: []rune("你好"), IsMarked: true},
		},
		selected:     1,
		screenHeight: 5,
	}

	tests := []struct {
		x, y          int
		expectedRune  rune
		expectedStyle tcell.Style
		expectedWidth int
	}{
		{0, 0, 'H', tcell.StyleDefault, 1},
		{1, 0, 'e', tcell.StyleDefault, 1},
		{0, 1, '你', tcell.StyleDefault.Reverse(true).Foreground(tcell.ColorGreen), 2},
		{1, 1, ' ', tcell.StyleDefault.Reverse(true).Foreground(tcell.ColorGreen), 1},
		{2, 1, '好', tcell.StyleDefault.Reverse(true).Foreground(tcell.ColorGreen), 2},
		{0, 2, ' ', tcell.StyleDefault, 1}, // out of bounds y
		{6, 0, ' ', tcell.StyleDefault, 1}, // out of bounds x
	}

	for _, tt := range tests {
		r, style, _, width := vs.GetCell(tt.x, tt.y)
		if r != tt.expectedRune {
			t.Errorf("GetCell(%d, %d) = %c; want %c", tt.x, tt.y, r, tt.expectedRune)
		}
		if style != tt.expectedStyle {
			t.Errorf("GetCell(%d, %d) style = %+v; want %+v", tt.x, tt.y, style, tt.expectedStyle)
		}
		if width != tt.expectedWidth {
			t.Errorf("GetCell(%d, %d) width = %d; want %d", tt.x, tt.y, width, tt.expectedWidth)
		}
	}
}

func TestNewVisualStateMeasuresBoundsInDisplayWidth(t *testing.T) {
	folder := files.NewTestFolder(".",
		files.NewTestFile("Hello", 0),
		files.NewTestFile("你好世界", 0),
	)
	vs := newVisualState(commands.State{Folder: folder}, 10)
	// each CJK rune occupies two columns, so "    0B 你好世界" is
	// 15 columns wide and its last addressable column is 14 —
	// measuring runes instead would report xbound of only 10
	xbound, ybound := vs.GetBounds()
	assert.Equal(t, 14, xbound)
	assert.Equal(t, 2, ybound)
}
