package payapi

import (
	"encoding/json"
	"fmt"
	"TFService/model"
	"time"

	"github.com/zhibingzhou/go_public/common"
)

type DFPAY struct {
	Pay_url    string
	Mer_code   string
	Notify_url string
	Key        string
	Back_url   string
	Uid        string
}

func (this *DFPAY) Init(pay_id string) {

	pay_conf := model.ApiConfigRedis(pay_id)
	p := map[string]string{}
	//判断是否存在用户名
	if len(pay_conf["id"]) > 0 {
		res := pay_conf["api_conf"]
		if len(res) > 0 {
			//将解密后的数据切割，并解析到公用参数map里面去
			err := json.Unmarshal([]byte(res), &p)
			if err != nil {
				common.LogsWithFileName(log_path, "dfpay_init_", "error->"+err.Error())
			}
		}
	}

	this.Pay_url = p["pay_url"]
	this.Mer_code = p["mer_code"]
	this.Notify_url = p["notify_url"]
	this.Key = p["key"]
	this.Back_url = p["back_url"]
	this.Uid = p["uid"]
}

func (api *DFPAY) CreatePay(p *PayData) (int, string, string, string, string, map[string]string) {
	api_method := "get"
	re_status := 200
	param_form := map[string]string{}
	re_msg := "success"

	now_time := time.Now().Unix()
	img_url := ""

	sign_str := fmt.Sprintf("cid=%s&uid=%s&time=%d&amount=%s&order_id=%s&ip=%s", api.Mer_code, api.Uid, now_time, p.Amount, p.Order_number, p.Ip)

	sign := common.BaseHmacSha(sign_str, api.Key)
	param_form["cid"] = api.Mer_code
	param_form["uid"] = api.Uid
	param_form["time"] = fmt.Sprintf("%d", now_time)
	param_form["amount"] = p.Amount
	param_form["order_id"] = p.Order_number
	param_form["ip"] = p.Ip
	param_form["type"] = p.Pay_bank
	param_form["sign"] = sign
	img_url = fmt.Sprintf("%s/dsdf/customer_pay/init_din", api.Pay_url)

	common.LogsWithFileName(log_path, "dfpay_create_", "img_url->"+img_url)

	return re_status, re_msg, api_method, img_url, img_url, param_form
}

//回调验证
func (api *DFPAY) CallBackPay(sign, sign_str string) int {
	result := 101

	verify_sign := common.BaseHmacSha(sign_str, api.Key)
	log_str := "sign_str->" + sign_str + "\nsign->" + sign

	if verify_sign == sign {
		result = 200
	}
	common.LogsWithFileName(log_path, "dfpay_sign_", log_str)
	return result
}

//代付
func (api *DFPAY) PayFor(pay_data map[string]string) (int, string, string) {
	//定义初始值
	api_status := 200
	pay_result := "processing"
	api_msg := "success"

	return api_status, api_msg, pay_result
}

func (api *DFPAY) PayQuery(pay_data map[string]string) (int, string, string) {
	//定义初始值
	api_status := 200
	pay_result := "processing"
	api_msg := "success"

	return api_status, api_msg, pay_result
}
