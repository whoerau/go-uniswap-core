package pair

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/whoerau/go-uniswap-core/utils/library"
	"github.com/whoerau/go-uniswap-core/utils/math"
	"github.com/whoerau/go-uniswap-core/utils/token"
	"math/big"
	"sync"
)

type UniswapV2Pair struct {
	Liquidity_token   *erc20.UniswapV2ERC20
	FeeTo             common.Address
	Token0            *erc20.UniswapV2ERC20
	Token1            *erc20.UniswapV2ERC20
	reserve0          *big.Int
	reserve1          *big.Int
	Klast             *big.Int
	rwmux             *sync.RWMutex
	Minimum_liquidity *big.Int
	Address           common.Address
}

func NewUniswapV2Pair(token0, token1 *erc20.UniswapV2ERC20) *UniswapV2Pair {
	return &UniswapV2Pair{
		Liquidity_token:   erc20.NewUniswapV2ERC20("Uniswap V2", "UNI-V2", 18, common.BigToAddress(big.NewInt(102))),
		FeeTo:             common.BigToAddress(new(big.Int)),
		Token0:            token0,
		Token1:            token1,
		reserve0:          new(big.Int),
		reserve1:          new(big.Int),
		Klast:             new(big.Int),
		rwmux:             new(sync.RWMutex),
		Minimum_liquidity: big.NewInt(1000),
		Address:           common.BigToAddress(big.NewInt(200)),
	}
}

func (pair *UniswapV2Pair) update(balance0, balance1, _reserve0, _reserve1 *big.Int) {
	// 省略计算 blockTimestamp， timeElapsed， price0CumulativeLast， price1CumulativeLast
	fmt.Println("update reserve", "balance0 =", balance0, "balance1 =", balance1, "_reserve0 =", _reserve0, "_reserve1 =", _reserve1)
	pair.reserve0 = balance0
	pair.reserve1 = balance1
}

func (pair *UniswapV2Pair) mintFee(_reserve0 *big.Int, _reserve1 *big.Int) (feeOn bool) {
	feeTo := pair.FeeTo
	feeOn = feeTo != common.Address{}
	_klast := pair.Klast
	if feeOn {
		if pair.Klast.Cmp(common.Big0) != 0 {
			rootK := math.Sqrt(new(big.Int).Mul(_reserve0, _reserve1).String())
			rootKlast := math.Sqrt(_klast.String())
			if rootK.Cmp(rootKlast) == 1 {
				numerator := new(big.Int).Mul(pair.Liquidity_token.TotalSupply, new(big.Int).Sub(rootK, rootKlast))
				denominator := new(big.Int).Add(new(big.Int).Mul(big.NewInt(5), rootK), rootKlast)
				liquidity := new(big.Int).Div(numerator, denominator)
				if liquidity.Cmp(common.Big0) == 1 {
					pair.Liquidity_token.Mint(feeTo, liquidity)
				}
			}
		}
	} else if _klast.Cmp(common.Big0) != 0 {
		_klast = common.Big0
	}
	return feeOn
}

func (pair *UniswapV2Pair) getReserves() (*big.Int, *big.Int) {
	return pair.reserve0, pair.reserve1
}

func (pair *UniswapV2Pair) Mint(to common.Address) (liquidity *big.Int) {
	_reserve0, _reserve1 := pair.getReserves()
	balance0 := pair.Token0.BalanceOf(pair.Address)
	balance1 := pair.Token1.BalanceOf(pair.Address)

	amount0 := new(big.Int).Sub(balance0, _reserve0)
	amount1 := new(big.Int).Sub(balance1, _reserve1)
	feeOn := pair.mintFee(_reserve0, _reserve1)
	_totalSupply := pair.Liquidity_token.TotalSupply
	fmt.Println("first add liquidity", "balance0 =", balance0, "amount0 = ", amount0, "_reserve0 =", _reserve0, "balance1 =", balance1, "amount1 = ", amount1, "_reserve1 =", _reserve1)
	if _totalSupply.Cmp(common.Big0) == 0 {
		liquidity = new(big.Int).Sub(math.Sqrt(new(big.Int).Mul(amount0, amount1).String()), pair.Minimum_liquidity)
		pair.Liquidity_token.Mint(common.Address{}, pair.Minimum_liquidity)
	} else {
		liquidity = math.Min(new(big.Int).Div(new(big.Int).Mul(amount0, _totalSupply), _reserve0), new(big.Int).Div(new(big.Int).Mul(amount1, _totalSupply), _reserve1))
	}
	if liquidity.Cmp(common.Big0) != 1 {
		fmt.Println("liquidity error", liquidity)
	}
	pair.Liquidity_token.Mint(to, liquidity)
	pair.update(balance0, balance1, _reserve0, _reserve1)
	if feeOn {
		pair.Klast = new(big.Int).Mul(pair.reserve0, pair.reserve1)
	}
	return liquidity
}

