package controller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"TFService/thread"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zhibingzhou/go_public/common"
)

func BackNewpay(c *gin.Context) {
	body := c.PostForm("body")
	sign := c.PostForm("sign")

	res := "error"
	body, _ = url.QueryUnescape(body)

	var json_res map[string]interface{}
	err := json.Unmarshal([]byte(body), &json_res)
	if err == nil {
		res = "SUCCESS"
		status := fmt.Sprintf("%v", json_res["status"])
		is_issued := fmt.Sprintf("%v", json_res["isIssued"])
		cash_status := "1"
		if status == "success" {
			if is_issued == "true" {
				cash_status = "3"
			}
		} else if status == "failure" {
			cash_status = "9"
		}
		if cash_status != "1" {
			amount := fmt.Sprintf("%v", json_res["amount"])
			amount_f, _ := strconv.ParseFloat(amount, 64)
			amount = fmt.Sprintf("%.2f", amount_f)
			trade_no := fmt.Sprintf("%v", json_res["outTradeNo"])
			sn := fmt.Sprintf("%v", json_res["sn"])
			sign_str := fmt.Sprintf(`body={"amount":%s,"isIssued":%s,"outTradeNo":"%s","payTime":"%v","sn":"%s","status":"%s"}`, amount, is_issued, trade_no, json_res["payTime"], sn, status)
			is_cent := 0
			thread.ThreadUpdateOrder(trade_no, sn, amount, cash_status, sign, sign_str, is_cent)
		}

	}

	common.LogsWithFileName(log_path, "newpay_back_", "sign->"+sign+"\nbody->"+body+"\nres->"+res)

	c.Writer.WriteString(res)
}

//c2c代付回调
func BackC2Cpay(c *gin.Context) {

	param := map[string]string{
		"tradeno":         c.PostForm("tradeno"),
		"tradedate":       c.PostForm("tradedate"),
		"tradedesc":       c.PostForm("tradedesc"),
		"totalamount":     c.PostForm("totalamount"),
		"bankaccountno":   c.PostForm("bankaccountno"),
		"bankaccountname": c.PostForm("bankaccountname"),
		"currency":        c.PostForm("currency"),
		"tradestatus":     c.PostForm("tradestatus"),
		"completedate":    c.PostForm("completedate"),
	}
	body, _ := json.Marshal(param)
	sign := c.PostForm("sign")
	res := "error"
	status := fmt.Sprintf("%v", param["tradestatus"])
	cash_status := "1"
	if status == "SUCCESS" {
		cash_status = "3"
	} else if status == "FAILED" || status == "NOT_EXIST" {
		cash_status = "9"
	}
	sign_str := common.MapCreatLinkSort(param, "&", true, true)
	if cash_status != "1" {
		amount := fmt.Sprintf("%v", param["totalamount"])
		trade_no := fmt.Sprintf("%v", param["tradeno"])

		is_cent := 0
		c_status, c_msg := thread.ThreadUpdateOrder(trade_no, trade_no, amount, cash_status, sign, sign_str, is_cent)
		res = c_msg
		if c_status == 200 {
			res = "SUCCESS"
		}
	}

	common.LogsWithFileName(log_path, "c2cpay_back_", "sign->"+sign+"\nsign_str->"+sign_str+"body"+string(body)+"\nres->"+res)

	c.Writer.WriteString(res)
}

