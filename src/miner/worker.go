package miner

import (
	"fmt"
	"github.com/go-vgo/robotgo"
	"rp-bot-client/src/window"
	"time"
)

const (
	holdTime = 6 * time.Second
	digTime  = 7 * time.Second
)

func NewWorker(mineBtn string, pid int32) *Worker {
	return &Worker{
		mineBtn:   mineBtn,
		pid:       pid,
		running:   true,
		stateChan: make(chan bool),
	}
}

func (w *Worker) Start(checkCaptchaCh chan<- struct{}) {
	go w.dig(checkCaptchaCh)
}

type Worker struct {
	mineBtn string
	pid     int32
	running bool

	stateChan chan bool
}

func (w *Worker) Interrupt() {
	fmt.Println("[*] Debug: Before interrupt")
	w.stateChan <- false
	fmt.Println("[*] Debug: After interrupt")
}

func (w *Worker) Resume() {
	fmt.Println("[*] Debug: Before resume")
	w.stateChan <- true
	fmt.Println("[*] Debug: After resume")
}

func (w *Worker) dig(checkCaptchaCh chan<- struct{}) {
	fmt.Println("[*] Debug: Dig once")

	digCh := make(chan struct{})

	w.holdMineBtn()
	timer := time.NewTimer(holdTime)

	for {
		select {
		case <-digCh:
			if !w.running {
				time.Sleep(100 * time.Millisecond)
				continue
			}
			w.holdMineBtn()
			timer.Reset(holdTime)
		// hold 6 sec to activate animation
		case <-timer.C:
			w.releaseMineBtn()
			// release and wait 7 sec until animation is finished
			time.Sleep(digTime)
			// miner finished to dig - time to check captcha
			checkCaptchaCh <- struct{}{}
			fmt.Println(fmt.Sprintf("[*] Debug: before send to dig"))
			go func() { digCh <- struct{}{} }()
			fmt.Println(fmt.Sprintf("[*] Debug: after send to dig"))
		case w.running = <-w.stateChan:
			if !w.running {
				w.releaseMineBtn()
				fmt.Println("[*] Debug: Miner was interrupted")
			} else {
				go func() { digCh <- struct{}{} }()
				fmt.Println("[*] Debug: Miner was resumed")
			}
		}
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
