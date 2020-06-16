package storage

import (
	"github.com/go-vgo/robotgo"
	"rp-bot-client/src/window"
	"time"
)

type pixel struct {
	x int
	y int
}

type Manipulator struct {
	pid int32
	inventoryFirstSlotPixel pixel
	storageFirstSlotPixel pixel
}

func NewManipulator(pid int32) *Manipulator {
	return &Manipulator{
		pid: pid,
		inventoryFirstSlotPixel: pixel{x: 435, y: 334},
		storageFirstSlotPixel: pixel{x: 1313, y: 365},
	}
}

func (m *Manipulator) ReplaceItemFromInventoryToStorage() {
	m.openInventory()
	m.replaceItem()
	m.closeInventory()
}

func (m *Manipulator) openInventory() {
	err := window.ActivatePidAndRun(m.pid, func() error {
		robotgo.KeyTap("i")
		time.Sleep(500 * time.Millisecond)

		return nil
	})

	if err != nil {
		panic(err)
	}
}

func (m *Manipulator) replaceItem()  {
	err := window.ActivatePidAndRun(m.pid, func() error {
		robotgo.Move(m.inventoryFirstSlotPixel.x, m.inventoryFirstSlotPixel.y)
		robotgo.DragSmooth(m.storageFirstSlotPixel.x, m.storageFirstSlotPixel.y)

		time.Sleep(500 * time.Millisecond)

		return nil
	})

	if err != nil {
		panic(err)
	}
}

func (m *Manipulator) closeInventory() {
	err := window.ActivatePidAndRun(m.pid, func() error {
		robotgo.KeyTap("esc")
		time.Sleep(500 * time.Millisecond)

		return nil
	})

	if err != nil {
		panic(err)
	}
}