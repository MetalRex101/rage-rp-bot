package window

import (
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/pkg/errors"
)

func FindGtaPid(name string) (int32, error) {
	var pid int32 = -1

	pc, err := robotgo.Process()
	if err != nil {
		return pid, err
	}

	for _, proc := range pc {
		if proc.Name == name {
			pid = proc.Pid
			fmt.Println(fmt.Sprintf("%+v", proc))
		}
	}

	if pid == -1 {
		return pid, errors.New(fmt.Sprintf("process not found, %s", name))
	}

	return pid, nil
}

