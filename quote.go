package blackscholes

import (
	"errors"
	"strings"
)

var errBadQuote = errors.New("an option quote needs a positive spot, strike and time to expiry, and a CALL or PUT right")

// A solved implied volatility outside [minVol, maxVol) (fractions) is garbage: the solver runs to ~0 for a price
// below intrinsic value, and to huge values for one no sane vol explains. 1000% still admits the very high IVs seen
// on 0DTE wings.
const (
	minVol = 5e-4
	maxVol = 10.0
)

// Right returns "CALL" or "PUT" for the short and long spellings providers use ("C", "call", "P", "Put", ...), in
// any case, and "" for anything else. GetDelta, GetImpliedVolatility and the others take only "CALL" and "PUT".
func Right(putCall string) string {
	switch {
	case strings.EqualFold(putCall, "C"), strings.EqualFold(putCall, "CALL"):
		return "CALL"
	case strings.EqualFold(putCall, "P"), strings.EqualFold(putCall, "PUT"):
		return "PUT"
	}
	return ""
}

func quoted(s, k, t float64, putCall string) bool {
	return Right(putCall) != "" && s > 0 && k > 0 && t > 0
}

// ImpliedVol is the implied volatility (a fraction, 0.20 == 20%) of an option's price, for the caller that has a
// quote rather than a clean input: putCall may be spelled any way Right accepts, and ok is false when the price
// identifies no sane vol, or the spot, strike or time to expiry is not positive, or the right is unknown.
func ImpliedVol(price, s, k, t, r float64, putCall string) (iv float64, ok bool) {
	if !quoted(s, k, t, putCall) {
		return 0, false
	}
	iv = GetImpliedVolatility(price, s, k, t, r, Right(putCall), 0)
	if !(iv >= minVol) || iv >= maxVol {
		return 0, false
	}
	return iv, true
}

// DeltaAtVol is GetDelta for the caller that has a quote: putCall may be spelled any way Right accepts, and it
// errors when the spot, strike or time to expiry is not positive or the right is unknown.
func DeltaAtVol(s, k, t, vol, r float64, putCall string) (float64, error) {
	if !quoted(s, k, t, putCall) {
		return 0, errBadQuote
	}
	return GetDelta(s, k, t, vol, r, Right(putCall))
}
