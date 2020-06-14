package window

import (
	"github.com/go-vgo/robotgo"
	"github.com/pkg/errors"
)

func ActivatePidAndRun(pid int32, callback func() error) error {
	// if we need more than 1 try to capture game window (if user clicked other window)
	for {
		if pid != robotgo.GetPID() {
			if err := robotgo.ActivePID(pid); err != nil {
				panic(errors.Wrap(err, "failed to activate"))
			}
		}
		if pid == robotgo.GetPID() {
			return callback()
		}
	}
}
