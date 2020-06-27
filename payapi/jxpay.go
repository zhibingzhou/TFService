package payapi

import (
	"encoding/json"
	"fmt"
	"TFService/model"

	"github.com/zhibingzhou/go_public/common"
)

type JXPAY struct {
	Pay_url    string
	Mer_code   string
	Notify_url string
	Key        string
	Return_url string
	Header     map[string]string
}

func (this *JXPAY) Init(pay_id string) {

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
	this.Header = make(map[string]string)
	this.Header["Content-type"] = "application/json; charset=UTF-8"
	this.Header["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.79 Safari/537.36"
}

func (api *JXPAY) CreatePay(p *PayData) (int, string, string, string, string, map[string]string) {
	api_method := "POST"
	re_status := 100
	param_form := map[string]string{}
	re_msg := "请求错误"
	img_url := ""
	sign_str := fmt.Sprintf("amount=%s&merchant_user_id=%s&notify_url=%s&out_trade_no=%s&pay_way=%s&return_url=%s", p.Amount, api.Mer_code, api.Notify_url, p.Order_number, p.Pay_bank, api.Return_url)

	sign := common.HexMd5(sign_str + api.Key)

	param := fmt.Sprintf(`{"merchant_user_id":"%s","sign_type":"MD5","sign":"%s","out_trade_no":"%s","pay_way":"%s","amount":"%s","notify_url":"%s","return_url":"%s"}`, api.Mer_code, sign, p.Order_number, p.Pay_bank, p.Amount, api.Notify_url, api.Return_url)
	post_url := fmt.Sprintf("%s/gateway/merchant/order", api.Pay_url)
	h_status, msg_b := common.HttpBody(post_url, api_method, param, api.Header)
	common.LogsWithFileName(log_path, "jxpay_create_", "post_url->"+post_url+"\nparam->"+param+"\nmsg->"+string(msg_b))
	if h_status != 200 {
		return re_status, re_msg, api_method, img_url, img_url, param_form
	}
	var json_res map[string]interface{}
	err := json.Unmarshal(msg_b, &json_res)
	if err != nil {
		re_msg = "json解析错误"
		return re_status, re_msg, api_method, img_url, img_url, param_form
	}
	code := fmt.Sprintf("%v", json_res["code"])
	re_msg = fmt.Sprintf("%v", json_res["msg"])
	if code != "200" || json_res["code"] == nil {
		return re_status, re_msg, api_method, img_url, img_url, param_form
	}

	data, ok := json_res["data"].(map[string]interface{})
	if !ok {
		return re_status, re_msg, api_method, img_url, img_url, param_form
	}
	re_status = 200
	re_msg = "success"
	img_url = fmt.Sprintf("%v", data["pay_url"])
	return re_status, re_msg, api_method, img_url, img_url, param_form
}

//回调验证
func (api *JXPAY) CallBackPay(sign, sign_str string) int {
	result := 100
	verify_sign := common.HexMd5(sign_str + api.Key)
	log_str := "sign_str->" + sign_str + "\nsign->" + sign + "\nverify_sign->" + verify_sign

	if verify_sign == sign {
		result = 200
	}
	common.LogsWithFileName(log_path, "jxpay_sign_", log_str)
	return result
}

//代付
func (api *JXPAY) PayFor(pay_data map[string]string) (int, string, string) {
	//定义初始值
	api_status := 200
	pay_result := "processing"
	api_msg := "success"

	return api_status, api_msg, pay_result
}

func (api *JXPAY) PayQuery(pay_data map[string]string) (int, string, string) {
	//定义初始值
	api_status := 200
	pay_result := "processing"
	api_msg := "success"

	return api_status, api_msg, pay_result
}
