package captcha

import "github.com/go-vgo/robotgo"

func NewMouseManipulator(pid int32) *MouseManipulator {
	return &MouseManipulator{
		pid: pid,
	}
}

type MouseManipulator struct {
	pid int32
}

func (m *MouseManipulator) Answer(answerNum int) error {
	if err := m.selectAnswer(answerNum); err != nil {
		return err
	}

	return m.clickAnswerButton()
}

func (m *MouseManipulator) selectAnswer(answerNum int) error {
	// if we need more than 1 try to capture game window (if user clicked other window)
	for {
		if err := robotgo.ActivePID(m.pid); err != nil {
			return err
		}
		if m.pid == robotgo.GetPID() {
			x, y := m.getAnswerCoordinates(answerNum)

			robotgo.MoveSmooth(x, y)
			robotgo.Click()

			return nil
		}
	}
}

func (m *MouseManipulator) clickAnswerButton() error {
	// answer button coordinates
	const x, y = 912, 702

	// if we need more than 1 try to capture game window (if user clicked other window)
	for {
		if err := robotgo.ActivePID(m.pid); err != nil {
			return err
		}
		if m.pid == robotgo.GetPID() {
			robotgo.MoveSmooth(x, y)
			robotgo.Click()

			return nil
		}
	}
}

func (m *MouseManipulator) getAnswerCoordinates(answerNum int) (x int, y int) {
	switch answerNum {
	case 1:
		x, y = 1011, 435
	case 2:
		x, y = 1026, 490
	case 3:
		x, y = 912, 702
	default:
		x, y = 0, 0
	}

	return
}
