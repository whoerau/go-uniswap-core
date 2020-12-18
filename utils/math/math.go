package math

import (
	"math/big"
)

func Sqrt(s string) *big.Int {
	var n, a, b, m, m2 big.Int
	n.SetString(s, 10)
	a.SetInt64(int64(1))
	b.Set(&n)

	for {
		m.Add(&a, &b).Div(&m, big.NewInt(2))
		if m.Cmp(&a) == 0 || m.Cmp(&b) == 0 {
			break
		}
		m2.Mul(&m, &m)
		if m2.Cmp(&n) > 0 {
			b.Set(&m)
		} else {
			a.Set(&m)
		}
	}
	return &m
}

func Min(x, y *big.Int) *big.Int {
	result := x.Cmp(y)
	var rsp *big.Int
	switch result {
	case 0:
		rsp = x
	case -1:
		rsp = x
	case 1:
		rsp = y
	}
	return rsp
}
