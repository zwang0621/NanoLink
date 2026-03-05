// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	ShortURLDB struct {
		DSN string
	}

	CacheRedis cache.CacheConf

	Sequence struct {
		DSN string
		//RedisAddr string
	}

	BaseString string

	ShortUrlBlackList []string

	ShortDomain string
}
