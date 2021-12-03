package utils

import (
	"math/rand"
	"strconv"
	"strings"
)

func doXOR(a [4]byte, b [4]byte) [4]byte {
	var c [4]byte

	for i := 0; i < 4; i++ {
		c[i] = a[i] ^ b[i]
	}

	return c
}

func hex2int(l string) int {
	switch l {
	case "0":
		return 0
	case "1":
		return 1
	case "2":
		return 2
	case "3":
		return 3
	case "4":
		return 4
	case "5":
		return 5
	case "6":
		return 6
	case "7":
		return 7
	case "8":
		return 8
	case "9":
		return 9
	case "a":
		return 10
	case "b":
		return 11
	case "c":
		return 12
	case "d":
		return 13
	case "e":
		return 14
	case "f":
		return 15
	}
	return 0
}

func CheckOverflow(a byte, b byte) bool {

	var result uint16 = uint16(a) * uint16(b)

	return result > 255
}

func MultiplicationWithOverflowCheck(a byte, b byte) byte {

	var result byte

	if a == 0x01 {
		result = b
	} else if a == 0x02 {
		result = a * b

		if CheckOverflow(a, b) {
			result ^= 0x1B
		}
	} else if a == 0x03 {
		result = 0x02 * b

		if CheckOverflow(0x02, b) {
			result ^= 0x1B
		}

		result ^= b

	}

	return result
}

func ConvertToArrayIndex(a byte) (int, int) {
	//Decimal to Hex Equivalent
	hexa := strconv.FormatInt(int64(a), 16)

	var firstIndex int
	var secondIndex int
	if len(hexa) > 1 {
		firstIndex = hex2int(string(hexa[0]))
		secondIndex = hex2int(string(hexa[1]))
	} else {
		firstIndex = 0
		secondIndex = hex2int(string(hexa[0]))
	}

	return firstIndex, secondIndex
}

func RandomKeyGenerator(passwordLength int, minSpecialChar int, minNum int, minUpperCase int) string {

	lowerCharSet := "abcdedfghijklmnopqrst"
	upperCharSet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialCharSet := "!@#$%&*"
	numberSet := "0123456789"
	allCharSet := lowerCharSet + upperCharSet + specialCharSet + numberSet

	var password strings.Builder

	//Set special character
	for i := 0; i < minSpecialChar; i++ {
		random := rand.Intn(len(specialCharSet))
		password.WriteString(string(specialCharSet[random]))
	}

	//Set numeric
	for i := 0; i < minNum; i++ {
		random := rand.Intn(len(numberSet))
		password.WriteString(string(numberSet[random]))
	}

	//Set uppercase
	for i := 0; i < minUpperCase; i++ {
		random := rand.Intn(len(upperCharSet))
		password.WriteString(string(upperCharSet[random]))
	}

	remainingLength := passwordLength - minSpecialChar - minNum - minUpperCase
	for i := 0; i < remainingLength; i++ {
		random := rand.Intn(len(allCharSet))
		password.WriteString(string(allCharSet[random]))
	}
	inRune := []rune(password.String())
	rand.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})

	return string(inRune)
}
