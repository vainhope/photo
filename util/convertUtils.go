package util

import "math"

func Wrap(num float64, retain int) int64 {
	return int64(num * math.Pow10(retain))
}
