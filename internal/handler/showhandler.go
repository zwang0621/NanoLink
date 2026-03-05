// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"net/http"

	"shortener/internal/logic"
	"shortener/internal/svc"
	"shortener/internal/types"

	"github.com/go-playground/validator/v10"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ShowHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ShowRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		//参数规则校验也写到这一层
		if err := validator.New().StructCtx(r.Context(), req); err != nil {
			logx.Errorw("validator check failed", logx.LogField{Key: "error", Value: err.Error()})
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewShowLogic(r.Context(), svcCtx)
		resp, err := l.Show(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			// 为便于在 Nginx access.log 中统计短链访问情况，这里通过响应头把短链和长链写出去
			// 之后可以在 Nginx 的 log_format 中使用 $sent_http_x_short_url / $sent_http_x_long_url 等变量进行统计
			w.Header().Set("X-Short-Url", req.ShortURL)
			w.Header().Set("X-Long-Url", resp.LongURL)
			w.Header().Set("X-Shortener-Service", "shortener-api")

			//httpx.OkJsonCtx(r.Context(), w, resp)
			http.Redirect(w, r, resp.LongURL, http.StatusFound)
		}
	}
}
