package captcha

import (
	"bytes"
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/pkg/errors"
	"golang.org/x/image/bmp"
	"image"
)

func NewSolver(pid int32, client *Recognizer, processor *ScreenshotProcessor) *Solver {
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

	client    *Recognizer
	processor *ScreenshotProcessor
}

func (s *Solver) Start(runCheckCh <-chan struct{}) <-chan int {
	go s.start(runCheckCh)

	return s.answerCh
}

func (s *Solver) start(runCheckCh <-chan struct{}) {
	for {
		select {
		case <-runCheckCh:
			fmt.Println("[*] Debug: Checking the captcha")
			answerNum, err := s.solve()
			if err == NoCaptchaAppearedErr {
				continue
			} else if err != nil {
				fmt.Println(fmt.Sprintf("failed to recognize captcha: %s. Continue", err))
				continue
			}

			s.answerCh <- answerNum
		case <-s.stopCh:
			return
		}
	}
}

func (s *Solver) solve() (int, error) {
	screenshot, err := s.takeScreenShot(s.pid)
	if err != nil {
		panic(fmt.Sprintf("failed to take screenshot: %s", err))
	}

	predictionId, err := s.processor.ProcessAndSave(screenshot)
	if err != nil {
		panic(fmt.Sprintf("failed to process screenshot: %s", err))
	}

	answerNum, err := s.client.recognizeAndSolve(predictionId)
	if err == NoCaptchaAppearedErr {
		fmt.Println("[*] Debug: no captcha appeared. Continue")

		if err := s.processor.CleanUp(predictionId); err != nil {
			fmt.Println(fmt.Sprintf("failed to clean up prediction [%d]: %s", predictionId, err))
		}

		return 0, NoCaptchaAppearedErr
	} else if err != nil {
		return 0, errors.Wrap(err, "failed to recognize captcha images")
	}

	fmt.Println("[*] Debug: captcha appeared and solved!")

	if err := s.processor.CleanUp(predictionId); err != nil {
		fmt.Println(fmt.Sprintf("failed to clean up prediction [%d]: %s", predictionId, err))
	}

	return answerNum, nil
}

func (s *Solver) takeScreenShot(pid int32) (image.Image, error) {
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

func (s *Solver) Stop() {
	s.stopCh <- struct{}{}
}
