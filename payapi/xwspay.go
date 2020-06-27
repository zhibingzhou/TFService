package payapi

import (
	"encoding/json"
	"fmt"
	"TFService/model"

	"github.com/zhibingzhou/go_public/common"
)

//新万顺支付
type XWSPAY struct {
	Notify_url string
	Fxback_url string
	Pay_url    string
	Mer_code   string
	Key        string
}

//新万顺支付初始化
func (this *XWSPAY) Init(pay_id string) {

	pay_conf := model.ApiConfigRedis(pay_id)
	p := map[string]string{}
	//判断是否存在用户名
	if len(pay_conf["id"]) > 0 {
		res := pay_conf["api_conf"]
		if len(res) > 0 {
			//将解密后的数据切割，并解析到公用参数map里面去
			err := json.Unmarshal([]byte(res), &p)
			if err != nil {
				common.LogsWithFileName(log_path, "xwspay_init_", "error->"+err.Error())
			}
		}
	}

	this.Notify_url = p["notify_url"]
	this.Fxback_url = p["fxback_url"]
	this.Pay_url = p["pay_url"]
	this.Mer_code = p["mer_code"]
	this.Key = p["key"]
}

func (api *XWSPAY) CreatePay(p *PayData) (int, string, string, string, string, map[string]string) {
	api_method := "POST"
	re_status := 100
	re_msg := "请求错误"

	img_url := ""
	param_form := map[string]string{
		"fxid":        api.Mer_code,
		"fxdesc":      p.Order_number,
		"fxddh":       p.Order_number,
		"fxnotifyurl": api.Notify_url,
		"fxbackurl":   api.Fxback_url,
		"fxfee":       p.Amount,
		"fxpay":       p.Pay_bank,
		"fxip":        p.Ip,
		"fxuserid":    api.Mer_code,
	}
	img_url = fmt.Sprintf("%s/Pay", api.Pay_url)
	//拼接
	result_url := common.MapCreatLink(param_form, "fxid,fxdesc,fxfee,fxnotifyurl", "", 2)
	result_url += fmt.Sprintf("%s", api.Key)
	sign := common.HexMd5(result_url)
	param_form["fxsign"] = sign

	//请求三方接口
	param := common.MapCreatLinkSort(param_form, "&", true, false)

	//把post form 表单提交 发送给目标服务器
	result, err := HttpPostForm(img_url, param_form)
	if err != nil {
		return re_status, re_msg, api_method, img_url, img_url, param_form
	}

	common.LogsWithFileName(log_path, "xwspay_create_", "param->"+param+"\nmsg->"+string(result))
	if err != nil {
		return re_status, re_msg, api_method, img_url, img_url, param_form
	}

	var json_res map[string]interface{}
	err = json.Unmarshal(result, &json_res)
	if err != nil {
		re_msg = "json错误"
		return re_status, re_msg, api_method, img_url, img_url, param_form
	}

	if fmt.Sprintf("%v", json_res["status"]) != "1" {
		re_msg = fmt.Sprintf("%v", json_res["error"])
		return re_status, re_msg, api_method, img_url, img_url, param_form
	}

	img_url = fmt.Sprintf("%v", json_res["payurl"])

	if img_url == "" {
		re_msg = "接口错误"
		return re_status, re_msg, api_method, img_url, img_url, param_form
	}
	fmt.Println(img_url)
	re_status = 200
	re_msg = "success"
	return re_status, re_msg, api_method, img_url, img_url, param_form
}

//回调验证
func (api *XWSPAY) CallBackPay(sign, sign_str string) int {
	result := 101

	verify_sign := common.HexMd5(sign_str + api.Key)
	log_str := "sign_str->" + sign_str + "\nsign->" + sign

	if verify_sign == sign {
		result = 200
	}
	common.LogsWithFileName(log_path, "xwspay_sign_", log_str)
	return result
}

//代付
func (api *XWSPAY) PayFor(pay_data map[string]string) (int, string, string) {
	//定义初始值
	api_status := 200
	pay_result := "processing"
	api_msg := "success"

	return api_status, api_msg, pay_result
}

func (api *XWSPAY) PayQuery(pay_data map[string]string) (int, string, string) {
	//定义初始值
	api_status := 200
	pay_result := "processing"
	api_msg := "success"

	return api_status, api_msg, pay_result
}
