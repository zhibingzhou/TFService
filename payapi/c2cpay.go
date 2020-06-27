package payapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"TFService/model"
	"time"

	"github.com/zhibingzhou/go_public/common"
)

//c2c
type C2CPAY struct {
	Notify_url string
	Pay_url    string
	Mer_code   string
	Key        string
}

func (this *C2CPAY) Init(pay_id string) {

	pay_conf := model.ApiConfigRedis(pay_id)
	p := map[string]string{}
	//判断是否存在用户名
	if len(pay_conf["id"]) > 0 {
		res := pay_conf["api_conf"]
		if len(res) > 0 {
			//将解密后的数据切割，并解析到公用参数map里面去
			err := json.Unmarshal([]byte(res), &p)
			if err != nil {
				common.LogsWithFileName(log_path, "C2CPAY_init_", "error->"+err.Error())
			}
		}
	}

	this.Pay_url = p["pay_url"]
	this.Mer_code = p["mer_code"]
	this.Notify_url = p["notify_url"]
	this.Key = p["key"]
}

func (api *C2CPAY) CreatePay(p *PayData) (int, string, string, string, string, map[string]string) {
	api_method := "POST"
	re_status := 100
	param_form := map[string]string{}
	re_msg := "请求错误"

	img_url := ""

	return re_status, re_msg, api_method, img_url, img_url, param_form
}

//代付回调验证
func (api *C2CPAY) CallBackPay(sign, sign_str string) int {
	result := 101

	sign_str += fmt.Sprintf("&paykey=%s", api.Key)

	verify_sign := common.HexMd5(sign_str)

	log_str := "sign_str->" + sign_str + "\nsign->" + sign

	if verify_sign == sign {
		result = 200
	}
	common.LogsWithFileName(log_path, "c2cpay_sign_", log_str)
	return result
}

//代付
func (api *C2CPAY) PayFor(pay_data map[string]string) (int, string, string) {
	//定义初始值
	api_status := 100
	pay_result := "error"
	api_msg := "代付失败"

	//请求参数
	param_form := map[string]string{
		"merchant":        api.Mer_code,
		"tradeno":         pay_data["order_number"],
		"tradedate":       fmt.Sprintf(time.Now().Format(s_format)), //时间
		"bankcode":        pay_data["bank_code"],
		"tradedesc":       api.Mer_code,
		"bankaccountno":   pay_data["card_number"],
		"bankaccountname": pay_data["card_name"],
		"bankaddress":     pay_data["bank_branch"],
		"currency":        "CNY",
		"totalamount":     pay_data["amount"],
		"notifyurl":       api.Notify_url,
	}

	//拼接
	sign_str := common.MapCreatLinkSort(param_form, "&", true, true)
	sign_str += fmt.Sprintf("&paykey=%s", api.Key)

	sign := common.HexMd5(sign_str)
	param_form["sign"] = sign
	api_url := fmt.Sprintf("%s/api/generateorder", api.Pay_url)
	param := common.MapCreatLinkSort(param_form, "&", true, false)

	result, err := HttpPostForm(api_url, param_form)
	if err != nil {
		return api_status, api_msg, pay_result
	}

	common.LogsWithFileName(log_path, "c2cpay_payfor_", "param->"+param+"\nmsg->"+string(result))
	if err != nil {
		return api_status, api_msg, pay_result
	}

	var json_res map[string]interface{}
	err = json.Unmarshal(result, &json_res)
	if err != nil {
		api_msg = "json错误"
		return api_status, api_msg, pay_result
	}

	if fmt.Sprintf("%v", json_res["resultCode"]) != "0000" && fmt.Sprintf("%v", json_res["tradestatus"]) != "SUCCESS" {
		api_msg = fmt.Sprintf("%v", json_res["errMsg"])
		return api_status, api_msg, pay_result
	}

	api_status = 200
	api_msg = "success"
	pay_result = "processing"
	return api_status, api_msg, pay_result
}

func (api *C2CPAY) PayQuery(pay_data map[string]string) (int, string, string) {
	//定义初始值
	api_status := 100
	pay_result := "processing"
	api_msg := "代付失败"

	//请求参数
	param_form := map[string]string{
		"merchant": api.Mer_code,
		"tradeno":  pay_data["order_number"],
	}

	//拼接
	sign_str := common.MapCreatLinkSort(param_form, "&", true, true)
	sign_str += fmt.Sprintf("&paykey=%s", api.Key)

	sign := common.HexMd5(sign_str)
	param_form["sign"] = sign
	api_url := fmt.Sprintf("%s/api/queryorder", api.Pay_url)
	param := common.MapCreatLinkSort(param_form, "&", true, true)
	//把post form 表单提交 发送给目标服务器
	result, err := HttpPostForm(api_url, param_form)
	if err != nil {
		return api_status, api_msg, pay_result
	}
	//post form 表单提交返回值 在body 里面

	common.LogsWithFileName(log_path, "c2cpay_payquery_", "param->"+param+"\nmsg->"+string(result))

	var json_res map[string]interface{}
	err = json.Unmarshal(result, &json_res)
	if err != nil {
		api_msg = "json错误"
		return api_status, api_msg, pay_result
	}

	if fmt.Sprintf("%v", json_res["resultCode"]) != "0000" {
		api_msg = fmt.Sprintf("%v", json_res["errMsg"])
		return api_status, api_msg, pay_result
	}

	order_status := fmt.Sprintf("%v", json_res["tradestatus"])

	switch order_status {
	case "SUCCESS":
		pay_result = "success"
		api_msg = "success"
	case "FAILED", "NOT_EXIST":
		pay_result = "fail"
	}
	api_status = 200
	return api_status, api_msg, pay_result
}

//Post from 表单提交
func HttpPostForm(post_url string, param_form map[string]string) ([]byte, error) {

	data := make(url.Values)
	for key, value := range param_form {
		data[key] = []string{value}
	}
	//把post form 表单提交 发送给目标服务器
	resp, err := http.PostForm(post_url, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
