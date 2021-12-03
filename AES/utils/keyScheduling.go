package utils

import (
	"strconv"
)

func CircularByteLeftShift(w [4]byte) [4]byte {
	var c [4]byte = w

	var temp byte = c[0]
	for i := 0; i < 3; i++ {
		c[i] = c[i+1]
	}
	c[3] = temp

	return c
}

func ByteSubstitution(w [4]byte) [4]byte {

	var c [4]byte = w

	for i := 0; i < 4; i++ {
		//Decimal to Hex Equivilant
		hexa := strconv.FormatInt(int64(c[i]), 16)

		var firstIndex int
		var secondIndex int
		if len(hexa) > 1 {
			firstIndex = hex2int(string(hexa[0]))
			secondIndex = hex2int(string(hexa[1]))
		} else {
			firstIndex = 0
			secondIndex = hex2int(string(hexa[0]))
		}

		//row * length + col
		c[i] = Sbox[firstIndex*16+secondIndex]
	}

	return c
}

func AddingRoundConstant(w [4]byte, round int) [4]byte {

	var c [4]byte = w

	c[0] = c[0] ^ RoundConstants[round]
	c[1] = c[1] ^ 0x00
	c[2] = c[2] ^ 0x00
	c[3] = c[3] ^ 0x00

	return c
}

func GW(w [4]byte, round int) [4]byte {
	var g [4]byte = w

	g = CircularByteLeftShift(g)
	g = ByteSubstitution(g)
	g = AddingRoundConstant(g, round)

	return g
}

func GenerateRoundKeys(prevRoundKey [16]byte, round int) [16]byte {
	var newRoundKey [16]byte

	//Seperating W Values
	var w [4][4]byte

	var k int = 0
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			w[i][j] = prevRoundKey[k]
			k++
		}
	}

	var gw3 [4]byte = GW(w[3], round)

	var newW [4][4]byte

	//Filling New W Values
	newW[0] = doXOR(w[0], gw3)
	index := 0

	for j := 0; j < 4; j++ {
		newRoundKey[index] = newW[0][j]
		index++
	}

	prevW := 1
	for i := 1; i < 4; i++ {
		newW[i] = doXOR(newW[i-1], w[prevW])
		prevW++

		for j := 0; j < 4; j++ {
			newRoundKey[index] = newW[i][j]
			index++
		}
	}

	return newRoundKey
}
