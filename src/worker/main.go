package worker

import (
	"github.com/pkg/errors"
	"rp-bot-client/src/captcha"
)

type Worker interface {
	Start()
	Restart()
	Resume()
	Interrupt()
	ToggleHoldTime()
}

func GetWorker(pid int32, botType, btn string, checker *captcha.Checker, solver *captcha.Solver) (Worker, error) {
	switch botType {
	case "oil":
		return NewOilMan(pid, checker, solver), nil
	case "mine":
		return NewMiner(btn, pid), nil
	}

	return nil, errors.Errorf("failed to get worker by bot type. Value given: %s", botType)
}