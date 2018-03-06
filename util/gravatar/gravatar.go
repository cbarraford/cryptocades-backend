package gravatar

// Code copied from https://github.com/zoonman/gravatar

import (
	"crypto/md5"
	"fmt"
	"log"
	"net/url"
	"strconv"

	_url "github.com/cbarraford/cryptocades-backend/util/url"
)

// Create Gravatar picture URL from email address
func Avatar(email string, size uint) string {
	hash := md5.Sum([]byte(email))
	sum := 0
	for _, i := range hash {
		sum += int(i)
	}
	// pick default image if gravatar doesn't have one
	defaultImg := _url.Get(fmt.Sprintf("/img/avatars/%d.png", sum%12))
	log.Printf("Default Image: %s", defaultImg)
	u, _ := url.Parse(fmt.Sprintf("https://www.gravatar.com/avatar/%x", hash))
	vals := url.Values{}
	vals.Set("s", strconv.FormatUint(uint64(size), 10))
	vals.Set("d", defaultImg.String())
	u.RawQuery = vals.Encode()

	return u.String()
}
