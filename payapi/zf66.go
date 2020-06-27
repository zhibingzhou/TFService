package payapi

import (
	"encoding/json"
	"fmt"
	"TFService/model"
	"strconv"
	"strings"

	"github.com/zhibingzhou/go_public/common"
)

type ZF66 struct {
	Pay_url    string
	Mer_code   string
	Notify_url string
	Key        string
	Back_url   string
}

func (this *ZF66) Init(pay_id string) {

	pay_conf := model.ApiConfigRedis(pay_id)
	p := map[string]string{}
	//判断是否存在用户名
	if len(pay_conf["id"]) > 0 {
		res := pay_conf["api_conf"]
		if len(res) > 0 {
			//将解密后的数据切割，并解析到公用参数map里面去
			err := json.Unmarshal([]byte(res), &p)
			if err != nil {
				common.LogsWithFileName(log_path, "ZF66_init_", "error->"+err.Error())
			}
		}
	}

	this.Pay_url = p["pay_url"]
	this.Mer_code = p["mer_code"]
	this.Notify_url = p["notify_url"]
	this.Key = p["key"]
	this.Back_url = p["back_url"]
}

func (api *ZF66) CreatePay(p *PayData) (int, string, string, string, string, map[string]string) {
	api_method := "get"
	re_status := 200
	param_form := map[string]string{}
	re_msg := "success"

	nonce := common.Random("smallnumber", 35)
	amount_f, _ := strconv.ParseFloat(p.Amount, 64)
	amount := fmt.Sprintf("%.0f", amount_f)

	sign_str := fmt.Sprintf("nonce=%s&ordercode=%s&paytype=%s&shopid=%s&total=%s&key=%s", nonce, p.Order_number, p.Pay_bank, api.Mer_code, amount, api.Key)

	sign := common.HexMd5(sign_str)
	sign = strings.ToUpper(sign)
	param_form["nonce"] = nonce
	param_form["ordercode"] = p.Order_number
	param_form["paytype"] = p.Pay_bank
	param_form["shopid"] = api.Mer_code
	param_form["total"] = amount
	param_form["sign"] = sign
	img_url := fmt.Sprintf("%s/pay.aspx", api.Pay_url)

	common.LogsWithFileName(log_path, "zf66_create_", "img_url->"+img_url+"\nsign_str->"+sign_str)

	return re_status, re_msg, api_method, img_url, img_url, param_form
}

//回调验证
func (api *ZF66) CallBackPay(sign, sign_str string) int {
	result := 101

	verify_sign := common.HexMd5(sign_str + "&key=" + api.Key)
	verify_sign = strings.ToUpper(verify_sign)
	log_str := "sign_str->" + sign_str + "\nsign->" + sign

	if verify_sign == sign {
		result = 200
	}
	common.LogsWithFileName(log_path, "zf66_sign_", log_str)
	return result
}

//代付
func (api *ZF66) PayFor(pay_data map[string]string) (int, string, string) {
	//定义初始值
	api_status := 200
	pay_result := "processing"
	api_msg := "success"

	return api_status, api_msg, pay_result
}

func (api *ZF66) PayQuery(pay_data map[string]string) (int, string, string) {
	//定义初始值
	api_status := 200
	pay_result := "processing"
	api_msg := "success"

	return api_status, api_msg, pay_result
}
