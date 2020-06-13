package bot

import (
	"errors"
	"fmt"
	"rp-bot-client/src/captcha"
	"rp-bot-client/src/event"
	"rp-bot-client/src/miner"
	"time"
)

var stopErr = errors.New("stop application")

func NewBot(
	pid int32,

	minerWorker *miner.Worker,
	captchaSolver *captcha.Solver,
	captchaMouseManipulator *captcha.MouseManipulator,
	eventListener *event.Listener,
) *Bot {
	return &Bot{
		pid: pid,

		minerWorker:             minerWorker,
		captchaSolver:           captchaSolver,
		captchaMouseManipulator: captchaMouseManipulator,
		eventListener:           eventListener,
	}
}

type Bot struct {
	pid     int32
	running bool

	minerWorker             *miner.Worker
	captchaSolver           *captcha.Solver
	captchaMouseManipulator *captcha.MouseManipulator
	eventListener           *event.Listener
}

func (b *Bot) Start() error {
	b.mainLoop()

	return nil
}

func (b *Bot) mainLoop() {
	checkCaptchaChan := make(chan struct{})

	eventCh := b.eventListener.Start()
	captchaSolvedCh := b.captchaSolver.Start(checkCaptchaChan)

	for {
		select {
		case e := <-eventCh:
			if err := b.handleEvent(e); err == stopErr {
				return
			} else if err != nil {
				panic(fmt.Sprintf("failed to handle event: %s", err))
			}
		case answerNum := <-captchaSolvedCh:
			if b.running {
				if err := b.captchaMouseManipulator.Answer(answerNum); err != nil {
					panic(fmt.Sprintf("captcha manipulator error: %s", err))
				}
			}
		default:
			if !b.running {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			<-b.minerWorker.DigOreOnce()
			checkCaptchaChan <- struct{}{}
		}
	}
}

func (b *Bot) handleEvent(e event.Event) error {
	if e.IsStop() {
		return stopErr
	}

	if e.IsPause() {
		b.running = false
		b.minerWorker.Interrupt()
	}

	if e.IsResume() {
		b.running = true
	}

	return nil
}
