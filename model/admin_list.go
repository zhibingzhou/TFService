package model

import (
	"github.com/zhibingzhou/go_public/common"
	"github.com/zhibingzhou/go_public/redis"
)

func AdminInfo(account string) AdminList {
	var a_info AdminList
	gdb.DB.Where("account=?", account).First(&a_info)
	return a_info
}

func AdminInfoRedis(account string) map[string]string {
	redis_key := "admin_list:" + account
	//优先查询redis
	a_map := redis.RediGo.HgetAll(redis_key)
	if len(a_map["account"]) < 1 {
		a_info := AdminInfo(account)
		if len(a_info.Account) > 0 {
			a_map = common.StructToMapSlow(a_info)
			redis.RediGo.Hmset(redis_key, a_map, redis_data_time)
			redis.RediGo.Sadd(Data_Redis_Key, redis_key, 0)
		}
	}
	return a_map
}

func AdminBySession(session_id string) AdminList {
	var d_list AdminList
	gdb.DB.Where("session_id=?", session_id).First(&d_list)
	return d_list
}

func AdminBySessionRedis(session_id string) map[string]string {
	redis_key := "admin_list:session_id:" + session_id
	//优先查询redis
	a_map := redis.RediGo.HgetAll(redis_key)
	if len(a_map["account"]) < 1 {
		a_info := AdminBySession(session_id)
		if len(a_info.Account) > 0 {
			a_map = common.StructToMapSlow(a_info)
			redis.RediGo.Hmset(redis_key, a_map, redis_data_time)
			redis.RediGo.Sadd(Data_Redis_Key, redis_key, 0)
		}
	}
	return a_map
}

func UpdatesAdminList(a_list AdminList, a_data map[string]interface{}) error {
	res := gdb.DB.Model(&a_list).UpdateColumns(a_data)
	return res.Error
}
