package worker

import (
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/pkg/errors"
	"rp-bot-client/src/captcha"
	"rp-bot-client/src/storage"
	"rp-bot-client/src/window"
	"time"
)

var (
	captchaSolveErr                   = errors.New("failed to solve captcha")
	captchaNotAppearedTooManyTimesErr = errors.New("captcha has not appeared too many times")
)

const (
	oilHoldShortTime           = 3500 * time.Millisecond
	oilHoldLongTime            = 4500 * time.Millisecond
	maxBarrelsCountInInventory = 1
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

func NewOilMan(pid int32, checker *captcha.Checker, solver *captcha.Solver, manipulator *storage.Manipulator) *OilMan {
	return &OilMan{
		pid:         pid,
		running:     true,
		holdTime:    oilHoldLongTime,
		withStorage: true,

		captchaChecker:     checker,
		captchaSolver:      solver,
		storageManipulator: manipulator,

		stateChan: make(chan bool),
	}
}

type OilMan struct {
	pid                     int32
	running                 bool
	captchaNotAppearedTimes int
	currentOil              int
	holdTime                time.Duration
	withStorage             bool
	barrelsCounter          int

	captchaChecker     *captcha.Checker
	captchaSolver      *captcha.Solver
	storageManipulator *storage.Manipulator

	stateChan chan bool
}

func (w *OilMan) Start() {
	go w.oil()
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

func (w *OilMan) ReEnterWindow() {
	w.pressEsc()

	<-time.After(time.Millisecond * 300)

	w.pressE()

	w.Restart()
}

func (w *OilMan) oil() {
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
			err := w.checkCaptchaAndSolveIfNeeded()
			if errors.Is(err, captchaSolveErr) {
				fmt.Println(fmt.Sprintf("[*] Debug: %s. Reentering window", err))
				go w.ReEnterWindow()
				continue
			} else if errors.Is(err, captchaNotAppearedTooManyTimesErr) {
				fmt.Println(fmt.Sprintf("[*] Debug: %s. Reentering window", err))
				go w.ReEnterWindow()
				continue
			}

			if err != nil {
				panic(fmt.Sprintf("unknown error: %s", err))
			}

			time.Sleep(100 * time.Millisecond)

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

func (w *OilMan) checkCaptchaAndSolveIfNeeded() error {
	if !w.isOilMiningIterationFinished() {
		return nil
	}

	defer func() { w.currentOil = 0 }()

	// worker finished to oil - time to check captcha
	if w.captchaChecker.IsCaptchaAppeared(w.pid) {
		fmt.Println(fmt.Sprintf("[*] Debug: captcha appeared: solving..."))
		w.captchaNotAppearedTimes = 0
		if err := w.captchaSolver.Solve(); err != nil {
			fmt.Println(fmt.Sprintf("[*] Debug: %s", err))

			return captchaSolveErr
		}

		return nil
	}

	w.moveBarrelsToStorageIfNeeded()

	if w.captchaNotAppearedTimes > 4 {
		return captchaNotAppearedTooManyTimesErr
	}

	w.captchaNotAppearedTimes++
	fmt.Println(fmt.Sprintf("[*] Debug: Captcha not appeared times: %d", w.captchaNotAppearedTimes))

	return nil
}

func (w *OilMan) moveBarrelsToStorageIfNeeded() {
	w.barrelsCounter++
	if w.barrelsCounter < maxBarrelsCountInInventory {
		return
	}
	w.barrelsCounter = 0

	w.pressEsc()
	w.storageManipulator.ReplaceItemFromInventoryToStorage()
	w.pressE()

	<-time.After(time.Second)
}

func (w *OilMan) isOilMiningIterationFinished() bool {
	return w.currentOil == 4
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

func (w *OilMan) pressEsc() {
	err := window.ActivatePidAndRun(w.pid, func() error {
		robotgo.KeyTap("esc")
		fmt.Println("[*] Debug: esc key tap")

		return nil
	})

	if err != nil {
		panic(err)
	}
}

func (w *OilMan) pressE() {
	err := window.ActivatePidAndRun(w.pid, func() error {
		robotgo.KeyTap("e")
		fmt.Println("[*] Debug: e key tap")

		return nil
	})

	if err != nil {
		panic(err)
	}
}
