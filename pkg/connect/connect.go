package connect

import (
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// client 全局的http客户端
var client = &http.Client{
	Transport: &http.Transport{
		DisableKeepAlives: true,
	},
	Timeout: time.Second * 2,
}

// Get 判断url能否请求通
func Get(url string) bool {
	resp, err := client.Get(url)
	if err != nil {
		logx.Errorw("connect client.Get failed", logx.LogField{Key: "err", Value: err.Error()})
		return false
	}
	err = resp.Body.Close()
	if err != nil {
		logx.Errorw("connect client.Get Close failed", logx.LogField{Key: "err", Value: err.Error()})
		return false
	}
	return resp.StatusCode == 200 //别人给我发一个跳转的响应也不算过
}
