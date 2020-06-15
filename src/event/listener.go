package event

import (
	"fmt"
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

func NewEventListener() *Listener {
	return &Listener{outCh: make(chan Event, 100)}
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
	robotgo.EventHook(hook.KeyDown, []string{"o", "ctrl"}, func(e hook.Event) {
		fmt.Println("ctrl-o: pause bot")
		l.outCh <- Event{T: pause}
		fmt.Println("ctrl-o: event sent")
	})

	robotgo.EventHook(hook.KeyDown, []string{"r", "ctrl"}, func(e hook.Event) {
		fmt.Println("ctrl-r: resume bot")
		l.outCh <- Event{T: resume}
		fmt.Println("ctrl-r: event sent")
	})

	robotgo.EventHook(hook.KeyDown, []string{"t", "ctrl"}, func(e hook.Event) {
		fmt.Println("ctrl-t: resume bot")
		l.outCh <- Event{T: restart}
		fmt.Println("ctrl-t: event sent")
	})

	robotgo.EventHook(hook.KeyDown, []string{"c", "ctrl"}, func(e hook.Event) {
		fmt.Println("ctrl-c: stop bot. Exiting...")
		l.outCh <- Event{T: stop}
		fmt.Println("ctrl-c: event sent")
		robotgo.EventEnd()
	})

	<-robotgo.EventProcess(robotgo.EventStart())
}

func writeHelpMessage() {
	shortcuts := []string{
		"ctrl+o to pause bot",
		"ctrl+r to resume bot",
		"ctrl+c to stop bot",
		"ctrl+t to restart bot",
	}

	fmt.Println("[*] Please use this keyboard shortcuts to control the bot: ")

	for _, msg := range shortcuts {
		fmt.Println(fmt.Sprintf("--- %s ---", msg))
	}
}
