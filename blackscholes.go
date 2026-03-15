package blackscholes

import (
	"math"
)

func isFinite(f float64) bool {
	return !math.IsInf(f, 0) && !math.IsNaN(f)
}

// StdNormCDF calculates the standard normal cumulative distribution function
func StdNormCDF(x float64) float64 {
	return 0.5 * math.Erfc(-x/math.Sqrt2)
}

// BlackScholes calculates the Black-Scholes option pricing formula
func BlackScholes(s, k, t, v, r float64, callPut string) float64 {
	var price float64
	w := (r*t + math.Pow(v, 2)*t/2 - math.Log(k/s)) / (v * math.Sqrt(t))
	if callPut == "CALL" {
		price = s*StdNormCDF(w) - k*math.Exp(-r*t)*StdNormCDF(w-v*math.Sqrt(t))
	} else {
		price = k*math.Exp(-r*t)*StdNormCDF(v*math.Sqrt(t)-w) - s*StdNormCDF(-w)
	}
	return price
}

// GetW calculates omega as defined in the Black-Scholes formula
func GetW(s, k, t, v, r float64) float64 {
	return (r*t + math.Pow(v, 2)*t/2 - math.Log(k/s)) / (v * math.Sqrt(t))
}
