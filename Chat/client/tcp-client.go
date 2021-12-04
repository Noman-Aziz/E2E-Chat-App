package client

import (
	"bytes"
	"io"
	"log"
	"math/big"
	"math/rand"
	"net"
	"strings"
	"time"

	"github.com/Noman-Aziz/E2E-Chat-App/AES"
	"github.com/Noman-Aziz/E2E-Chat-App/Chat/protocol"
	"github.com/Noman-Aziz/E2E-Chat-App/RSA"
)

type TcpChatClient struct {
	conn      net.Conn
	cmdReader *protocol.CommandReader
	cmdWriter *protocol.CommandWriter
	name      string
	incoming  chan protocol.MessageCommand
	first     bool
}

var (
	MyRsaPubKey        RSA.PubKey
	MyRsaPrivKey       RSA.PrivKey
	OtherPeerRsaPubKey RSA.PubKey
	AesKey             AES.Key
	AesKeyExchanged    bool
	RsaKeySent         bool
	RsaKeyReceived     bool
	RsaKeyAcknowledged bool
)

func NewClient(first bool) *TcpChatClient {
	temp := &TcpChatClient{
		incoming: make(chan protocol.MessageCommand),
	}

	temp.first = first

	return temp
}

func (c *TcpChatClient) Dial(address string) error {
	conn, err := net.Dial("tcp", address)

	if err == nil {
		c.conn = conn
	}

	c.cmdReader = protocol.NewCommandReader(conn)
	c.cmdWriter = protocol.NewCommandWriter(conn)

	return err
}

