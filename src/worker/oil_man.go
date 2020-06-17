package worker

import (
	"github.com/go-vgo/robotgo"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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
	maxBarrelsCountInInventory = 100
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

func NewOilMan(pid int32, checker *captcha.Checker, solver *captcha.Solver, manipulator *storage.Manipulator, withStorage bool) *OilMan {
	return &OilMan{
		pid:         pid,
		running:     true,
		holdTime:    oilHoldLongTime,
		withStorage: withStorage,

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
	log.Debug("Before interrupt")
	w.stateChan <- false
	log.Debug("After interrupt")
}

func (w *OilMan) Resume() {
	log.Debug("Before resume")
	w.stateChan <- true
	log.Debug("After resume")
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

	// todo move to config to fit into client system requirements, make 1 sec as default
	<-time.After(2 * time.Second)

	w.Restart()
}

func (w *OilMan) oil() {
	log.Debug("Starting to oil")

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
		case <-timer.C:
			if err := w.releaseOilAndCheckCaptcha(); err != nil {
				log.WithError(err).Error("failed to check of solve captcha")
				go w.ReEnterWindow()
				continue
			}

			log.Debug("before send to oil ch")
			go func() { oilCh <- struct{}{} }()
			log.Debug("after send to oil ch")
		case w.running = <-w.stateChan:
			if !w.running {
				w.releaseOil()
				log.Debug("Oilman was interrupted")
			} else {
				go func() { oilCh <- struct{}{} }()
				log.Debug("Oilman was resumed")
			}
		}
	}
}

func (w *OilMan) releaseOilAndCheckCaptcha() error {
	log.Debugf("current oil: %d", w.currentOil)
	w.currentOil++
	w.releaseOil()
	err := w.checkCaptchaAndSolveIfNeeded()
	if errors.Is(err, captchaSolveErr) || errors.Is(err, captchaNotAppearedTooManyTimesErr) {
		return err
	} else if err != nil {
		log.WithError(err).Fatalf("unknown error")
	}

	time.Sleep(100 * time.Millisecond)

	return nil
}

func (w *OilMan) checkCaptchaAndSolveIfNeeded() error {
	if !w.isOilMiningIterationFinished() {
		return nil
	}

	defer func() { w.currentOil = 0 }()

	// worker finished to oil - time to check captcha
	if w.captchaChecker.IsCaptchaAppeared(w.pid) {
		log.Info("captcha appeared: solving...")
		w.captchaNotAppearedTimes = 0
		if err := w.captchaSolver.Solve(); err != nil {
			log.WithError(err).Error("failed to solve captcha")

			return captchaSolveErr
		}

		return nil
	}

	w.moveBarrelsToStorageIfNeeded()

	if w.captchaNotAppearedTimes > 4 {
		return captchaNotAppearedTooManyTimesErr
	}

	w.captchaNotAppearedTimes++
	log.Debugf("Captcha not appeared times: %d", w.captchaNotAppearedTimes)

	return nil
}

func (w *OilMan) moveBarrelsToStorageIfNeeded() {
	if !w.withStorage {
		return
	}

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
		log.Debug("oil mouse key down")

		return nil
	})

	if err != nil {
		panic(err)
	}
}

func (w *OilMan) releaseOil() {
	err := window.ActivatePidAndRun(w.pid, func() error {
		robotgo.MouseToggle("up")
		log.Debug("oil mouse key up")

		return nil
	})

	if err != nil {
		panic(err)
	}
}

func (w *OilMan) pressEsc() {
	err := window.ActivatePidAndRun(w.pid, func() error {
		robotgo.KeyTap("esc")
		log.Debug("esc key tap")

		return nil
	})

	if err != nil {
		panic(err)
	}
}

func (w *OilMan) pressE() {
	err := window.ActivatePidAndRun(w.pid, func() error {
		robotgo.KeyTap("e")
		log.Debug("e key tap")

		return nil
	})

	if err != nil {
		panic(err)
	}
}
