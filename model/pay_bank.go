package model

import (
	"fmt"

	"github.com/zhibingzhou/go_public/common"
	"github.com/zhibingzhou/go_public/redis"
)

func PayBankInfo(is_mobile, pay_code, class_code, bank_code string) PayBank {
	var pay_bank PayBank
	gdb.DB.Model(&PayBank{}).Where("is_mobile=? and pay_code=? and class_code=? and bank_code=?", is_mobile, pay_code, class_code, bank_code).First(&pay_bank)
	return pay_bank
}

func PayBankRedis(is_mobile, pay_code, class_code, bank_code string) map[string]string {
	redis_key := fmt.Sprintf("pay_bank:%s_%s_%s_%s", is_mobile, pay_code, class_code, bank_code)
	//优先查询redis
	pay_map := redis.RediGo.HgetAll(redis_key)
	if len(pay_map["Api_conf"]) < 1 {
		conf_info := PayBankInfo(is_mobile, pay_code, class_code, bank_code)
		if len(conf_info.Bank_code) > 0 {
			pay_map = common.StructToMapSlow(conf_info)
			redis.RediGo.Hmset(redis_key, pay_map, redis_max_time)
			redis.RediGo.Sadd(Conf_Redis_Key, redis_key, 0)
		}
	}
	return pay_map
}

func BankList(is_mobile, pay_code string) []PayBank {
	var b_list []PayBank
	gdb.DB.Model(&PayBank{}).Where("is_mobile=? and pay_code=?", is_mobile,
		pay_code).Find(&b_list)

	return b_list
}
