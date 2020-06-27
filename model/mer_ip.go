package model

import (
	"fmt"

	"github.com/zhibingzhou/go_public/common"
	"github.com/zhibingzhou/go_public/redis"
)

func MerIpAuth(mer_code, ip string) MerIp {
	var mer_ip MerIp
	gdb.DB.Where("mer_code=? and ip=?", mer_code, ip).First(&mer_ip)
	return mer_ip
}

func MerIpAuthRedis(mer_code, ip string) map[string]string {
	redis_key := fmt.Sprintf("mer_ip:%s:%s", mer_code, ip)
	//优先查询redis
	ip_map := redis.RediGo.HgetAll(redis_key)
	if len(ip_map["ip"]) < 1 {
		ip_info := MerIpAuth(mer_code, ip)
		if len(ip_info.Ip) > 0 {
			ip_map = common.StructToMapSlow(ip_info)
			redis.RediGo.Hmset(redis_key, ip_map, redis_max_time)
			redis.RediGo.Sadd(Data_Redis_Key, redis_key, 0)
		}
	}
	return ip_map
}

func MerIpByIp(ip string) MerIp {
	var mer_ip MerIp
	gdb.DB.Where("ip=?", ip).First(&mer_ip)
	return mer_ip
}

func MerIpByIpRedis(ip string) map[string]string {
	redis_key := fmt.Sprintf("mer_ip:%s", ip)
	//优先查询redis
	ip_map := redis.RediGo.HgetAll(redis_key)
	if len(ip_map["ip"]) < 1 {
		ip_info := MerIpByIp(ip)
		if len(ip_info.Ip) > 0 {
			ip_map = common.StructToMapSlow(ip_info)
			redis.RediGo.Hmset(redis_key, ip_map, redis_max_time)
			redis.RediGo.Sadd(Data_Redis_Key, redis_key, 0)
		}
	}
	return ip_map
}
