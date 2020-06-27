package model

import (
	"fmt"

	"github.com/zhibingzhou/go_public/common"
	"github.com/zhibingzhou/go_public/redis"
)

/**
*  根据网站订单查询
 */
func CashByWeb(web_order string) CashList {
	var p_list CashList
	gdb.DB.Model(&PayList{}).Where("order_number = ?", web_order).First(&p_list)
	return p_list
}

/**
* 根据网站订单查询
*  有redis缓存
 */
func CashByWebRedis(web_order string) map[string]string {
	redis_key := fmt.Sprintf("cash_list:order_number:%s", web_order)
	//优先查询redis
	pay_map := redis.RediGo.HgetAll(redis_key)
	if len(pay_map["order_number"]) < 1 {
		p_info := CashByWeb(web_order)
		if len(p_info.Id) > 0 {
			pay_map = common.StructToMapSlow(p_info)
			redis.RediGo.Hmset(redis_key, pay_map, redis_max_time)
			redis.RediGo.Sadd(Data_Redis_Key, redis_key, 0)
		}
	}
	return pay_map
}

/**
*  根据系统订单查询
 */
func CashById(order_number string) CashList {
	var p_list CashList
	gdb.DB.Model(&CashList{}).Where("id = ?", order_number).First(&p_list)
	return p_list
}

func UpdatesCashList(p_list CashList, p_data map[string]interface{}) error {
	res := gdb.DB.Model(&p_list).UpdateColumns(p_data)
	return res.Error
}

func PushCashList(c_time string) []CashList {
	var cash_lists []CashList

	gdb.DB.Where("push_status<>3 and status=3 and push_num<15 and create_time>?",
		c_time).Order("create_time").Limit(100).Find(&cash_lists)

	return cash_lists
}
