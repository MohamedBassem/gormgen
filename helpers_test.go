package gormgen

import (
	"bytes"
	"math/rand"
)

func generateRandomString(l int) string {
	seed := []rune("abcdefghijklmnobqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	var buf bytes.Buffer
	for i := 0; i < l; i++ {
		buf.WriteRune(seed[rand.Intn(len(seed))])
	}
	return buf.String()
}
