package payapi

import (
	"encoding/json"
	"fmt"
	"net/url"
	"TFService/model"

	"github.com/zhibingzhou/go_public/common"
)

type BFPAY struct {
	Mer_code   string //商户API识别码
	Return_url string //同步跳转
	Notify_url string //异步回调
	Key        string //商户key
	Pay_url    string //支付地址
}

/**
* 对象初始化
 */

func (this *BFPAY) Init(pay_id string) {
	pay_conf := model.ApiConfigRedis(pay_id)
	p := map[string]string{}
	//判断是否存在用户名
	if len(pay_conf["id"]) > 0 {
		res := pay_conf["api_conf"]
		if len(res) > 0 {
			//将解密后的数据切割，并解析到公用参数map里面去
			err := json.Unmarshal([]byte(res), &p)
			if err != nil {
				common.LogsWithFileName(log_path, "jxpay_init_", "error->"+err.Error())
			}
		}
	}

	this.Pay_url = p["pay_url"]
	this.Mer_code = p["mer_code"]
	this.Notify_url = p["notify_url"]
	this.Key = p["key"]
	this.Return_url = p["return_url"]
}

/**
* 发出支付请求
* @param	*PayData	支付信息的指针值
* return	string	需要提交的参数
* return	map	需要用于表单的内容
 */
func (api *BFPAY) CreatePay(p *PayData) (int, string, string, string, string, map[string]string) {
	re_status := 100
	re_msg := "支付请求错误"
	api_method := "POST"
	img_url := ""
	param_form := map[string]string{}
	notify_nurl := url.QueryEscape(api.Notify_url)
	sign_str := fmt.Sprintf("amount=%s&cus_code=%s&cus_order_sn=%s&notify_url=%s&payment_flag=%s", p.Amount, api.Mer_code, p.Order_number, notify_nurl, p.Pay_bank)

	//sign值
	sign := common.HexMd5(sign_str + "&key=" + api.Key)

	param_result := fmt.Sprintf("%s&sign=%s", sign_str, sign)

	post_url := fmt.Sprintf("%s/api/payment/deposit", api.Pay_url)
	h_status, msg_b := common.HttpBody(post_url, api_method, param_result, http_header)
	common.LogsWithFileName(log_path, "bfpay_create_", "param_result->"+param_result+"\npost_url->"+post_url+"\nmsg->"+string(msg_b))
	if h_status != 200 {
		return re_status, re_msg, api_method, img_url, img_url, param_form
	}
	var json_res map[string]interface{}
	err := json.Unmarshal(msg_b, &json_res)
	if err != nil {
		re_msg = "JSON解析失败"
		return re_status, re_msg, api_method, img_url, img_url, param_form
	}

	re_msg = fmt.Sprintf("%v", json_res["message"])
	//捞取结果，赋值到变量
	result := fmt.Sprintf("%v", json_res["result"])
	if result != "success" {
		return re_status, re_msg, api_method, img_url, img_url, param_form

	}
	order_info, ok := json_res["order_info"].(map[string]interface{})
	if !ok {
		return re_status, re_msg, api_method, img_url, img_url, param_form
	}
	re_status = 200
	re_msg = "success"
	img_url = fmt.Sprintf("%v", order_info["payment_uri"])
	return re_status, re_msg, api_method, api.Pay_url, img_url, param_form
}

//回调验证
func (api *BFPAY) CallBackPay(sign, sign_str string) int {
	result := 101

	verify_sign := common.HexMd5(sign_str + "&key=" + api.Key)
	log_str := "sign_str->" + sign_str + "\nsign->" + sign

	if verify_sign == sign {
		result = 200
	}
	common.LogsWithFileName(log_path, "bfpay_sign_", log_str)
	return result
}

//代付
func (api *BFPAY) PayFor(pay_data map[string]string) (int, string, string) {
	//定义初始值
	api_status := 100
	pay_result := "error"
	api_msg := "代付失败"
	//请求参数
	param_form := map[string]string{
		"cus_code":     api.Mer_code,
		"cus_order_sn": pay_data["order_number"],
		"payment_flag": "pay_webbk",
		"amount":       pay_data["amount"],
		"bank_code":    pay_data["bank_code"],
		"bank_account": pay_data["card_number"],
		"account_name": pay_data["card_name"],
		"notify_url":   api.Notify_url,
	}
	encode_form := map[string]string{}
	for key, value := range param_form {
		encode_form[key] = url.QueryEscape(value)
	}
	//拼接
	sign_str := common.MapCreatLinkSort(encode_form, "&", true, false)
	sign_str += fmt.Sprintf("&key=%s", api.Key)

	sign := common.HexMd5(sign_str)
	param_form["sign"] = sign

	param_result := common.MapCreatLinkSort(param_form, "&", true, false)

	api_url := fmt.Sprintf("%s/api/charge/receive", api.Pay_url)

	msg_b, err := HttpPostForm(api_url, param_form)
	common.LogsWithFileName(log_path, "bfpay_payfor_", "param_result->"+param_result+"\napi_url->"+api_url+"\nmsg->"+string(msg_b))
	if err != nil {
		return api_status, api_msg, pay_result
	}
	var json_res map[string]interface{}
	err = json.Unmarshal([]byte(msg_b), &json_res)
	if err != nil {
		api_msg = "json错误"
		return api_status, api_msg, pay_result
	}

	if fmt.Sprintf("%v", json_res["result"]) != "success" && fmt.Sprintf("%v", json_res["status"]) != "200" {
		api_msg = fmt.Sprintf("%v", json_res["message"])
		return api_status, api_msg, pay_result
	}

	api_status = 200
	api_msg = "success"
	pay_result = "processing"
	return api_status, api_msg, pay_result
}

func (api *BFPAY) PayQuery(pay_data map[string]string) (int, string, string) {

	//定义初始值
	api_status := 100
	pay_result := "processing"
	api_msg := "代付失败"

	//请求参数
	param_form := map[string]string{
		"cus_code": api.Mer_code,
		"order_sn": pay_data["order_number"], //这个要用第三方的订单号去查
	}

	//拼接
	sign_str := common.MapCreatLinkSort(param_form, "&", true, false)
	sign_str += fmt.Sprintf("&key=%s", api.Key)
	sign := common.HexMd5(sign_str)
	param_form["sign"] = sign

	//写log
	param_result := common.MapCreatLinkSort(param_form, "&", true, false)

	api_url := fmt.Sprintf("%s/api/charge/info", api.Pay_url)

	msg_b, err := HttpPostForm(api_url, param_form)
	common.LogsWithFileName(log_path, "bfpay_payquery_", "param_result->"+param_result+"\napi_url->"+api_url+"\nmsg->"+string(msg_b))
	if err != nil {
		return api_status, api_msg, pay_result
	}

	var json_res map[string]interface{}
	err = json.Unmarshal([]byte(msg_b), &json_res)
	if err != nil {
		api_msg = "json错误"
		return api_status, api_msg, pay_result
	}

	if fmt.Sprintf("%v", json_res["result"]) != "success" && fmt.Sprintf("%v", json_res["status"]) != "200" {
		api_msg = fmt.Sprintf("%v", json_res["message"])
		return api_status, api_msg, pay_result
	}
	order_status := map[string]interface{}{}
	order_status = json_res["order_info"].(map[string]interface{})

	switch fmt.Sprintf("%v", order_status["order_status"]) {
	case "success":
		pay_result = "success"
		api_msg = "success"
	case "fail", "cancel":
		pay_result = "fail"
	}
	api_status = 200

	return api_status, api_msg, pay_result
}
