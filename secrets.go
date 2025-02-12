package main

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"time"
)

// generate secret for github
func generateSecret() string {
	// generate new password hash
	byteArray := bytes.Join(
		[][]byte{[]byte(time.Now().GoString())},
		[]byte("github&*^$%"),
	) // generate random complex byte array
	secret := fmt.Sprintf("%x", sha1.Sum(byteArray)) // complex byte hash
	return secret[0:15]                              // Make secret shorter
}
