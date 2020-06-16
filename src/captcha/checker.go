package captcha

import (
	"github.com/go-vgo/robotgo"
	"rp-bot-client/src/window"
)

const answerBtnColorHex = "ffcc30"

type pixel struct {
	x int
	y int
}

type Checker struct {
	answerBtnPixels []pixel
}

func NewChecker() *Checker {
	return &Checker{answerBtnPixels: []pixel{{840,688}, {940,688}, {844,709}, {937,709}}}
}

func (c *Checker) IsCaptchaAppeared(pid int32) bool {
	for _, p := range c.answerBtnPixels {
		if !c.isPixelColorEqualToAnswerBtnColor(pid, p) {
			return false
		}
	}

	return true
}

func (c *Checker) isPixelColorEqualToAnswerBtnColor(pid int32, p pixel) bool {
	captchaAppeared := true

	err := window.ActivatePidAndRun(pid, func() error {
		captchaColor := robotgo.GetPixelColor(p.x, p.y)

		if captchaColor != answerBtnColorHex {
			captchaAppeared = false
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	return captchaAppeared
}