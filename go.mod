module github.com/sf8193/blackscholes

go 1.23

// A fast and accurate Black-Scholes options pricing library
// Includes implied volatility calculation using binary search method
retract (
	v1.4.0 // changed GetImpliedVolatility's signature
	v1.5.0 // added SolveImpliedVolatility, since removed
)
