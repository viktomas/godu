package main

import (
	"sync"

	"github.com/gdamore/tcell"
	"github.com/viktomas/godu/commands"
)

func parseCommand(s tcell.Screen, commandsChan chan commands.Executer, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape, tcell.KeyCtrlC:
				close(commandsChan)
				return
			case tcell.KeyEnter, tcell.KeyRight:
				commandsChan <- commands.Enter{}
			case tcell.KeyDown:
				commandsChan <- commands.Down{}
			case tcell.KeyUp:
				commandsChan <- commands.Up{}
			case tcell.KeyBackspace, tcell.KeyLeft:
				commandsChan <- commands.GoBack{}
			case tcell.KeyCtrlL:
				s.Sync()
			case tcell.KeyRune:
				switch ev.Rune() {
				case ' ':
					commandsChan <- commands.Mark{}
				case 'q':
					close(commandsChan)
					return
				case 'h':
					commandsChan <- commands.GoBack{}
				case 'j':
					commandsChan <- commands.Down{}
				case 'k':
					commandsChan <- commands.Up{}
				case 'l':
					commandsChan <- commands.Enter{}
				}

			}
		case *tcell.EventResize:
			s.Sync()
		}
	}
}
