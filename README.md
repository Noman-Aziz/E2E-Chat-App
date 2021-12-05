# E2E-Chat-App
End to End Encrypted Private Chat Application implemented in Golang

## Features
1. It uses RSA-3072 bit algorithm for Symmetric Key Exchange.
2. It uses AES-128 bit ECB Mode for Secure Communication after AES Key Exchange
3. It uses SHA-512 algorithm for Hashing

## Working
1. Rsa Public Key Exchange Happens
2. Sender Authentication and Non Repudiation is done using RSA Signature
3. AES Key Exchange happens using RSA
4. All communication is done from here using AES
5. Message integrity is also checked at each communication.
6. Base64 encoding is used for after encrypting message with AES to avoid any errors from the server end.

## To Run
### 1. Server
```
go run ./Chat/server/cmd/main.go
```
### 2. Client 1
```
go run ./Chat/tui/cmd/main.go -server :3333 -first true
```

### 3. Client 2
```
go run ./Chat/tui/cmd/main.go -server :3333
```