func (pair *UniswapV2Pair) Burn(to common.Address) (amount0 *big.Int, amount1 *big.Int) {
	_reserve0, _reserve1 := pair.getReserves()
	_token0 := pair.Token0
	_token1 := pair.Token1
	balance0 := _token0.BalanceOf(pair.Address)
	balance1 := _token1.BalanceOf(pair.Address)
	liquidity := pair.Liquidity_token.BalanceOf(pair.Address)

	feeOn := pair.mintFee(_reserve0, _reserve1)
	_totalSupply := pair.Liquidity_token.TotalSupply
	amount0 = new(big.Int).Div(new(big.Int).Mul(liquidity, balance0), _totalSupply)
	amount1 = new(big.Int).Div(new(big.Int).Mul(liquidity, balance1), _totalSupply)
	if amount0.Cmp(common.Big0) != 1 || amount1.Cmp(common.Big0) != 1 {
		fmt.Println("amount error", "amount0", amount0, "amount1", amount1)
	}
	pair.Liquidity_token.Burn(pair.Address, liquidity)
	pair.Token0.TransferFrom(pair.Address, to, amount0)
	pair.Token1.TransferFrom(pair.Address, to, amount1)
	balance0 = pair.Token0.BalanceOf(pair.Address)
	balance1 = pair.Token1.BalanceOf(pair.Address)
	pair.update(balance0, balance1, _reserve0, _reserve1)
	if feeOn {
		pair.Klast = new(big.Int).Mul(pair.reserve0, pair.reserve1)
	}
	fmt.Println("Burn", "amount0", amount0, "amount1", amount1, "to", to.Hex())
	return amount0, amount1
}

// 忽略 data 字段
func (pair *UniswapV2Pair) Swap(amount0Out *big.Int, amount1Out *big.Int, to common.Address, data []byte) {
	_reserve0, _reserve1 := pair.getReserves()
	var balance0 *big.Int
	var balance1 *big.Int
	_token0 := pair.Token0
	_token1 := pair.Token1
	if amount0Out.Cmp(common.Big0) == 1 {
		_token0.TransferFrom(pair.Address, to, amount0Out)
	}
	if amount1Out.Cmp(common.Big0) == 1 {
		_token1.TransferFrom(pair.Address, to, amount1Out)
	}
	balance0 = _token0.BalanceOf(pair.Address)
	balance1 = _token1.BalanceOf(pair.Address)

	// 二次验证是否已经扣除手续费
	var amount0In *big.Int
	var amount1In *big.Int
	if balance0.Cmp(new(big.Int).Sub(_reserve0, amount0Out)) == 1 {
		amount0In = new(big.Int).Sub(balance0, new(big.Int).Sub(_reserve0, amount0Out))
	} else {
		amount0In = common.Big0
	}
	if balance1.Cmp(new(big.Int).Sub(_reserve1, amount1Out)) == 1 {
		amount1In = new(big.Int).Sub(balance1, new(big.Int).Sub(_reserve1, amount1Out))
	} else {
		amount1In = common.Big0
	}
	balance0Adjusted := new(big.Int).Sub(new(big.Int).Mul(balance0, big.NewInt(1000)), new(big.Int).Mul(amount0In, big.NewInt(3)))
	balance1Adjusted := new(big.Int).Sub(new(big.Int).Mul(balance1, big.NewInt(1000)), new(big.Int).Mul(amount1In, big.NewInt(3)))
	if new(big.Int).Mul(balance0Adjusted, balance1Adjusted).Cmp(new(big.Int).Mul(new(big.Int).Mul(_reserve0, _reserve1), big.NewInt(1000000))) == -1 {
		fmt.Println("UniswapV2: K error", "balance0Adjusted", balance0Adjusted, "balance1Adjusted", balance1Adjusted, "_reserve0", _reserve0, "_reserve1", _reserve1)
	}
	pair.update(balance0, balance1, _reserve0, _reserve1)
	fmt.Println("Swap", "amount0In =", amount0In, "amount1In =", amount1In, "amount0Out =", amount0Out, "amount1Out =", amount1Out)
}

