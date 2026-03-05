package url

import (
	"errors"
	"net/url"
	"path"
)

// GetBasePath 获取url路径的最后一节
func GetBasePath(targetUrl string) (string, error) {
	myUrl, err := url.Parse(targetUrl)
	if err != nil {
		return "", err
	}
	if len(myUrl.Host) == 0 {
		return "", errors.New("no host")
	}
	return path.Base(myUrl.Path), nil
}
