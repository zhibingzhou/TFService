package payapi

import (
	"encoding/json"
	"fmt"
	"TFService/model"
	"strconv"
	"strings"
	"time"

	"github.com/zhibingzhou/go_public/common"
)

//超凡支付
type CFPAY struct {
	Pay_url    string
	Mer_code   string
	Notify_url string
	Key        string
	Return_url string
}

//超凡支付
func (this *CFPAY) Init(pay_id string) {

	pay_conf := model.ApiConfigRedis(pay_id)
	p := map[string]string{}
	//判断是否存在用户名
	if len(pay_conf["id"]) > 0 {
		res := pay_conf["api_conf"]
		if len(res) > 0 {
			//将解密后的数据切割，并解析到公用参数map里面去
			err := json.Unmarshal([]byte(res), &p)
			if err != nil {
				common.LogsWithFileName(log_path, "CFPAY_init_", "error->"+err.Error())
			}
		}
	}

	this.Pay_url = p["pay_url"]
	this.Mer_code = p["mer_code"]
	this.Notify_url = p["notify_url"]
	this.Return_url = p["return_url"]
	this.Key = p["key"]
}

func (api *CFPAY) CreatePay(p *PayData) (int, string, string, string, string, map[string]string) {
	api_method := "POST"
	re_status := 100
	re_msg := "请求错误"

	img_url := ""
	param_form := map[string]string{}

	amount, err := RmbTranfer(p.Amount, false)
	if err != nil {
		return re_status, re_msg, api_method, img_url, img_url, param_form
	}

	param_form = map[string]string{
		"amount":     amount,
		"orderNo":    p.Order_number,
		"notifyUrl":  api.Notify_url,
		"merchantId": api.Mer_code,
		"version":    "1.0",
		"clientIp":   p.Ip,
		"service":    p.Pay_bank,
		"key":        api.Key,
		"tradeDate":  fmt.Sprintf(time.Now().Format("20060102")),
		"tradeTime":  fmt.Sprintf(time.Now().Format("150405")),
	}

	//拼接
	result_url := common.MapCreatLinkSort(param_form, "&", true, false)

	sign := common.HexMd5(result_url)
	param_form["sign"] = sign
	delete(param_form, "key")

	param := common.MapCreatLinkSort(param_form, "&", true, false)
	h_status, msg_b := common.HttpBody(api.Pay_url, api_method, param, http_header)

	common.LogsWithFileName(log_path, "cfpay_create_", "param->"+param+"\nmsg->"+string(msg_b))
	if h_status != 200 {
		return re_status, re_msg, api_method, img_url, img_url, param_form
	}

	var json_res map[string]string
	json_res = UrlToMap(string(msg_b))
	if fmt.Sprintf("%v", json_res["repCode"]) != "0001" {
		re_msg = fmt.Sprintf("%v", json_res["repMsg"])
		return re_status, re_msg, api_method, img_url, img_url, param_form
	}

	img_url = fmt.Sprintf("%v", json_res["resultUrl"])

	if img_url == "" {
		re_msg = "接口错误"
		return re_status, re_msg, api_method, img_url, img_url, param_form
	}

	re_status = 200
	re_msg = "success"
	return re_status, re_msg, api_method, img_url, img_url, param_form
}

//回调验证
func (api *CFPAY) CallBackPay(sign, sign_str string) int {
	result := 101

	var json_res map[string]string
	err := json.Unmarshal([]byte(sign_str), &json_res)
	if err != nil {
		return result
	}
	json_res["key"] = api.Key

	param := common.MapCreatLinkSort(json_res, "&", true, false)

	verify_sign := common.HexMd5(param)

	log_str := "sign_str->" + param + "\nsign->" + sign

	if verify_sign == sign {
		result = 200
	}
	common.LogsWithFileName(log_path, "cfpay_sign_", log_str)
	return result
}

//代付
func (api *CFPAY) PayFor(pay_data map[string]string) (int, string, string) {
	//定义初始值
	api_status := 200
	pay_result := "processing"
	api_msg := "success"

	return api_status, api_msg, pay_result
}

func (api *CFPAY) PayQuery(pay_data map[string]string) (int, string, string) {
	//定义初始值
	api_status := 200
	pay_result := "processing"
	api_msg := "success"

	return api_status, api_msg, pay_result
}

//url参数转map
func UrlToMap(url string) map[string]string {
	result := make(map[string]string)
	parm := strings.Split(url, "&")
	for _, value := range parm {
		keys := strings.Split(value, "=")
		result[keys[0]] = keys[1]
	}
	return result
}

//人民转换
//dic 分转元  传true ，元转分  传false
func RmbTranfer(amount string, dic bool) (string, error) {
	ramount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return "", err
	}
	if dic == false {
		ramount = ramount * 1000 // 元转分
		return strconv.FormatFloat(float64(ramount), 'f', 0, 64), nil
	} else {
		ramount = ramount / 1000 // 分转元,保留两们小数
		return strconv.FormatFloat(float64(ramount), 'f', 2, 64), nil
	}

}
