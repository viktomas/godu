package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/gdamore/tcell"
	"github.com/viktomas/godu/core"
	"github.com/viktomas/godu/interactive"
)

func main() {
	limit := flag.Int64("l", 10, "show only files larger than limit (in MB)")
	flag.Parse()
	args := flag.Args()
	root := "."
	if len(args) > 0 {
		root = args[0]
	}
	root, err := filepath.Abs(root)
	if err != nil {
		log.Println(err.Error())
		os.Exit(0)
	}
	log.Printf("godu will walk through `%s` that might take up to few minutes\n", root)
	tree := core.WalkFolder(root, ioutil.ReadDir, getIgnoredFolders())
	tree.Name = root
	err = core.PrepareTree(tree, *limit*core.MEGABYTE)
	if err != nil {
		log.Println(err.Error())
		os.Exit(0)
	}
	s := initScreen()
	commands := make(chan core.Executer)
	states := make(chan core.State)
	lastStateChan := make(chan *core.State, 1)
	var wg sync.WaitGroup
	wg.Add(3)
	go core.StartProcessing(tree, commands, states, lastStateChan, &wg)
	go InteractiveTree(s, states, &wg)
	go ParseCommand(s, commands, &wg)
	wg.Wait()
	s.Fini()
	lastState := <-lastStateChan
	printMarkedFiles(lastState)
}

func printMarkedFiles(lastState *core.State) {
	markedFiles := interactive.QuoteMarkedFiles(lastState.MarkedFiles)
	for _, quotedFile := range markedFiles {
		fmt.Println(quotedFile)
	}
}

func initScreen() tcell.Screen {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	s, e := tcell.NewScreen()
	if e != nil {
		log.Printf("%v\n", e)
		os.Exit(1)
	}
	if e = s.Init(); e != nil {
		log.Printf("%v\n", e)
		os.Exit(1)
	}
	s.Clear()
	return s
}
