package captcha

import (
	"bytes"
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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
		log.WithError(err).Fatalf("failed to take screenshot")
	}

	predictionId, err := s.processor.ProcessAndSave(screenshot)
	if err != nil {
		log.WithError(err).Fatalf("failed to process screenshot")
	}

	answerNum, err := s.client.recognizeAndSolve(predictionId)
	if err != nil {
		//if err := s.processor.CleanUp(predictionId); err != nil {
		//	log.WithError(err).WithField("prediction_id", predictionId).Errorf("failed to clean up prediction")
		//}

		return errors.Wrap(err, "failed to recognize captcha images")
	}

	if err := s.manipulator.Answer(answerNum); err != nil {
		return errors.Wrap(err, "failed to answer to captcha with manipulator")
	}

	log.Info("captcha appeared and solved!")

	if err := s.processor.CleanUp(predictionId); err != nil {
		log.WithError(err).WithField("prediction_id", predictionId).Errorf("failed to clean up prediction")
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
