package util

import (
	"math/rand"
	"os"
	"time"

	"github.com/stvp/rollbar"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	rollbar.Token = os.Getenv("ROLLBAR_ACCESS_TOKEN")
	rollbar.Environment = os.Getenv("ENVIRONMENT")
}

var LowerLetters = []rune("abcdefghijklmnopqrstuvwxyz")
var Letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var AlphaNumeric = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var LowerAlphaNumeric = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

func RandSeq(n int, seq []rune) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = seq[rand.Intn(len(seq))]
	}
	return string(b)
}
