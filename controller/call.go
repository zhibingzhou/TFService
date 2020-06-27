package controller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"TFService/thread"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zhibingzhou/go_public/common"
)

func CallDypay(c *gin.Context) {
	res := "error"
	json_res := make(map[string]string)
	body_str := ""
	sign := ""
	sign_str := ""
	if c.Request.Body != nil {
		body := make([]byte, c.Request.ContentLength)
		c.Request.Body.Read(body)

		body_str = string(body)
		//解析body
		body_arr := strings.Split(body_str, "&")
		if len(body_arr) > 0 {
			for _, b_val := range body_arr {
				val_arr := strings.Split(b_val, "=")
				if len(val_arr) > 1 {
					json_res[val_arr[0]] = val_arr[1]
				}
			}
		}

		if len(json_res) > 1 {
			amount := json_res["Amount"]
			order_line := json_res["OrderLine"]
			merchant_id := json_res["MerchantId"]
			merchant_order_id := json_res["MerchantOrderId"]
			timestamp := json_res["Timestamp"]
			sign := json_res["Sign"]
			call_status := 101
			sign_str = fmt.Sprintf("Amount=%s&OrderLine=%s&MerchantId=%s&MerchantOrderId=%s&Timestamp=%s", amount, order_line, merchant_id, merchant_order_id, timestamp)

			is_cent := 0
			call_status, res = thread.ThreadUpdatePay(merchant_order_id, order_line, amount, sign, sign_str, is_cent)

			if call_status == 200 {
				res = "SUCCESS"
			}
		}
	}

	common.LogsWithFileName(log_path, "dypay_call_",
		"body->"+body_str+"\nsign->"+sign+"\nsign_str->"+sign_str+"\nres->"+res)

	c.Writer.WriteString(res)
}

func CallZf66(c *gin.Context) {
	nonce := c.Query("nonce")
	ordercode := c.Query("ordercode")
	returncode := c.Query("returncode")
	total := c.Query("total")
	sign := c.Query("sign")
	res := returncode
	sign_str := ""
	if returncode == "SUCCESS" {
		call_status := 101
		sign_str = fmt.Sprintf("nonce=%s&ordercode=%s&returncode=%s&total=%s", nonce, ordercode, returncode, total)
		is_cent := 0
		call_status, res = thread.ThreadUpdatePay(ordercode, ordercode, total, sign, sign_str, is_cent)
		if call_status == 200 {
			res = "SUCCESS"
		}
	}

	common.LogsWithFileName(log_path, "zf66_call_", "sign->"+sign+"\nsign_str->"+sign_str+"\nres->"+res)

	c.Writer.WriteString(res)
}

func CallIpay(c *gin.Context) {
	pid := c.Query("pid")
	cid := c.Query("cid")
	oid := c.Query("oid")
	sid := c.Query("sid")
	uid := c.Query("uid")
	amount := c.Query("amount")
	ramount := c.Query("ramount")
	stime := c.Query("stime")
	code := c.Query("code")
	sign := c.Query("sign")
	res := code
	sign_str := ""
	if code == "101" {
		call_status := 101
		sign_str = fmt.Sprintf("pid=%s&cid=%s&oid=%s&sid=%s&uid=%s&amount=%s&stime=%s&code=%s", pid, cid, oid, sid, uid, amount, stime, code)
		is_cent := 1
		call_status, res = thread.ThreadUpdatePay(oid, sid, ramount, sign, sign_str, is_cent)
		if call_status == 200 {
			res = "Success"
		}
	}

	common.LogsWithFileName(log_path, "ipay_call_", "sign->"+sign+"\nsign_str->"+sign_str+"\nres->"+res)

	c.Writer.WriteString(res)
}

