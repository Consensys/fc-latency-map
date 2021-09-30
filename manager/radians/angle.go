package radians

import "math"

const degrees180 = 180

// Radians convert degrees into radians
func Radians(degree float64) float64 {
	return degree * (math.Pi / degrees180)
}
