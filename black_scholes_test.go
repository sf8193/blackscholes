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

// TestGetThetaCallVsPut guards against the case-routing bug where GetTheta
// compared callPut against lowercase "call" instead of the package convention
// "CALL", which silently returned the put theta for CALL options.
func TestGetThetaCallVsPut(t *testing.T) {
	s, k, ti, v, r := 100.0, 100.0, 0.5, 0.2, 0.05
	scale := 365

	callTheta := GetTheta(s, k, ti, v, r, "CALL", scale)
	putTheta := GetTheta(s, k, ti, v, r, "PUT", scale)

	if callTheta == putTheta {
		t.Fatalf("GetTheta returned identical theta for CALL and PUT (%f); callPut case routing is broken", callTheta)
	}
	if !(callTheta < putTheta) {
		t.Errorf("expected call theta (%f) < put theta (%f) for r=%f", callTheta, putTheta, r)
	}
	t.Logf("call theta: %f, put theta: %f", callTheta, putTheta)
}

// TestGetRhoCallVsPut guards against the same case-routing bug in GetRho.
func TestGetRhoCallVsPut(t *testing.T) {
	s, k, ti, v, r := 100.0, 100.0, 0.5, 0.2, 0.05
	scale := 100

	callRho := GetRho(s, k, ti, v, r, "CALL", scale)
	putRho := GetRho(s, k, ti, v, r, "PUT", scale)

	if callRho == putRho {
		t.Fatalf("GetRho returned identical rho for CALL and PUT (%f); callPut case routing is broken", callRho)
	}
	if !(callRho > putRho) {
		t.Errorf("expected call rho (%f) > put rho (%f)", callRho, putRho)
	}
	t.Logf("call rho: %f, put rho: %f", callRho, putRho)
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
			iv, ok := GetImpliedVolatility(tt.premium, tt.strike, tt.underlyingValue, tt.timeToExpiry, tt.riskFreeRate, tt.contractType)
			if !ok {
				t.Fatalf("no vol identified")
			}
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
