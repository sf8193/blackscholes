package blackscholes

import (
	"math"
	"testing"
)

// A 0DTE SPX wing is worth cents. The solver must reproduce its price to a tenth of a cent, not a cent.
func TestImpliedVolatilityMatchesAWingPriceToATenthOfACent(t *testing.T) {
	const s, t0, r = 6000.0, 0.0015, 0.045
	for _, tc := range []struct {
		k       float64
		vol     float64
		callPut string
	}{
		{5950, 0.14, "PUT"},
		{5900, 0.17, "PUT"},
		{5850, 0.19, "PUT"},
		{6050, 0.12, "CALL"},
		{6100, 0.15, "CALL"},
		{6150, 0.18, "CALL"},
	} {
		price := BlackScholes(s, tc.k, t0, tc.vol, r, tc.callPut)
		iv := GetImpliedVolatility(price, s, tc.k, t0, r, tc.callPut, 0)
		if diff := math.Abs(BlackScholes(s, tc.k, t0, iv, r, tc.callPut) - price); diff >= 0.001 {
			t.Errorf("%s %.0f: price %.4f, solved vol %.6f reprices %.6f off", tc.callPut, tc.k, price, iv, diff)
		}
	}
}
