package bot

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"rp-bot-client/src/event"
	"rp-bot-client/src/worker"
)

var stopErr = errors.New("stop application")

func NewBot(
	pid int32,

	minerWorker worker.Worker,
	eventListener *event.Listener,
) *Bot {
	return &Bot{
		pid:     pid,
		running: true,

		worker:                  minerWorker,
		eventListener:           eventListener,
	}
}

type Bot struct {
	pid     int32
	running bool

	worker                  worker.Worker
	eventListener           *event.Listener
}

func (b *Bot) Start() error {
	b.mainLoop()

	return nil
}

func (b *Bot) mainLoop() {
	eventCh := b.eventListener.Start()
	b.worker.Start()

	log.Info("Bot have started")
	for {
		select {
		case e := <-eventCh:
			if err := b.handleEvent(e); err == stopErr {
				return
			} else if err != nil {
				panic(fmt.Sprintf("failed to handle event: %s", err))
			}
		}
	}
}

func (b *Bot) handleEvent(e event.Event) error {
	if e.IsStop() {
		return stopErr
	}

	if e.IsPause() {
		b.running = false
		b.worker.Interrupt()
	}

	if e.IsResume() {
		b.running = true
		b.worker.Resume()
	}

	if e.IsRestart() {
		b.running = true
		b.worker.Restart()
	}

	return nil
}
