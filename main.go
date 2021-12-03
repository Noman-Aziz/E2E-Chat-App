package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Noman-Aziz/E2E-Chat-App/RSA"
)

func main() {

	rand.Seed(time.Now().UTC().UnixNano())

	var pubKey RSA.PubKey
	var privKey RSA.PrivKey

	pubKey, privKey = RSA.Initialization(3072)

	fmt.Println("Pub Key", pubKey, "\n")
	fmt.Println("Priv Key", privKey, "\n")

	fmt.Println(RSA.Decrypt(privKey, RSA.Encrypt(pubKey, "HELLO")))
}
