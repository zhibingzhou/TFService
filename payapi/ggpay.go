package payapi

import (
	"encoding/json"
	"fmt"
	"TFService/model"

	"github.com/zhibingzhou/go_public/common"
)

//硅谷支付
type GGPAY struct {
	Pay_url    string
	Mer_code   string
	Notify_url string
	Key        string
}

//硅谷初始化
func (this *GGPAY) Init(pay_id string) {

	pay_conf := model.ApiConfigRedis(pay_id)
	p := map[string]string{}
	//判断是否存在用户名
	if len(pay_conf["id"]) > 0 {
		res := pay_conf["api_conf"]
		if len(res) > 0 {
			//将解密后的数据切割，并解析到公用参数map里面去
			err := json.Unmarshal([]byte(res), &p)
			if err != nil {
				common.LogsWithFileName(log_path, "GGPAY_init_", "error->"+err.Error())
			}
		}
	}

	this.Notify_url = p["notify_url"]
	this.Pay_url = p["pay_url"]
	this.Mer_code = p["mer_code"]
	this.Key = p["key"]
}

func (api *GGPAY) CreatePay(p *PayData) (int, string, string, string, string, map[string]string) {
	api_method := "POST"
	re_status := 100
	re_msg := "请求错误"

	img_url := api.Pay_url
	param_form := map[string]string{
		"partner":   api.Mer_code,
		"amount":    p.Amount,
		"tradeNo":   p.Order_number,
		"notifyUrl": api.Notify_url,
		"service":   p.Pay_bank,
	}

	//拼接
	result_url := common.MapCreatLinkSort(param_form, "&", true, false)
	result_url += fmt.Sprintf("&key=%s", api.Key)
	sign := common.HexMd5(result_url)
	param_form["sign"] = sign

	re_status = 200
	re_msg = "success"
	return re_status, re_msg, api_method, img_url, img_url, param_form
}

//回调验证
func (api *GGPAY) CallBackPay(sign, sign_str string) int {
	result := 101

	sign_str += fmt.Sprintf("&key=%s", api.Key)
	verify_sign := common.HexMd5(sign_str)

	log_str := "sign_str->" + sign_str + "\nsign->" + sign

	if verify_sign == sign {
		result = 200
	}
	common.LogsWithFileName(log_path, "GGPAY_sign_", log_str)
	return result
}

//代付
func (api *GGPAY) PayFor(pay_data map[string]string) (int, string, string) {
	//定义初始值
	api_status := 200
	pay_result := "processing"
	api_msg := "success"

	return api_status, api_msg, pay_result
}

func (api *GGPAY) PayQuery(pay_data map[string]string) (int, string, string) {
	//定义初始值
	api_status := 200
	pay_result := "processing"
	api_msg := "success"

	return api_status, api_msg, pay_result
}
