package captcha

import (
	"image"
	"rp-bot-client/src/cropper"
	"rp-bot-client/src/saver"
	"time"
)

func NewScreenshotProcessor(c *cropper.Img, s *saver.Img,) *ScreenshotProcessor {
	return &ScreenshotProcessor{c: c, s: s}
}

type ScreenshotProcessor struct {
	c *cropper.Img
	s *saver.Img
}

// extracts captcha question and answers from screenshot,
// filters images and saves it to separate dir
func (p *ScreenshotProcessor) ProcessAndSave (img image.Image) (int64, error) {
	predictionId := time.Now().UnixNano()

	captcha, err := p.c.CropCaptcha(img)
	if err != nil {
		return 0, err
	}

	question, err := p.c.CropQuestion(captcha)
	if err != nil {
		return 0, err
	}
	_, err = p.s.SaveQuestion(question, predictionId)
	if err != nil {
		return 0, err
	}

	for i:= 1; i < 4; i++ {
		answer, err := p.c.CropAnswer(captcha, i)
		if err != nil {
			return 0, err
		}

		_, err = p.s.SaveAnswer(answer, i, predictionId)
		if err != nil {
			return 0, err
		}
	}

	return predictionId, nil
}

func (p *ScreenshotProcessor) CleanUp(predictionId int64) error {
	return p.s.CleanUpPredictionFiles(predictionId)
}