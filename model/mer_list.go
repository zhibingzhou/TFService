package model

import (
	"fmt"

	"github.com/zhibingzhou/go_public/common"
	"github.com/zhibingzhou/go_public/redis"
)

func MerInfo(mer_code string) MerList {
	var m_info MerList
	gdb.DB.Where("code=?", mer_code).First(&m_info)
	return m_info
}

func MerInfoRedis(mer_code string) map[string]string {
	redis_key := fmt.Sprintf("mer_list:%s", mer_code)
	//优先查询redis
	m_map := redis.RediGo.HgetAll(redis_key)
	if len(m_map["code"]) < 1 {
		m_info := MerInfo(mer_code)
		if len(m_info.Code) > 0 {
			m_map = common.StructToMapSlow(m_info)
			redis.RediGo.Hmset(redis_key, m_map, redis_max_time)
			redis.RediGo.Sadd(Data_Redis_Key, redis_key, 0)
		}
	}
	return m_map
}

func UpdatesMer(m_info MerList, p_data map[string]interface{}) error {
	res := gdb.DB.Model(&m_info).UpdateColumns(p_data)
	return res.Error
}
