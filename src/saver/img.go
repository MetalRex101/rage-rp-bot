package saver

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"strings"
)

const predictionsBasePathTpl = "resources/predictions"

func NewImg() *Img {
	return &Img{}
}

type Img struct{}

func (s *Img) saveToFile(img image.Image, path string) (string, error) {
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// return path, png.Encode(f, img)
	return path, jpeg.Encode(f, img, &jpeg.Options{Quality: 100})
}

func (s *Img) SaveAnswerToDataset(img image.Image, datasetNum, answerNum int, screenName string) (string, error) {
	screenName = strings.ReplaceAll(screenName, ".png", "")
	path := fmt.Sprintf("resources/dataset_%d/answers/%s_%d.jpg", datasetNum, screenName, answerNum)

	path, err := s.saveToFile(img, path)
	if err != nil {
		return "", err
	}

	gtPath := strings.ReplaceAll(path, ".jpg", ".gt.txt")

	_, err = os.Create(gtPath)
	if err != nil {
		return "", err
	}

	return path, nil
}

func (s *Img) SaveAnswer(img image.Image, answerNum int, predictionId int64) (string, error) {
	path := fmt.Sprintf("%s/%d/answer_%d.jpg", predictionsBasePathTpl, predictionId, answerNum)

	return s.saveToFile(img, path)
}

func (s *Img) SaveQuestion(img image.Image, predictionId int64) (string, error) {
	path := fmt.Sprintf("%s/%d/question.jpg", predictionsBasePathTpl, predictionId)

	return s.saveToFile(img, path)
}

func (s *Img) SaveScreenshot(img image.Image, name string) (string, error) {
	path := fmt.Sprintf("resources/screenshots/%s.jpg", name)

	return s.saveToFile(img, path)
}