func CallDay88(c *gin.Context) {
	pid := c.Query("pid")
	trade_no := c.Query("trade_no")
	out_trade_no := c.Query("out_trade_no")
	type_str := c.Query("type")
	name := c.Query("name")
	money := c.Query("money")
	trade_status := c.Query("trade_status")
	sign := c.Query("sign")
	res := trade_status
	sign_str := ""
	if trade_status == "TRADE_SUCCESS" {
		call_status := 101
		sign_str = fmt.Sprintf("money=%s&name=%s&out_trade_no=%s&pid=%s&trade_no=%s&trade_status=%s&type=%s", money, name, out_trade_no, pid, trade_no, trade_status, type_str)
		is_cent := 0
		call_status, res = thread.ThreadUpdatePay(out_trade_no, trade_no, money, sign, sign_str, is_cent)
		if call_status == 200 {
			res = "success"
		}
	}

	common.LogsWithFileName(log_path, "day88_call_", "sign->"+sign+"\nsign_str->"+sign_str+"\nres->"+res)

	c.Writer.WriteString(res)
}

func CallHytpay(c *gin.Context) {
	amount := c.PostForm("Amount")
	merchant_id := c.PostForm("MerchantId")
	o_id := c.PostForm("MerchantUniqueOrderId")
	timestamp := c.PostForm("Timestamp")
	sign := c.PostForm("Sign")

	res := "error"

	call_status := 101
	sign_str := fmt.Sprintf("Amount=%s&MerchantId=%s&MerchantUniqueOrderId=%s&Timestamp=%s", amount, merchant_id, o_id, timestamp)
	is_cent := 0
	call_status, res = thread.ThreadUpdatePay(o_id, o_id, amount, sign, sign_str, is_cent)
	if call_status == 200 {
		res = "SUCCESS"
	}

	common.LogsWithFileName(log_path, "hytpay_call_", "sign->"+sign+"\nsign_str->"+sign_str+"\nres->"+res)

	c.Writer.WriteString(res)
}

func CallNewpay(c *gin.Context) {
	body := c.PostForm("body")
	sign := c.PostForm("sign")

	res := "error"

	var json_res map[string]interface{}
	err := json.Unmarshal([]byte(body), &json_res)
	if err == nil {
		res = "SUCCESS"
		status := fmt.Sprintf("%v", json_res["status"])
		if status == "success" {
			actual_amount := fmt.Sprintf("%v", json_res["actualAmount"])
			if !strings.Contains(actual_amount, ".") {
				actual_amount = actual_amount + ".00"
			}
			amount := fmt.Sprintf("%v", json_res["amount"])
			amount_f, _ := strconv.ParseFloat(amount, 64)
			amount = fmt.Sprintf("%.2f", amount_f)
			trade_no := fmt.Sprintf("%v", json_res["outTradeNo"])
			sn := fmt.Sprintf("%v", json_res["sn"])
			sign_str := fmt.Sprintf(`body={"actualAmount":%s,"amount":%s,"outTradeNo":"%s","payTime":"%v","sn":"%s","status":"%s"}`, actual_amount, amount, trade_no, json_res["payTime"], sn, status)
			is_cent := 0
			thread.ThreadUpdatePay(trade_no, sn, amount, sign, sign_str, is_cent)
		}
	}

	common.LogsWithFileName(log_path, "newpay_call_", "sign->"+sign+"\nbody->"+body+"\nres->"+res)

	c.Writer.WriteString(res)
}

func CallJxpay(c *gin.Context) {
	out_trade_no := c.PostForm("out_trade_no")
	input_charset := c.PostForm("input_charset")
	trade_no := c.PostForm("trade_no")
	amount := c.PostForm("amount")
	notify_time := c.PostForm("notify_time")
	notify_id := c.PostForm("notify_id")
	status := c.PostForm("status")
	sign := c.PostForm("sign")

	res := "success"

	sign_str := fmt.Sprintf("amount=%s&input_charset=%s&notify_id=%s&notify_time=%s&out_trade_no=%s&status=%s&trade_no=%s", amount, input_charset, notify_id, notify_time, out_trade_no, status, trade_no)
	if status == "SUCCESS" {
		is_cent := 0
		c_status, c_msg := thread.ThreadUpdatePay(out_trade_no, trade_no, amount, sign, sign_str, is_cent)
		if c_status != 200 {
			res = c_msg
		}
	}

	common.LogsWithFileName(log_path, "jxpay_call_", "sign->"+sign+"\nsign_str>"+sign_str+"\nres->"+res)

	c.Writer.WriteString(res)
}

