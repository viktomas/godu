package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gdamore/tcell"
	"github.com/gosuri/uilive"
	"github.com/viktomas/godu/commands"
	"github.com/viktomas/godu/files"

	"github.com/viktomas/godu/interactive"
)

// the correct version is injected by `go build` command in release.sh script
var goduVersion = "master"

func main() {
	limit := flag.Int64("l", 10, "show only files larger than limit (in MB)")
	nullTerminate := flag.Bool("print0", false, "print null-terminated strings")
	version := flag.Bool("v", false, "show version")
	flag.Parse()
	if *version {
		fmt.Printf("godu %s\n", goduVersion)
		os.Exit(0)
	}
	args := flag.Args()
	rootFolderName := "."
	if len(args) > 0 {
		rootFolderName = args[0]
	}
	rootFolderName, err := filepath.Abs(rootFolderName)
	if err != nil {
		log.Fatalln(err.Error())
	}
	progress := make(chan int)
	go reportProgress(progress)
	rootFolder := files.WalkFolder(rootFolderName, ioutil.ReadDir, ignoreBasedOnIgnoreFile(readIgnoreFile()), progress)
	rootFolder.Name = rootFolderName
	err = commands.ProcessFolder(rootFolder, *limit*files.MEGABYTE)
	if err != nil {
		log.Fatalln(err.Error())
	}
	s := initScreen()
	commandsChan := make(chan commands.Executer)
	states := make(chan commands.State)
	lastStateChan := make(chan *commands.State, 1)
	var wg sync.WaitGroup
	wg.Add(3)
	go commands.StartProcessing(rootFolder, commandsChan, states, lastStateChan, &wg)
	go interactiveFolder(s, states, &wg)
	go parseCommand(s, commandsChan, &wg)
	wg.Wait()
	s.Fini()
	lastState := <-lastStateChan
	printMarkedFiles(lastState, *nullTerminate)
}

func reportProgress(progress <-chan int) {
	const interval = 50 * time.Millisecond
	writer := uilive.New()
	writer.Out = os.Stderr
	writer.Start()
	defer writer.Stop()
	total := 0
	ticker := time.NewTicker(interval)
	for {
		select {
		case c, ok := <-progress:
			if !ok {
				return
			}
			total += c
		case <-ticker.C:
			fmt.Fprintf(writer, "Walked through %d folders\n", total)
		}
	}
}

func printMarkedFiles(lastState *commands.State, nullTerminate bool) {
	markedFiles := interactive.FilesAsSlice(lastState.MarkedFiles)
	var printFunc func(string)
	if nullTerminate {
		printFunc = func(s string) {
			fmt.Printf("%s\x00", s)
		}
	} else {
		printFunc = func(s string) {
			fmt.Println(s)
		}
	}
	for _, f := range markedFiles {
		printFunc(f)
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
