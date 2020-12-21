package route

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/whoerau/go-uniswap-core/factory"
	"github.com/whoerau/go-uniswap-core/utils/library"
	erc20 "github.com/whoerau/go-uniswap-core/utils/token"
	"math/big"
)

type UniswapV2Router02 struct {
	Factory *factory.UniswapV2Factory
	Address common.Address
}

func NewUniswapV2Router02(factory *factory.UniswapV2Factory, address common.Address) *UniswapV2Router02 {
	return &UniswapV2Router02{
		Factory: factory,
		Address: address,
	}
}

func (route *UniswapV2Router02) swap(amounts []*big.Int, path []common.Address, _to common.Address) {
	len_path := len(path)
	for i := 0; i < len_path-1; i++ {
		input, output := path[i], path[i+1]
		token0, _ := library.SortToken(input, output)
		amountOut := amounts[i+1]

		var amount0Out *big.Int
		var amount1Out *big.Int
		if input == token0 {
			amount0Out, amount1Out = common.Big0, amountOut
		} else {
			amount0Out, amount1Out = amountOut, common.Big0
		}

		var to common.Address
		if i < len_path-2 {
			to = route.Factory.GetPair(output, path[i+2]).Address
		} else {
			to = _to
		}

		route.Factory.GetPair(input, output).Swap(amount0Out, amount1Out, to, nil)
	}
}

func (route *UniswapV2Router02) addLiquidity(tokenA common.Address, tokenB common.Address, amountADesired, amountBDesired, amountAMin, amountBMin *big.Int) (amountA, amountB *big.Int) {
	reserveA, reserveB := library.GetReserves(route.Factory, tokenA, tokenB)
	if reserveA.Cmp(common.Big0) == 0 && reserveB.Cmp(common.Big0) == 0 {
		amountA, amountB = amountADesired, amountBDesired
	} else {
		amountBOptimal := library.Quote(amountADesired, reserveA, reserveB)
		if amountBOptimal.Cmp(amountBDesired) != 1 {
			amountA, amountB = amountADesired, amountBOptimal
		} else {
			amountAOptimal := library.Quote(amountBDesired, reserveB, reserveA)
			if amountBOptimal.Cmp(amountBDesired) == 1 {
				return
			} else {
				amountA, amountB = amountAOptimal, amountBDesired
			}
		}
	}
	return
}

func (route *UniswapV2Router02) AddLiquidity(tokenA common.Address, tokenB common.Address, amountADesired, amountBDesired, amountAMin, amountBMin *big.Int, to common.Address) (amountA, amountB, liquidity *big.Int) {
	amountA, amountB = route.addLiquidity(tokenA, tokenB, amountADesired, amountBDesired, amountAMin, amountBMin)
	pair := route.Factory.GetPair(tokenA, tokenB)
	token0, _ := route.Factory.AddressToInstance.Load(tokenA)
	token1, _ := route.Factory.AddressToInstance.Load(tokenB)
	token0.(*erc20.UniswapV2ERC20).TransferFrom(to, pair.Address, amountA)
	token1.(*erc20.UniswapV2ERC20).TransferFrom(to, pair.Address, amountB)
	liquidity = pair.Mint(to)
	return
}

func (route *UniswapV2Router02) RemoveLiquidity(tokenA common.Address, tokenB common.Address, liquidity, amountAmin, AMountBmin *big.Int, to common.Address) (amountA, amountB *big.Int) {
	pair := route.Factory.GetPair(tokenA, tokenB)
	pair.Liquidity_token.TransferFrom(to, pair.Address, liquidity)
	amount0, amount1 := pair.Burn(to)
	token0, _ := library.SortToken(tokenA, tokenB)
	if tokenA == token0 {
		amountA, amountB = amount0, amount1
	} else {
		amountA, amountB = amount1, amount0
	}
	return
}

func (route *UniswapV2Router02) SwapExactTokensForTokens(amountIn *big.Int, amountOutMin *big.Int, path []common.Address, to common.Address) (amounts []*big.Int) {
	amounts = library.GetAmoutsOut(route.Factory, amountIn, path)
	if amounts[len(amounts)-1].Cmp(amountOutMin) == -1 {
		fmt.Println("error")
	}
	token, _ := route.Factory.AddressToInstance.Load(path[0])
	pair := route.Factory.GetPair(path[0], path[1])
	token.(*erc20.UniswapV2ERC20).TransferFrom(to, pair.Address, amounts[0])
	route.swap(amounts, path, to)
	return
}

func (route *UniswapV2Router02) SwapTokensForExactTokens(amountOut *big.Int, amountInMax *big.Int, path []common.Address, to common.Address) (amounts []*big.Int) {
	amounts = library.GetAmoutsIn(route.Factory, amountOut, path)
	if amounts[0].Cmp(amountInMax) == 1 {
		fmt.Println("error")
	}
	token, _ := route.Factory.AddressToInstance.Load(path[0])
	pair := route.Factory.GetPair(path[0], path[1])
	token.(*erc20.UniswapV2ERC20).TransferFrom(to, pair.Address, amounts[0])
	route.swap(amounts, path, to)
	return
}
