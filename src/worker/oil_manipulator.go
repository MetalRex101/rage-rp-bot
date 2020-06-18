package worker

import (
	"github.com/go-vgo/robotgo"
	log "github.com/sirupsen/logrus"
	"rp-bot-client/src/window"
	"time"
)

const oilDoneHexColor = ""

type coordinates struct {
	x int
	y int
}

type OilManipulator struct {
	pid            int32
	oilCoordinates map[int]coordinates
}

func newOilManipulator(pid int32) *OilManipulator {
	return &OilManipulator{
		pid: pid,
		oilCoordinates: map[int]coordinates{
			0: {283, 675},
			1: {388, 672},
			2: {494, 673},
			3: {595, 678},
		},
	}
}

func (m *OilManipulator) holdOil(currentOil int) {
	err := window.ActivatePidAndRun(m.pid, func() error {
		<-time.After(100 * time.Millisecond)

		coordinates := m.oilCoordinates[currentOil]
		robotgo.Move(coordinates.x, coordinates.y)

		<-time.After(100 * time.Millisecond)

		robotgo.MouseToggle("down")

		log.Debug("oil mouse key down")

		return nil
	})

	if err != nil {
		panic(err)
	}
}

func (m *OilManipulator) releaseOilOnDone(currentOil int) {
	coordinates := m.oilCoordinates[currentOil]

	for {
		color := robotgo.GetPixelColor(coordinates.x, coordinates.y)
		if color == oilDoneHexColor {
			m.releaseOil()
			return
		}

		<-time.After(300 * time.Millisecond)
	}
}

func (m *OilManipulator) releaseOil() {
	err := window.ActivatePidAndRun(m.pid, func() error {
		robotgo.MouseToggle("up")

		<-time.After(100 * time.Millisecond)

		log.Debug("oil mouse key up")

		return nil
	})

	if err != nil {
		panic(err)
	}
}

func (m *OilManipulator) ReOpenWindow() {
	m.pressEsc()

	<-time.After(300 * time.Millisecond)

	m.pressE()

	// todo move to config to fit into client system requirements, make 1 sec as default
	<-time.After(2 * time.Second)
}

func (m *OilManipulator) pressEsc() {
	err := window.ActivatePidAndRun(m.pid, func() error {
		robotgo.KeyTap("esc")
		log.Debug("esc key tap")

		return nil
	})

	if err != nil {
		panic(err)
	}
}

func (m *OilManipulator) pressE() {
	err := window.ActivatePidAndRun(m.pid, func() error {
		robotgo.KeyTap("e")
		log.Debug("e key tap")

		return nil
	})

	if err != nil {
		panic(err)
	}
}