func (pair *UniswapV2Pair) addLiquidity(amountADesired, amountBDesired, amountAMin, amountBMin *big.Int) (amountA, amountB *big.Int) {
	reserveA, reserveB := pair.getReserves()
	if reserveA.Cmp(common.Big0) == 0 && reserveB.Cmp(common.Big0) == 0 {
		amountA, amountB = amountADesired, amountBDesired
	} else {
		amountBOptimal := library.Quote(amountADesired, reserveA, reserveB)
		if amountBOptimal.Cmp(amountBDesired) != 1 {
			if amountBOptimal.Cmp(amountBMin) == -1 {
				fmt.Println("error, UniswapV2Router: INSUFFICIENT_B_AMOUNT")
			}
			amountA, amountB = amountADesired, amountBOptimal
		} else {
			amountAOptimal := library.Quote(amountBDesired, reserveB, reserveA)
			if amountAOptimal.Cmp(amountADesired) == 1 || amountAOptimal.Cmp(amountAMin) == -1 {
				fmt.Println("error, UniswapV2Router: INSUFFICIENT_A_AMOUNT")
			}
			amountA, amountB = amountAOptimal, amountADesired
		}
	}
	return amountA, amountB
}

// token0, token1 的是顺序就是 A, B
func (pair *UniswapV2Pair) AddLiquidity(amountADesired, amountBDesired, amountAMin, amountBMin *big.Int, to common.Address) (amountA, amountB, liquidity *big.Int) {
	amountA, amountB = pair.addLiquidity(amountADesired, amountBDesired, amountAMin, amountBMin)
	pair.Token0.TransferFrom(to, pair.Address, amountA)
	pair.Token1.TransferFrom(to, pair.Address, amountB)
	liquidity = pair.Mint(to)
	return amountA, amountB, liquidity
}

// token0, token1 的是顺序就是 A, B
func (pair *UniswapV2Pair) RemoveLiquidity(liquidity, amountAMin, amountBMin *big.Int, to common.Address) (amountA, amountB *big.Int) {
	pair.Liquidity_token.TransferFrom(to, pair.Address, liquidity)
	amount0, amount1 := pair.Burn(to)
	amountA, amountB = amount0, amount1
	return amountA, amountB
}

// token0 -> token1
func (pair *UniswapV2Pair) SwapExactTokensForTokens(amount0In, amount1OutMIn *big.Int, to common.Address) *big.Int {
	reserveIn, reserveOut := pair.getReserves()
	amount1Out := library.GetAmountOut(amount0In, reserveIn, reserveOut)
	if amount1Out.Cmp(amount1OutMIn) == -1 {
		fmt.Println("UniswapV2Router: INSUFFICIENT_OUTPUT_AMOUNT")
	}
	pair.Token0.TransferFrom(to, pair.Address, amount1Out)
	pair.Swap(common.Big0, amount1Out, to, nil)
	return amount1Out
}

// token0 -> token1
func (pair *UniswapV2Pair) SwapTokensForExactTokens(amount1Out, amount0InMax *big.Int, to common.Address) *big.Int {
	reserveIn, reserveOut := pair.getReserves()
	amount0In := library.GetAmountIn(amount1Out, reserveIn, reserveOut)
	if amount0In.Cmp(amount0InMax) == 1 {
		fmt.Println("UniswapV2Router: EXCESSIVE_INPUT_AMOUNT")
	}
	pair.Token0.TransferFrom(to, pair.Address, amount0In)
	pair.Swap(common.Big0, amount1Out, to, nil)
	return amount0In
}
