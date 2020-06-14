package captcha

import (
	"bytes"
	"fmt"
	"github.com/go-vgo/robotgo"
	"golang.org/x/image/bmp"
	"image"
)

func NewSolver(pid int32, client *RecognizerClient, processor *ScreenshotProcessor) *Solver {
	return &Solver{
		pid:       pid,
		stopCh:    make(chan struct{}),
		answerCh:  make(chan int),
		client:    client,
		processor: processor,
	}
}

type Solver struct {
	pid      int32
	stopCh   chan struct{}
	answerCh chan int

	client    *RecognizerClient
	processor *ScreenshotProcessor
}

func (c *Solver) Start(runCheckCh <-chan struct{}) <-chan int {
	go c.start(runCheckCh)

	return c.answerCh
}

func (c *Solver) start(runCheckCh <-chan struct{}) {
	for {
		select {
		case <-runCheckCh:
			screenshot, err := c.takeScreenShot(c.pid)
			if err != nil {
				panic(fmt.Sprintf("failed to take screenshot: %s", err))
			}

			predictionId, err := c.processor.ProcessAndSave(screenshot)
			if err != nil {
				panic(fmt.Sprintf("failed to process screenshot: %s", err))
			}

			answerNum, err := c.client.recognizeAndSolve(predictionId)
			if err == NoCaptchaAppearedErr {
				fmt.Println("[*] Debug: no captcha appeared. Continue")
				continue
			}

			if err != nil {
				panic(fmt.Sprintf("failed to recognize captcha images: %s", err))
			}

			fmt.Println("[*] Debug: captcha appeared and solved!")

			c.answerCh <- answerNum
		case <-c.stopCh:
			return
		}
	}
}

func (c *Solver) takeScreenShot(pid int32) (image.Image, error) {
	// if we need more than 1 try to capture game window (if user clicked other window)
	for {
		if err := robotgo.ActivePID(pid); err != nil {
			return nil, err
		}
		if pid == robotgo.GetPID() {
			bitmap := robotgo.CaptureScreen(0, 0, 1920, 1080)

			buf := bytes.NewBuffer(robotgo.ToBitmapBytes(bitmap))
			return bmp.Decode(buf)
		}
	}
}

func (c *Solver) Stop() {
	c.stopCh <- struct{}{}
}
