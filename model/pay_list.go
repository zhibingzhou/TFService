package model

import (
	"fmt"

	"github.com/zhibingzhou/go_public/common"
	"github.com/zhibingzhou/go_public/redis"
)

func CreatePayList(pay_list PayList) error {
	res := gdb.DB.Create(&pay_list)
	return res.Error
}

/**
*  根据网站订单查询
 */
func OrderByWeb(web_order string) PayList {
	var p_list PayList
	gdb.DB.Model(&PayList{}).Where("order_number = ?", web_order).First(&p_list)
	return p_list
}

/**
* 根据网站订单查询
*  有redis缓存
 */
func OrderByWebRedis(web_order string) map[string]string {
	redis_key := fmt.Sprintf("pay_list:order_number:%s", web_order)
	//优先查询redis
	pay_map := redis.RediGo.HgetAll(redis_key)
	if len(pay_map["order_number"]) < 1 {
		p_info := OrderByWeb(web_order)
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
func OrderById(order_number string) PayList {
	var p_list PayList
	gdb.DB.Model(&PayList{}).Where("id = ?", order_number).First(&p_list)
	return p_list
}

func UpdatesPayList(p_list PayList, p_data map[string]interface{}) error {
	res := gdb.DB.Model(&p_list).UpdateColumns(p_data)
	return res.Error
}

func PushPayList(c_time string) []PayList {
	var pay_lists []PayList

	gdb.DB.Where("push_status<>3 and status=3 and push_num<15 and create_time>?",
		c_time).Order("create_time").Limit(100).Find(&pay_lists)

	return pay_lists
}
