package RSA

import "math/big"

type Parameters struct {
	P  *big.Int
	Q  *big.Int
	N  *big.Int
	ON *big.Int
}

type PubKey struct {
	E *big.Int
	N *big.Int
}

type PrivKey struct {
	D *big.Int
	N *big.Int
}
