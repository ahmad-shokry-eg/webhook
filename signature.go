package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// verify github signature
func verifySignature(body []byte, signature string, secret string) bool {
	if signature == "" {
		return false
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expectedMAC := mac.Sum(nil)
	expectedSignature := "sha256=" + hex.EncodeToString(expectedMAC)

	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
