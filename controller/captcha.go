package controller

import (
	"errors"

	"github.com/dchest/captcha"
)

var ErrInvalidCaptchaSolution = errors.New("invalid captcha solution")

type Captcha struct {
	ID string `json:"captchaId"`
}

func NewCaptcha() *Captcha {
	return &Captcha{ID: captcha.New()}
}
