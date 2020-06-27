package payapi

import (
	"encoding/json"
	"fmt"
	"TFService/model"

	"github.com/zhibingzhou/go_public/common"
)

type DAY88 struct {
	Pay_url    string
	Mer_code   string
	Notify_url string
	Key        string
	Back_url   string
	Sign_type  string
}

func (this *DAY88) Init(pay_id string) {

	pay_conf := model.ApiConfigRedis(pay_id)
	p := map[string]string{}
	//判断是否存在用户名
	if len(pay_conf["id"]) > 0 {
		res := pay_conf["api_conf"]
		if len(res) > 0 {
			//将解密后的数据切割，并解析到公用参数map里面去
			err := json.Unmarshal([]byte(res), &p)
			if err != nil {
				common.LogsWithFileName(log_path, "day88_init_", "error->"+err.Error())
			}
		}
	}

	this.Pay_url = p["pay_url"]
	this.Mer_code = p["mer_code"]
	this.Notify_url = p["notify_url"]
	this.Key = p["key"]
	this.Back_url = p["back_url"]
	this.Sign_type = p["sign_type"]

}

func (api *DAY88) CreatePay(p *PayData) (int, string, string, string, string, map[string]string) {
	api_method := "get"
	re_status := 200
	param_form := map[string]string{}
	re_msg := "success"

	nonce := common.Random("smallnumber", 8)
	sitename := common.Random("smallnumber", 9)
	img_url := ""

	sign_str := fmt.Sprintf("money=%s&name=%s&notify_url=%s&out_trade_no=%s&pid=%s&return_url=%s&sitename=%s&type=%s%s", p.Amount, nonce, api.Notify_url, p.Order_number, api.Mer_code, api.Back_url, sitename, p.Pay_bank, api.Key)

	sign := common.HexMd5(sign_str)
	param_form["money"] = p.Amount
	param_form["name"] = nonce
	param_form["notify_url"] = api.Notify_url
	param_form["out_trade_no"] = p.Order_number
	param_form["pid"] = api.Mer_code
	param_form["return_url"] = api.Back_url
	param_form["sitename"] = sitename
	param_form["type"] = p.Pay_bank
	param_form["sign_type"] = api.Sign_type
	param_form["sign"] = sign
	img_url = fmt.Sprintf("%s/submit.php", api.Pay_url)

	common.LogsWithFileName(log_path, "day88_create_", "img_url->"+img_url)

	return re_status, re_msg, api_method, img_url, img_url, param_form
}

//回调验证
func (api *DAY88) CallBackPay(sign, sign_str string) int {
	result := 101

	verify_sign := common.HexMd5(sign_str + api.Key)
	log_str := "sign_str->" + sign_str + "\nsign->" + sign

	if verify_sign == sign {
		result = 200
	}
	common.LogsWithFileName(log_path, "day88_sign_", log_str)
	return result
}

//代付
func (api *DAY88) PayFor(pay_data map[string]string) (int, string, string) {
	//定义初始值
	api_status := 200
	pay_result := "processing"
	api_msg := "success"
	return api_status, api_msg, pay_result
}

func (api *DAY88) PayQuery(pay_data map[string]string) (int, string, string) {
	//定义初始值
	api_status := 200
	pay_result := "processing"
	api_msg := "success"

	return api_status, api_msg, pay_result
}
