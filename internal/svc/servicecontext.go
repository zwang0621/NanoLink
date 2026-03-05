// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"shortener/internal/config"
	"shortener/model"
	"shortener/sequence"

	"github.com/zeromicro/go-zero/core/bloom"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config        config.Config
	ShortURLModel model.ShortUrlMapModel
	Sequence      sequence.Sequence

	ShortURLBlackList map[string]struct{}
	Filter            *bloom.Filter
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.ShortURLDB.DSN)

	// 把配置文件中的黑名单切片加载到map中，方便后续直接定位
	m := make(map[string]struct{}, len(c.ShortUrlBlackList))
	for _, v := range c.ShortUrlBlackList {
		m[v] = struct{}{}
	}

	// 初始化布隆过滤器
	store := redis.MustNewRedis(redis.RedisConf{
		Host: c.CacheRedis[0].Host,
		Type: redis.NodeType,
	})
	filter := bloom.New(store, "bloom_filter", 20*(1<<20))

	return &ServiceContext{
		Config:        c,
		ShortURLModel: model.NewShortUrlMapModel(conn, c.CacheRedis),
		Sequence:      sequence.NewMySQL(c.Sequence.DSN),
		//Sequence: sequence.NewRedis(c.Sequence.RedisAddr),
		ShortURLBlackList: m, //短链接黑名单map
		Filter:            filter,
	}
}
