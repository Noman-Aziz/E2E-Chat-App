package utils

func AddRoundKey(stateMatrix [16]byte, roundKey [16]byte) [16]byte {
	var newStateMatrix [16]byte

	for i := 0; i < 16; i++ {
		newStateMatrix[i] = stateMatrix[i] ^ roundKey[i]
	}

	return newStateMatrix
}

func BytesSubstitution(stateMatrix [16]byte) [16]byte {

	var newStateMatrix [16]byte = stateMatrix

	for i := 0; i < 16; i++ {
		firstIndex, secondIndex := ConvertToArrayIndex(newStateMatrix[i])

		//row * length + col
		newStateMatrix[i] = Sbox[firstIndex*16+secondIndex]
	}

	return newStateMatrix
}

func ShiftRows(stateMatrix [16]byte) [16]byte {
	var newStateMatrix [16]byte = stateMatrix

	//Converting Col Major into Row Major
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			stateMatrix[i*4+j] = newStateMatrix[j*4+i]
		}
	}

	//Offset 0
	newStateMatrix[0] = stateMatrix[0]
	newStateMatrix[1] = stateMatrix[1]
	newStateMatrix[2] = stateMatrix[2]
	newStateMatrix[3] = stateMatrix[3]

	//Offset 1
	newStateMatrix[4] = stateMatrix[5]
	newStateMatrix[5] = stateMatrix[6]
	newStateMatrix[6] = stateMatrix[7]
	newStateMatrix[7] = stateMatrix[4]

	//Offset 2
	newStateMatrix[8] = stateMatrix[10]
	newStateMatrix[9] = stateMatrix[11]
	newStateMatrix[10] = stateMatrix[8]
	newStateMatrix[11] = stateMatrix[9]

	//Offset 3
	newStateMatrix[12] = stateMatrix[15]
	newStateMatrix[13] = stateMatrix[12]
	newStateMatrix[14] = stateMatrix[13]
	newStateMatrix[15] = stateMatrix[14]

	//Converting Row Major back to Col Major
	var temp [16]byte = newStateMatrix
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			newStateMatrix[i*4+j] = temp[j*4+i]
		}
	}

	return newStateMatrix
}

func MixColumns(stateMatrix [16]byte) [16]byte {
	var newStateMatrix [16]byte = stateMatrix

	var M [16]byte = [16]byte{0x02, 0x03, 0x01, 0x01, 0x01, 0x02, 0x03, 0x01, 0x01, 0x01, 0x02, 0x03, 0x03, 0x01, 0x01, 0x02}

	//Converting Col Major into Row Major
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			stateMatrix[i*4+j] = newStateMatrix[j*4+i]
		}
	}

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			newStateMatrix[i*4+j] = 0
			for k := 0; k < 4; k++ {
				newStateMatrix[i*4+j] ^= MultiplicationWithOverflowCheck(M[i*4+k], stateMatrix[k*4+j])
			}
		}
	}

	//Converting Row Major back to Col Major
	var temp [16]byte = newStateMatrix
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			newStateMatrix[i*4+j] = temp[j*4+i]
		}
	}

	return newStateMatrix
}

func Encrypt(stateMatrix [16]byte, rounds int, roundKeys [][16]byte) [16]byte {
	var tempStateMatrix [16]byte

	//Round 0
	tempStateMatrix = AddRoundKey(stateMatrix, roundKeys[0])

	//Rest Rounds
	for i := 1; i <= rounds; i++ {

		//Substituting Bytes
		tempStateMatrix = BytesSubstitution(tempStateMatrix)

		//Shifting Rows
		tempStateMatrix = ShiftRows(tempStateMatrix)

		//Except Last Round
		if i != rounds {
			//Mixing Columns
			tempStateMatrix = MixColumns(tempStateMatrix)
		}

		//Adding Round Key
		tempStateMatrix = AddRoundKey(tempStateMatrix, roundKeys[i])
	}

	return tempStateMatrix
}
