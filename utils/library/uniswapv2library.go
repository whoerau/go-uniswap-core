package library

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/whoerau/go-uniswap-core/factory"
	"math/big"
)

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

func SortToken(tokenA common.Address, tokenB common.Address) (token0 common.Address, token1 common.Address) {
	if tokenA.Hex() < tokenB.Hex() {
		return tokenA, tokenB
	} else {
		return tokenB, tokenA
	}
}

func GetReserves(factory *factory.UniswapV2Factory, tokenA common.Address, tokenB common.Address) (reserveA *big.Int, reserveB *big.Int) {
	token0, token1 := SortToken(tokenA, tokenB)
	pair := factory.GetPair(token0, token1)
	reserve0, reserve1 := pair.GetReserves()
	if tokenA == token0 {
		return reserve0, reserve1
	} else {
		return reserve1, reserve0
	}
}

func GetAmoutsOut(factory *factory.UniswapV2Factory, amountIn *big.Int, path []common.Address) []*big.Int {
	lenPath := len(path)
	amounts := make([]*big.Int, lenPath)
	amounts[0] = amountIn
	for i := 0; i < lenPath-1; i++ {
		reserveIn, reserveOut := GetReserves(factory, path[i], path[i+1])
		amounts[i+1] = GetAmountOut(amounts[i], reserveIn, reserveOut)
	}
	return amounts
}

func GetAmoutsIn(factory *factory.UniswapV2Factory, amountOut *big.Int, path []common.Address) []*big.Int {
	lenPath := len(path)
	amounts := make([]*big.Int, lenPath)
	amounts[lenPath-1] = amountOut
	for i := lenPath - 1; i > 0; i-- {
		reserveIn, reserveOut := GetReserves(factory, path[i-1], path[i])
		amounts[i-1] = GetAmountIn(amounts[i], reserveIn, reserveOut)
	}
	return amounts
}
