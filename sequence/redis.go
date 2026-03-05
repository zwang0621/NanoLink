package sequence

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
)

//基于redis实现取号器

type Redis struct {
	//redis连接
	conn *redis.Redis
}

func NewRedis(redisAddr string) Sequence {
	r := redis.MustNewRedis(redis.RedisConf{
		Host: redisAddr, // "127.0.0.1:6379"
		Type: "node",    //单机模式
	})
	return &Redis{
		conn: r,
	}
}

func (r *Redis) Next() (u uint64, err error) {
	//使用redis实现取号器的思路
	//INCR原子递增
	id, err := r.conn.Incr("sequence:global")
	if err != nil {
		return
	}
	return uint64(id), nil
}
