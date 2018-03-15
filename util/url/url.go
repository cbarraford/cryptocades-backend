package url

import (
	"fmt"
	"net/url"
	"os"
)

func Get(path string) *url.URL {
	u, err := url.Parse(fmt.Sprintf("%s%s", os.Getenv("BASE_URL"), path))
	if err != nil {
		return nil
	}
	return u
}
