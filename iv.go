package blackscholes

import (
	"math"
)

// GetImpliedVolatility calculates a close estimate of implied volatility given an option price
func GetImpliedVolatility(expectedCost, s, k, t, r float64, callPut string, estimate float64) float64 {
	if estimate == 0 {
		estimate = 0.1
	}
	low := 0.0
	high := math.Inf(1)
	// perform 100 iterations max
	for i := 0; i < 100; i++ {
		actualCost := BlackScholes(s, k, t, estimate, r, callPut)
		if int(expectedCost*100) == int(actualCost*100) {
			break
		} else if actualCost > expectedCost {
			high = estimate
			estimate = (estimate-low)/2 + low
		} else {
			low = estimate
			estimate = (high-estimate)/2 + estimate
			if !isFinite(estimate) {
				estimate = low * 2
			}
		}
	}
	return estimate
}
