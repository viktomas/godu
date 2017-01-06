package interactive

import (
	"github.com/gdamore/tcell"
	"github.com/viktomas/godu/core"
)

func ParseCommand(s tcell.Screen, commands chan core.Executer, quit chan struct{}) {
	for {
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape, tcell.KeyCtrlC:
				close(quit)
				return
			case tcell.KeyEnter:
				commands <- core.Enter{1}
			case tcell.KeyCtrlL:
				s.Sync()
			}
		case *tcell.EventResize:
			s.Sync()
		}
	}
}
