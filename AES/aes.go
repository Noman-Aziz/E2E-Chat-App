package AES

import (
	"math"
	"strings"

	"github.com/Noman-Aziz/E2E-Chat-App/AES/utils"
)

func Initialization() Key {

	//Initial Variables
	var keys Key

	keys.Rounds = 10

	//Fixing Sizes for Dynamic Array
	keys.RoundKeys = make([][16]byte, keys.Rounds+1)

	//Generating Random 16 Bytes Key Input
	buffer := utils.RandomKeyGenerator(16, 4, 4, 4)

	//Read initial key
	for i := 0; i < 16; i++ {
		keys.RoundKeys[0][i] = buffer[i]
	}

	for i := 0; i < keys.Rounds; i++ {
		keys.RoundKeys[i+1] = utils.GenerateRoundKeys(keys.RoundKeys[i], i)
	}

	return keys
}

func Encryption(keys Key, message string) [][16]byte {

	var plainText PlainText
	var CipherTexts [][16]byte

	plainText.Text = []byte(message)

	//Selecting and Displaying Padding Character
	plainText.PaddingCharacter = '`'

	//Determining Length of Plain Text and Allocating Memory Accordingly
	var temp float64 = float64(len(plainText.Text)) / 16.0
	plainText.NumChunks = int(math.Ceil(temp))
	if plainText.NumChunks == 0 {
		plainText.NumChunks = 1
	}

	plainText.StateMatrix = make([][16]byte, plainText.NumChunks)
	CipherTexts = make([][16]byte, plainText.NumChunks)

	//Seperating the chunks from PlainText
	var index int = 0
	for i := 0; i < plainText.NumChunks; i++ {
		for j := 0; j < 16; j++ {
			//Padding
			if index >= len(plainText.Text) {
				plainText.StateMatrix[i][j] = plainText.PaddingCharacter
			} else {
				plainText.StateMatrix[i][j] = plainText.Text[index]
				index++
			}
		}
	}

	//Perform Encryption in ECB Mode
	for i := 0; i < plainText.NumChunks; i++ {
		CipherTexts[i] = utils.Encrypt(plainText.StateMatrix[i], keys.Rounds, keys.RoundKeys)
	}

	return CipherTexts
}

func Decryption(keys Key, CipherTexts [][16]byte) string {
	var plainText PlainText

	//Perform Decryption in ECB Mode

	for i := 0; i < len(CipherTexts); i++ {
		temp := utils.Decrypt(CipherTexts[i], keys.Rounds, keys.RoundKeys)

		for j := 0; j < 16; j++ {
			plainText.Text = append(plainText.Text, temp[j])
		}
	}

	//Removing Padding and Sending Text Back
	return strings.Trim(string(plainText.Text), "`")
}
