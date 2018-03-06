package gravatar

// Code copied from https://github.com/zoonman/gravatar

import (
	"crypto/md5"
	"fmt"
	"net/url"
	"strconv"
)

// Create Gravatar picture URL from email address
func Avatar(email string, size uint) string {
	hash := md5.Sum([]byte(email))
	sum := 0
	for _, i := range hash {
		sum += int(i)
	}
	// pick default image if gravatar doesn't have one
	defaultImg := fmt.Sprintf("/img/avatars/%d.png", sum%12)
	u, _ := url.Parse(fmt.Sprintf("https://www.gravatar.com/avatar/%x", hash))
	vals := url.Values{}
	vals.Set("s", strconv.FormatUint(uint64(size), 10))
	vals.Set("d", defaultImg)
	u.RawQuery = vals.Encode()

	return u.String()
}
