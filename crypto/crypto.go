package crypto

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func Hash(origin, salt []byte) []byte {
	h := hmac.New(sha256.New, salt)
	h.Write(origin)
	data := h.Sum(nil)

	return data
}

func HashString(origin, salt string) string {
	data := Hash([]byte(origin), []byte(salt))
	return hex.EncodeToString(data)
}

func SaltGen(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		fmt.Println(err)
	}
	return base64.RawStdEncoding.EncodeToString(b)
}
