package blackscholes

import (
	"math"
	"testing"
)

// oldStdNormCDF is the previous Taylor series implementation, kept for benchmarking
func oldStdNormCDF(x float64) float64 {
	var probability float64
	if x >= 8 {
		probability = 1
	} else if x <= -8 {
		probability = 0
	} else {
		var cache [100]float64
		for i := 0; i < 100; i++ {
			n := float64(2*i + 1)
			val := 1.0
			for j := n; j > 1; j -= 2 {
				val *= j
			}
			cache[i] = val
		}
		for i := 0; i < 100; i++ {
			probability += (math.Pow(x, float64(2*i+1)) / cache[i])
		}
		probability *= math.Exp(-0.5 * math.Pow(x, 2))
		probability /= math.Sqrt(2 * math.Pi)
		probability += 0.5
	}
	return probability
}

// oldBlackScholes uses the old CDF for comparison
func oldBlackScholes(s, k, t, v, r float64, callPut string) float64 {
	w := (r*t + math.Pow(v, 2)*t/2 - math.Log(k/s)) / (v * math.Sqrt(t))
	if callPut == "CALL" {
		return s*oldStdNormCDF(w) - k*math.Exp(-r*t)*oldStdNormCDF(w-v*math.Sqrt(t))
	}
	return k*math.Exp(-r*t)*oldStdNormCDF(v*math.Sqrt(t)-w) - s*oldStdNormCDF(-w)
}

func oldGetImpliedVolatility(expectedCost, s, k, t, r float64, callPut string, estimate float64) float64 {
	if estimate == 0 {
		estimate = 0.1
	}
	low := 0.0
	high := math.Inf(1)
	for i := 0; i < 100; i++ {
		actualCost := oldBlackScholes(s, k, t, estimate, r, callPut)
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

// BenchmarkDeltaSearch simulates FindByDelta's workload: IV solve + delta calc for N strikes
// This mirrors the production hot path when entering a trade by delta
func BenchmarkDeltaSearch(b *testing.B) {
	// Simulate SPX option chain: 200 strikes, 5 points apart
	const numStrikes = 200
	const underlying = 6000.0
	const timeToExpiry = 0.003 // ~1 DTE
	const riskFreeRate = 0.045
	const baseVol = 0.15

	type strike struct {
		price   float64
		premium float64
		vol     float64
	}

	strikes := make([]strike, numStrikes)
	for i := range strikes {
		k := underlying - 500 + float64(i)*5
		vol := baseVol + 0.001*float64(i) // slight vol skew
		premium := BlackScholes(underlying, k, timeToExpiry, vol, riskFreeRate, "PUT")
		strikes[i] = strike{price: k, premium: premium, vol: vol}
	}

	b.Run("old_taylor_series", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			for _, s := range strikes {
				iv := oldGetImpliedVolatility(s.premium, underlying, s.price, timeToExpiry, riskFreeRate, "PUT", s.vol)
				oldStdNormCDF(GetW(underlying, s.price, timeToExpiry, iv, riskFreeRate))
			}
		}
	})

	b.Run("new_erfc", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			for _, s := range strikes {
				iv, _ := GetImpliedVolatility(s.premium, underlying, s.price, timeToExpiry, riskFreeRate, "PUT")
				StdNormCDF(GetW(underlying, s.price, timeToExpiry, iv, riskFreeRate))
			}
		}
	})
}
