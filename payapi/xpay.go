package payapi

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"TFService/model"

	"github.com/zhibingzhou/go_public/common"
)

//xpay
type XPAY struct {
	Notify_url string
	Pay_url    string
	Mer_code   string
	Aes_key    string
	Key        string
	Header     map[string]string
}

//xpay初始化
func (this *XPAY) Init(pay_id string) {

	pay_conf := model.ApiConfigRedis(pay_id)
	p := map[string]string{}
	//判断是否存在用户名
	if len(pay_conf["id"]) > 0 {
		res := pay_conf["api_conf"]
		if len(res) > 0 {
			//将解密后的数据切割，并解析到公用参数map里面去
			err := json.Unmarshal([]byte(res), &p)
			if err != nil {
				common.LogsWithFileName(log_path, "XPAY_init_", "error->"+err.Error())
			}
		}
	}

	this.Pay_url = p["pay_url"]
	this.Mer_code = p["mer_code"]
	this.Notify_url = p["notify_url"]
	this.Aes_key = p["aes_key"]
	this.Key = p["key"]
	this.Header = make(map[string]string)
	this.Header["Content-type"] = "application/json; charset=UTF-8"
	this.Header["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.79 Safari/537.36"
}

func (api *XPAY) CreatePay(p *PayData) (int, string, string, string, string, map[string]string) {
	api_method := "POST"
	re_status := 100
	re_msg := "请求错误"

	img_url := ""
	//业务参数
	yparam_form := map[string]string{
		"merchantNo":        api.Mer_code,
		"orderNo":           p.Order_number,
		"productCode":       p.Pay_bank,
		"tradeProfitType":   "ReceivableProduct", //类型支付
		"orderAmount":       p.Amount,
		"serverCallbackUrl": api.Notify_url,
		"goodsName":         "test",
		"bankCode":          "ICBC",
		"bankBusinessType":  "B2C",
		"bankCardType":      "DEBIT",
	}

	//请求参数
	param_form := map[string]string{
		"merchantNo":       api.Mer_code,
		"orderNo":          p.Order_number,
		"tradeProfitType":  "ReceivableProduct", //类型支付
		"productCode":      p.Pay_bank,
		"bankBusinessType": "B2C",
		"bankCardType":     "DEBIT",
	}

	//base64 解码 aeskey
	aeskey, err := base64.StdEncoding.DecodeString(api.Aes_key)

	if err != nil {
		return re_status, re_msg, api_method, img_url, img_url, param_form
	}

	result_url, _ := json.Marshal(yparam_form)

	mer_ecb := common.SetAESECB(string(aeskey), "", "", "", 16)
	contents := mer_ecb.AesEncryptString(string(result_url))
	param_form["content"] = contents

	//拼接
	param := common.MapCreatLinkSort(yparam_form, ",", true, false)
	//首尾拼接key
	param = api.Key + "," + param + "," + api.Key

	//sha 256 加密
	h := sha256.New()
	h.Write([]byte(param))
	bs := fmt.Sprintf("%x", h.Sum(nil))

	param_form["sign"] = string(bs)
	params, _ := json.Marshal(param_form)
	rep := string(params)

	//请求三方接口
	h_status, msg_b := common.HttpBody(api.Pay_url, api_method, rep, api.Header)

	common.LogsWithFileName(log_path, "xpay_create_", "param->"+rep+"\nmsg->"+string(msg_b))
	if h_status != 200 {
		return re_status, re_msg, api_method, img_url, img_url, param_form
	}

	var json_res map[string]interface{}

	err = json.Unmarshal(msg_b, &json_res)
	if err != nil {
		re_msg = "json错误"
		return re_status, re_msg, api_method, img_url, img_url, param_form
	}

	jsonresult, _ := json_res["content"].(string)

	get_ecb := common.SetAESECB(string(aeskey), "", "", "", 16)
	url_body := get_ecb.AesDecryptString(jsonresult)

	url_res := make(map[string]interface{})
	err = json.Unmarshal([]byte(url_body), &url_res)

	if fmt.Sprintf("%v", url_res["responseCode"]) != "0000" {
		re_msg = fmt.Sprintf("%v", url_res["errorMsg"])
		return re_status, re_msg, api_method, img_url, img_url, param_form
	}

	img_url = fmt.Sprintf("%v", url_res["payUrl"])

	if img_url == "" {
		re_msg = "接口错误"
		return re_status, re_msg, api_method, img_url, img_url, param_form
	}

	re_status = 200
	re_msg = "success"
	return re_status, re_msg, api_method, img_url, img_url, param_form

}

//回调验证
func (api *XPAY) CallBackPay(sign, sign_str string) int {
	result := 101
	//首尾拼接key
	sign_str = api.Key + "," + sign_str + "," + api.Key
	//sha 256 加密
	h := sha256.New()
	h.Write([]byte(sign_str))
	verify_sign := fmt.Sprintf("%x", h.Sum(nil))

	log_str := "sign_str->" + sign_str + "\nsign->" + sign

	if verify_sign == sign {
		result = 200
	}
	common.LogsWithFileName(log_path, "XPAY_sign_", log_str)
	return result
}

//代付
func (api *XPAY) PayFor(pay_data map[string]string) (int, string, string) {
	//定义初始值
	api_status := 100
	pay_result := "error"
	api_msg := "代付失败"
	api_method := "POST"

	//业务参数
	yparam_form := map[string]string{
		"merchantNo":        api.Mer_code,
		"orderNo":           pay_data["order_number"],
		"orderAmount":       pay_data["amount"],
		"bankCode":          pay_data["bank_code"],
		"accountName":       pay_data["card_name"],
		"bankBusinessType":  "PRIVATE",
		"accountNo":         pay_data["card_number"],
		"productCode":       "PAYAPI_PRIVATE", //产品编码
		"tradeProfitType":   "PayProduct",     //业务类型
		"serverCallbackUrl": api.Notify_url,
	}

	//请求参数
	param_form := map[string]string{
		"merchantNo":      api.Mer_code,
		"orderNo":         pay_data["order_number"],
		"tradeProfitType": "PayProduct", //类型代付
		"productCode":     "PAYAPI_PRIVATE",
	}

	aeskey, err := base64.StdEncoding.DecodeString(api.Aes_key)

	if err != nil {
		api_msg = "base4 解密aeskey出错"
		return api_status, api_msg, pay_result
	}
	//json
	result_url, _ := json.Marshal(yparam_form)

	mer_ecb := common.SetAESECB(string(aeskey), "", "", "", 16)
	contents := mer_ecb.AesEncryptString(string(result_url))
	param_form["content"] = contents

	//拼接
	params := common.MapCreatLinkSort(yparam_form, ",", true, false)
	//首尾拼接key
	params = api.Key + "," + params + "," + api.Key

	//sha 256 加密
	h := sha256.New()
	h.Write([]byte(params))
	bs := fmt.Sprintf("%x", h.Sum(nil))

	param_form["sign"] = string(bs)
	rep, _ := json.Marshal(param_form)
	param := string(rep)

	api_status, api_b := common.HttpBody(api.Pay_url, api_method, param, api.Header)

	common.LogsWithFileName(log_path, "xpay_payfor_", "param->"+param+"\nmsg->"+string(api_b))
	if api_status != 200 {
		return api_status, api_msg, pay_result
	}

	var json_res map[string]interface{}

	err = json.Unmarshal(api_b, &json_res)
	if err != nil {
		api_msg = "json错误"
		return api_status, api_msg, pay_result
	}

	//接收参数

	jsonresult, _ := json_res["content"].(string)

	get_ecb := common.SetAESECB(string(aeskey), "", "", "", 16)
	url_body := get_ecb.AesDecryptString(jsonresult)

	url_res := make(map[string]interface{})
	err = json.Unmarshal([]byte(url_body), &url_res)

	if fmt.Sprintf("%v", url_res["responseCode"]) != "0000" {
		api_msg = fmt.Sprintf("%v", url_res["errorMsg"])
		return api_status, api_msg, pay_result
	}

	api_status = 200
	api_msg = "success"
	pay_result = "processing"
	return api_status, api_msg, pay_result
}

func (api *XPAY) PayQuery(pay_data map[string]string) (int, string, string) {
	//定义初始值
	api_status := 100
	pay_result := "processing"
	api_msg := "代付失败"
	api_method := "POST"

	//业务参数
	yparam_form := map[string]string{
		"merchantNo":      api.Mer_code,
		"orderNo":         pay_data["order_number"],
		"productCode":     "PAYAPI_PRIVATE", //产品编码
		"tradeProfitType": "ONLINE",         //类型代付
	}

	//请求参数
	param_form := map[string]string{
		"merchantNo":      api.Mer_code,
		"orderNo":         pay_data["order_number"],
		"tradeProfitType": "ONLINE", //类型代付
		"productCode":     "PAYAPI_PRIVATE",
	}

	aeskey, err := base64.StdEncoding.DecodeString(api.Aes_key)

	if err != nil {
		api_msg = "base4 解密aeskey出错"
		return api_status, api_msg, pay_result
	}
	//json
	result_url, _ := json.Marshal(yparam_form)

	//AES 加密
	mer_ecb := common.SetAESECB(string(aeskey), "", "", "", 16)
	contents := mer_ecb.AesEncryptString(string(result_url))
	param_form["content"] = contents

	//拼接
	params := common.MapCreatLinkSort(yparam_form, ",", true, false)
	//首尾拼接key
	params = api.Key + "," + params + "," + api.Key

	//sha 256 加密
	h := sha256.New()
	h.Write([]byte(params))
	bs := fmt.Sprintf("%x", h.Sum(nil))

	param_form["sign"] = string(bs)
	rep, _ := json.Marshal(param_form)
	param := string(rep)

	api_status, api_b := common.HttpBody(api.Pay_url, api_method, param, api.Header)

	common.LogsWithFileName(log_path, "xpay_payfor_", "param->"+param+"\nmsg->"+string(api_b))
	if api_status != 200 {
		return api_status, api_msg, pay_result
	}

	var json_res map[string]interface{}

	err = json.Unmarshal(api_b, &json_res)
	if err != nil {
		api_msg = "json错误"
		return api_status, api_msg, pay_result
	}

	jsonresult, _ := json_res["content"].(string)

	get_ecb := common.SetAESECB(string(aeskey), "", "", "", 16)
	url_body := get_ecb.AesDecryptString(jsonresult)

	url_res := make(map[string]interface{})
	err = json.Unmarshal([]byte(url_body), &url_res)

	if fmt.Sprintf("%v", url_res["responseCode"]) != "0000" {
		api_msg = fmt.Sprintf("%v", url_res["errorMsg"])
		return api_status, api_msg, pay_result
	}
	order_status := fmt.Sprintf("%v", url_res["orderStatus"])

	switch order_status {
	case "SUCCESS":
		pay_result = "success"
	case "FAILED", "BACK":
		pay_result = "fail"
	}

	api_status = 200
	api_msg = "success"
	return api_status, api_msg, pay_result
}
