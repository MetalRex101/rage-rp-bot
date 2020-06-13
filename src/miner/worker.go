package miner

import (
	"github.com/go-vgo/robotgo"
	"time"
)

func NewWorker(mineBtn string) *Worker {
	return &Worker{
		mineBtn: mineBtn,
		finished: make(chan struct{}),
		interruptCh: make(chan struct{}),
	}
}

type Worker struct {
	mineBtn     string
	finished    chan struct{}
	interruptCh chan struct{}
}

func (w *Worker) DigOreOnce() <-chan struct{} {
	go w.dig()

	return w.finished
}

func (w *Worker) Interrupt() {
	w.interruptCh <- struct{}{}
}

func (w *Worker) dig() {
	defer robotgo.KeyToggle(w.mineBtn, "up")

	robotgo.KeyToggle(w.mineBtn, "down")

	toggleTimeout := time.NewTimer(time.Second * 6)

	select {
	case <-toggleTimeout.C:
		robotgo.KeyToggle(w.mineBtn, "up")
		time.Sleep(time.Second * 7)

		<- w.finished
	case <-w.interruptCh:
		return
	}
}
