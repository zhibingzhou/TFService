package hook

import (
	"TFService/model"
)

func AuthIp(mer_code, ip string) (int, string) {
	h_status := 403
	h_msg := "IP->" + ip + "不在白名单中"
	if mer_code == "1005" || len(ip) > 15 {
		h_status = 200
		h_msg = "success"
		return h_status, h_msg
	}
	ip_info := model.MerIpAuthRedis(mer_code, ip)
	if len(ip_info["ip"]) > 1 {
		h_status = 200
		h_msg = "success"
	}
	return h_status, h_msg
}
