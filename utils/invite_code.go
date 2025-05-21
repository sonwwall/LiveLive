package utils

import (
	"math/rand"
	"time"
)

const inviteCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func GenerateInviteCode(length int) string {
	code := make([]byte, length)
	for i := range code {
		code[i] = inviteCharset[seededRand.Intn(len(inviteCharset))]
	}
	return string(code)
}
