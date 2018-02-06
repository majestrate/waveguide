package util

import (
	"math/rand"
)

func RandBool(rate float64) bool {
	return rand.Float64() < rate
}
