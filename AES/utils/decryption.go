package utils

func InverseBytesSubstitution(stateMatrix [16]byte) [16]byte {

	var newStateMatrix [16]byte = stateMatrix

	for i := 0; i < 16; i++ {
		firstIndex, secondIndex := ConvertToArrayIndex(newStateMatrix[i])
		//row * length + col
		newStateMatrix[i] = SboxInv[firstIndex*16+secondIndex]
	}

	return newStateMatrix
}

func InverseShiftRows(stateMatrix [16]byte) [16]byte {
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
	newStateMatrix[4] = stateMatrix[7]
	newStateMatrix[5] = stateMatrix[4]
	newStateMatrix[6] = stateMatrix[5]
	newStateMatrix[7] = stateMatrix[6]

	//Offset 2
	newStateMatrix[8] = stateMatrix[10]
	newStateMatrix[9] = stateMatrix[11]
	newStateMatrix[10] = stateMatrix[8]
	newStateMatrix[11] = stateMatrix[9]

	//Offset 3
	newStateMatrix[12] = stateMatrix[13]
	newStateMatrix[13] = stateMatrix[14]
	newStateMatrix[14] = stateMatrix[15]
	newStateMatrix[15] = stateMatrix[12]

	//Converting Row Major back to Col Major
	var temp [16]byte = newStateMatrix
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			newStateMatrix[i*4+j] = temp[j*4+i]
		}
	}

	return newStateMatrix
}

func InverseMixColumns(stateMatrix [16]byte) [16]byte {
	var newStateMatrix [16]byte = stateMatrix

	var MInv [16]byte = [16]byte{0x0e, 0x0b, 0x0d, 0x09, 0x09, 0x0e, 0x0b, 0x0d, 0x0d, 0x09, 0x0e, 0x0b, 0x0b, 0x0d, 0x09, 0x0e}

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

				if MInv[i*4+k] == 0x01 {
					newStateMatrix[i*4+j] ^= stateMatrix[k*4+j]
				} else if stateMatrix[k*4+j] == 0x01 {
					newStateMatrix[i*4+j] ^= MInv[i*4+k]
				} else if MInv[i*4+k] == 0x00 || stateMatrix[k*4+j] == 0x00 {
					newStateMatrix[i*4+j] ^= 0x00
				} else {
					i1, i2 := ConvertToArrayIndex(MInv[i*4+k])
					i3, i4 := ConvertToArrayIndex(stateMatrix[k*4+j])

					var temp uint16 = uint16(L[i1*16+i2]) + uint16(L[i3*16+i4])

					if temp > 0xff {
						temp -= 0xff
					}
					i1, i2 = ConvertToArrayIndex(byte(temp))
					newStateMatrix[i*4+j] ^= E[i1*16+i2]
				}

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

func Decrypt(stateMatrix [16]byte, rounds int, roundKeys [][16]byte) [16]byte {
	var tempStateMatrix [16]byte

	//Round 10
	tempStateMatrix = AddRoundKey(stateMatrix, roundKeys[rounds])

	//Rest Rounds
	for i := rounds - 1; i >= 0; i-- {

		//Inverse Shifting Rows
		tempStateMatrix = InverseShiftRows(tempStateMatrix)

		//Inverse Substituting Bytes
		tempStateMatrix = InverseBytesSubstitution(tempStateMatrix)

		//Adding Round Key
		tempStateMatrix = AddRoundKey(tempStateMatrix, roundKeys[i])

		//Except First Round
		if i != 0 {
			//Mixing Columns
			tempStateMatrix = InverseMixColumns(tempStateMatrix)
		}
	}

	return tempStateMatrix
}
