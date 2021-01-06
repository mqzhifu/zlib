package zlib

import (
	"math/rand"
	"time"
)

func getRandIntNum(max int) int{
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max)
}
