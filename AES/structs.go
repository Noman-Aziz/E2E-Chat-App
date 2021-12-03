package AES

type Key struct {
	Rounds    int
	RoundKeys [][16]byte
}

type PlainText struct {
	StateMatrix      [][16]byte
	NumChunks        int
	PaddingCharacter byte
	Text             []byte
}