func CallRlpay(c *gin.Context) {
	customer_no := c.PostForm("customer_no")
	customer_order := c.PostForm("customer_order")
	amount := c.PostForm("amount")
	trading_num := c.PostForm("trading_num")
	trading_time := c.PostForm("trading_time")
	trading_code := c.PostForm("trading_code")
	sign_md5 := c.PostForm("sign_md5")

	sign_str := fmt.Sprintf("amount=%s&customer_no=%s&customer_order=%s&trading_code=%s&trading_num=%s&trading_time=%s", amount, customer_no, customer_order, trading_code, trading_num, trading_time)
	res := "OK"
	if trading_code == "00" {
		is_cent := 0
		c_status, c_msg := thread.ThreadUpdatePay(customer_order, trading_num, amount, sign_md5, sign_str, is_cent)
		if c_status != 200 {
			res = c_msg
		}
	}

	common.LogsWithFileName(log_path, "rlpay_call_", "sign->"+sign_md5+"\nsign_str>"+sign_str+"\nres->"+res)

	c.Writer.WriteString(res)
}

/**
* yeo的支付回调接口
 */
func CallYeopay(c *gin.Context) {

	customer_no := c.PostForm("customer_no")       //商户编号
	customer_order := c.PostForm("customer_order") //订单号
	amount := c.PostForm("amount")                 //订单金额
	trading_num := c.PostForm("trading_num")       //交易流水号
	trading_time := c.PostForm("trading_time")     // 交易时间
	trading_code := c.PostForm("trading_code")     //交易状态
	sign_md5 := c.PostForm("sign_md5")             //签名

	res := "error"

	//需要验证签名的参数
	sign_str := fmt.Sprintf("amount=%s&customer_no=%s&customer_order=%s&trading_code=%s&trading_num=%s&trading_time=%s", amount, customer_no, customer_order, trading_code, trading_num, trading_time)

	if trading_code == "00" {
		res = "OK"
		is_cent := 0
		c_status, c_msg := thread.ThreadUpdatePay(customer_order, trading_num, amount, sign_md5, sign_str, is_cent)
		if c_status != 200 {
			res = c_msg
		}
	}

	common.LogsWithFileName(log_path, "yeopay_call_", "sign_md5->"+sign_md5+"\nsign_str->"+sign_str+"\nres->"+res)

	c.Writer.WriteString(res)
}

/**
* stpay的回调接口
* 隐式回调接口，就是服务器对服务器端的推送
 */
func CallStpay(c *gin.Context) {
	res := "error"
	json_res := make(map[string]interface{})
	body_str := ""
	sign := ""
	sign_str := ""
	if c.Request.Body != nil {
		body := make([]byte, c.Request.ContentLength)
		body, err := ioutil.ReadAll(c.Request.Body)

		body_str = string(body)
		err = json.Unmarshal(body, &json_res)
		if err == nil {

			mer_id := fmt.Sprintf("%v", json_res["MerchantId"])              //商户对应的密钥账号
			req_id := fmt.Sprintf("%v", json_res["RequestId"])               //用户订单号
			order_no := fmt.Sprintf("%v", json_res["OrderNo"])               //平台订单号
			money := fmt.Sprintf("%v", json_res["Money"])                    //用户订单金额
			code := fmt.Sprintf("%v", json_res["Code"])                      //请求返回码（0:成功，其他:失败）
			time_str := url.QueryEscape(fmt.Sprintf("%v", json_res["Time"])) //平台传送的当前时间
			msg := fmt.Sprintf("%v", json_res["Message"])                    //附加信息
			req_type := fmt.Sprintf("%v", json_res["Type"])                  //请求类型
			fee := fmt.Sprintf("%v", json_res["Fee"])                        //手续费
			sign = fmt.Sprintf("%v", json_res["Sign"])                       //签名

			//需要验证签名的参数
			sign_str = fmt.Sprintf("Code=%s&Fee=%s&MerchantId=%s&Message=%s&Money=%s&OrderNo=%s&RequestId=%s&Time=%s&Type=%s", code, fee, mer_id, msg, money, order_no, req_id, time_str, req_type)
			if code == "0" {
				res = "SUCCESS"
				is_cent := 0
				c_status, c_msg := thread.ThreadUpdatePay(req_id, order_no, money, sign, sign_str, is_cent)
				if c_status != 200 {
					res = c_msg
				}
			}
		}
	}
	common.LogsWithFileName(log_path, "stpay_callback_", "body->"+body_str+"\nsign->"+sign+"\nsign_str->"+sign_str+"\nres->"+res)

	c.Writer.WriteString(res)
}

