package recognizer

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func NewGOCR() *gocr {
	return &gocr{}
}

type gocr struct{}

// returns correct answer number: from 1 to 3
func (r *gocr) RecognizeAndSolve(predictionId int64) (int, error) {
	t1, t2, err := r.parserQuestion(predictionId)
	if err != nil {
		return 0, errors.Wrap(err, "failed to parse question")
	}

	calculatedAnswer := t1 + t2

	answers, err := r.parseAnswers(predictionId)
	if err != nil {
		return 0, errors.Wrap(err, "failed to parse answers")
	}

	for i, answer := range answers {
		if answer == calculatedAnswer {
			return i + 1, nil
		}
	}

	return 0, errors.New("no correct answer found")
}

func (r *gocr) GetEngine() Engine {
	return GOCR
}

func (r *gocr) parserQuestion(predictionId int64) (int, int, error) {
	path := fmt.Sprintf("resources/predictions/%d/question/question.pbm", predictionId)

	questionStr, err := r.predict(path)
	if err != nil {
		return 0, 0, err
	}

	questionStr = strings.ReplaceAll(questionStr, " ", "")
	questionStr = strings.TrimSpace(questionStr)

	var t1, t2 int
	regex, err := regexp.Compile(`[0-9]\+[0-9]`)
	if err != nil {
		return 0, 0, err
	}

	if !regex.MatchString(questionStr) {
		return 0, 0, QuestionValidationErr
	}

	runes := []rune(questionStr)
	q1, q2 := string(runes[0]), string(runes[2])

	t1, err = strconv.Atoi(q1)
	if err != nil {
		return 0, 0, err
	}
	t2, err = strconv.Atoi(q2)
	if err != nil {
		return 0, 0, err
	}

	return t1, t2, nil
}

func (r *gocr) parseAnswers(predictionId int64) ([]int, error) {
	path := fmt.Sprintf("resources/predictions/%d/answers", predictionId)

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var answers = make([]int, 0, 3)
	for _, f := range files {
		answerStr, err := r.predict(fmt.Sprintf("%s/%s", path, f.Name()))
		if err != nil {
			return nil, err
		}

		answerStr = strings.TrimSpace(answerStr)

		answer, err := strconv.Atoi(answerStr)
		if err != nil {
			return nil, err
		}

		answers = append(answers, answer)
	}

	return answers, nil
}

func (r *gocr) predict(path string) (string, error) {
	exePath := "gocr.exe"

	fileToPredictPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	cmd := exec.Command(
		exePath,
		"-C", "0-9+", // set characters whitelist
		"-i", fileToPredictPath,
	)

	var stdErr bytes.Buffer
	cmd.Stderr = &stdErr

	var stdOut bytes.Buffer
	cmd.Stdout = &stdOut

	err = cmd.Run()

	return stdOut.String(), errors.Wrap(err, stdErr.String())
}
