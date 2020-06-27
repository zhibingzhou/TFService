package payapi

import (
	"encoding/json"
	"fmt"
	"TFService/model"
	"strconv"
	"strings"

	"github.com/zhibingzhou/go_public/common"
)

type IPAY struct {
	Pay_url    string
	Mer_code   string
	Notify_url string
	Key        string
	Back_url   string
	Stype      string
	Uid        string
	Call_url   string
	Pay_pwd    string
}

func (this *IPAY) Init(pay_id string) {

	pay_conf := model.ApiConfigRedis(pay_id)
	p := map[string]string{}
	//判断是否存在用户名
	if len(pay_conf["id"]) > 0 {
		res := pay_conf["api_conf"]
		if len(res) > 0 {
			//将解密后的数据切割，并解析到公用参数map里面去
			err := json.Unmarshal([]byte(res), &p)
			if err != nil {
				common.LogsWithFileName(log_path, "ipay_init_", "error->"+err.Error())
			}
		}
	}

	this.Pay_url = p["pay_url"]
	this.Mer_code = p["mer_code"]
	this.Notify_url = p["notify_url"]
	this.Key = p["key"]
	this.Back_url = p["back_url"]
	this.Stype = p["stype"]
	this.Uid = p["uid"]
	this.Pay_pwd = p["pay_pwd"]
	this.Call_url = p["call_url"]
}

func (api *IPAY) CreatePay(p *PayData) (int, string, string, string, string, map[string]string) {
	api_method := "get"
	re_status := 200
	param_form := map[string]string{}
	re_msg := "success"

	nonce := common.Random("smallnumber", 17)
	img_url := ""
	amount_f, _ := strconv.ParseFloat(p.Amount, 64)
	if amount_f < 1.00 {
		re_status = 100
		re_msg = "付款金额错误"
		return re_status, re_msg, api_method, img_url, img_url, param_form
	}
	amount_f = amount_f * 100.00

	sign_str := fmt.Sprintf("amount=%.0f&burl=%s&cid=%s&eparam=%s&ip=%s&nurl=%s&oid=%s&pid=%s&stype=%s&type=&uid=%s&key=%s", amount_f, api.Back_url, p.Pay_bank, nonce, p.Ip, api.Notify_url, p.Order_number, api.Mer_code, api.Stype, api.Uid, api.Key)

	sign := common.HexMd5(sign_str)
	sign = strings.ToUpper(sign)
	param_form["amount"] = fmt.Sprintf("%.0f", amount_f)
	param_form["burl"] = p.Order_number
	param_form["cid"] = p.Pay_bank
	param_form["eparam"] = api.Mer_code
	param_form["ip"] = p.Ip
	param_form["nurl"] = api.Notify_url
	param_form["oid"] = p.Order_number
	param_form["pid"] = api.Mer_code
	param_form["stype"] = api.Stype
	param_form["uid"] = api.Uid
	param_form["type"] = ""
	param_form["sign"] = sign
	img_url = fmt.Sprintf("%s/go/gateway.go", api.Pay_url)

	common.LogsWithFileName(log_path, "ipay_create_", "img_url->"+img_url)

	return re_status, re_msg, api_method, img_url, img_url, param_form
}

//回调验证
func (api *IPAY) CallBackPay(sign, sign_str string) int {
	result := 101

	verify_sign := common.HexMd5(sign_str + "&key=" + api.Key)
	log_str := "sign_str->" + sign_str + "\nsign->" + sign

	if verify_sign == sign {
		result = 200
	}
	common.LogsWithFileName(log_path, "ipay_sign_", log_str)
	return result
}

//代付
func (api *IPAY) PayFor(pay_data map[string]string) (int, string, string) {
	//定义初始值
	api_status := 100
	pay_result := "processing"
	api_msg := "请求失败"

	pay_pwd := common.HexMd5(api.Pay_pwd)
	sign_str := fmt.Sprintf("amount=%s&anumber=%s&bbranch=%s&bcity=%s&bname=%s&bprovince=%s&ctype=%s&nurl=%s&oid=%s&paypwd=%s&pid=%s&uname=%s", pay_data["amount"], pay_data["card_number"], pay_data["bank_branch"], pay_data["bank_city"], pay_data["bank_title"], pay_data["bank_province"], "1", api.Call_url, pay_data["order_number"], pay_pwd, api.Mer_code, pay_data["card_name"])

	sign := common.HexMd5(sign_str + "&key=" + api.Key)
	param := fmt.Sprintf("%s&sign=%s", sign_str, sign)
	pay_url := fmt.Sprintf("%s/go/withdraw/found.go", api.Pay_url)
	p_statu, p_msg := common.HttpBody(pay_url, "POST", param, http_header)
	if p_statu != 200 {
		return api_status, api_msg, pay_result
	}
	var json_res map[string]interface{}
	err := json.Unmarshal(p_msg, &json_res)
	if err != nil {
		api_msg = "json解析错误"
		return api_status, api_msg, pay_result
	}
	if fmt.Sprintf("%v", json_res["code"]) != "101" {
		api_msg = fmt.Sprintf("%v", json_res["msg"])
		pay_result = "fail"
		return api_status, api_msg, pay_result
	}
	api_status = 200
	api_msg = "success"

	return api_status, api_msg, pay_result
}

func (api *IPAY) PayQuery(pay_data map[string]string) (int, string, string) {
	//定义初始值
	api_status := 200
	pay_result := "processing"
	api_msg := "success"

	return api_status, api_msg, pay_result
}
