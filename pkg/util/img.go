package util

import (
	"bytes"
	"io"

	"github.com/lifei6671/gocaptcha"
)

const (
	width  = 144
	height = 48
)

func NewImg(content string) ([]byte, error) {
	return NewImgWithSize(content, width, height)
}

func NewImgWithSize(content string, w, h int) ([]byte, error) {
	captchaImage := gocaptcha.NewCaptchaImage(w, h, gocaptcha.RandLightColor())
	//画边框
	captchaImage.DrawBorder(gocaptcha.ColorToRGB(0x17A7A7A))

	//噪点
	captchaImage.DrawLine(3)
	captchaImage.DrawNoise(gocaptcha.CaptchaComplexHigh)

	captchaImage.DrawText(content)

	var b bytes.Buffer
	err := captchaImage.SaveImage(io.Writer(&b), gocaptcha.ImageFormatJpeg)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
