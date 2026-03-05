// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"shortener/internal/svc"
	"shortener/internal/types"
	"shortener/model"
	"shortener/pkg/base62"
	"shortener/pkg/connect"
	"shortener/pkg/md5"
	"shortener/pkg/url"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ConvertLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConvertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConvertLogic {
	return &ConvertLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConvertLogic) Convert(req *types.ConvertRequest) (resp *types.ConvertResponse, err error) {
	// 1.参数校验
	// 1.1 判断是否为空（已在handler中完成，使用第三方库validator实现）
	// 1.2 判断请求的链接是否能正常访问
	if connect.Get(req.LongURL) == false {
		return nil, errors.New("无效的连接")
	}

	// 1.3 判断请求的链接是否已经被转为短链接（得用MD5值去查，而不是直接拿长连接去查）
	md5Value := md5.Sum([]byte(req.LongURL))
	u, err := l.svcCtx.ShortURLModel.FindOneByMd5(l.ctx, sql.NullString{String: md5Value, Valid: true})
	if !errors.Is(err, sqlx.ErrNotFound) {
		if err == nil {
			return nil, fmt.Errorf("该链接已经被转化为%s", u.Surl.String)
		}
		logx.Errorw("ShortURLModel.FindOneByMd5 failed", logx.LogField{
			Key:   "error",
			Value: err.Error(),
		})
		return nil, err
	}

	// 1.4 判断请求的链接是否已经是短链接
	basePath, err := url.GetBasePath(req.LongURL)
	if err != nil {
		logx.Errorw("GetBasePath failed", logx.LogField{
			Key:   "lurl",
			Value: req.LongURL,
		}, logx.LogField{
			Key:   "error",
			Value: err.Error(),
		})
		return nil, err
	}

	u, err = l.svcCtx.ShortURLModel.FindOneBySurl(l.ctx, sql.NullString{
		String: basePath,
		Valid:  true,
	})
	if !errors.Is(err, sqlx.ErrNotFound) {
		if err == nil {
			return nil, errors.New("该链接已经被转化为短链接了")
		}
		logx.Errorw("ShortURLModel.FindOneBySurl failed", logx.LogField{
			Key:   "error",
			Value: err.Error(),
		})
		return nil, err
	}

	var short string
	for {
		// 2.取号器取号（用mysql自增主键配合REPLACE INTO实现）
		seq, err := l.svcCtx.Sequence.Next()
		if err != nil {
			return nil, err
		}

		// 3.转短链接（62进制转换实现）
		// 3.1 考虑安全性，避免被人恶意请求
		// 3.2 避免生成特殊的词，比如version，stupid
		short = base62.Base62Encode(seq)
		if _, ok := l.svcCtx.ShortURLBlackList[short]; !ok {
			break //生成不在黑名单的short就跳出循环
		}
	}

	// 4.短链接记录到redis版本布隆过滤器中
	err = l.svcCtx.Filter.Add([]byte(short))
	if err != nil {
		logx.Errorw("Filter.Add failed", logx.LogField{
			Key:   "error",
			Value: err.Error(),
		})
		return nil, err
	}

	// 4.长短连接映射关系记录到数据库表中
	if _, err := l.svcCtx.ShortURLModel.Insert(
		l.ctx,
		&model.ShortUrlMap{
			Url:  sql.NullString{String: req.LongURL, Valid: true},
			Md5:  sql.NullString{String: md5Value, Valid: true},
			Surl: sql.NullString{String: short, Valid: true},
		}); err != nil {
		logx.Errorw("ShortURLModel.Insert failed", logx.LogField{
			Key:   "error",
			Value: err.Error(),
		})
		return nil, err
	}

	// 5.返回响应
	shortURL := l.svcCtx.Config.ShortDomain + "/" + short
	return &types.ConvertResponse{
		ShortURL: shortURL,
	}, nil
}
