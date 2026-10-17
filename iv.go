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

const (
	minVol = 5e-4
	maxVol = 10.0
	// ivIterations halves a 10 wide bracket to ~9e-18, past float64 resolution.
	ivIterations = 60
	// vegaFloor is the dollars of price per 1.00 of vol below which a price
	// cannot identify a vol (living SPX rows measure vega in the tens).
	vegaFloor = 0.05
)

// SolveImpliedVolatility inverts an option price to a volatility (decimal, 0.20 == 20%)
// by bisection over [minVol, maxVol]; price is monotone in vol. ok is false when the
// price cannot identify a vol: a non-positive price, a price outside the band the price
// function can reach, or vega at the solution below vegaFloor (deep ITM near expiry,
// where any vol reproduces the price).
//
// s, k and t must be finite and positive; callers hold a live quote and a time to expiry.
// Outside that, ok is false but iv is not meaningful.
//
// A caller that needs the IV itself (skew, an IV ratio) must skip the row when !ok. A
// caller that needs only delta may use iv as returned: for a real quote with no time
// value, GetDelta at iv 0 gives +/-1 in the money and 0 out of the money.
func SolveImpliedVolatility(price, s, k, t, r float64, callPut string) (iv float64, ok bool) {
	if price <= 0 || (callPut != "CALL" && callPut != "PUT") {
		return 0, false
	}
	lo, hi := minVol, maxVol
	// Written as the positive test so a NaN price function (negative spot or strike) is rejected too.
	if !(price >= BlackScholes(s, k, t, lo, r, callPut) && price <= BlackScholes(s, k, t, hi, r, callPut)) {
		return 0, false
	}
	for i := 0; i < ivIterations; i++ {
		mid := (lo + hi) / 2
		if BlackScholes(s, k, t, mid, r, callPut) < price {
			lo = mid
		} else {
			hi = mid
		}
	}
	iv = (lo + hi) / 2
	return iv, GetVega(s, k, t, iv, r)*100 >= vegaFloor
}
