package util

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestNewImg(t *testing.T) {
	src, err := NewImg("hello world!")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	str := base64.StdEncoding.EncodeToString(src)
	fmt.Println(str)
}
