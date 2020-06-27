package model

import (
	"fmt"

	"github.com/zhibingzhou/go_public/common"
	"github.com/zhibingzhou/go_public/redis"
)

func CashBankInfo(pay_code, bank_code string) CashBank {
	var c_bank CashBank
	gdb.DB.Where("pay_code=? and bank_code=?", pay_code, bank_code).First(&c_bank)
	return c_bank
}

func CashBankRedis(pay_code, bank_code string) map[string]string {
	redis_key := fmt.Sprintf("cash_bank:%s_%s", pay_code, bank_code)
	//优先查询redis
	m_map := redis.RediGo.HgetAll(redis_key)
	if len(m_map["code"]) < 1 {
		m_info := CashBankInfo(pay_code, bank_code)
		if len(m_info.Bank_code) > 0 {
			m_map = common.StructToMapSlow(m_info)
			redis.RediGo.Hmset(redis_key, m_map, redis_max_time)
			redis.RediGo.Sadd(Data_Redis_Key, redis_key, 0)
		}
	}
	return m_map
}
