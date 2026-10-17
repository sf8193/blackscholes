package blackscholes

import (
	"math"
	"testing"
)

// A 0DTE SPX wing is worth cents; the old solver stopped at a one-cent price match.
func TestImpliedVolWingIsExact(t *testing.T) {
	const s, t0, r = 6000.0, 0.0015, 0.045
	for _, tc := range []struct {
		k       float64
		vol     float64
		callPut string
	}{
		{5950, 0.14, "PUT"},
		{5900, 0.17, "PUT"},
		{6050, 0.12, "CALL"},
		{6100, 0.15, "CALL"},
	} {
		price := BlackScholes(s, tc.k, t0, tc.vol, r, tc.callPut)
		iv, ok := SolveImpliedVolatility(price, s, tc.k, t0, r, tc.callPut)
		if !ok {
			t.Fatalf("%s %.0f: no vol identified for price %.4f", tc.callPut, tc.k, price)
		}
		if math.Abs(iv-tc.vol) > 1e-13 {
			t.Errorf("%s %.0f: iv %.12f, want %.12f", tc.callPut, tc.k, iv, tc.vol)
		}
	}
}

// Inputs that identify no vol return iv 0, so a delta-only caller gets +/-1 or 0.
func TestImpliedVolNotIdentified(t *testing.T) {
	good := BlackScholes(100, 100, 0.25, 0.2, 0.05, "CALL")
	atMin := BlackScholes(100, 100, 0.25, 5e-4, 0.05, "CALL")
	atMax := BlackScholes(100, 100, 0.25, 10.0, 0.05, "CALL")
	nan := math.NaN()
	for _, tc := range []struct {
		name              string
		price, s, k, t, r float64
		callPut           string
	}{
		{"zero price", 0, 100, 100, 0.25, 0.05, "CALL"},
		{"zero price far out of the money", 0, 100, 200, 0.25, 0.05, "CALL"},
		{"negative price", -1, 100, 100, 0.25, 0.05, "CALL"},
		{"NaN price", nan, 100, 100, 0.25, 0.05, "CALL"},
		{"zero spot", good, 0, 100, 0.25, 0.05, "CALL"},
		{"zero strike", good, 100, 0, 0.25, 0.05, "CALL"},
		{"expired", good, 100, 100, 0, 0.05, "CALL"},
		{"NaN rate", good, 100, 100, 0.25, nan, "CALL"},
		{"+Inf price", math.Inf(1), 100, 100, 0.25, 0.05, "CALL"},
		{"NaN spot", good, nan, 100, 0.25, 0.05, "CALL"},
		{"+Inf spot", good, math.Inf(1), 100, 0.25, 0.05, "CALL"},
		{"+Inf time", good, 100, 100, math.Inf(1), 0.05, "CALL"},
		{"negative time", good, 100, 100, -0.25, 0.05, "CALL"},
		{"NaN strike", good, 100, nan, 0.25, 0.05, "CALL"},
		{"NaN time", good, 100, 100, nan, 0.05, "CALL"},
		{"negative spot", good, -1, 100, 0.25, 0.05, "CALL"},
		{"negative strike", good, 100, -1, 0.25, 0.05, "PUT"},
		{"negative spot and strike", good, -1, -1, 0.25, 0.05, "CALL"},
		{"bad right", good, 100, 100, 0.25, 0.05, "X"},
		{"below intrinsic", 1, 100, 50, 0.25, 0.05, "CALL"},
		{"just under the min-vol price", atMin * (1 - 1e-9), 100, 100, 0.25, 0.05, "CALL"},
		{"just over the max-vol price", atMax * (1 + 1e-9), 100, 100, 0.25, 0.05, "CALL"},
	} {
		iv, ok := SolveImpliedVolatility(tc.price, tc.s, tc.k, tc.t, tc.r, tc.callPut)
		if ok || iv != 0 {
			t.Errorf("%s: got (%v, %v), want (0, false)", tc.name, iv, ok)
		}
	}
	// Deep in the money near expiry: a vol solves but vega is too small to trust it.
	if _, ok := SolveImpliedVolatility(100-50*math.Exp(-0.05*0.001), 100, 50, 0.001, 0.05, "CALL"); ok {
		t.Error("deep ITM near expiry: ok = true, want false")
	}
}

// The edges of the vol range are identified, not rejected.
func TestImpliedVolBandEdges(t *testing.T) {
	for _, vol := range []float64{5e-4, 10.0} {
		price := BlackScholes(100, 100, 0.25, vol, 0.05, "CALL")
		iv, _ := SolveImpliedVolatility(price, 100, 100, 0.25, 0.05, "CALL")
		if math.Abs(iv-vol) > 1e-9 {
			t.Errorf("price at vol %v: iv %v", vol, iv)
		}
	}
}

// Vega is measured in dollars per 1.00 of vol; the floor is 0.05. At S=100, T=0.01, r=0.05,
// vol 0.20, a call struck at 94.6 has vega 0.077 and one struck at 94.2 has 0.041.
func TestImpliedVolVegaFloor(t *testing.T) {
	for _, tc := range []struct {
		k      float64
		wantOk bool
	}{{94.6, true}, {94.2, false}} {
		price := BlackScholes(100, tc.k, 0.01, 0.2, 0.05, "CALL")
		iv, ok := SolveImpliedVolatility(price, 100, tc.k, 0.01, 0.05, "CALL")
		if ok != tc.wantOk || math.Abs(iv-0.2) > 1e-9 {
			t.Errorf("K %.1f: got (%v, %v), want (0.2, %v)", tc.k, iv, ok, tc.wantOk)
		}
	}
}
