package RSA

import (
	"crypto/rand"
	"math/big"
)

func CreateRandomPrime(bits int) *big.Int {

	num, err := rand.Prime(rand.Reader, bits)

	for err != nil {
		num, err = rand.Prime(rand.Reader, bits)
	}

	return num
}

func CreateRandomInt(p *big.Int) *big.Int {

	for {
		num, err := rand.Int(rand.Reader, p)

		for err != nil {
			num, err = rand.Int(rand.Reader, p)
		}

		// Num > 1
		if num.Cmp(big.NewInt(1)) == 1 {
			return num
		}
	}
}

func ModInv(p *big.Int, n *big.Int) *big.Int {

	bigP := new(big.Int)
	bigP = bigP.Set(p)

	bigBase := new(big.Int)
	bigBase = bigBase.Set(n)

	// Compute inverse of bigP modulo bigBase
	bigGcd := big.NewInt(0)
	bigX := big.NewInt(0)
	bigGcd.GCD(bigX, nil, bigP, bigBase)

	// x*bigP+y*bigBase=1
	// => x*bigP = 1 modulo bigBase

	if bigGcd.Cmp(big.NewInt(1)) != 0 {
		str := "GCD of " + bigP.String() + " wrt " + bigBase.String() + " != 1"
		panic(str)
	}

	return Mod(bigX, n)
}

func CheckGcdOne(p *big.Int, n *big.Int) bool {

	bigP := new(big.Int)
	bigP = bigP.Set(p)

	bigBase := new(big.Int)
	bigBase = bigBase.Set(n)

	// Compute inverse of bigP modulo bigBase
	bigGcd := big.NewInt(0)
	bigX := big.NewInt(0)
	bigGcd.GCD(bigX, nil, bigP, bigBase)

	//GCD is 1
	return bigGcd.Cmp(big.NewInt(1)) == 0
}

func Mod(n *big.Int, p *big.Int) *big.Int {
	if n.Cmp(big.NewInt(0)) == -1 {
		for n.Cmp(big.NewInt(0)) == -1 {
			n = n.Add(n, p)
		}
		return n
	} else {
		return n.Mod(n, p)
	}
}
