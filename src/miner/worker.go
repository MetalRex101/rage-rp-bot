package miner

import (
	"fmt"
	"github.com/go-vgo/robotgo"
	"rp-bot-client/src/window"
	"time"
)

func NewWorker(mineBtn string, pid int32) *Worker {
	return &Worker{
		mineBtn:     mineBtn,
		pid:         pid,
		interruptCh: make(chan struct{}),
	}
}

type Worker struct {
	mineBtn string
	pid     int32

	interruptCh chan struct{}
}

func (w *Worker) DigOreOnce() {
	w.dig()
}

func (w *Worker) Interrupt() {
	w.interruptCh <- struct{}{}
}

func (w *Worker) dig() {
	fmt.Println("[*] Debug: Dig once")

	w.holdMineBtn()
	// hold 6 sec to activate animation
	toggleTimeout := time.NewTimer(time.Second * 6)

	select {
	case <-toggleTimeout.C:
		w.releaseMineBtn()
		// release and wait 7 sec until animation is finished
		time.Sleep(time.Second * 7)
	case <-w.interruptCh:
		w.releaseMineBtn()
		return
	}
}

func (w *Worker) holdMineBtn() {
	err := window.ActivatePidAndRun(w.pid, func() error {
		robotgo.KeyToggle(w.mineBtn, "down")
		fmt.Println("[*] Debug: mine key down")
		return nil
	})

	if err != nil {
		panic(err)
	}
}

func (w *Worker) releaseMineBtn() {
	err := window.ActivatePidAndRun(w.pid, func() error {
		robotgo.KeyToggle(w.mineBtn, "up")
		fmt.Println("[*] Debug: mine key up")
		return nil
	})

	if err != nil {
		panic(err)
	}
}
