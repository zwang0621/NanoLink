package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// HealthHandler 提供给 Nginx / 负载均衡做健康检查使用
// 仅返回 200 OK，便于在 access.log 中区分业务访问与健康检查流量
func HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 标记为健康检查请求，方便 Nginx access.log 单独统计或过滤
		w.Header().Set("X-Shortener-Service", "shortener-api")
		w.Header().Set("X-Shortener-Health", "1")

		httpx.OkJsonCtx(r.Context(), w, map[string]string{
			"status": "ok",
		})
	}
}

