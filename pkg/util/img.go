package util

import (
	"bytes"
	"image/png"
	"io"

	"github.com/afocus/captcha"
)

var (
	cap *captcha.Captcha
)

const (
	width  = 120
	height = 40
)

func NewImg(text string) ([]byte, error) {
	return newImg(text)
}

func newImg(text string) ([]byte, error) {
	img := cap.CreateCustom(text)
	var b bytes.Buffer
	if err := png.Encode(io.Writer(&b), img); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func init() {
	cap = captcha.New()
	cap.SetSize(width, height)
	_ = cap.AddFontFromBytes(COMICSAN)
}
