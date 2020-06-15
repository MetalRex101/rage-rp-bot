package worker

import (
	"fmt"
	"github.com/go-vgo/robotgo"
	"rp-bot-client/src/window"
	"time"
)

const (
	oilHoldShortTime = 3500 * time.Millisecond
	oilHoldLongTime  = 4500 * time.Millisecond
)

type coordinates struct {
	x int
	y int
}

var oilCoordinates = map[int]coordinates{
	0: {283, 675},
	1: {388, 672},
	2: {494, 673},
	3: {595, 678},
}

func NewOilMan(pid int32) *OilMan {
	return &OilMan{
		pid:      pid,
		running:  true,
		holdTime: oilHoldLongTime,

		stateChan: make(chan bool),
	}
}

func (w *OilMan) Start(checkCaptchaCh chan<- struct{}) {
	go w.oil(checkCaptchaCh)
}

type OilMan struct {
	mineBtn    string
	pid        int32
	running    bool
	currentOil int
	holdTime   time.Duration

	stateChan chan bool
}

func (w *OilMan) Interrupt() {
	fmt.Println("[*] Debug: Before interrupt")
	w.stateChan <- false
	fmt.Println("[*] Debug: After interrupt")
}

func (w *OilMan) Resume() {
	fmt.Println("[*] Debug: Before resume")
	w.stateChan <- true
	fmt.Println("[*] Debug: After resume")
}

func (w *OilMan) Restart() {
	w.Interrupt()
	w.currentOil = 0
	time.Sleep(100 * time.Millisecond)
	w.Resume()
}

func (w *OilMan) ToggleHoldTime() {
	if w.holdTime == oilHoldLongTime {
		w.holdTime = oilHoldShortTime
	} else {
		w.holdTime = oilHoldLongTime
	}
	w.Restart()
}

func (w *OilMan) oil(checkCaptchaCh chan<- struct{}) {
	fmt.Println("[*] Debug: Starting to oil")

	oilCh := make(chan struct{})

	w.holdOil()
	timer := time.NewTimer(w.holdTime)

	for {
		select {
		case <-oilCh:
			if !w.running {
				time.Sleep(100 * time.Millisecond)
				continue
			}
			w.holdOil()
			timer.Reset(w.holdTime)
		// hold 6 sec to activate animation
		case <-timer.C:
			fmt.Println(fmt.Sprintf("[*] Debug: current oil: %d", w.currentOil))
			w.currentOil++
			w.releaseOil()

			// check captcha after last oil
			if w.currentOil == 4 {
				// worker finished to oil - time to check captcha
				checkCaptchaCh <- struct{}{}

				w.currentOil = 0
			}

			fmt.Println(fmt.Sprintf("[*] Debug: before send to oil ch"))
			go func() { oilCh <- struct{}{} }()
			fmt.Println(fmt.Sprintf("[*] Debug: after send to oil ch"))
		case w.running = <-w.stateChan:
			if !w.running {
				w.releaseOil()
				fmt.Println("[*] Debug: Oilman was interrupted")
			} else {
				go func() { oilCh <- struct{}{} }()
				fmt.Println("[*] Debug: Oilman was resumed")
			}
		}
	}
}

func (w *OilMan) holdOil() {
	err := window.ActivatePidAndRun(w.pid, func() error {
		coord := oilCoordinates[w.currentOil]
		robotgo.Move(coord.x, coord.y)
		robotgo.MouseToggle("down")
		fmt.Println("[*] Debug: oil mouse key down")

		return nil
	})

	if err != nil {
		panic(err)
	}
}

func (w *OilMan) releaseOil() {
	err := window.ActivatePidAndRun(w.pid, func() error {
		robotgo.MouseToggle("up")
		fmt.Println("[*] Debug: oil mouse key up")

		return nil
	})

	if err != nil {
		panic(err)
	}
}
