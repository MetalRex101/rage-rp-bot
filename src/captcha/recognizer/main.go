package recognizer

import "github.com/pkg/errors"

var (
	AnswerValidationErr   = errors.New("answer validation error: answer value should be between 1 and 3")
	QuestionValidationErr = errors.New("question regexp validation failed")
)

type Recognizer interface {
	RecognizeAndSolve(predictionId int64) (int, error)
	GetEngine() Engine
}

type Engine string

const (
	Calamari Engine = "calamari"
	GOCR Engine = "gocr"
)

func Get(engine Engine) (Recognizer, error) {
	switch engine {
	case Calamari:
		return NewCalamari(), nil
	case GOCR:
		return NewGOCR(), nil
	default:
		return nil, errors.Errorf("recognizer engine '%s' not found", engine)
	}
}