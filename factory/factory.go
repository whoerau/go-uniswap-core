package factory

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/whoerau/go-uniswap-core/pair"
	"github.com/whoerau/go-uniswap-core/utils/token"
	"sync"
)

type UniswapV2Factory struct {
	Pairs             *sync.Map
	AddressToInstance *sync.Map
	Address           common.Address
}

func NewUniswapV2Factory(address common.Address) *UniswapV2Factory {
	return &UniswapV2Factory{
		Pairs:             new(sync.Map),
		AddressToInstance: new(sync.Map),
		Address:           address,
	}
}

func (uf *UniswapV2Factory) CreatePair(tokenA, tokenB *erc20.UniswapV2ERC20) *pair.UniswapV2Pair {
	pr := pair.NewUniswapV2Pair(tokenA, tokenB)
	var p = new(sync.Map)
	var q = new(sync.Map)
	p.Store(tokenA.Address, pr)
	q.Store(tokenB.Address, pr)
	uf.Pairs.Store(tokenB.Address, p)
	uf.Pairs.Store(tokenA.Address, q)
	uf.AddressToInstance.Store(tokenA.Address, tokenA)
	uf.AddressToInstance.Store(tokenB.Address, tokenB)
	return pr
}

func (uf *UniswapV2Factory) GetPair(tokenA, tokenB common.Address) *pair.UniswapV2Pair {
	p, ok := uf.Pairs.Load(tokenA)
	if !ok {
		fmt.Println("error a")
		return nil
	} else {
		q, ok := p.(*sync.Map).Load(tokenB)
		if !ok {
			fmt.Println("error b")
			return nil
		} else {
			return q.(*pair.UniswapV2Pair)
		}
	}
}
