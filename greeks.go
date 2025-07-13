package blackscholes

import (
	"fmt"
	"math"
)

// Black-Scholes delta calculation for a European call and put option
// S := 100.0   // Current stock price
// K := 100.0   // Strike price
// T := 1.0     // Time to maturity (in years)
// r := 0.05    // Risk-free interest rate
// sigma := 0.2 // Volatility

// StdNormDensity calculates the standard normal density function
func StdNormDensity(x float64) float64 {
	return math.Exp(-math.Pow(x, 2)/2) / math.Sqrt(2*math.Pi)
}

// GetDelta calculates the delta of an option
func GetDelta(s, k, t, v, r float64, callPut string) (float64, error) {
	if callPut == "CALL" {
		return callDelta(s, k, t, v, r), nil
	} else if callPut == "PUT" {
		return putDelta(s, k, t, v, r), nil
	}
	return 0, fmt.Errorf("callput is not of type CALL or PUT %s", callPut)
}

// CallDelta calculates the delta of a call option
func callDelta(s, k, t, v, r float64) float64 {
	w := GetW(s, k, t, v, r)
	var delta float64
	if !isFinite(w) {
		if s > k {
			delta = 1
		} else {
			delta = 0
		}
	} else {
		delta = StdNormCDF(w)
	}
	return delta
}

// PutDelta calculates the delta of a put option
func putDelta(s, k, t, v, r float64) float64 {
	delta := callDelta(s, k, t, v, r) - 1
	if delta == -1 && k == s {
		return 0
	}
	return delta
}

// GetRho calculates the rho of an option
func GetRho(s, k, t, v, r float64, callPut string, scale int) float64 {
	if scale == 0 {
		scale = 100
	}
	if callPut == "call" {
		return callRho(s, k, t, v, r) / float64(scale)
	}
	return putRho(s, k, t, v, r) / float64(scale)
}

// CallRho calculates the rho of a call option
func callRho(s, k, t, v, r float64) float64 {
	w := GetW(s, k, t, v, r)
	if !math.IsNaN(w) {
		return k * t * math.Exp(-r*t) * StdNormCDF(w-v*math.Sqrt(t))
	}
	return 0
}

// PutRho calculates the rho of a put option
func putRho(s, k, t, v, r float64) float64 {
	w := GetW(s, k, t, v, r)
	if !math.IsNaN(w) {
		return -k * t * math.Exp(-r*t) * StdNormCDF(v*math.Sqrt(t)-w)
	}
	return 0
}

// GetVega calculates the vega of a call and put option
func GetVega(s, k, t, v, r float64) float64 {
	w := GetW(s, k, t, v, r)
	if isFinite(w) {
		return s * math.Sqrt(t) * StdNormDensity(w) / 100
	}
	return 0
}

// GetTheta calculates the theta of an option
func GetTheta(s, k, t, v, r float64, callPut string, scale int) float64 {
	if scale == 0 {
		scale = 365
	}
	if callPut == "call" {
		return callTheta(s, k, t, v, r) / float64(scale)
	}
	return putTheta(s, k, t, v, r) / float64(scale)
}

// CallTheta calculates the theta of a call option
func callTheta(s, k, t, v, r float64) float64 {
	w := GetW(s, k, t, v, r)
	if isFinite(w) {
		return -v*s*StdNormDensity(w)/(2*math.Sqrt(t)) - k*r*math.Exp(-r*t)*StdNormCDF(w-v*math.Sqrt(t))
	}
	return 0
}

// PutTheta calculates the theta of a put option
func putTheta(s, k, t, v, r float64) float64 {
	w := GetW(s, k, t, v, r)
	if isFinite(w) {
		return -v*s*StdNormDensity(w)/(2*math.Sqrt(t)) + k*r*math.Exp(-r*t)*StdNormCDF(v*math.Sqrt(t)-w)
	}
	return 0
}

// GetGamma calculates the gamma of a call and put option
func GetGamma(s, k, t, v, r float64) float64 {
	w := GetW(s, k, t, v, r)
	if isFinite(w) {
		return StdNormDensity(w) / (s * v * math.Sqrt(t))
	}
	return 0
}
