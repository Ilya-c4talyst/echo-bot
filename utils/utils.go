package utils

import (
	"math/rand"
	"time"
)

func Random(min, max int) int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	if min > max {
		return min
	} else {
		return rand.Intn(max-min) + min
	}
}
