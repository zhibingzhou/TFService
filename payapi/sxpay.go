package payapi

import (
	"encoding/json"
	"fmt"
	"TFService/model"

	"github.com/zhibingzhou/go_public/common"
)

type SXPAY struct {
	Pay_url    string
	Mer_code   string
	Notice_url string
	Return_url string
	Key        string
}

//舒信初始化
func (this *SXPAY) Init(pay_id string) {

	pay_conf := model.ApiConfigRedis(pay_id)
	p := map[string]string{}
	//判断是否存在用户名
	if len(pay_conf["id"]) > 0 {
		res := pay_conf["api_conf"]
		if len(res) > 0 {
			//将解密后的数据切割，并解析到公用参数map里面去
			err := json.Unmarshal([]byte(res), &p)
			if err != nil {
				common.LogsWithFileName(log_path, "SXPAY_init_", "error->"+err.Error())
			}
		}
	}

	this.Pay_url = p["pay_url"]
	this.Mer_code = p["mer_code"]
	this.Notice_url = p["notice_url"]
	this.Return_url = p["return_url"]
	this.Key = p["key"]
}

func (api *SXPAY) CreatePay(p *PayData) (int, string, string, string, string, map[string]string) {
	api_method := "POST"
	re_status := 100
	re_msg := "请求错误"

	img_url := ""
	param_form := map[string]string{
		"merchantName": api.Mer_code,
		"amount":       p.Amount,
		"orderId":      p.Order_number,
		"channelId":    p.Pay_bank,
		"orderIp":      p.Ip,
		"returnUrl":    api.Return_url,
		"noticeUrl":    api.Notice_url,
	}

	//拼接
	result_url := common.MapCreatLinkSort(param_form, "&", true, false)
	result_url += fmt.Sprintf("&key=%s", api.Key)

	sign := common.HexMd5(result_url)
	param_form["sign"] = sign
	param_form["signType"] = "MD5"

	//拼接并请求三方接口
	param := common.MapCreatLinkSort(param_form, "&", true, false)
	h_status, msg_b := common.HttpBody(api.Pay_url, api_method, param, http_header)

	common.LogsWithFileName(log_path, "sxpay_create_", "param->"+param+"\nmsg->"+string(msg_b))
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

	//多层map取值
	data, _ := json_res["data"].(map[string]interface{})
	url, _ := data["url"].(string)
	img_url = fmt.Sprintf("%v", url)

	if img_url == "" {
		re_msg = "接口错误"
		return re_status, re_msg, api_method, img_url, img_url, param_form
	}

	re_status = 200
	re_msg = "success"
	return re_status, re_msg, api_method, img_url, img_url, param_form
}

//回调验证
func (api *SXPAY) CallBackPay(sign, sign_str string) int {
	result := 101

	sign_str += fmt.Sprintf("&key=%s", api.Key)
	verify_sign := common.HexMd5(sign_str)

	log_str := "sign_str->" + sign_str + "\nsign->" + sign

	if verify_sign == sign {
		result = 200
	}
	common.LogsWithFileName(log_path, "SXPAY_sign_", log_str)
	return result
}

//代付
func (api *SXPAY) PayFor(pay_data map[string]string) (int, string, string) {
	//定义初始值
	api_status := 200
	pay_result := "processing"
	api_msg := "success"

	return api_status, api_msg, pay_result
}

func (api *SXPAY) PayQuery(pay_data map[string]string) (int, string, string) {
	//定义初始值
	api_status := 200
	pay_result := "processing"
	api_msg := "success"

	return api_status, api_msg, pay_result
}
