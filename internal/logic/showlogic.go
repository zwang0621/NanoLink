// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"database/sql"
	"errors"
	"shortener/internal/svc"
	"shortener/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShowLogic {
	return &ShowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShowLogic) Show(req *types.ShowRequest) (resp *types.ShowResponse, err error) {
	// 为避免缓存穿透问题，加布隆过滤器
	ok, err := l.svcCtx.Filter.Exists([]byte(req.ShortURL))
	if err != nil {
		logx.Errorw("Filter.Exists failed", logx.LogField{
			Key:   "error",
			Value: err.Error(),
		})
	}
	if !ok {
		return nil, errors.New("404布隆过滤器中不存在该缓存")
	}
	logx.Infof("bloom check short=%s, ok=%v", req.ShortURL, ok)

	// 查看短链接对应的长连接,查数据库之前可以加缓存
	// go-zero的缓存本身就支持singleflight
	u, err := l.svcCtx.ShortURLModel.FindOneBySurl(l.ctx, sql.NullString{
		String: req.ShortURL,
		Valid:  true,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("404")
		}
		logx.Errorw("ShortURLModel.FindOneBySurl failed", logx.LogField{
			Key:   "error",
			Value: err.Error(),
		})
		return nil, err
	}
	// 返回查询到的长连接，在handler层返回重定向的响应
	return &types.ShowResponse{LongURL: u.Url.String}, nil
}
