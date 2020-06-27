package model

import (
	"github.com/zhibingzhou/go_public/common"
	"github.com/zhibingzhou/go_public/redis"
)

func UserPower(account, p_code string) AdminPower {
	var a_info AdminPower
	gdb.DB.Where("account=? and power_code=?", account, p_code).First(&a_info)
	return a_info
}

func UserPowerRedis(account, p_code string) map[string]string {
	redis_key := "admin_list:" + account + ":" + p_code
	//优先查询redis
	a_map := redis.RediGo.HgetAll(redis_key)
	if len(a_map["account"]) < 1 {
		a_info := UserPower(account, p_code)
		if a_info.Id > 0 {
			a_map = common.StructToMapSlow(a_info)
			redis.RediGo.Hmset(redis_key, a_map, redis_data_time)
			redis.RediGo.Sadd(Data_Redis_Key, redis_key, 0)
		}
	}
	return a_map
}
