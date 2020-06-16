package worker

import (
	"fmt"
	"github.com/go-vgo/robotgo"
	"rp-bot-client/src/window"
	"time"
)

const (
	mineHoldTime = 6 * time.Second
	mineDigTime  = 7 * time.Second
)

func NewMiner(mineBtn string, pid int32) *Miner {
	return &Miner{
		mineBtn:   mineBtn,
		pid:       pid,
		running:   true,
		stateChan: make(chan bool),
	}
}

func (w *Miner) Start() {
	go w.dig()
}

type Miner struct {
	mineBtn string
	pid     int32
	running bool

	stateChan chan bool
}

func (w *Miner) Interrupt() {
	fmt.Println("[*] Debug: Before interrupt")
	w.stateChan <- false
	fmt.Println("[*] Debug: After interrupt")
}

func (w *Miner) Resume() {
	fmt.Println("[*] Debug: Before resume")
	w.stateChan <- true
	fmt.Println("[*] Debug: After resume")
}

func (w *Miner) Restart() {
	w.Resume()
}

func (w *Miner) ToggleHoldTime() {}

func (w *Miner) dig() {
	fmt.Println("[*] Debug: Dig once")

	digCh := make(chan struct{})

	w.holdMineBtn()
	timer := time.NewTimer(mineHoldTime)

	for {
		select {
		case <-digCh:
			if !w.running {
				time.Sleep(100 * time.Millisecond)
				continue
			}
			w.holdMineBtn()
			timer.Reset(mineHoldTime)
		// hold 6 sec to activate animation
		case <-timer.C:
			w.releaseMineBtn()
			// release and wait 7 sec until animation is finished
			time.Sleep(mineDigTime)
			// worker finished to dig - time to check captcha
			// checkCaptchaCh <- struct{}{}
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

func (w *Miner) holdMineBtn() {
	err := window.ActivatePidAndRun(w.pid, func() error {
		robotgo.KeyToggle(w.mineBtn, "down")
		fmt.Println("[*] Debug: mine key down")
		return nil
	})

	if err != nil {
		panic(err)
	}
}

func (w *Miner) releaseMineBtn() {
	err := window.ActivatePidAndRun(w.pid, func() error {
		robotgo.KeyToggle(w.mineBtn, "up")
		fmt.Println("[*] Debug: mine key up")
		return nil
	})

	if err != nil {
		panic(err)
	}
}
