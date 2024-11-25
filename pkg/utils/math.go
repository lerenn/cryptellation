package utils

import "math"

// Round rounds a float64 number to a given precision.
func Round(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	round := int(num*output + math.Copysign(0.5, num))
	return float64(round) / output
}
