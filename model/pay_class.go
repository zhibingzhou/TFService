package model

import (
	"fmt"

	"github.com/zhibingzhou/go_public/common"
	"github.com/zhibingzhou/go_public/redis"
)

func ClassInfo(class_code string) PayClass {
	var m_info PayClass
	gdb.DB.Where("code=?", class_code).First(&m_info)
	return m_info
}

func ClassInfoRedis(class_code string) map[string]string {
	redis_key := fmt.Sprintf("pay_class:%s", class_code)
	//优先查询redis
	m_map := redis.RediGo.HgetAll(redis_key)
	if len(m_map["code"]) < 1 {
		m_info := ClassInfo(class_code)
		if len(m_info.Code) > 0 {
			m_map = common.StructToMapSlow(m_info)
			redis.RediGo.Hmset(redis_key, m_map, redis_max_time)
			redis.RediGo.Sadd(Data_Redis_Key, redis_key, 0)
		}
	}
	return m_map
}