func (c *TcpChatClient) Start() {

	//Generating Self RSA Key Pair
	rand.Seed(time.Now().UTC().UnixNano())
	MyRsaPubKey, MyRsaPrivKey = RSA.Initialization(3072)

	time.Sleep(2 * time.Second)

	//If it is the First Client, it will generate AES Key
	if c.first {
		AesKey = AES.Initialization(false, "")
	}

	AesKeyExchanged = false

	RsaKeySent = false
	RsaKeyReceived = false
	RsaKeyAcknowledged = false

	for {
		cmd, err := c.cmdReader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("Read error %v", err)
		}

		if cmd != nil {
			switch v := cmd.(type) {
			case protocol.MessageCommand:

				//Server has asked C1 to start sharing RSA Key
				if v.Message == "Exchange RSA Public Keys" && c.first {
					//Showing Received Message
					c.incoming <- v

					//Sending Public Key
					temp := "RSAKEYINCOMING,"
					temp = temp + MyRsaPubKey.E.String() + "," + MyRsaPubKey.N.String()
					c.SendMessage(temp)

					RsaKeySent = true

				} else if strings.Contains(v.Message, "RSAKEYINCOMING") && !c.first && !RsaKeySent && !RsaKeyReceived {
					//C2 Received C1 RSA Key
					temp := strings.Split(v.Message, ",")

					//Allocating Memory
					OtherPeerRsaPubKey.E = new(big.Int)
					OtherPeerRsaPubKey.N = new(big.Int)

					//Assigning Values
					OtherPeerRsaPubKey.E.SetString(temp[1], 10)
					OtherPeerRsaPubKey.N.SetString(temp[2], 10)

					//Changing Message Content
					v.Message = "RSA Public Key Received"

					//Showing Received Message
					c.incoming <- v

					//C2 Send Ack to C1
					c.SendMessage("RSAKEYACK")

					//Sleep for 1s before sending another message
					time.Sleep(1 * time.Second)

					//C2 Send its RSA Key to C1
					temp1 := "RSAKEYINCOMING,"
					temp1 = temp1 + MyRsaPubKey.E.String() + "," + MyRsaPubKey.N.String()
					c.SendMessage(temp1)

					RsaKeySent = true
					RsaKeyReceived = true

				} else if strings.Contains(v.Message, "RSAKEYACK") && c.first && !RsaKeyAcknowledged {
					//C1 Receives Ack of C2
					RsaKeyAcknowledged = true
				} else if strings.Contains(v.Message, "RSAKEYINCOMING") && c.first && RsaKeySent && RsaKeyAcknowledged && !RsaKeyReceived {
					//C1 Received C2 RSA Key
					temp := strings.Split(v.Message, ",")

					//Allocating Memory
					OtherPeerRsaPubKey.E = new(big.Int)
					OtherPeerRsaPubKey.N = new(big.Int)

					//Assigning Values
					OtherPeerRsaPubKey.E.SetString(temp[1], 10)
					OtherPeerRsaPubKey.N.SetString(temp[2], 10)

					//Changing Message Content
					v.Message = "RSA Public Key Received"

					//Showing Received Message
					c.incoming <- v

					//Send ACK Response to C2
					c.SendMessage("RSAKEYACK")

					RsaKeyReceived = true
				} else if strings.Contains(v.Message, "RSAKEYACK") && !c.first && !RsaKeyAcknowledged {
					//C2 Receives Ack of C1
					RsaKeyAcknowledged = true
				} else if c.first && RsaKeyAcknowledged && !AesKeyExchanged {
					//C1 Send AES Key to C2
					temp := string(AesKey.RoundKeys[0][:])

					encAesKey := RSA.Encrypt(OtherPeerRsaPubKey, temp)

					//Converting Big Int Array to String
					temp = "AESKEYINCOMING"
					for _, val := range encAesKey {
						temp = temp + "," + val.String()
					}

					c.SendMessage(temp)

					AesKeyExchanged = true
				} else if strings.Contains(v.Message, "AESKEYINCOMING") && !c.first && !AesKeyExchanged {
					//C2 Receives AES Key from C1
					temp := strings.Split(v.Message, "AESKEYINCOMING")

					//COnverting Message to Apropiate Format
					temp = strings.Split(temp[1], ",")
					var encAesKey []big.Int
					for _, val := range temp {
						var crypt *big.Int = new(big.Int)
						crypt.SetString(val, 10)
						encAesKey = append(encAesKey, *crypt)
					}

					//Decrypting the AES Key with RSA
					decAesKey := RSA.Decrypt(MyRsaPrivKey, encAesKey)

					AesKey = AES.Initialization(true, decAesKey)

					//Changing Message Content
					v.Message = "AES Key Exchanged " + string(AesKey.RoundKeys[0][:])

					AesKeyExchanged = true

					//Showing Received Message
					c.incoming <- v
				} else if AesKeyExchanged && RsaKeyAcknowledged && !strings.Contains(v.Message, "AESKEYINCOMING") {
					//Normal Communication

					//Showing Received Message after Decrypting with AES Key
					var cipherText [][16]byte

					var dynamicStr []byte = bytes.NewBufferString(v.Message).Bytes()
					var staticStr [16]byte

					index := 0
					for i := 0; i < len(dynamicStr); i++ {
						staticStr[index] = dynamicStr[i]

						if index == 15 {
							cipherText = append(cipherText, staticStr)
							index = 0
							continue
						}

						index++
					}

					plainText := AES.Decryption(AesKey, cipherText)

					v.Message = plainText
					c.incoming <- v
				}

			default:
				log.Printf("Unknown command: %v", v)
			}
		}
	}
}

func (c *TcpChatClient) Close() {
	c.conn.Close()
}

func (c *TcpChatClient) Incoming() chan protocol.MessageCommand {

	return c.incoming
}

func (c *TcpChatClient) Send(command interface{}) error {
	return c.cmdWriter.Write(command)
}

func (c *TcpChatClient) SetName(name string) error {
	c.name = name
	return c.Send(protocol.NameCommand{name})
}

func (c *TcpChatClient) SendMessage(message string) error {

	//If AES Key have been exchanged, encrypt message and send
	if AesKeyExchanged {
		cipherText := AES.Encryption(AesKey, message)

		message = ""

		for i := 0; i < len(cipherText); i++ {
			message += bytes.NewBuffer(cipherText[i][:]).String()
		}
	}

	return c.Send(protocol.SendCommand{
		Message: message,
	})
}
