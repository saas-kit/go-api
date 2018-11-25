package random

import (
	"math/rand"
	"strings"
	"time"
)

type (
	// Random structure
	Random struct{}
)

// Charsets
const (
	Uppercase    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Lowercase    = "abcdefghijklmnopqrstuvwxyz"
	Alphabetic   = Uppercase + Lowercase
	Numeric      = "0123456789"
	Alphanumeric = Alphabetic + Numeric
	Symbols      = "`" + `~!@#$%^&*()-_+={}[]|\;:"<>,./?`
	Hex          = Numeric + "abcdef"
)

var (
	global = New()
)

// New Random structure
func New() *Random {
	rand.Seed(time.Now().UnixNano())
	return new(Random)
}

// String function
func (r *Random) String(length uint8, charsets ...string) string {
	charset := strings.Join(charsets, "")
	if charset == "" {
		charset = Alphanumeric
	}
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Int63()%int64(len(charset))]
	}
	return string(b)
}

// String is a helper to get random string
func String(length uint8, charsets ...string) string {
	return global.String(length, charsets...)
}
