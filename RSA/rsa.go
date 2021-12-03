package RSA

import (
	"math/big"
)

func Initialization(b int) (PubKey, PrivKey) {

	//Variables

	var pubKey PubKey
	var privKey PrivKey

	var p Parameters

	//Parameters Initialization

	//Both should not be same
	for {
		p.P = CreateRandomPrime(b)
		p.Q = CreateRandomPrime(b)

		if p.P.Cmp(p.Q) != 0 {
			break
		}
	}

	p.N = big.NewInt(0).Mul(p.P, p.Q)

	temp1 := big.NewInt(0).Sub(p.P, big.NewInt(1))
	temp2 := big.NewInt(0).Sub(p.Q, big.NewInt(1))
	p.ON = big.NewInt(0).Mul(temp1, temp2)

	//Public Key Generation
	for {
		pubKey.E = CreateRandomInt(p.ON)

		//GCD(E, ON) = 1
		if CheckGcdOne(pubKey.E, p.ON) {
			pubKey.N = p.N
			break
		}
	}

	//Private Key Generation
	privKey.D = ModInv(pubKey.E, p.ON)
	privKey.N = p.N

	return pubKey, privKey
}

func Encrypt(pubKey PubKey, message string) []big.Int {

	var asciiText []byte = []byte(message)
	var cipherText []big.Int

	for _, m := range asciiText {
		//Converting int to big int
		temp := big.NewInt(int64(m))

		//Encryption Formula
		temp = temp.Exp(temp, pubKey.E, pubKey.N)

		cipherText = append(cipherText, *temp)
	}

	return cipherText
}

func Decrypt(privKey PrivKey, cipherText []big.Int) string {
	var asciiText []byte

	for _, c := range cipherText {
		//Decryption Formula
		temp := c.Exp(&c, privKey.D, privKey.N)

		//Converting Big int to Byte
		asciiText = append(asciiText, temp.Bytes()...)
	}

	return string(asciiText)
}
