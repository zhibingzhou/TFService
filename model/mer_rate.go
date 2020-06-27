package model

import (
	"fmt"

	"github.com/zhibingzhou/go_public/common"
	"github.com/zhibingzhou/go_public/redis"
)

func MerRateInfo(mer_code, pay_code, class_code, bank_code string) MerRate {
	var m_info MerRate
	gdb.DB.Where("mer_code=? and pay_code=? and class_code=? and bank_code=?", mer_code, pay_code, class_code, bank_code).First(&m_info)
	return m_info
}

func RateList(mer_code, class_code, bank_code string) []MerRate {
	var m_info []MerRate
	gdb.DB.Where("mer_code=? and class_code=? and bank_code=? and limit_amount>=day_amount", mer_code, class_code, bank_code).Find(&m_info)
	return m_info
}

func MerRateRedis(mer_code, pay_code, class_code, bank_code string) map[string]string {
	redis_key := fmt.Sprintf("mer_rate:%s_%s_%s_%s", mer_code, pay_code, class_code, bank_code)
	//优先查询redis
	m_map := redis.RediGo.HgetAll(redis_key)
	if len(m_map["mer_code"]) < 1 {
		m_info := MerRateInfo(mer_code, pay_code, class_code, bank_code)
		if len(m_info.Mer_code) > 0 {
			m_map = common.StructToMapSlow(m_info)
			redis.RediGo.Hmset(redis_key, m_map, redis_max_time)
			redis.RediGo.Sadd(Data_Redis_Key, redis_key, 0)
		}
	}
	return m_map
}

func MerRateById(rate_id string) MerRate {
	var m_info MerRate
	gdb.DB.Where("id=?", rate_id).First(&m_info)
	return m_info
}

func MerRateByIdRedis(rate_id string) map[string]string {
	redis_key := fmt.Sprintf("mer_rate:id:%s", rate_id)
	//优先查询redis
	m_map := redis.RediGo.HgetAll(redis_key)
	if len(m_map["mer_code"]) < 1 {
		m_info := MerRateById(rate_id)
		if len(m_info.Mer_code) > 0 {
			m_map = common.StructToMapSlow(m_info)
			redis.RediGo.Hmset(redis_key, m_map, redis_max_time)
			redis.RediGo.Sadd(Data_Redis_Key, redis_key, 0)
		}
	}
	return m_map
}

func UpdatesMerRate(a_list MerRate, a_data map[string]interface{}) error {
	res := gdb.DB.Model(&a_list).UpdateColumns(a_data)
	return res.Error
}
