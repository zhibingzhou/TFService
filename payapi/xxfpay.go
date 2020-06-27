package payapi

import (
	"encoding/json"
	"fmt"
	"TFService/model"
	"strings"
	"time"

	"github.com/zhibingzhou/go_public/common"
)

//新先锋支付
type XXFPAY struct {
	Return_url string
	Pay_url    string
	Mer_code   string
	Token      string
	Notify_url string
}

//新先锋支付
func (this *XXFPAY) Init(pay_id string) {

	pay_conf := model.ApiConfigRedis(pay_id)
	p := map[string]string{}
	//判断是否存在用户名
	if len(pay_conf["id"]) > 0 {
		res := pay_conf["api_conf"]
		if len(res) > 0 {
			//将解密后的数据切割，并解析到公用参数map里面去
			err := json.Unmarshal([]byte(res), &p)
			if err != nil {
				common.LogsWithFileName(log_path, "XXFPAY_init_", "error->"+err.Error())
			}
		}
	}

	this.Pay_url = p["pay_url"]
	this.Mer_code = p["mer_code"]
	this.Notify_url = p["notify_url"]
	this.Return_url = p["return_url"]
	this.Token = p["token"]
}

func (api *XXFPAY) CreatePay(p *PayData) (int, string, string, string, string, map[string]string) {
	api_method := "POST"
	re_status := 100
	re_msg := "请求错误"

	img_url := ""
	param_form := map[string]string{
		"token":      api.Token,
		"money":      p.Amount,
		"outTradeNo": p.Order_number,
		"channel":    p.Pay_bank,
		"timestamp":  fmt.Sprintf("%v", time.Now().UnixNano()/1e6),
		"uid":        api.Mer_code,
		"returnUrl":  api.Return_url,
		"notifyUrl":  api.Notify_url,
	}

	//拼接
	result_url := common.MapCreatLinkSort(param_form, "&", true, false)

	sign := common.HexMd5(result_url)
	sign = strings.ToUpper(sign)
	param_form["sign"] = sign
	delete(param_form, "token")

	//请求三方接口
	param := common.MapCreatLinkSort(param_form, "&", true, false)
	h_status, msg_b := common.HttpBody(api.Pay_url, api_method, param, http_header)
	fmt.Println(param)
	common.LogsWithFileName(log_path, "xxfpay_create_", "param->"+param+"\nmsg->"+string(msg_b))
	if h_status != 200 {
		return re_status, re_msg, api_method, img_url, img_url, param_form
	}

	var json_res map[string]interface{}
	err := json.Unmarshal(msg_b, &json_res)
	if err != nil {
		re_msg = "json错误"
		return re_status, re_msg, api_method, img_url, img_url, param_form
	}

	if fmt.Sprintf("%v", json_res["code"]) != "0" {
		re_msg = fmt.Sprintf("%v", json_res["msg"])
		return re_status, re_msg, api_method, img_url, img_url, param_form
	}

	order_info, ok := json_res["data"].(map[string]interface{})
	if !ok {
		return re_status, re_msg, api_method, img_url, img_url, param_form
	}

	img_url = fmt.Sprintf("%v", order_info["payUrl"])

	if img_url == "" {
		re_msg = "接口错误"
		return re_status, re_msg, api_method, img_url, img_url, param_form
	}

	re_status = 200
	re_msg = "success"
	return re_status, re_msg, api_method, img_url, img_url, param_form
}

//回调验证
func (api *XXFPAY) CallBackPay(sign, sign_str string) int {
	result := 101
	//传过来json数据
	var json_res map[string]string
	err := json.Unmarshal([]byte(sign_str), &json_res)
	if err != nil {
		return result
	}
	json_res["token"] = api.Token
	str := common.MapCreatLinkSort(json_res, "&", true, false)

	verify_sign := common.HexMd5(str)
	verify_sign = strings.ToUpper(verify_sign)

	log_str := "sign_str->" + sign_str + "\nsign->" + sign

	if verify_sign == sign {
		result = 200
	}
	common.LogsWithFileName(log_path, "XXFPAY_sign_", log_str)
	return result
}

//代付
func (api *XXFPAY) PayFor(pay_data map[string]string) (int, string, string) {
	//定义初始值
	api_status := 200
	pay_result := "processing"
	api_msg := "success"

	return api_status, api_msg, pay_result
}

func (api *XXFPAY) PayQuery(pay_data map[string]string) (int, string, string) {
	//定义初始值
	api_status := 200
	pay_result := "processing"
	api_msg := "success"

	return api_status, api_msg, pay_result
}
