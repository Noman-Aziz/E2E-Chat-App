package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Noman-Aziz/E2E-Chat-App/AES"
	"github.com/Noman-Aziz/E2E-Chat-App/RSA"
)

func main() {

	rand.Seed(time.Now().UTC().UnixNano())

	var rsaPubKey RSA.PubKey
	var rsaPrivKey RSA.PrivKey

	rsaPubKey, rsaPrivKey = RSA.Initialization(3072)

	fmt.Println("Pub Key", rsaPubKey, "\n")
	fmt.Println("Priv Key", rsaPrivKey, "\n")

	var aesKey AES.Key = AES.Initialization(false, "")

	fmt.Println("Generating Random AES Key and Doing Key Change using RSA\n")

	encAesKey := RSA.Encrypt(rsaPubKey, string(aesKey.RoundKeys[0][:]))
	decAesKey := RSA.Decrypt(rsaPrivKey, encAesKey)

	fmt.Println("AES Key", decAesKey, "\n")

	cipherText := AES.Encryption(aesKey, "HELLO")

	plainText := AES.Decryption(aesKey, cipherText)

	fmt.Println("AES Decrypted Text", plainText)
}
