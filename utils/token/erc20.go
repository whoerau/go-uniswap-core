package erc20

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"sync"
)

type UniswapV2ERC20 struct {
	Name             string
	Symbol           string
	Decimals         uint8
	TotalSupply      *big.Int
	rwMuxTotalSupply *sync.RWMutex
	Balance          *sync.Map
	Allowance        *sync.Map
	Nonces           *sync.Map
	Address          common.Address
}

func NewUniswapV2ERC20(name string, symbol string, decimals uint8, address common.Address) *UniswapV2ERC20 {
	return &UniswapV2ERC20{
		Name:             name,
		Symbol:           symbol,
		Decimals:         decimals,
		TotalSupply:      new(big.Int),
		rwMuxTotalSupply: new(sync.RWMutex),
		Balance:          new(sync.Map),
		Allowance:        new(sync.Map),
		Nonces:           new(sync.Map),
		Address:          address,
	}
}

func (erc20 *UniswapV2ERC20) Mint(to common.Address, value *big.Int) {
	erc20.rwMuxTotalSupply.Lock()
	defer erc20.rwMuxTotalSupply.Unlock()
	balance := erc20.BalanceOf(to)
	erc20.TotalSupply = new(big.Int).Add(erc20.TotalSupply, value)
	erc20.Balance.Store(to, new(big.Int).Add(balance, value))
}

func (erc20 *UniswapV2ERC20) Burn(to common.Address, value *big.Int) {
	erc20.rwMuxTotalSupply.Lock()
	defer erc20.rwMuxTotalSupply.Unlock()
	balance := erc20.BalanceOf(to)
	erc20.Balance.Store(to, new(big.Int).Sub(balance, value))
	erc20.TotalSupply = new(big.Int).Sub(erc20.TotalSupply, value)

}

func (erc20 *UniswapV2ERC20) BalanceOf(from common.Address) *big.Int {
	balanceFrom, ok := erc20.Balance.Load(from)
	if !ok {
		return common.Big0
	} else {
		balance := new(big.Int)
		balance = balanceFrom.(*big.Int)
		return balance
	}
}

func (erc20 *UniswapV2ERC20) TransferFrom(from common.Address, to common.Address, value *big.Int) {
	balanceFrom := erc20.BalanceOf(from)
	balanceTo := erc20.BalanceOf(to)
	erc20.Balance.Store(from, new(big.Int).Sub(balanceFrom, value))
	erc20.Balance.Store(to, new(big.Int).Add(balanceTo, value))
}
