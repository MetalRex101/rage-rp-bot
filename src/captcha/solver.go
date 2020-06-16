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
		pid:         pid,
		client:      client,
		processor:   processor,
		manipulator: NewMouseManipulator(pid),
	}
}

type Solver struct {
	pid      int32
	stopCh   chan struct{}
	answerCh chan int

	client      *Recognizer
	manipulator *MouseManipulator
	processor   *ScreenshotProcessor
}

func (s *Solver) Solve() error {
	screenshot, err := s.takeScreenShot(s.pid)
	if err != nil {
		panic(fmt.Sprintf("failed to take screenshot: %s", err))
	}

	predictionId, err := s.processor.ProcessAndSave(screenshot)
	if err != nil {
		panic(fmt.Sprintf("failed to process screenshot: %s", err))
	}

	answerNum, err := s.client.recognizeAndSolve(predictionId)
	if err != nil {
		//if err := s.processor.CleanUp(predictionId); err != nil {
		//	fmt.Println(fmt.Sprintf("failed to clean up prediction [%d]: %s", predictionId, err))
		//}

		return errors.Wrap(err, "failed to recognize captcha images")
	}

	if err := s.manipulator.Answer(answerNum); err != nil {
		return errors.Wrap(err, "failed to answer to captcha with manipulator")
	}

	fmt.Println("[*] Debug: captcha appeared and solved!")

	if err := s.processor.CleanUp(predictionId); err != nil {
		fmt.Println(fmt.Sprintf("failed to clean up prediction [%d]: %s", predictionId, err))
	}

	return nil
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
