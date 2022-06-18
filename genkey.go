package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
)

func genkeyCommand() {
	keyb := make([]byte, 32)
	_, err := rand.Read(keyb)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err.Error())
		os.Exit(1)
	}

	fmt.Printf("%s\n", hex.EncodeToString(keyb))
}
