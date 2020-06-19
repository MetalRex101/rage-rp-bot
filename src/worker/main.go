package worker

import (
	"github.com/pkg/errors"
	"rp-bot-client/src/captcha"
	"rp-bot-client/src/storage"
)

type Worker interface {
	Start()
	Restart()
	Resume()
	Interrupt()
}

func GetWorker(
	pid int32,
	botType,
	btn string,
	checker *captcha.Checker,
	solver *captcha.Solver,
	manipulator *storage.Manipulator,
	withStorage bool,
) (Worker, error) {
	switch botType {
	case "oil":
		return NewOilMan(pid, checker, solver, manipulator, withStorage), nil
	case "mine":
		return NewMiner(btn, pid), nil
	}

	return nil, errors.Errorf("failed to get worker by bot type. Value given: %s", botType)
}
