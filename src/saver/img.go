package saver

import (
	"fmt"
	pnm "github.com/jbuchbinder/gopnm"
	"image"
	"os"
	"path/filepath"
	"strings"
)

const predictionsBasePathTpl = "resources/predictions"

func NewImg() *Img {
	return &Img{}
}

type Img struct{}

func (s *Img) mkdir(path string) {
	dir := filepath.Dir(path)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			panic(err)
		}
	}
}

func (s *Img) saveToFile(img image.Image, path string) (string, error) {
	s.mkdir(path)

	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// return path, png.Encode(f, img)

	return path, pnm.Encode(f, img, 0)
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
	path := fmt.Sprintf("%s/%d/answers/answer_%d.pbm", predictionsBasePathTpl, predictionId, answerNum)

	return s.saveToFile(img, path)
}

func (s *Img) SaveQuestion(img image.Image, predictionId int64) (string, error) {
	path := fmt.Sprintf("%s/%d/question/question.pbm", predictionsBasePathTpl, predictionId)

	return s.saveToFile(img, path)
}

func (s *Img) SaveScreenshot(img image.Image, name string) (string, error) {
	path := fmt.Sprintf("resources/screenshots/%s.jpg", name)

	return s.saveToFile(img, path)
}

func (s *Img) CleanUpPredictionFiles(predictionId int64) error {
	return os.RemoveAll(fmt.Sprintf("%s/%d", predictionsBasePathTpl, predictionId))
}
