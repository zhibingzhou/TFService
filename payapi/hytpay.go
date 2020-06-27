package payapi

import (
	"encoding/json"
	"fmt"
	"TFService/model"
	"time"

	"github.com/zhibingzhou/go_public/common"
)

type HYTPAY struct {
	Pay_url    string
	Mer_code   string
	Notify_url string
	Key        string
	Back_url   string
	Call_url   string
	Header     map[string]string
}

func (this *HYTPAY) Init(pay_id string) {

	pay_conf := model.ApiConfigRedis(pay_id)
	p := map[string]string{}
	//判断是否存在用户名
	if len(pay_conf["id"]) > 0 {
		res := pay_conf["api_conf"]
		if len(res) > 0 {
			//将解密后的数据切割，并解析到公用参数map里面去
			err := json.Unmarshal([]byte(res), &p)
			if err != nil {
				common.LogsWithFileName(log_path, "hytpay_init_", "error->"+err.Error())
			}
		}
	}

	this.Pay_url = p["pay_url"]
	this.Mer_code = p["mer_code"]
	this.Notify_url = p["notify_url"]
	this.Key = p["key"]
	this.Back_url = p["back_url"]
	this.Call_url = p["call_url"]
	this.Header = make(map[string]string)
	this.Header["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36"
	this.Header["Content-Type"] = "application/x-www-form-urlencoded; charset=UTF-8"
}

func (api *HYTPAY) CreatePay(p *PayData) (int, string, string, string, string, map[string]string) {
	api_method := "POST"
	re_status := 100
	param_form := map[string]string{}
	re_msg := "请求错误"

	img_url := ""

	param := fmt.Sprintf("Amount=%s&Ip=%s&MerchantId=%s&MerchantUniqueOrderId=%s&NotifyUrl=%s&PayTypeId=%s&ReturnUrl=%s", p.Amount, p.Ip, api.Mer_code, p.Order_number, api.Notify_url, p.Pay_bank, api.Back_url)

	api_url := fmt.Sprintf("%s/InterfaceV4/CreatePayOrder/", api.Pay_url)

	api_status, api_b := common.HttpBody(api_url, api_method, param, api.Header)
	common.LogsWithFileName(log_path, "hytpay_create_", "param->"+param+"\nmsg->"+string(api_b))
	if api_status != 200 {
		return re_status, re_msg, api_method, img_url, img_url, param_form
	}

	var json_res map[string]interface{}
	err := json.Unmarshal(api_b, &json_res)
	if err != nil {
		re_msg = "json错误"
		return re_status, re_msg, api_method, img_url, img_url, param_form
	}

	if fmt.Sprintf("%v", json_res["Code"]) != "0" {
		re_msg = fmt.Sprintf("%v", json_res["MessageForUser"])
		return re_status, re_msg, api_method, img_url, img_url, param_form
	}

	img_url = fmt.Sprintf("%v", json_res["Url"])

	if img_url == "" {
		re_msg = "接口错误"
		return re_status, re_msg, api_method, img_url, img_url, param_form
	}
	re_status = 200
	re_msg = "success"
	return re_status, re_msg, api_method, img_url, img_url, param_form
}

//回调验证
func (api *HYTPAY) CallBackPay(sign, sign_str string) int {
	result := 101

	verify_sign := common.HexMd5(sign_str + api.Key)
	log_str := "sign_str->" + sign_str + "\nsign->" + sign

	if verify_sign == sign {
		result = 200
	}
	common.LogsWithFileName(log_path, "hytpay_sign_", log_str)
	return result
}

//代付
func (api *HYTPAY) PayFor(pay_data map[string]string) (int, string, string) {
	//定义初始值
	api_status := 100
	pay_result := "error"
	api_msg := "代付失败"
	api_method := "POST"
	timestamp := time.Now().Format(day_f)

	sign_str := fmt.Sprintf("Amount=%s&BankCardBankName=%s&BankCardNumber=%s&BankCardRealName=%s&MerchantId=%s&MerchantUniqueOrderId=%s&NotifyUrl=%s&Timestamp=%s&WithdrawTypeId=0", pay_data["amount"], pay_data["bank_title"], pay_data["card_number"], pay_data["card_name"], api.Mer_code, pay_data["order_number"], api.Call_url, timestamp)

	sign := common.HexMd5(sign_str + api.Key)

	param := fmt.Sprintf("%s&Sign=%s", sign_str, sign)

	api_url := fmt.Sprintf("%s/WithdrawOrder/CreateWithdrawOrderV3/", api.Pay_url)

	api_status, api_b := common.HttpBody(api_url, api_method, param, api.Header)
	common.LogsWithFileName(log_path, "hytpay_payfor_", "param->"+param+"\nmsg->"+string(api_b))
	if api_status != 200 {
		return api_status, api_msg, pay_result
	}

	var json_res map[string]interface{}
	err := json.Unmarshal(api_b, &json_res)
	if err != nil {
		api_msg = "json错误"
		return api_status, api_msg, pay_result
	}

	if fmt.Sprintf("%v", json_res["Code"]) != "0" {
		api_msg = fmt.Sprintf("%v", json_res["Message"])
		return api_status, api_msg, pay_result
	}

	api_status = 200
	api_msg = "success"
	pay_result = "processing"
	return api_status, api_msg, pay_result
}

func (api *HYTPAY) PayQuery(pay_data map[string]string) (int, string, string) {
	//定义初始值
	api_status := 100
	pay_result := "processing"
	api_msg := "代付失败"
	api_method := "POST"
	timestamp := time.Now().Format(day_f)

	sign_str := fmt.Sprintf("MerchantId=%s&MerchantUniqueOrderId=%s&Timestamp=%s", api.Mer_code, pay_data["order_number"], timestamp)

	sign := common.HexMd5(sign_str + api.Key)

	param := fmt.Sprintf("%s&Sign=%s", sign_str, sign)

	api_url := fmt.Sprintf("%s/WithdrawOrder/QueryWithdrawOrderV3/", api.Pay_url)

	api_status, api_b := common.HttpBody(api_url, api_method, param, api.Header)
	common.LogsWithFileName(log_path, "hytpay_query_", "param->"+param+"\nmsg->"+string(api_b))
	if api_status != 200 {
		return api_status, api_msg, pay_result
	}

	var json_res map[string]interface{}
	err := json.Unmarshal(api_b, &json_res)
	if err != nil {
		api_msg = "json错误"
		return api_status, api_msg, pay_result
	}

	if fmt.Sprintf("%v", json_res["Code"]) != "0" {
		api_msg = fmt.Sprintf("%v", json_res["Message"])
		return api_status, api_msg, pay_result
	}

	order_status := fmt.Sprintf("%v", json_res["WithdrawOrderStatus"])
	switch order_status {
	case "100":
		pay_result = "success"
	case "-90":
		pay_result = "fail"
	}

	api_status = 200
	api_msg = "success"
	return api_status, api_msg, pay_result
}
