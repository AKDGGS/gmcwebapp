package model

import (
	"crypto/rand"
	"math/big"
	"strings"
)

const dict = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890$%*+-./:"

type Token struct {
	ID          int
	Token       string
	Description string
}

func (tk *Token) Generate() {
	var t strings.Builder
	for i := 0; i < 128; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(dict))))
		t.WriteByte(dict[int(n.Int64())])
	}
	tk.Token = t.String()
}
