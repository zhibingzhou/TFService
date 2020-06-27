package model

import (
	"fmt"

	"github.com/zhibingzhou/go_public/common"
	"github.com/zhibingzhou/go_public/redis"
)

func MerBankInfo(card_number, mer_code string) MerBank {
	var m_info MerBank
	gdb.DB.Where("card_number=? and mer_code=?", card_number, mer_code).First(&m_info)
	return m_info
}

func MerBankRedis(card_number, mer_code string) map[string]string {
	redis_key := fmt.Sprintf("mer_bank:card_number:mer_code:%s_%s", card_number, mer_code)
	//优先查询redis
	m_map := redis.RediGo.HgetAll(redis_key)
	if len(m_map["code"]) < 1 {
		m_info := MerBankInfo(card_number, mer_code)
		if len(m_info.Bank_code) > 0 {
			m_map = common.StructToMapSlow(m_info)
			redis.RediGo.Hmset(redis_key, m_map, redis_max_time)
			redis.RediGo.Sadd(Data_Redis_Key, redis_key, 0)
		}
	}
	return m_map
}

func MerBankById(b_id string) MerBank {
	var m_info MerBank
	gdb.DB.Where("id=?", b_id).First(&m_info)
	return m_info
}

func MerBankByIdRedis(b_id string) map[string]string {
	redis_key := fmt.Sprintf("mer_bank:%s", b_id)
	//优先查询redis
	m_map := redis.RediGo.HgetAll(redis_key)
	if len(m_map["code"]) < 1 {
		m_info := MerBankById(b_id)
		if len(m_info.Bank_code) > 0 {
			m_map = common.StructToMapSlow(m_info)
			redis.RediGo.Hmset(redis_key, m_map, redis_max_time)
			redis.RediGo.Sadd(Data_Redis_Key, redis_key, 0)
		}
	}
	return m_map
}
