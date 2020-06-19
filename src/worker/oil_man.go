package worker

import (
	"context"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"os"
	"rp-bot-client/src/captcha"
	"rp-bot-client/src/storage"
	"time"
)

var (
	captchaSolveErr                   = errors.New("failed to solve captcha")
	captchaNotAppearedTooManyTimesErr = errors.New("captcha has not appeared too many times")
)

const (
	maxBarrelsCountInInventory = 100
)

func NewOilMan(
	pid int32,
	checker *captcha.Checker,
	solver *captcha.Solver,
	storageManipulator *storage.Manipulator,
	withStorage bool,
) *OilMan {
	return &OilMan{
		pid:         pid,
		running:     true,
		withStorage: withStorage,

		oilManipulator:     newOilManipulator(pid),
		captchaChecker:     checker,
		captchaSolver:      solver,
		storageManipulator: storageManipulator,

		stateChan: make(chan bool),
	}
}

type OilMan struct {
	pid                     int32
	running                 bool
	captchaNotAppearedTimes int
	withStorage             bool
	barrelsCounter          int

	oilManipulator     *OilManipulator
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
	go func() { w.stateChan <- false }()
	log.Debug("After interrupt")
}

func (w *OilMan) Resume() {
	log.Debug("Before resume")
	go func() { w.stateChan <- true }()
	log.Debug("After resume")
}

func (w *OilMan) Restart() {
	w.Interrupt()
	//w.currentOil = 0
	time.Sleep(100 * time.Millisecond)
	w.Resume()
}

func (w *OilMan) oil() {
	log.Debug("Starting to oil")

	oilCh := make(chan struct{})

	for {
		select {
		case <-oilCh:
			if !w.running {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			heldOilCoordinates := w.oilManipulator.holdOil()
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			isAllOilsDone := w.oilManipulator.releaseOilOnDone(heldOilCoordinates, ctx)

			if isAllOilsDone {
				solved, err := w.solveCaptchaIfNeeded()
				if err != nil {
					log.WithError(err).Error("failed to check or solve captcha")
					w.oilManipulator.ReOpenWindow()
					w.Restart()
					continue
				}
				if solved {
					w.Restart()
					continue
				}

				<-time.After(time.Second)
			}

			log.Debug("before send to oil ch")
			go func() { oilCh <- struct{}{} }()
			log.Debug("after send to oil ch")
		case w.running = <-w.stateChan:
			if !w.running {
				w.oilManipulator.releaseOil()
				log.Debug("Oilman was interrupted")
			} else {
				go func() { oilCh <- struct{}{} }()
				log.Debug("Oilman was resumed")
			}
		}
	}
}

func (w *OilMan) solveCaptchaIfNeeded() (bool, error) {
	solved, err := w.checkCaptchaAndSolveIfNeeded()
	if errors.Is(err, captchaSolveErr) || errors.Is(err, captchaNotAppearedTooManyTimesErr) {
		return false, err
	} else if err != nil {
		log.WithError(err).Fatalf("unknown error")
	}

	time.Sleep(100 * time.Millisecond)

	return solved, nil
}

func (w *OilMan) checkCaptchaAndSolveIfNeeded() (bool, error) {
	defer func() {
		log.Debug("Oil mining iteration has finished")
		w.moveBarrelsToStorageIfNeeded()
	}()

	// worker finished to oil - time to check captcha
	if w.captchaChecker.IsCaptchaAppeared(w.pid) {
		log.Info("captcha appeared: solving...")
		w.captchaNotAppearedTimes = 0
		if err := w.captchaSolver.Solve(); err != nil {
			log.WithError(err).Error("failed to solve captcha")

			return false, captchaSolveErr
		}

		return true, nil
	}

	// 5 iterations
	if w.captchaNotAppearedTimes > 100 {
		os.Exit(2)
		return false, captchaNotAppearedTooManyTimesErr
	}

	w.captchaNotAppearedTimes++
	log.Debugf("Captcha not appeared times: %d", w.captchaNotAppearedTimes)

	return false, nil
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

	w.oilManipulator.pressEsc()
	w.storageManipulator.ReplaceItemFromInventoryToStorage()
	w.oilManipulator.pressE()

	<-time.After(time.Second)
}
