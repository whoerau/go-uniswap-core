package pair

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/whoerau/go-uniswap-core/utils/token"
	"math/big"
	"testing"
)

func TestPair(t *testing.T) {
	address_token0 := common.BigToAddress(big.NewInt(100))
	address_token1 := common.BigToAddress(big.NewInt(101))
	address_user0 := common.BigToAddress(big.NewInt(300))
	address_user1 := common.BigToAddress(big.NewInt(301))

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

	//token0.Balance.Range(func(k, v interface{}) bool {
	//	t.Log("iterate_0:\n", k, v)
	//	return true
	//})
	//token1.Balance.Range(func(k, v interface{}) bool {
	//	t.Log("iterate_1:\n", k, v)
	//	return true
	//})

	pair := NewUniswapV2Pair(token0, token1)
	t.Log("balance Liquidity_token", "user0 =", pair.Liquidity_token.BalanceOf(address_user1), "user1 =", pair.Liquidity_token.BalanceOf(address_user1))

	tmp1, _ := new(big.Int).SetString("1000000000000000000000", 10)
	tmp2, _ := new(big.Int).SetString("2000000000000000000000", 10)
	amountA, amountB, liquidity := pair.AddLiquidity(tmp1, tmp2, common.Big0, common.Big0, address_user0)
	t.Log("AddLiquidity", "amountA =", amountA, "amountB =", amountB, "liquidity =", liquidity)

	swap0In, _ := new(big.Int).SetString("100000000000000000000", 10)
	t.Log("balance user1 before swap", "user0_token0 =", token0.BalanceOf(address_user1), "user0_token1 =", token1.BalanceOf(address_user1))
	pair.SwapExactTokensForTokens(swap0In, common.Big0, address_user1)
	t.Log("balance user1 after swap", "user0_token0 = ", token0.BalanceOf(address_user1), "user0_token1 =", token1.BalanceOf(address_user1))

	//swap1Out, _ := new(big.Int).SetString("181322178776029826316", 10)
	//t.Log("balance user1 before swap", "user0_token0 =", token0.BalanceOf(address_user1), "user0_token1 =", token1.BalanceOf(address_user1))
	//pair.SwapTokensForExactTokens(swap1Out, common.Big0, address_user1)
	//t.Log("balance user1 after swap", "user0_token0 = ", token0.BalanceOf(address_user1), "user0_token1 =", token1.BalanceOf(address_user1))

	liquidity_0 := pair.Liquidity_token.BalanceOf(address_user0)
	pair.RemoveLiquidity(liquidity_0, common.Big0, common.Big0, address_user0)
	t.Log("balance token0", "user0 =", token0.BalanceOf(address_user0), "user1 =", token0.BalanceOf(address_user1))
	t.Log("balance token1", "user0 =", token1.BalanceOf(address_user0), "user1 =", token1.BalanceOf(address_user1))
	t.Log("liquidity_0 =", pair.Liquidity_token.BalanceOf(address_user0))
}
