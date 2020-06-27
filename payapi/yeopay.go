package payapi

import (
	"encoding/json"
	"fmt"
	"TFService/model"
	"strings"
	"time"

	"github.com/zhibingzhou/go_public/common"
)

type YEOPAY struct {
	Mer_code   string //商家号
	Notify_url string //服务端通知
	Return_url string //页面跳转通知
	Key        string //秘钥
	Pay_url    string
}

/**
* 对象初始化
* @param	customer_no 		商家号
* @param	pay_url		支付地址
* @param	notify_url		服务端通知
* @param	callback_url		页面跳转通知
* @param	key	秘钥
 */
func (api *YEOPAY) Init(pay_id string) {
	pay_conf := model.ApiConfigRedis(pay_id)
	p := map[string]string{}
	//判断是否存在用户名
	if len(pay_conf["id"]) > 0 {
		res := pay_conf["api_conf"]
		if len(res) > 0 {
			//将解密后的数据切割，并解析到公用参数map里面去
			err := json.Unmarshal([]byte(res), &p)
			if err != nil {
				common.LogsWithFileName(log_path, "yeopay_init_", "error->"+err.Error())
			}
		}
	}
	api.Mer_code = p["mer_code"]
	api.Notify_url = p["notify_url"]
	api.Return_url = p["return_url"]
	api.Key = p["key"]
	api.Pay_url = p["pay_url"]

}

/**
* 发出支付请求
* @param	*PayData	支付信息的指针值
*
* return	string	需要提交的参数
* return	map	需要用于表单的内容
 */
func (api *YEOPAY) CreatePay(p *PayData) (int, string, string, string, string, map[string]string) {
	api_method := "post"
	re_status := 200
	img_url := ""
	re_msg := "success"

	goods_name := "pay"
	produce_date := time.Now().Format(date_format)
	//sign参数
	input_str := fmt.Sprintf("amount=%s&bank_code=%s&callback_url=%s&customer_no=%s&customer_order=%s&notify_url=%s&produce_date=%s", p.Amount, p.Pay_bank, api.Return_url, api.Mer_code, p.Order_number, api.Notify_url, produce_date)
	//sign+key
	sign_str := fmt.Sprintf("%s&key=%s", input_str, api.Key)
	//md5加密
	sign := strings.ToUpper(common.HexMd5(sign_str))

	//请求参数
	param_result := fmt.Sprintf("%s&sign_md5=%s&goods_name=%s", input_str, sign, goods_name)
	//将支付请求写入日志
	common.LogsWithFileName(log_path, "yeopay_create_", "param_result->"+param_result+"\nsign_str->"+sign_str)

	//html
	param_form := map[string]string{}
	param_form["amount"] = p.Amount
	param_form["bank_code"] = p.Pay_bank
	param_form["callback_url"] = api.Return_url
	param_form["customer_no"] = api.Mer_code
	param_form["customer_order"] = p.Order_number
	param_form["notify_url"] = api.Notify_url
	param_form["produce_date"] = produce_date
	param_form["sign_md5"] = sign
	param_form["goods_name"] = goods_name

	img_url = fmt.Sprintf("%s/Pay_Index.html", api.Pay_url)
	return re_status, re_msg, api_method, img_url, img_url, param_form
}

/**
* 被动接收返回值
 */
func (api *YEOPAY) CallBackPay(sign, sign_str string) int {
	//sign是返回的签名
	//sign_str是需要验证签名的参数
	result := 100
	md5_str := sign_str + "&key=" + api.Key
	verify_sign := strings.ToUpper(common.HexMd5(md5_str))

	log_str := "sign_str->" + md5_str + "\nsign->" + sign + "\nverify_sign->" + verify_sign
	//如果签名相同
	if verify_sign == sign {
		result = 200
	}
	common.LogsWithFileName(log_path, "yeopay_sign_", log_str)
	return result
}

/**
*  代付
 */
func (api *YEOPAY) PayFor(pay_data map[string]string) (int, string, string) {
	//定义初始值
	api_status := 200
	pay_result := "processing"
	api_msg := "success"

	return api_status, api_msg, pay_result
}

func (api *YEOPAY) PayQuery(pay_data map[string]string) (int, string, string) {
	//定义初始值
	api_status := 200
	pay_result := "processing"
	api_msg := "success"

	return api_status, api_msg, pay_result
}
