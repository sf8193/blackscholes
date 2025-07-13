package blackscholes

import (
	"math"
)

func isFinite(f float64) bool {
	return !math.IsInf(f, 0) && !math.IsNaN(f)
}

// StdNormCDF calculates the standard normal cumulative distribution function
func StdNormCDF(x float64) float64 {
	var probability float64
	if x >= 8 {
		probability = 1
	} else if x <= -8 {
		probability = 0
	} else {
		for i := 0; i < 100; i++ {
			probability += (math.Pow(x, float64(2*i+1)) / doubleFactorial(float64(2*i+1)))
		}
		probability *= math.Exp(-0.5 * math.Pow(x, 2))
		probability /= math.Sqrt(2 * math.Pi)
		probability += 0.5
	}
	return probability
}

// DoubleFactorial calculates the double factorial of n
func doubleFactorial(n float64) float64 {
	val := 1.0
	for i := n; i > 1; i -= 2 {
		val *= i
	}
	return val
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
