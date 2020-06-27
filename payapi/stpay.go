package payapi

import (
	"encoding/json"
	"fmt"
	"TFService/model"
	"strings"

	"github.com/zhibingzhou/go_public/common"
)

type STPAY struct {
	Mer_code   string //商户账号
	Return_url string //同步跳转
	Notify_url string //异步回调
	Aes_key    string //AES key
	Key        string //商户key
	Pay_url    string //支付地址
}

/**
* 对象初始化
 */

func (api *STPAY) Init(pay_id string) {
	pay_conf := model.ApiConfigRedis(pay_id)
	p := map[string]string{}
	//判断是否存在用户名
	if len(pay_conf["id"]) > 0 {
		res := pay_conf["api_conf"]
		if len(res) > 0 {
			//将解密后的数据切割，并解析到公用参数map里面去
			err := json.Unmarshal([]byte(res), &p)
			if err != nil {
				common.LogsWithFileName(log_path, "stpay_init_", "error->"+err.Error())
			}
		}
	}
	api.Mer_code = p["mer_code"]
	api.Notify_url = p["notify_url"]
	api.Return_url = p["return_url"]
	api.Key = p["key"]
	api.Aes_key = p["aes_key"]
	api.Pay_url = p["pay_url"]
}

/**
* 发出支付请求
* @param	*PayData	支付信息的指针值
* return	string	需要提交的参数
* return	map	需要用于表单的内容
 */
func (api *STPAY) CreatePay(p *PayData) (int, string, string, string, string, map[string]string) {
	re_status := 100
	re_msg := "支付失败"
	api_method := "POST"
	img_url := ""
	param_form := map[string]string{}

	action := "deposit"
	data_str := fmt.Sprintf(`{"RequestOrderNo":"%s","BankCode":"%s","Money":"%s","CallBackUrl":"%s","UserLevel":"1"}`, p.Order_number, p.Pay_bank, p.Amount, api.Notify_url)
	aes := common.SetAESECB(api.Aes_key, "", "", "hex", 32)
	//data值
	data := strings.ToUpper(aes.AesEncryptString(data_str))

	sign_str := fmt.Sprintf("Action=%s&Data=%s&MerchantId=%s&Key=%s", action, data, api.Mer_code, api.Key)
	//sign值
	sign := common.HexMd5(sign_str)

	param_result := fmt.Sprintf("MerchantId=%s&Action=%s&Data=%s&Sign=%s", api.Mer_code, action, data, sign)

	pay_url := fmt.Sprintf("%s/action", api.Pay_url)
	t_status, msg_b := common.HttpBody(pay_url, api_method, param_result, http_header)
	common.LogsWithFileName(log_path, "stpay_create_", "param_result->"+param_result+"\ndata_str->"+data_str+"\nsign_str"+sign_str+"\nmsg->"+string(msg_b))
	if t_status != 200 {
		re_msg = "支付请求错误"
		return re_status, re_msg, api_method, api.Pay_url, img_url, param_form
	}
	json_res := make(map[string]interface{})
	err := json.Unmarshal(msg_b, &json_res)
	if err != nil {
		re_msg = err.Error()
		return re_status, re_msg, api_method, api.Pay_url, img_url, param_form
	}
	//捞取结果，赋值到变量
	Code := fmt.Sprintf("%v", json_res["Code"])
	result := fmt.Sprintf("%v", json_res["Result"])
	re_msg = fmt.Sprintf("%v", json_res["ErrMsg"])
	if Code != "0" {
		return re_status, re_msg, api_method, api.Pay_url, img_url, param_form
	}
	url_body := aes.AesDecryptString(result)
	url_res := make(map[string]interface{})
	err = json.Unmarshal([]byte(url_body), &url_res)

	if err != nil {
		return re_status, re_msg, api_method, api.Pay_url, img_url, param_form
	}
	img_url = url_res["Url"].(string)
	re_status = 200
	re_msg = "success"
	return re_status, re_msg, api_method, api.Pay_url, img_url, param_form
}

/**
* 被动接收返回值
 */
func (api *STPAY) CallBackPay(sign, sign_str string) int {
	result := 100
	sign_code := sign_str + "&Key=" + api.Key
	// //对字符串进行URLEncode
	// sign_code := url.QueryEscape(md5_str)
	check_sign := common.HexMd5(sign_code)
	log_str := "sign_str->" + sign_str + "\nsign_code" + sign_code + "\nsign->" + sign + "\ncheck_sign" + check_sign
	//如果签名相同
	if check_sign == sign {
		result = 200
	}
	common.LogsWithFileName(log_path, "stpay_sign_", log_str)
	return result
}

/**
*  代付
 */
func (api *STPAY) PayFor(pay_data map[string]string) (int, string, string) {
	//定义初始值
	api_status := 200
	pay_result := "processing"
	api_msg := "success"

	return api_status, api_msg, pay_result
}

func (api *STPAY) PayQuery(pay_data map[string]string) (int, string, string) {
	//定义初始值
	api_status := 200
	pay_result := "processing"
	api_msg := "success"

	return api_status, api_msg, pay_result
}
