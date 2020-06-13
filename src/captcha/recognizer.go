package captcha

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

var (
	StatusCodeIsNotOkErr = errors.New("status code != 2**")
	AnswerValidationErr  = errors.New("answer validation error: answer value should be between 1 and 3")
	NoCaptchaAppearedErr = errors.New("no captcha appeared")
)

type recognitionResp struct {
	Answer int    `json:"answer"`
	Error  string `json:"error"`
}

func NewRecognizerClient() *RecognizerClient {
	return &RecognizerClient {}
}

type RecognizerClient struct {
	*http.Client
}

// returns correct answer number: from 1 to 3
func (c *RecognizerClient) recognizeAndSolve() (int, error) {
	resp, err := c.Get("localhost:8118/recognize-and-solve")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode == http.StatusNoContent {
		return 0, NoCaptchaAppearedErr
	}

	if resp.StatusCode > 299 {
		fmt.Println(fmt.Sprintf("Response body: %s", respBody))

		return 0, StatusCodeIsNotOkErr
	}

	recognitionResp := recognitionResp{}

	if err := json.Unmarshal(respBody, &recognitionResp); err != nil {
		return 0, err
	}

	if recognitionResp.Answer < 1 || recognitionResp.Answer > 3 {
		return 0, AnswerValidationErr
	}

	return recognitionResp.Answer, nil
}
