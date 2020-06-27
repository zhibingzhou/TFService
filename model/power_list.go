package model

import (
	"fmt"

	"github.com/zhibingzhou/go_public/common"
	"github.com/zhibingzhou/go_public/redis"
)

func PowerInfo(code string) PowerList {
	var m_info PowerList
	gdb.DB.Where("code=?", code).First(&m_info)
	return m_info
}

func PowerInfoRedis(code string) map[string]string {
	redis_key := fmt.Sprintf("power_list:%s", code)
	//优先查询redis
	m_map := redis.RediGo.HgetAll(redis_key)
	if len(m_map["code"]) < 1 {
		m_info := PowerInfo(code)
		if len(m_info.Code) > 0 {
			m_map = common.StructToMapSlow(m_info)
			redis.RediGo.Hmset(redis_key, m_map, redis_max_time)
			redis.RediGo.Sadd(Data_Redis_Key, redis_key, 0)
		}
	}
	return m_map
}

func PowerByUrl(url string) PowerList {
	var m_info PowerList
	gdb.DB.Where("url=?", url).First(&m_info)
	return m_info
}

func PowerByUrlRedis(url string) map[string]string {
	redis_key := fmt.Sprintf("power_list:url:%s", url)
	//优先查询redis
	m_map := redis.RediGo.HgetAll(redis_key)
	if len(m_map["code"]) < 1 {
		m_info := PowerByUrl(url)
		if len(m_info.Code) > 0 {
			m_map = common.StructToMapSlow(m_info)
			redis.RediGo.Hmset(redis_key, m_map, redis_max_time)
			redis.RediGo.Sadd(Data_Redis_Key, redis_key, 0)
		}
	}
	return m_map
}
