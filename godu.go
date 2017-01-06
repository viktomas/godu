package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/gdamore/tcell"
	"github.com/viktomas/godu/core"
	"github.com/viktomas/godu/interactive"
)

func main() {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e = s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	s.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorBlack).
		Background(tcell.ColorWhite))
	s.Clear()
	limit := flag.Int64("l", 10, "show only files larger than limit (in MB)")
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}
	tree := core.GetSubTree(roots[0], ioutil.ReadDir, getIgnoredFolders())
	commands := make(chan core.Executer)
	quit := make(chan struct{})
	go interactive.InteractiveTree(&tree, s, commands, quit, *limit*core.MEGABYTE)
	go interactive.ParseCommand(s, commands, quit)

	cnt := 0
	dur := time.Duration(0)
loop:
	for {
		select {
		case <-quit:
			break loop
		case <-time.After(time.Millisecond * 50):
		}
		start := time.Now()
		//makebox(s)
		cnt++
		dur += time.Now().Sub(start)
	}

	s.Fini()
	fmt.Printf("Finished %d boxes in %s\n", cnt, dur)
	fmt.Printf("Average is %0.3f ms / box\n", (float64(dur)/float64(cnt))/1000000.0)
}
