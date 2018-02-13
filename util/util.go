package util

import (
	"errors"
	"math/rand"
	"os"
	"regexp"
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

//// REGEX ////
var (
	BTCRegex    = regexp.MustCompile("^[13][a-km-zA-HJ-NP-Z0-9]{26,33}$")
	emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

///////////////

func RandSeq(n int, seq []rune) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = seq[rand.Intn(len(seq))]
	}
	return string(b)
}

func ValidateEmail(email string) error {
	if !emailRegexp.MatchString(email) {
		return errors.New("Invalid email format")
	}
	return nil
}