//bfpay代付回调
func Backbfpay(c *gin.Context) {
	json_res := make(map[string]interface{})
	sign := ""
	sign_str := ""
	res := "fail"
	bodystr := ""
	if c.Request.Body != nil {
		body := make([]byte, c.Request.ContentLength)
		body, err := ioutil.ReadAll(c.Request.Body)
		bodystr = string(body)
		err = json.Unmarshal(body, &json_res)
		if err == nil {

			encode_form := map[string]string{}
			for key, value := range json_res {
				encode_form[key] = url.QueryEscape(fmt.Sprintf("%v", value))
			}
			delete(encode_form, "sign")
			//拼接
			sign_str = common.MapCreatLinkSort(encode_form, "&", true, false)
			status := fmt.Sprintf("%v", json_res["status"])
			cash_status := "1"
			if status == "success" {
				cash_status = "3"
			} else if status == "failed" {
				cash_status = "9"
			}
			if cash_status != "1" {
				amount := fmt.Sprintf("%v", json_res["order_amount"])
				trade_no := fmt.Sprintf("%v", json_res["cus_order_sn"])
				pay_order := fmt.Sprintf("%v", json_res["order_sn"])
				sign = fmt.Sprintf("%v", json_res["sign"])
				is_cent := 0
				c_status, c_msg := thread.ThreadUpdateOrder(trade_no, pay_order, amount, cash_status, sign, sign_str, is_cent)
				res = c_msg
				if c_status == 200 {
					res = "SUCCESS"
				}
			}
		}
	}

	common.LogsWithFileName(log_path, "bfpay_back_", "sign->"+sign+"\nsign_str->"+sign_str+"body"+bodystr+"\nres->"+res)

	c.Writer.WriteString(res)
}

//xpay代付回调
func Backxpay(c *gin.Context) {

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
			result, msg := thread.GetBackConfig(param["orderNo"], false)
			if msg == "success" && result != "" {
				err := json.Unmarshal([]byte(result), &p_result)
				if err != nil {
					common.LogsWithFileName(log_path, "Backxpay", "error->"+err.Error()+"\nresult->"+result)
				}

				//解密content
				aeskey, _ := base64.StdEncoding.DecodeString(p_result["aes_key"])
				get_ecb := common.SetAESECB(string(aeskey), "", "", "", 16)
				url_body := get_ecb.AesDecryptString(param["content"])

				content := make(map[string]interface{})
				err = json.Unmarshal([]byte(url_body), &content)
				if err != nil {
					common.LogsWithFileName(log_path, "Backxpay", "error->"+err.Error())
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

				cash_status := "1"
				status := fmt.Sprintf("%v", content["orderStatus"])
				if status == "SUCCESS" {
					cash_status = "3"
				} else if status == "FAILED" || status == "BACK" {
					cash_status = "9"
				}
				if cash_status != "1" {
					is_cent := 0
					c_status, c_msg := thread.ThreadUpdateOrder(param["orderNo"], param["orderNo"], contents["orderAmount"], cash_status, sign, sign_str, is_cent)
					res = c_msg
					if c_status == 200 {
						res = "SUCCESS"
					}
				}
			}
		}
	}

	common.LogsWithFileName(log_path, "xpay_back_", "sign->"+sign+"\nsign_str->"+sign_str+"body"+body_str+"\nres->"+res)

	c.Writer.WriteString(res)
}

//mbpay代付回调
func Backmbpay(c *gin.Context) {

	param := map[string]string{
		"TransactionCode":   c.PostForm("TransactionCode"),
		"TransactionAmount": c.PostForm("TransactionAmount"),
		"Status":            c.PostForm("Status"),
		"SerialNumber":      c.PostForm("SerialNumber"),
	}
	body, _ := json.Marshal(param)
	sign := c.PostForm("sign")
	res := "error"
	status := fmt.Sprintf("%v", param["Status"])
	cash_status := "1"
	if status == "3" {
		cash_status = "3"
	} else if status == "4" {
		cash_status = "9"
	}
	sign_str := common.MapCreatLinkSort(param, "&", true, true)
	if cash_status != "1" {
		trade_no := fmt.Sprintf("%v", param["TransactionCode"])
		is_cent := 0
		c_status, c_msg := thread.ThreadUpdateOrder(trade_no, param["SerialNumber"], param["TransactionAmount"], cash_status, sign, sign_str, is_cent)
		res = c_msg
		if c_status == 200 {
			res = "SUCCESS"
		}
	}

	common.LogsWithFileName(log_path, "mbpay_back_", "sign->"+sign+"\nsign_str->"+sign_str+"body"+string(body)+"\nres->"+res)

	c.Writer.WriteString(res)
}
