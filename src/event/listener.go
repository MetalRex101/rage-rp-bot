package event

import (
	"fmt"
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

func NewEventListener() *Listener {
	return &Listener{outCh: make(chan Event)}
}

type Listener struct {
	outCh chan Event
}

func (l *Listener) Start () <-chan Event {
	go l.start()

	return l.outCh
}

func (l *Listener) start() {
	writeHelpMessage()
	robotgo.EventHook(hook.KeyDown, []string{"p", "ctrl"}, func(e hook.Event) {
		fmt.Println("ctrl-p: pause bot")
		l.outCh <- Event{T: pause}
	})

	robotgo.EventHook(hook.KeyDown, []string{"r", "ctrl"}, func(e hook.Event) {
		fmt.Println("ctrl-r: resume bot")
		l.outCh <- Event{T: resume}
	})

	robotgo.EventHook(hook.KeyDown, []string{"c", "ctrl"}, func(e hook.Event) {
		fmt.Println("ctrl-c: stop bot. Exiting...")
		l.outCh <- Event{T: stop}
		robotgo.EventEnd()
	})

	<-robotgo.EventProcess(robotgo.EventStart())
}

func writeHelpMessage() {
	shortcuts := []string{
		"ctrl+p to pause bot",
		"ctrl+r to resume bot",
		"ctrl+c tp stop bot",
	}

	fmt.Println("[*] Please use this keyboard shortcuts to control the bot: ")

	for _, msg := range shortcuts {
		fmt.Println(fmt.Sprintf("--- %s ---", msg))
	}
}
