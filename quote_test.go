package blackscholes

import (
	"math"
	"testing"
)

const (
	qSpot = 7600.0
	qTte  = 0.0015
	qRate = 0.045
)

func TestRight(t *testing.T) {
	for in, want := range map[string]string{
		"C": "CALL", "c": "CALL", "CALL": "CALL", "Call": "CALL", "call": "CALL",
		"P": "PUT", "p": "PUT", "PUT": "PUT", "Put": "PUT", "put": "PUT",
		"": "", "X": "", "CALLS": "", " C": "",
	} {
		if got := Right(in); got != want {
			t.Errorf("Right(%q) = %q, want %q", in, got, want)
		}
	}
}

// A price made at a known vol gives that vol back (to the solver's precision), whichever spelling of the right.
func TestImpliedVolSpellings(t *testing.T) {
	for _, tc := range []struct {
		strike  float64
		right   string
		spelled string
	}{
		{7550, "PUT", "P"}, {7550, "PUT", "put"}, {7650, "CALL", "C"}, {7650, "CALL", "Call"},
	} {
		price := BlackScholes(qSpot, tc.strike, qTte, 0.14, qRate, tc.right)
		iv, ok := ImpliedVol(price, qSpot, tc.strike, qTte, qRate, tc.spelled)
		if !ok || math.Abs(iv-0.14) > 1e-3 {
			t.Errorf("%s: got (%v, %v), want about 0.14", tc.spelled, iv, ok)
		}
	}
}

// The vol range is [0.05%, 1000%): the edges a real 0DTE wing reaches are inside it.
func TestImpliedVolRange(t *testing.T) {
	for _, vol := range []float64{0.02, 2.0, 6.0} {
		price := BlackScholes(qSpot, 7600, qTte, vol, qRate, "CALL")
		if iv, ok := ImpliedVol(price, qSpot, 7600, qTte, qRate, "CALL"); !ok || math.Abs(iv-vol)/vol > 0.05 {
			t.Errorf("vol %v: got (%v, %v)", vol, iv, ok)
		}
	}
}

// An input the model cannot be asked about, or a price no sane vol explains, is not identified.
func TestImpliedVolNotIdentified(t *testing.T) {
	price := BlackScholes(qSpot, 7600, qTte, 0.14, qRate, "CALL")
	insane := BlackScholes(qSpot, 7600, qTte, 12, qRate, "CALL") // a 1200% vol
	for name, got := range map[string]bool{
		"zero spot":       okOf(ImpliedVol(price, 0, 7600, qTte, qRate, "CALL")),
		"zero strike":     okOf(ImpliedVol(price, qSpot, 0, qTte, qRate, "CALL")),
		"expired":         okOf(ImpliedVol(price, qSpot, 7600, 0, qRate, "CALL")),
		"past expiry":     okOf(ImpliedVol(price, qSpot, 7600, -0.001, qRate, "CALL")),
		"unknown right":   okOf(ImpliedVol(price, qSpot, 7600, qTte, qRate, "X")),
		"no price":        okOf(ImpliedVol(0, qSpot, 7600, qTte, qRate, "CALL")),
		"below intrinsic": okOf(ImpliedVol(qSpot-7200-3, qSpot, 7200, qTte, qRate, "CALL")),
		"vol over 1000":   okOf(ImpliedVol(insane, qSpot, 7600, qTte, qRate, "CALL")),
		"vol under 0.05%": okOf(ImpliedVol(BlackScholes(qSpot, 7600, qTte, 2e-4, qRate, "CALL"), qSpot, 7600, qTte, qRate, "CALL")),
	} {
		if got {
			t.Errorf("%s: ok = true, want false", name)
		}
	}
}

func okOf(_ float64, ok bool) bool { return ok }

func TestDeltaAtVol(t *testing.T) {
	for _, tc := range []struct{ right, spelled string }{{"CALL", "C"}, {"PUT", "p"}} {
		want, err := GetDelta(qSpot, 7600, qTte, 0.14, qRate, tc.right)
		if err != nil {
			t.Fatal(err)
		}
		got, err := DeltaAtVol(qSpot, 7600, qTte, 0.14, qRate, tc.spelled)
		if err != nil || got != want {
			t.Errorf("%s: got (%v, %v), want %v", tc.spelled, got, err, want)
		}
	}
	for name, args := range map[string][4]float64{
		"expired":     {qSpot, 7600, 0, 0.14},
		"past expiry": {qSpot, 7600, -0.001, 0.14},
		"zero spot":   {0, 7600, qTte, 0.14},
		"zero strike": {qSpot, 0, qTte, 0.14},
	} {
		if _, err := DeltaAtVol(args[0], args[1], args[2], args[3], qRate, "CALL"); err == nil {
			t.Errorf("%s: no error", name)
		}
	}
	if _, err := DeltaAtVol(qSpot, 7600, qTte, 0.14, qRate, "X"); err == nil {
		t.Error("unknown right: no error")
	}
}
