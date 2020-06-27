package payapi

import (
	"encoding/json"
	"fmt"
	"TFService/model"

	"github.com/zhibingzhou/go_public/common"
)

//mbpay
type MBPAY struct {
	Notify_url string
	Pay_url    string
	SN         string
	Key        string
	Secret_key string
	Header     map[string]string
}

//mbpay初始化
func (this *MBPAY) Init(pay_id string) {

	pay_conf := model.ApiConfigRedis(pay_id)
	p := map[string]string{}
	//判断是否存在用户名
	if len(pay_conf["id"]) > 0 {
		res := pay_conf["api_conf"]
		if len(res) > 0 {
			//将解密后的数据切割，并解析到公用参数map里面去
			err := json.Unmarshal([]byte(res), &p)
			if err != nil {
				common.LogsWithFileName(log_path, "MBPAY_init_", "error->"+err.Error())
			}
		}
	}

	this.Pay_url = p["pay_url"]
	this.Notify_url = p["notify_url"]
	this.Key = p["key"]
	this.SN = p["sn"]
	this.Secret_key = p["secret_key"]
	this.Header = make(map[string]string)
	this.Header["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36"
	this.Header["Content-Type"] = "application/x-www-form-urlencoded; charset=UTF-8"
}

func (api *MBPAY) CreatePay(p *PayData) (int, string, string, string, string, map[string]string) {
	api_method := "POST"
	re_status := 100
	param_form := map[string]string{}
	re_msg := "请求错误"

	img_url := ""

	return re_status, re_msg, api_method, img_url, img_url, param_form
}

//回调验证
func (api *MBPAY) CallBackPay(sign, sign_str string) int {
	result := 101
	log_str := "sign_str->" + sign_str + "\nsign->" + sign
	result = 200
	common.LogsWithFileName(log_path, "MBPAY_sign_", log_str)
	return result
}

//代付
func (api *MBPAY) PayFor(pay_data map[string]string) (int, string, string) {
	//定义初始值
	api_status := 100
	pay_result := "error"
	api_msg := "代付失败"
	api_method := "POST"

	param_form := map[string]string{
		"TransactionCode":   pay_data["order_number"],
		"AccountNumber":     pay_data["card_number"],
		"AccountName":       pay_data["card_name"],
		"TransactionAmount": pay_data["amount"],
		"BankName":          pay_data["bank_title"],
		"Callback":          api.Notify_url,
		"Version":           "0",
	}

	api.Header["Authorization"] = api.Key
	pay_url := api.Pay_url + fmt.Sprintf("/order/withdraw?sn=%s", api.SN)
	param := common.MapCreatLinkSort(param_form, "&", true, false)

	api_status, api_b := common.HttpBody(pay_url, api_method, param, api.Header)
	common.LogsWithFileName(log_path, "mbpay_payfor_", "param->"+param+"\nmsg->"+string(api_b))

	if api_status != 200 {
		return api_status, api_msg, pay_result
	}

	var json_res map[string]interface{}

	err := json.Unmarshal(api_b, &json_res)
	if err != nil {
		api_msg = "json错误"
		return api_status, api_msg, pay_result
	}

	if fmt.Sprintf("%v", json_res["status"]) != "true" {
		api_msg = fmt.Sprintf("%v", json_res["error"])
		return api_status, api_msg, pay_result
	}

	api_status = 200
	api_msg = "success"
	pay_result = "processing"
	return api_status, api_msg, pay_result
}

func (api *MBPAY) PayQuery(pay_data map[string]string) (int, string, string) {
	//定义初始值
	api_status := 100
	pay_result := "processing"
	api_msg := "代付失败"
	api_method := "POST"

	//请求参数
	param_form := map[string]string{
		"OrderType ":      "withdraw",
		"TransactionCode": pay_data["order_number"],
		"SerialNumber":    pay_data["order_number"],
	}

	api.Header["Authorization"] = api.Key
	pay_url := api.Pay_url + fmt.Sprintf("/order/query?sn=%s", api.SN)
	param := common.MapCreatLinkSort(param_form, "&", true, false)

	api_status, api_b := common.HttpBody(pay_url, api_method, param, api.Header)
	common.LogsWithFileName(log_path, "mbpay_payfor_", "param->"+param+"\nmsg->"+string(api_b))

	if api_status != 200 {
		return api_status, api_msg, pay_result
	}

	var json_res map[string]interface{}

	err := json.Unmarshal(api_b, &json_res)
	if err != nil {
		api_msg = "json错误"
		return api_status, api_msg, pay_result
	}

	if fmt.Sprintf("%v", json_res["status"]) != "true" {
		api_msg = fmt.Sprintf("%v", "下单失败")
		return api_status, api_msg, pay_result
	}

	url_res := make(map[string]interface{})
	err = json.Unmarshal([]byte(fmt.Sprintf("%v", json_res["data"])), &url_res)

	order_status := fmt.Sprintf("%v", url_res["Status"])

	pay_result = "processing"
	switch order_status {
	case "3":
		pay_result = "success"
		api_msg = "success"
	case "4":
		pay_result = "fail"
	}

	api_status = 200
	return api_status, api_msg, pay_result
}
