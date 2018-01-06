package url

import (
	"net/url"
	"os"
)

func Get(path string) *url.URL {
	u, err := url.Parse(os.Getenv("BASE_URL"))
	if err != nil {
		return nil
	}
	u.Path = path
	return u
}
