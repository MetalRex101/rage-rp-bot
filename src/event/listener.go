package event

import (
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
	log "github.com/sirupsen/logrus"
)

func NewEventListener() *Listener {
	return &Listener{outCh: make(chan Event, 100)}
}

type Listener struct {
	outCh chan Event
}

func (l *Listener) Start() <-chan Event {
	go l.start()

	return l.outCh
}

func (l *Listener) start() {
	writeHelpMessage()
	robotgo.EventHook(hook.KeyDown, []string{"o", "ctrl"}, func(e hook.Event) {
		log.Debug("ctrl-o: pause bot")
		l.outCh <- Event{T: pause}
		log.Debug("ctrl-o: event sent")
	})

	robotgo.EventHook(hook.KeyDown, []string{"r", "ctrl"}, func(e hook.Event) {
		log.Debug("ctrl-r: resume bot")
		l.outCh <- Event{T: resume}
		log.Debug("ctrl-r: event sent")
	})

	robotgo.EventHook(hook.KeyDown, []string{"t", "ctrl"}, func(e hook.Event) {
		log.Debug("ctrl-t: restart bot")
		l.outCh <- Event{T: restart}
		log.Debug("ctrl-t: event sent")
	})

	robotgo.EventHook(hook.KeyDown, []string{"y", "ctrl"}, func(e hook.Event) {
		log.Debug("ctrl-y: toggle speed")
		l.outCh <- Event{T: toggleHoldTime}
		log.Debug("ctrl-y: event sent")
	})

	robotgo.EventHook(hook.KeyDown, []string{"c", "ctrl"}, func(e hook.Event) {
		log.Debug("ctrl-c: stop bot. Exiting...")
		l.outCh <- Event{T: stop}
		log.Debug("ctrl-c: event sent")
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

	log.Info("Please use this keyboard shortcuts to control the bot: ")

	for _, msg := range shortcuts {
		log.Infof("--- %s ---", msg)
	}
}
