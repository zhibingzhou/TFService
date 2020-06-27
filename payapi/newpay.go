package payapi

import (
	"encoding/json"
	"fmt"
	"TFService/model"
	"strings"

	"github.com/zhibingzhou/go_public/common"
)

type NEWPAY struct {
	Pay_url    string
	Mer_code   string
	Notify_url string
	Key        string
	Back_url   string
	Call_url   string
	Pay_pwd    string
}

func (this *NEWPAY) Init(pay_id string) {

	pay_conf := model.ApiConfigRedis(pay_id)
	p := map[string]string{}
	//判断是否存在用户名
	if len(pay_conf["id"]) > 0 {
		res := pay_conf["api_conf"]
		if len(res) > 0 {
			//将解密后的数据切割，并解析到公用参数map里面去
			err := json.Unmarshal([]byte(res), &p)
			if err != nil {
				common.LogsWithFileName(log_path, "NEWPAY_init_", "error->"+err.Error())
			}
		}
	}

	this.Pay_url = p["pay_url"]
	this.Mer_code = p["mer_code"]
	this.Notify_url = p["notify_url"]
	this.Back_url = p["back_url"]
	this.Key = p["key"]
	this.Pay_pwd = p["pay_pwd"]
	this.Call_url = p["call_url"]
}

func (api *NEWPAY) CreatePay(p *PayData) (int, string, string, string, string, map[string]string) {
	api_method := "POST"
	re_status := 100
	param_form := map[string]string{}
	re_msg := "请求错误"

	nonce := common.Random("smallnumber", 8)
	img_url := ""

	sign_str := fmt.Sprintf("amount=%s&appId=%s&asyncUrl=%s&ip=%s&nonceStr=%s&outTradeNo=%s&payType=%s&returnUrl=%s", p.Amount, api.Mer_code, api.Notify_url, p.Ip, nonce, p.Order_number, p.Pay_bank, api.Back_url)
	sign := common.HexMd5(sign_str + "&key=" + api.Key)
	sign = strings.ToUpper(sign)
	param := fmt.Sprintf("%s&sign=%s", sign_str, sign)

	post_url := fmt.Sprintf("%s/api/v1/pay", api.Pay_url)
	h_status, msg_b := common.HttpBody(post_url, api_method, param, http_header)
	common.LogsWithFileName(log_path, "newpay_create_", "post_url->"+post_url+"\nmsg->"+string(msg_b))
	if h_status != 200 {
		return re_status, re_msg, api_method, img_url, img_url, param_form
	}

	var json_res map[string]interface{}
	err := json.Unmarshal(msg_b, &json_res)
	if err != nil {
		re_msg = "JSON错误"
		return re_status, re_msg, api_method, img_url, img_url, param_form
	}

	if fmt.Sprintf("%v", json_res["code"]) != "200" {
		re_msg = fmt.Sprintf("%v", json_res["message"])
		return re_status, re_msg, api_method, img_url, img_url, param_form
	}

	if json_res["data"] == nil {
		re_msg = fmt.Sprintf("%v", json_res["message"])
		return re_status, re_msg, api_method, img_url, img_url, param_form
	}

	data, ok := json_res["data"].(map[string]interface{})
	if !ok {
		re_msg = "类型异常"
		return re_status, re_msg, api_method, img_url, img_url, param_form
	}
	img_url = fmt.Sprintf("%v", data["payUrl"])
	re_status = 200
	re_msg = "success"
	return re_status, re_msg, api_method, img_url, img_url, param_form
}

//回调验证
func (api *NEWPAY) CallBackPay(sign, sign_str string) int {
	result := 100

	verify_sign := common.HexMd5(sign_str + "&key=" + api.Key)
	verify_sign = strings.ToUpper(verify_sign)

	if verify_sign == sign {
		result = 200
	}
	log_str := "sign_str->" + sign_str + "\nsign->" + sign + "\nverify_sign->" + verify_sign
	common.LogsWithFileName(log_path, "newpay_sign_", log_str)
	return result
}

//代付
func (api *NEWPAY) PayFor(pay_data map[string]string) (int, string, string) {
	//定义初始值
	api_status := 200
	pay_result := "processing"
	api_method := "POST"
	api_msg := "success"
	nonce := common.Random("smallnumber", 8)
	pwd := common.HexMd5(api.Pay_pwd)
	sign_str := fmt.Sprintf("amount=%s&appId=%s&asyncUrl=%s&bankBranch=%s&bankCode=%s&card=%s&name=%s&nonceStr=%s&outTradeNo=%s&password=%s&returnUrl=%s", pay_data["amount"], api.Mer_code, api.Call_url, pay_data["bank_branch"], pay_data["bank_code"], pay_data["card_number"], pay_data["card_name"], nonce, pay_data["order_number"], pwd, api.Back_url)

	sign := common.HexMd5(sign_str + "&key=" + api.Key)
	sign = strings.ToUpper(sign)
	param := fmt.Sprintf("%s&sign=%s", sign_str, sign)

	post_url := fmt.Sprintf("%s/api/v1/issued", api.Pay_url)
	h_status, msg_b := common.HttpBody(post_url, api_method, param, http_header)
	common.LogsWithFileName(log_path, "newpay_payfor_", "post_url->"+post_url+"\nparam->"+param+"\nmsg->"+string(msg_b))
	if h_status != 200 {
		return api_status, api_msg, pay_result
	}

	var json_res map[string]interface{}
	err := json.Unmarshal(msg_b, &json_res)
	if err != nil {
		api_msg = "JSON错误"
		return api_status, api_msg, pay_result
	}

	if fmt.Sprintf("%v", json_res["code"]) != "200" {
		api_msg = fmt.Sprintf("%v", json_res["message"])
		pay_result = "fail"
		return api_status, api_msg, pay_result
	}
	return api_status, api_msg, pay_result
}

func (api *NEWPAY) PayQuery(pay_data map[string]string) (int, string, string) {
	//定义初始值
	api_status := 200
	pay_result := "processing"
	api_msg := "success"

	return api_status, api_msg, pay_result
}
