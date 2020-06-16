package captcha

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

var (
	AnswerValidationErr  = errors.New("answer validation error: answer value should be between 1 and 3")
	regexpValidationErr  = errors.New("question regexp validation failed")
)

type recognitionType string

const (
	question recognitionType = "question"
	answer   recognitionType = "answer"
)

func NewRecognizer() *Recognizer {
	return &Recognizer{}
}

type Recognizer struct {
}

// returns correct answer number: from 1 to 3
func (r *Recognizer) recognizeAndSolve(predictionId int64) (int, error) {
	t1, t2, err := r.parserQuestion(predictionId)
	if err != nil {
		return 0, err
	}

	calculatedAnswer := t1 + t2

	answers, err := r.parseAnswers(predictionId)
	if err != nil {
		return 0, err
	}

	for i, answer := range answers {
		if answer == calculatedAnswer {
			return i + 1, nil
		}
	}

	return 0, errors.New("no correct answer found")
}

func (r *Recognizer) parserQuestion(predictionId int64) (int, int, error) {
	path := fmt.Sprintf("resources/predictions/%d/question", predictionId)

	if err := r.generatePredictions(path, question, 1); err != nil {
		return 0, 0, err
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return 0, 0, err
	}

	var t1, t2 int
	for _, f := range files {
		if !strings.Contains(f.Name(), ".pred.txt") {
			continue
		}

		question, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", path, f.Name()))
		if err != nil {
			return 0, 0, err
		}

		r, err := regexp.Compile(`[0-9]\+[0-9]`)
		if err != nil {
			return 0, 0, err
		}

		if !r.MatchString(string(question)) {
			return 0, 0, regexpValidationErr
		}

		runes := []rune(string(question))
		q1, q2 := string(runes[0]), string(runes[2])

		t1, err = strconv.Atoi(q1)
		if err != nil {
			return 0, 0, err
		}
		t2, err = strconv.Atoi(q2)
		if err != nil {
			return 0, 0, err
		}
	}

	return t1, t2, nil
}

func (r *Recognizer) parseAnswers(predictionId int64) ([]int, error) {
	path := fmt.Sprintf("resources/predictions/%d/answers", predictionId)

	if err := r.generatePredictions(path, answer, 3); err != nil {
		return nil, err
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var answers = make([]int, 0, 3)
	for _, f := range files {
		if !strings.Contains(f.Name(), ".pred.txt") {
			continue
		}

		dat, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", path, f.Name()))
		if err != nil {
			return nil, err
		}

		answer, err := strconv.Atoi(string(dat))
		if err != nil {
			return nil, err
		}

		answers = append(answers, answer)
	}

	return answers, nil
}

func (r *Recognizer) generatePredictions(path string, t recognitionType, procNum int) error {
	exePath := "calamari-predict"

	modelPath, err := filepath.Abs(fmt.Sprintf("resources/ml_models_new/%s/model_last.ckpt", t))
	if err != nil {
		return err
	}

	filesToPredictPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	filesToPredictPath = fmt.Sprintf("%s/*.jpg", filesToPredictPath)

	cmd := exec.Command(
		exePath,
		"--checkpoint", modelPath,
		"--files", filesToPredictPath,
		"--processes", strconv.Itoa(procNum),
	)

	var out bytes.Buffer
	cmd.Stderr = &out

	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, out.String())
	}

	return nil
}
