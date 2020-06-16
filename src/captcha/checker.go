package captcha

import (
	"github.com/go-vgo/robotgo"
	"rp-bot-client/src/window"
)

const answerBtnColorHex = ""

type pixel struct {
	x int
	y int
}

type Checker struct {
	answerBtnPixels []pixel
}

func NewChecker() *Checker {
	return &Checker{answerBtnPixels: []pixel{{1,1}, {2,2}, {3,3}, {4,4}}}
}

func (c *Checker) isCaptchaAppeared (pid int32) bool {
	for _, p := range c.answerBtnPixels {
		if !c.isPixelColorEqualToAnswerBtnColor(pid, p) {
			return false
		}
	}

	return true
}

func (c *Checker) isPixelColorEqualToAnswerBtnColor (pid int32, p pixel) bool {
	captchaAppeared := true

	err := window.ActivatePidAndRun(pid, func() error {
		if robotgo.GetPixelColor(p.x, p.y) != answerBtnColorHex {
			captchaAppeared = false
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	return captchaAppeared
}