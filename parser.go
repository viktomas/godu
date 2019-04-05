package main

import (
	"sync"

	"github.com/gdamore/tcell"
	"github.com/viktomas/godu/core"
)

func parseCommand(s tcell.Screen, commands chan core.Executer, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape, tcell.KeyCtrlC:
				close(commands)
				return
			case tcell.KeyEnter, tcell.KeyRight:
				commands <- core.Enter{}
			case tcell.KeyDown:
				commands <- core.Down{}
			case tcell.KeyUp:
				commands <- core.Up{}
			case tcell.KeyBackspace, tcell.KeyLeft:
				commands <- core.GoBack{}
			case tcell.KeyCtrlL:
				s.Sync()
			case tcell.KeyRune:
				switch ev.Rune() {
				case ' ':
					commands <- core.Mark{}
				case 'q':
					close(commands)
					return
				case 'h':
					commands <- core.GoBack{}
				case 'j':
					commands <- core.Down{}
				case 'k':
					commands <- core.Up{}
				case 'l':
					commands <- core.Enter{}
				}

			}
		case *tcell.EventResize:
			s.Sync()
		}
	}
}
