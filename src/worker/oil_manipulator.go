package worker

import (
	"context"
	"fmt"
	"github.com/go-vgo/robotgo"
	log "github.com/sirupsen/logrus"
	"rp-bot-client/src/window"
	"time"
)

const oilDoneHexColorFirst = "c1c1c1"
const oilDoneHexColorSecond = "ffffff"

type coordinates struct {
	x int
	y int
}

type OilManipulator struct {
	pid            int32
	oilCoordinates []coordinates
}

func newOilManipulator(pid int32) *OilManipulator {
	return &OilManipulator{
		pid: pid,
		oilCoordinates: []coordinates{
			{283, 675},
			{388, 672},
			{494, 673},
			{595, 678},
		},
	}
}

func (m *OilManipulator) holdOil() coordinates {
	var oilToHoldCoordinates coordinates

	err := window.ActivatePidAndRun(m.pid, func() error {
		<-time.After(100 * time.Millisecond)

		for _, coordinates := range m.oilCoordinates {
			if !m.oilHasDoneColor(robotgo.GetPixelColor(coordinates.x, coordinates.y)) {
				oilToHoldCoordinates = coordinates
				break
			}
		}

		robotgo.Move(oilToHoldCoordinates.x, oilToHoldCoordinates.y)

		<-time.After(100 * time.Millisecond)

		robotgo.MouseToggle("down")

		log.Debug(fmt.Sprintf("oil mouse key down on coordinates %v", oilToHoldCoordinates))

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return oilToHoldCoordinates
}

func (m *OilManipulator) releaseOilOnDone(heldOilCoordinates coordinates, ctx context.Context) bool {
	if m.getDoneOilsCount() == 3 {
		m.releaseLastOilOnDone(ctx)
		return true
	}

	log.Debug("Waiting to release oil")

	err := window.ActivatePidAndRun(m.pid, func() error {
		for {
			select {
			case <-ctx.Done():
				m.releaseOil()
				log.Debug(fmt.Sprintf("oil with coordinates: %v has been released by timeout", heldOilCoordinates))
				return nil
			default:
				color := robotgo.GetPixelColor(heldOilCoordinates.x, heldOilCoordinates.y)
				log.Debug(fmt.Sprintf("current oil: coordinates: %v color: %s", heldOilCoordinates, color))
				if m.oilHasDoneColor(color) {
					m.releaseOil()
					return nil
				}
			}

			<-time.After(300 * time.Millisecond)
		}
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Debug("Oil released")

	return false
}

func (m *OilManipulator) releaseLastOilOnDone(ctx context.Context) {
	log.Debug("Waiting to release last oil")

	err := window.ActivatePidAndRun(m.pid, func() error {
		for {
			select {
			case <-ctx.Done():
				log.Debug("last oil has been released by timeout")
				m.releaseOil()
				return nil
			default:
				allFinished := true
				for _, coordinates := range m.oilCoordinates {
					newColor := robotgo.GetPixelColor(coordinates.x, coordinates.y)
					log.Debug(fmt.Sprintf("Oil with coordinates %d:%d has color %s", coordinates.x, coordinates.y, newColor))
					if m.oilHasDoneColor(newColor) {
						allFinished = false
					}
				}
				if allFinished {
					m.releaseOil()
					return nil
				}
			}

			<-time.After(300 * time.Millisecond)
		}
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Debug("last oil released")
}

func (m *OilManipulator) releaseOil() {
	err := window.ActivatePidAndRun(m.pid, func() error {
		robotgo.MouseToggle("up")

		<-time.After(100 * time.Millisecond)

		log.Debug("oil mouse key up")

		return nil
	})

	if err != nil {
		log.Debug(err)
	}
}

func (m *OilManipulator) ReOpenWindow() {
	m.pressEsc()

	<-time.After(time.Second)

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
		log.Fatal(err)
	}
}

func (m *OilManipulator) getDoneOilsCount() int {
	doneOilsCount := 0

	err := window.ActivatePidAndRun(m.pid, func() error {
		for _, coordinates := range m.oilCoordinates {
			if m.oilHasDoneColor(robotgo.GetPixelColor(coordinates.x, coordinates.y)) {
				doneOilsCount++
			}
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return doneOilsCount
}

func (m *OilManipulator) oilHasDoneColor(color string) bool {
	return color == oilDoneHexColorFirst || color == oilDoneHexColorSecond
}
