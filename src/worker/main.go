package worker

import "github.com/pkg/errors"

type Worker interface {
	Start(checkCaptchaCh chan<- struct{})
	Restart()
	Resume()
	Interrupt()
}

func GetWorker(pid int32, botType string, btn string) (Worker, error) {
	switch botType {
	case "oil":
		return NewOilMan(pid), nil
	case "mine":
		return NewMiner(btn, pid), nil
	}

	return nil, errors.Errorf("failed to get worker by bot type. Value given: %s", botType)
}