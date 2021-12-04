# E2E-Chat-App
End to End Encrypted Private Chat Application implemented in Golang

## Key Features
1. It uses RSA-3072 bit algorithm for Symmetric Key Exchange.
2. It uses AES-128 bit ECB Mode for Secure Communication after AES Key Exchange

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