/**
* baifu的回调接口
* 隐式回调接口，就是服务器对服务器端的推送
 */
func CallBfpay(c *gin.Context) {
	res := "error"
	json_res := make(map[string]interface{})
	body_str := ""
	sign_str := ""
	sign := ""
	if c.Request.Body != nil {
		body := make([]byte, c.Request.ContentLength)
		body, err := ioutil.ReadAll(c.Request.Body)
		body_str = string(body)
		err = json.Unmarshal(body, &json_res)
		if err == nil {
			status := fmt.Sprintf("%v", json_res["status"])                 //交易状态
			order_sn := fmt.Sprintf("%v", json_res["order_sn"])             //订单号
			cus_order_sn := fmt.Sprintf("%v", json_res["cus_order_sn"])     //商户订单号
			order_amount := fmt.Sprintf("%v", json_res["order_amount"])     //用户订单金额
			receive_amount := fmt.Sprintf("%v", json_res["receive_amount"]) //实收金额
			//attach_data := fmt.Sprintf("%v",json_res["attach_data"])       //商戶請求交易時額外寄存資料
			sign = fmt.Sprintf("%v", json_res["sign"]) //签名
			message := url.QueryEscape(json_res["message"].(string))
			//需要验证签名的参数
			sign_str = fmt.Sprintf("cus_order_sn=%s&message=%s&order_amount=%s&order_sn=%s&receive_amount=%s&status=%s", cus_order_sn, message, order_amount, order_sn, receive_amount, status)

			res = "success"
			if status == "success" {
				res = "success"
				is_cent := 0
				c_status, c_msg := thread.ThreadUpdatePay(cus_order_sn, order_sn, receive_amount, sign, sign_str, is_cent)
				if c_status != 200 {
					res = c_msg
				}
			}
		}
	}
	common.LogsWithFileName(log_path, "bfpay_call_", "body->"+body_str+"\nsign->"+sign+"\nsign_str->"+sign_str+"\nres->"+res)
	c.Writer.WriteString(res)
}

/**
* yf的回调接口
* 隐式回调接口，就是服务器对服务器端的推送
 */
