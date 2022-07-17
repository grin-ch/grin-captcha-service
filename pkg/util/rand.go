package util

import (
	"math/rand"
	"time"
)

var (
	lower = []byte{
		'a', 'b', 'c', 'd', 'e', 'f', 'g',
		'h', 'i', 'j', 'k', 'l', 'm', 'n',
		'o', 'p', 'q', 'r', 's', 't', 'u',
		'v', 'w', 'x', 'y', 'z',
	}
	upper = []byte{
		'A', 'B', 'C', 'D', 'E', 'F', 'G',
		'H', 'I', 'J', 'K', 'L', 'M', 'N',
		'O', 'P', 'Q', 'R', 'S', 'T', 'U',
		'V', 'W', 'X', 'Y', 'Z',
	}
	nums = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
)

func NumSet() []int {
	return nums
}

func LowerSet() []byte {
	return lower
}

func UpperSet() []byte {
	return upper
}

func GenFormSet[T any](l int, src []T) []T {
	m := len(src)
	if l <= 0 || m == 0 {
		return nil
	}
	list := make([]T, l)
	r := rand.New(rand.NewSource(time.Now().UnixMilli()))
	for i := 0; i < l; i++ {
		list[i] = src[r.Int()%m]
	}

	return list
}
