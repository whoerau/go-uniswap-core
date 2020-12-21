package route

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/whoerau/go-uniswap-core/factory"
	erc20 "github.com/whoerau/go-uniswap-core/utils/token"
	"math/big"
	"testing"
)

func TestRoute(t *testing.T) {
	address_token0 := common.BigToAddress(big.NewInt(100))
	address_token1 := common.BigToAddress(big.NewInt(101))
	address_user0 := common.BigToAddress(big.NewInt(300))
	address_user1 := common.BigToAddress(big.NewInt(301))
	address_factory := common.BigToAddress(big.NewInt(400))
	address_route := common.BigToAddress(big.NewInt(500))

	token0 := erc20.NewUniswapV2ERC20("token-0", "TOKEN-0", 18, address_token0)
	token1 := erc20.NewUniswapV2ERC20("token-1", "TOKEN-1", 18, address_token1)

	total0, _ := new(big.Int).SetString("1000000000000000000000000", 10)
	total1, _ := new(big.Int).SetString("2000000000000000000000000", 10)
	total2, _ := new(big.Int).SetString("1000000000000000000000000", 10)
	total3, _ := new(big.Int).SetString("2000000000000000000000000", 10)

	token0.Mint(address_user0, total0)
	token1.Mint(address_user0, total1)
	token0.Mint(address_user1, total2)
	token1.Mint(address_user1, total3)

	t.Log("balance token0", "user0 =", token0.BalanceOf(address_user0), "user1 =", token0.BalanceOf(address_user1))
	t.Log("balance token1", "user0 =", token1.BalanceOf(address_user0), "user1 =", token1.BalanceOf(address_user1))

	f := factory.NewUniswapV2Factory(address_factory)
	r := NewUniswapV2Router02(f, address_route)
	r.Factory.CreatePair(token0, token1)

	tmp1, _ := new(big.Int).SetString("1000000000000000000000", 10)
	tmp2, _ := new(big.Int).SetString("2000000000000000000000", 10)
	r.AddLiquidity(address_token0, address_token1, tmp1, tmp2, common.Big0, common.Big0, address_user0)

	//swap0In, _ := new(big.Int).SetString("100000000000000000000", 10)
	//pathTmp := []common.Address{address_token0, address_token1}
	//r.SwapExactTokensForTokens(swap0In, common.Big0, pathTmp, address_user1)

	swap0InMax, _ := new(big.Int).SetString("100000000000000000000", 10)
	swap1Out, _ := new(big.Int).SetString("181322178776029826316", 10)
	pathTmp := []common.Address{address_token0, address_token1}
	r.SwapTokensForExactTokens(swap1Out, swap0InMax, pathTmp, address_user1)

	liquidity := r.Factory.GetPair(address_token0, address_token1).Liquidity_token.BalanceOf(address_user0)
	r.RemoveLiquidity(address_token0, address_token1, liquidity, common.Big0, common.Big0, address_user0)
}
