package blackscholes

import (
	"math"
	"testing"
)

// TestRoundTripAccuracy tests that GetPrice -> GetImpliedVolatility -> GetPrice gives the same result
func TestRoundTripAccuracy(t *testing.T) {
	testCases := []struct {
		name         string
		underlying   float64
		strike       float64
		timeToExpiry float64
		volatility   float64
		riskFreeRate float64
		contractType string
	}{
		{
			name:         "ATM Call - 30 days",
			underlying:   100.0,
			strike:       100.0,
			timeToExpiry: 30.0 / 365.0,
			volatility:   0.20,
			riskFreeRate: 0.05,
			contractType: "CALL",
		},
		{
			name:         "ATM Put - 30 days",
			underlying:   100.0,
			strike:       100.0,
			timeToExpiry: 30.0 / 365.0,
			volatility:   0.20,
			riskFreeRate: 0.05,
			contractType: "PUT",
		},
		{
			name:         "Deep OTM Put",
			underlying:   100.0,
			strike:       120.0,
			timeToExpiry: 90.0 / 365.0,
			volatility:   0.40,
			riskFreeRate: 0.05,
			contractType: "PUT",
		},
		{
			name:         "Long expiry - 2 years",
			underlying:   100.0,
			strike:       100.0,
			timeToExpiry: 2.0,
			volatility:   0.15,
			riskFreeRate: 0.03,
			contractType: "PUT",
		},
		{
			name:         "High volatility",
			underlying:   100.0,
			strike:       100.0,
			timeToExpiry: 0.25,
			volatility:   1.0, // 100% IV
			riskFreeRate: 0.05,
			contractType: "CALL",
		},
		{
			name:         "Low volatility",
			underlying:   100.0,
			strike:       100.0,
			timeToExpiry: 0.25,
			volatility:   0.05, // 5% IV
			riskFreeRate: 0.05,
			contractType: "PUT",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Step 1: Calculate theoretical price
			originalPrice := BlackScholes(tc.underlying, tc.strike, tc.timeToExpiry, tc.volatility, tc.riskFreeRate, tc.contractType)
			
			// Step 2: Calculate implied volatility from that price
			impliedVol := GetImpliedVolatility(originalPrice, tc.underlying, tc.strike, tc.timeToExpiry, tc.riskFreeRate, tc.contractType, 0.0)
			
			// Step 3: Calculate price again using implied volatility
			roundTripPrice := BlackScholes(tc.underlying, tc.strike, tc.timeToExpiry, impliedVol, tc.riskFreeRate, tc.contractType)
			
			// Check accuracy
			priceDiff := math.Abs(originalPrice - roundTripPrice)
			volDiff := math.Abs(tc.volatility - impliedVol)
			
			t.Logf("Original price: $%.6f", originalPrice)
			t.Logf("Original vol:   %.6f (%.2f%%)", tc.volatility, tc.volatility*100)
			t.Logf("Implied vol:    %.6f (%.2f%%)", impliedVol, impliedVol*100)
			t.Logf("Round-trip price: $%.6f", roundTripPrice)
			t.Logf("Price difference: $%.8f", priceDiff)
			t.Logf("Vol difference:   %.8f (%.4f%%)", volDiff, volDiff*100)
			
			// Price should match within 0.01 (1 cent) - realistic tolerance
			if priceDiff > 0.01 {
				t.Errorf("Round-trip price difference too large: $%.6f > $0.01", priceDiff)
			}
			
			// Volatility should match within 0.001 (0.1%) - realistic tolerance  
			if volDiff > 0.001 {
				t.Errorf("Volatility difference too large: %.6f > 0.001", volDiff)
			}
			
			// Sanity checks
			if originalPrice <= 0 {
				t.Errorf("Original price should be positive: $%.6f", originalPrice)
			}
			if impliedVol <= 0 {
				t.Errorf("Implied volatility should be positive: %.6f", impliedVol)
			}
		})
	}
}

// TestRoundTripEdgeCases tests round-trip accuracy for challenging cases
func TestRoundTripEdgeCases(t *testing.T) {
	edgeCases := []struct {
		name         string
		underlying   float64
		strike       float64
		timeToExpiry float64
		volatility   float64
		riskFreeRate float64
		contractType string
		tolerance    float64 // Custom tolerance for difficult cases
	}{
		{
			name:         "Very deep ITM call",
			underlying:   100.0,
			strike:       50.0,
			timeToExpiry: 0.25,
			volatility:   0.20,
			riskFreeRate: 0.05,
			contractType: "CALL",
			tolerance:    0.001, // Slightly higher tolerance
		},
		{
			name:         "Very deep OTM put",
			underlying:   100.0,
			strike:       150.0,
			timeToExpiry: 0.25,
			volatility:   0.20,
			riskFreeRate: 0.05,
			contractType: "PUT",
			tolerance:    0.001,
		},
		{
			name:         "Near-zero volatility",
			underlying:   100.0,
			strike:       100.0,
			timeToExpiry: 0.25,
			volatility:   0.01, // 1% IV
			riskFreeRate: 0.05,
			contractType: "PUT",
			tolerance:    0.001,
		},
		{
			name:         "Extreme high volatility",
			underlying:   100.0,
			strike:       100.0,
			timeToExpiry: 0.25,
			volatility:   2.0, // 200% IV
			riskFreeRate: 0.05,
			contractType: "CALL",
			tolerance:    0.01, // Higher tolerance for extreme vol
		},
	}

	for _, tc := range edgeCases {
		t.Run(tc.name, func(t *testing.T) {
			originalPrice := BlackScholes(tc.underlying, tc.strike, tc.timeToExpiry, tc.volatility, tc.riskFreeRate, tc.contractType)
			
			// Skip if price is too small (numerical issues expected)
			if originalPrice < 0.001 {
				t.Skipf("Skipping case with very small price: $%.8f", originalPrice)
			}
			
			impliedVol := GetImpliedVolatility(originalPrice, tc.underlying, tc.strike, tc.timeToExpiry, tc.riskFreeRate, tc.contractType, 0.0)
			roundTripPrice := BlackScholes(tc.underlying, tc.strike, tc.timeToExpiry, impliedVol, tc.riskFreeRate, tc.contractType)
			
			priceDiff := math.Abs(originalPrice - roundTripPrice)
			volDiff := math.Abs(tc.volatility - impliedVol)
			
			t.Logf("Original price: $%.6f", originalPrice)
			t.Logf("Round-trip price: $%.6f", roundTripPrice)
			t.Logf("Price difference: $%.8f (tolerance: $%.3f)", priceDiff, tc.tolerance)
			t.Logf("Vol difference: %.6f", volDiff)
			
			if priceDiff > tc.tolerance {
				t.Errorf("Price difference %.8f exceeds tolerance %.8f", priceDiff, tc.tolerance)
			}
		})
	}
}