func CallYfpay(c *gin.Context) {
	res := "fail"
	sign_str := ""
	sign := c.Query("sign")
	param := map[string]string{
		"payKey":      c.Query("payKey"),
		"orderPrice":  c.Query("orderPrice"),
		"outTradeNo":  c.Query("outTradeNo"),
		"productType": c.Query("productType"),
		"orderTime":   c.Query("orderTime"),
		"productName": c.Query("productName"),
		"tradeStatus": c.Query("tradeStatus"),
		"successTime": c.Query("successTime"),
		"remark":      c.Query("remark"),
		"trxNo":       c.Query("trxNo"),
	}

	if param["tradeStatus"] == "SUCCESS" {
		//排序拼接
		var keys []string
		for key := range param {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, k := range keys {
			if param[k] != "" {
				sign_str += k + "=" + fmt.Sprintf("%s", param[k]) + "&"
			}
		}
		sign_str = strings.TrimRight(sign_str, "&")
		is_cent := 0
		call_status, _ := thread.ThreadUpdatePay(param["outTradeNo"], param["outTradeNo"], param["orderPrice"], sign, sign_str, is_cent)
		if call_status == 200 {
			res = "SUCCESS"
		}

	}

	common.LogsWithFileName(log_path, "yfay_call_", "sign->"+sign+"\nsign_str->"+sign_str+"\nres->"+res)

	c.Writer.WriteString(res)
}

/**
* jpay的回调接口
* 隐式回调接口，就是服务器对服务器端的推送
 */
func CallJpay(c *gin.Context) {
	res := "fail"
	sign_str := ""
	sign := c.Query("sign")
	param := map[string]string{
		"amount":     c.PostForm("amount"),
		"statusStr":  c.PostForm("statusStr"),
		"outTradeNo": c.PostForm("outTradeNo"),
		"status":     c.PostForm("status"),
	}

	if param["status"] == "1" {
		//排序拼接
		var keys []string
		for key := range param {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, k := range keys {
			if param[k] != "" {
				sign_str += k + "=" + fmt.Sprintf("%s", param[k]) + "&"
			}
		}
		sign_str = strings.TrimRight(sign_str, "&")
		is_cent := 0
		call_status, _ := thread.ThreadUpdatePay(param["outTradeNo"], param["outTradeNo"], param["amount"], sign, sign_str, is_cent)
		if call_status == 200 {
			res = "success"
		}

	}

	common.LogsWithFileName(log_path, "jpay_call_", "sign->"+sign+"\nsign_str->"+sign_str+"\nres->"+res)

	c.Writer.WriteString(res)
}

/**
* sxpay的回调接口
* 隐式回调接口，就是服务器对服务器端的推送
 */
func CallSxpay(c *gin.Context) {
	res := "fail"
	sign_str := ""
	sign := c.PostForm("sign")
	param := map[string]string{
		"amount":    c.PostForm("amount"),
		"sq_amount": c.PostForm("sq_amount"),
		"orderId":   c.PostForm("orderId"),
		"state":     c.PostForm("state"),
	}

	if param["state"] == "收款成功" {
		sign_str = fmt.Sprintf("amount=%s&orderId=%s", param["amount"], param["orderId"])
		is_cent := 0
		call_status, _ := thread.ThreadUpdatePay(param["orderId"], param["orderId"], param["amount"], sign, sign_str, is_cent)
		if call_status == 200 {
			res = "success"
		}

	}

	common.LogsWithFileName(log_path, "sxpay_call_", "sign->"+sign+"\nsign_str->"+sign_str+"\nres->"+res)

	c.Writer.WriteString(res)
}

/**
* thpay的回调接口
* 隐式回调接口，就是服务器对服务器端的推送
 */
func Callthpay(c *gin.Context) {

	res := "fail"
	sign_str := ""
	param := map[string]string{
		"fxid":          c.PostForm("fxid"),
		"fxorderid":     c.PostForm("fxorderid"),
		"fxtranid":      c.PostForm("fxtranid"),
		"fxamount":      c.PostForm("fxamount"),
		"fxamount_succ": c.PostForm("fxamount_succ"),
		"fxstatus":      c.PostForm("fxstatus"),
		"fxremark":      c.PostForm("fxremark"),
		"fxsign":        c.PostForm("fxsign"),
	}
	sign := param["fxsign"]
	if param["fxstatus"] == "succ" {
		sign_str := fmt.Sprintf("%s%s%s%s%s", param["fxstatus"], param["fxid"], param["fxamount_succ"], param["fxorderid"], param["fxamount"])
		is_cent := 0
		c_status, c_msg := thread.ThreadUpdatePay(param["fxorderid"], param["fxtranid"], param["fxamount_succ"], sign, sign_str, is_cent)
		res = "success"
		if c_status != 200 {
			res = c_msg
		}

	}
	body, _ := json.Marshal(param)

	common.LogsWithFileName(log_path, "thpay_call_", "body->"+string(body)+"\n"+"sign->"+sign+"\nsign_str->"+sign_str+"\nres->"+res)

	c.Writer.WriteString(res)
}

/**
* ggpay的回调接口
* 隐式回调接口，就是服务器对服务器端的推送
 */
func Callggpay(c *gin.Context) {
	res := "fail"
	sign_str := ""
	sign := c.PostForm("sign")
	param := map[string]string{
		"amount":     c.PostForm("amount"),
		"statusStr":  c.PostForm("statusStr"),
		"outTradeNo": c.PostForm("outTradeNo"),
		"status":     c.PostForm("status"),
	}

	if param["status"] == "1" {
		//排序拼接
		sign_str = common.MapCreatLinkSort(param, "&", true, false)
		is_cent := 0
		call_status, _ := thread.ThreadUpdatePay(param["outTradeNo"], param["outTradeNo"], param["amount"], sign, sign_str, is_cent)
		if call_status == 200 {
			res = "success"
		}

	}

	common.LogsWithFileName(log_path, "jpay_call_", "sign->"+sign+"\nsign_str->"+sign_str+"\nres->"+res)

	c.Writer.WriteString(res)
}

/**
* xxfpay的回调接口
* 隐式回调接口，就是服务器对服务器端的推送
 */
func Callxxfpay(c *gin.Context) {

	res := "fail"
	sign_str := ""
	param := map[string]string{
		"channel":    c.PostForm("channel"),
		"tradeNo":    c.PostForm("tradeNo"),
		"outTradeNo": c.PostForm("outTradeNo"),
		"money":      c.PostForm("money"),
		"realMoney":  c.PostForm("realMoney"),
		"uid":        c.PostForm("uid"),
		"outUserId":  c.PostForm("outUserId"),
		"outBody":    c.PostForm("outBody"),
	}
	sign := param["sign"]

	//传json过去，因为要和token 一起加密
	str, _ := json.Marshal(param)
	sign_str = string(str)
	is_cent := 0
	call_status, _ := thread.ThreadUpdatePay(param["outTradeNo"], param["outTradeNo"], param["realMoney"], sign, sign_str, is_cent)
	if call_status == 200 {
		res = "SUCCESS"
	}

	common.LogsWithFileName(log_path, "xxfpay_call_", "sign->"+sign+"\nsign_str->"+sign_str+"\nres->"+res)

	c.Writer.WriteString(res)
}

/**
* zofpay的回调接口
* 隐式回调接口，就是服务器对服务器端的推送
 */
func Callzofpay(c *gin.Context) {

	res := "fail"
	sign_str := ""
	param := map[string]string{
		"fxid":     c.PostForm("fxid"),
		"fxddh":    c.PostForm("fxddh"),
		"fxorder":  c.PostForm("fxorder"),
		"fxdesc":   c.PostForm("fxdesc"),
		"fxfee":    c.PostForm("fxfee"),
		"fxattch":  c.PostForm("fxattch"),
		"fxstatus": c.PostForm("fxstatus"),
		"fxtime":   c.PostForm("fxtime"),
	}
	sign := c.PostForm("fxsign")
	if param["fxstatus"] == "1" {
		sign_str := common.MapCreatLink(param, "fxstatus,fxid,fxddh,fxfee", "", 2)
		is_cent := 0
		call_status, _ := thread.ThreadUpdatePay(param["fxddh"], param["fxorder"], param["fxfee"], sign, sign_str, is_cent)
		if call_status == 200 {
			res = "success"
		}

	}

	common.LogsWithFileName(log_path, "zofpay_call_", "sign->"+sign+"\nsign_str->"+sign_str+"\nres->"+res)

	c.Writer.WriteString(res)
}

/**
* cfpay的回调接口
* 隐式回调接口，就是服务器对服务器端的推送
 */
func Callcfpay(c *gin.Context) {

	res := "fail"
	sign_str := ""
	param := map[string]string{
		"version":    c.PostForm("version"),
		"merchantId": c.PostForm("merchantId"),
		"orderNo":    c.PostForm("orderNo"),
		"tradeDate":  c.PostForm("tradeDate"),
		"tradeTime":  c.PostForm("tradeTime"),
		"amount":     c.PostForm("amount"),
		"resultCode": c.PostForm("resultCode"),
		"attach":     c.PostForm("attach"),
	}
	sign := c.PostForm("sign")
	str, _ := json.Marshal(param)
	sign_str = string(str)
	if param["resultCode"] == "0" {
		//传json过去，因为要和key 一起排序加密
		is_cent := 0
		res = "success"

		ramount, _ := strconv.ParseFloat(param["amount"], 64)
		ramount = ramount / 1000 // 分转元,保留两们小数
		amount := strconv.FormatFloat(float64(ramount), 'f', 2, 64)

		call_status, c_msg := thread.ThreadUpdatePay(param["orderNo"], param["orderNo"], amount, sign, sign_str, is_cent)
		if call_status != 200 {
			res = c_msg
		}
	}

	common.LogsWithFileName(log_path, "cfpay_call_", "body->"+sign_str+"\nsign->"+sign+"\nsign_str->"+sign_str+"\nres->"+res)

	c.Writer.WriteString(res)
}

/**
* xpay的回调接口
* 隐式回调接口，就是服务器对服务器端的推送
 */
func Callxpay(c *gin.Context) {
	res := "fail"
	sign_str := ""
	body_str := ""
	sign := ""
	param := map[string]string{}

	if c.Request.Body != nil {
		body := make([]byte, c.Request.ContentLength)
		body, err := ioutil.ReadAll(c.Request.Body)
		body_str = string(body)
		err = json.Unmarshal(body, &param)
		if err == nil {
			p_result := map[string]string{}
			//获取 pay_config 参数
			result, msg := thread.GetBackConfig(param["orderNo"], true)
			if msg == "success" && result != "" {
				err := json.Unmarshal([]byte(result), &p_result)
				if err != nil {
					common.LogsWithFileName(log_path, "Callxpay", "error->"+err.Error()+"\nresult->"+result)
				}

				//解密content
				aeskey, _ := base64.StdEncoding.DecodeString(p_result["aes_key"])
				get_ecb := common.SetAESECB(string(aeskey), "", "", "", 16)
				url_body := get_ecb.AesDecryptString(param["content"])

				content := make(map[string]interface{})
				err = json.Unmarshal([]byte(url_body), &content)
				if err != nil {
					common.LogsWithFileName(log_path, "Callxpay", "error->"+err.Error())
				}
				contents := make(map[string]string)
				for key, value := range content {
					contents[key] = fmt.Sprintf("%v", value)
				}
				//钱要保留两们小数
				ramount, _ := strconv.ParseFloat(contents["orderAmount"], 64)
				contents["orderAmount"] = fmt.Sprintf("%.2f", ramount)
				//拼接 业务参数
				sign_str = common.MapCreatLinkSort(contents, ",", true, false)

				sign = fmt.Sprintf("%v", param["sign"])

				if fmt.Sprintf("%v", content["orderStatus"]) == "SUCCESS" {
					is_cent := 0
					res = "SUCCESS"
					call_status, c_msg := thread.ThreadUpdatePay(param["orderNo"], param["orderNo"], contents["orderAmount"], sign, sign_str, is_cent)
					if call_status != 200 {
						res = c_msg
					}
				}
			}
		}
	}
	common.LogsWithFileName(log_path, "xpay_call_", "body->"+body_str+"\nsign->"+sign+"\nsign_str->"+sign_str+"\nres->"+res)

	c.Writer.WriteString(res)
}

/**
* xwspay的回调接口
* 隐式回调接口，就是服务器对服务器端的推送
 */
func Callxwspay(c *gin.Context) {

	res := "fail"
	sign_str := ""
	param := map[string]string{
		"fxid":     c.PostForm("fxid"),
		"fxddh":    c.PostForm("fxddh"),
		"fxorder":  c.PostForm("fxorder"),
		"fxdesc":   c.PostForm("fxdesc"),
		"fxfee":    c.PostForm("fxfee"),
		"fxattch":  c.PostForm("fxattch"),
		"fxstatus": c.PostForm("fxstatus"),
		"fxtime":   c.PostForm("fxtime"),
	}
	body, _ := json.Marshal(param)

	sign := c.PostForm("fxsign")
	if param["fxstatus"] == "1" {
		sign_str := common.MapCreatLink(param, "fxstatus,fxid,fxddh,fxfee", "", 2)
		is_cent := 0
		call_status, _ := thread.ThreadUpdatePay(param["fxddh"], param["fxorder"], param["fxfee"], sign, sign_str, is_cent)
		if call_status == 200 {
			res = "success"
		}

	}

	common.LogsWithFileName(log_path, "xwspay_call_", "body ->"+string(body)+"\nsign->"+sign+"\nsign_str->"+sign_str+"\nres->"+res)

	c.Writer.WriteString(res)
}
