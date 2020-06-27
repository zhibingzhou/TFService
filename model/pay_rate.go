package model

import (
	"fmt"

	"github.com/zhibingzhou/go_public/common"
	"github.com/zhibingzhou/go_public/redis"
)

func PayRateInfo(pay_code, class_code, bank_code string) PayRate {
	var m_info PayRate
	gdb.DB.Where("pay_code=? and class_code=? and bank_code=?", pay_code, class_code, bank_code).First(&m_info)
	return m_info
}

func PayRateById(r_id string) PayRate {
	var m_info PayRate
	gdb.DB.Where("id=?", r_id).First(&m_info)
	return m_info
}

func PayRateRedis(pay_code, class_code, bank_code string) map[string]string {
	redis_key := fmt.Sprintf("pay_rate:%s_%s_%s", pay_code, class_code, bank_code)
	//优先查询redis
	m_map := redis.RediGo.HgetAll(redis_key)
	if len(m_map["pay_code"]) < 1 {
		m_info := PayRateInfo(pay_code, class_code, bank_code)
		if len(m_info.Pay_code) > 0 {
			m_map = common.StructToMapSlow(m_info)
			redis.RediGo.Hmset(redis_key, m_map, redis_max_time)
			redis.RediGo.Sadd(Data_Redis_Key, redis_key, 0)
		}
	}
	return m_map
}

func UpdatesPayRate(a_list PayRate, a_data map[string]interface{}) error {
	res := gdb.DB.Model(&a_list).UpdateColumns(a_data)
	return res.Error
}

func PayRateList(code_arr []string) []PayRate {
	var m_info []PayRate
	gdb.DB.Where("limit_amount>=day_amount and pay_code in (?)", code_arr).Find(&m_info)
	return m_info
}
