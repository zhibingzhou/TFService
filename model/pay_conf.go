package model

import (
	"fmt"

	"github.com/zhibingzhou/go_public/common"
	"github.com/zhibingzhou/go_public/redis"
)

func PayConf(pay_id string) PayConfig {
	var api_conf PayConfig
	gdb.DB.Where("id=?", pay_id).First(&api_conf)
	return api_conf
}

func ApiConfigRedis(pay_id string) map[string]string {
	redis_key := fmt.Sprintf("pay_config:%s", pay_id)
	//优先查询redis
	pay_map := redis.RediGo.HgetAll(redis_key)
	if len(pay_map["api_conf"]) < 1 {
		conf_info := PayConf(pay_id)
		if len(conf_info.Api_conf) > 0 {
			pay_map = common.StructToMapSlow(conf_info)
			redis.RediGo.Hmset(redis_key, pay_map, redis_max_time)
			redis.RediGo.Sadd(Conf_Redis_Key, redis_key, 0)
		}
	}
	return pay_map
}
