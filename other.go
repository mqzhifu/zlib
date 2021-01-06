package zlib

import (
	"math/rand"
	"strings"
	"time"
)

func GetRandIntNum(max int) int{
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max)
}

func CheckStrEmpty(str string)bool{
	if str == ""{
		return true
	}
	str = strings.Trim(str," ")
	if str == ""{
		return true
	}
	return false
}