package model

import (
	"fmt"

	"github.com/zhibingzhou/go_public/common"
	"github.com/zhibingzhou/go_public/redis"
)

func MerPayInfo(mer_code, pay_id string) MerPay {
	var m_info MerPay
	gdb.DB.Where("mer_code=? and pay_id=?", mer_code, pay_id).First(&m_info)
	return m_info
}

func MerPayRedis(mer_code, pay_id string) map[string]string {
	redis_key := fmt.Sprintf("mer_pay:%s_%s", mer_code, pay_id)
	//优先查询redis
	m_map := redis.RediGo.HgetAll(redis_key)
	if len(m_map["mer_code"]) < 1 {
		m_info := MerPayInfo(mer_code, pay_id)
		if len(m_info.Mer_code) > 0 {
			m_map = common.StructToMapSlow(m_info)
			redis.RediGo.Hmset(redis_key, m_map, redis_max_time)
			redis.RediGo.Sadd(Data_Redis_Key, redis_key, 0)
		}
	}
	return m_map
}

func MerPayList(mer_code string) []MerPay {
	var m_pay []MerPay
	p_status := 1
	gdb.DB.Where("mer_code=? and status=?", mer_code, p_status).Find(&m_pay)
	return m_pay
}

func MerPayById(p_id string) MerPay {
	var m_info MerPay
	gdb.DB.Where("id=?", p_id).First(&m_info)
	return m_info
}

func MerPayByIdRedis(p_id string) map[string]string {
	redis_key := fmt.Sprintf("mer_pay:id_%s", p_id)
	//优先查询redis
	m_map := redis.RediGo.HgetAll(redis_key)
	if len(m_map["mer_code"]) < 1 {
		m_info := MerPayById(p_id)
		if len(m_info.Mer_code) > 0 {
			m_map = common.StructToMapSlow(m_info)
			redis.RediGo.Hmset(redis_key, m_map, redis_max_time)
			redis.RediGo.Sadd(Data_Redis_Key, redis_key, 0)
		}
	}
	return m_map
}

func UpdatesMerPay(a_list MerPay, a_data map[string]interface{}) error {
	res := gdb.DB.Model(&a_list).UpdateColumns(a_data)
	return res.Error
}
