package payapi

import (
	"encoding/json"
	"fmt"
	"TFService/model"
	"strings"
	"time"

	"github.com/zhibingzhou/go_public/common"
)

type RLPAY struct {
	Pay_url    string
	Mer_code   string
	Notify_url string
	Key        string
	Return_url string
	Header     map[string]string
}

func (this *RLPAY) Init(pay_id string) {

	pay_conf := model.ApiConfigRedis(pay_id)
	p := map[string]string{}
	//判断是否存在用户名
	if len(pay_conf["id"]) > 0 {
		res := pay_conf["api_conf"]
		if len(res) > 0 {
			//将解密后的数据切割，并解析到公用参数map里面去
			err := json.Unmarshal([]byte(res), &p)
			if err != nil {
				common.LogsWithFileName(log_path, "RLPAY_init_", "error->"+err.Error())
			}
		}
	}

	this.Pay_url = p["pay_url"]
	this.Mer_code = p["mer_code"]
	this.Notify_url = p["notify_url"]
	this.Key = p["key"]
	this.Return_url = p["return_url"]
	this.Header = make(map[string]string)
	this.Header["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36"
	this.Header["Content-Type"] = "application/x-www-form-urlencoded; charset=UTF-8"
}

func (api *RLPAY) CreatePay(p *PayData) (int, string, string, string, string, map[string]string) {
	api_method := "POST"
	re_status := 200
	param_form := map[string]string{}
	re_msg := "success"

	img_url := ""
	now_date := time.Now().Format(date_format)
	goods_name := common.Random("smallnumber", 9)
	sign_str := fmt.Sprintf("amount=%s&bank_code=%s&callback_url=%s&customer_no=%s&customer_order=%s&notify_url=%s&produce_date=%s", p.Amount, p.Pay_bank, api.Return_url, api.Mer_code, p.Order_number, api.Notify_url, now_date)
	sign := common.HexMd5(sign_str + "&key=" + api.Key)
	sign = strings.ToUpper(sign)
	param := fmt.Sprintf("goods_name=%s&%s&sign_md5=%s", goods_name, sign_str, sign)
	param_form["amount"] = p.Amount
	param_form["bank_code"] = p.Pay_bank
	param_form["callback_url"] = api.Return_url
	param_form["customer_no"] = api.Mer_code
	param_form["customer_order"] = p.Order_number
	param_form["notify_url"] = api.Notify_url
	param_form["produce_date"] = now_date
	param_form["goods_name"] = goods_name
	param_form["sign_md5"] = sign
	img_url = fmt.Sprintf("%s/Pay_Defray.html", api.Pay_url)

	common.LogsWithFileName(log_path, "rlpay_create_", "param->"+param+"\nurl->"+img_url)

	return re_status, re_msg, api_method, img_url, img_url, param_form
}

//回调验证
func (api *RLPAY) CallBackPay(sign, sign_str string) int {
	result := 101

	verify_sign := common.HexMd5(sign_str + "&key=" + api.Key)
	verify_sign = strings.ToUpper(verify_sign)
	log_str := "sign_str->" + sign_str + "\nsign->" + sign + "\nverify_sign->" + verify_sign

	if verify_sign == sign {
		result = 200
	}
	common.LogsWithFileName(log_path, "RLPAY_sign_", log_str)
	return result
}

//代付
func (api *RLPAY) PayFor(pay_data map[string]string) (int, string, string) {
	//定义初始值
	api_status := 200
	pay_result := "processing"
	api_msg := "代付失败"

	return api_status, api_msg, pay_result
}

func (api *RLPAY) PayQuery(pay_data map[string]string) (int, string, string) {
	//定义初始值
	api_status := 200
	pay_result := "processing"
	api_msg := "success"

	return api_status, api_msg, pay_result
}
