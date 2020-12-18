package library

import "math/big"

func Quote(amountA, reserveA, reserveB *big.Int) (amountB *big.Int) {
	amountB = new(big.Int).Div(new(big.Int).Mul(amountA, reserveB), reserveA)
	return amountB
}

func GetAmountOut(amountIn, reserveIn, reserveOut *big.Int) (amountOut *big.Int) {
	amountInWithFee := new(big.Int).Mul(amountIn, big.NewInt(997))
	numerator := new(big.Int).Mul(amountInWithFee, reserveOut)
	denominator := new(big.Int).Add(new(big.Int).Mul(reserveIn, big.NewInt(1000)), amountInWithFee)
	amountOut = new(big.Int).Div(numerator, denominator)
	return amountOut
}

func GetAmountIn(amountOut, reserveIn, reserveOut *big.Int) (amountIn *big.Int) {
	numerator := new(big.Int).Mul(new(big.Int).Mul(reserveIn, amountOut), big.NewInt(1000))
	denominator := new(big.Int).Mul(new(big.Int).Sub(reserveOut, amountOut), big.NewInt(997))
	amountIn = new(big.Int).Add(new(big.Int).Div(numerator, denominator), big.NewInt(1))
	return amountIn
}
