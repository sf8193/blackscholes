package blackscholes

import (
	"math"
	"testing"
)

func TestCalculateGetDelta(t *testing.T) {
	tests := []struct {
		name            string
		timeToExpiry    float64
		underlyingValue float64
		strike          float64
		contractType    string
		expectedDelta   float64
		riskFreeRate    float64
		volatility      float64
	}{
		{
			name:            "putTest",
			timeToExpiry:    0.3846,
			underlyingValue: 49.0,
			strike:          50.0,
			contractType:    "PUT",
			expectedDelta:   -0.4783983660284239,
			riskFreeRate:    0.05,
			volatility:      0.2,
		},
		{
			name:            "callTest",
			timeToExpiry:    0.3846,
			underlyingValue: 49.0,
			strike:          50.0,
			contractType:    "CALL",
			expectedDelta:   0.522, // Placeholder value for expected delta
			riskFreeRate:    0.05,
			volatility:      0.2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delta, _ := GetDelta(tt.underlyingValue, tt.strike, tt.timeToExpiry, tt.volatility, tt.riskFreeRate, tt.contractType)
			if delta == 0 {
				t.Errorf("calculateDelta() returned zero delta")
			}
			if !(math.Abs(delta-tt.expectedDelta) < 0.1) {
				t.Fatalf("Delta value %f is out of expected range", delta)
			}
			t.Logf("Calculated Delta: %f", delta)
		})
	}
}

func TestCalculateIV(t *testing.T) {
	tests := []struct {
		name            string
		timeToExpiry    float64
		underlyingValue float64
		strike          float64
		contractType    string
		expectedVol     float64
		riskFreeRate    float64
		premium         float64
	}{
		{
			name:            "callTest2",
			timeToExpiry:    0.5,
			underlyingValue: 100.00,
			strike:          100.00,
			contractType:    "CALL",
			expectedVol:     0.20,
			riskFreeRate:    0.01,
			premium:         5.87602423383,
		},
	}

	// func BSImpliedVol(callType bool, lastTradedPrice float64, underlying float64, strike float64, timeToExpiration float64, startAnchorVolatility float64, riskFreeInterest float64, dividend float64) float64 {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iv := GetImpliedVolatility(tt.premium, tt.strike, tt.underlyingValue, tt.timeToExpiry, tt.riskFreeRate, tt.contractType, 0)
			if iv == 0 {
				t.Errorf("calculateDelta() returned zero delta")
			}
			if !(math.Abs(iv-tt.expectedVol) < 0.1) {
				t.Fatalf("iv value %f is out of expected range", iv)
			}
			t.Logf("Calculated Delta: %f", iv)
		})
	}
}
