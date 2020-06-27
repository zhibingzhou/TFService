package controller

import (
	"TFService/hook"
	"TFService/thread"

	"github.com/gin-gonic/gin"
)

/**
* 支付的请求
 */
func PayCreate(c *gin.Context) {
	//定义需要输出的结果
	c_status := 100
	c_msg := "请求完成"
	pay_map := map[string]string{}

	api_jump := 0
	form_param := map[string]string{}
	tpl_param := map[string]string{}
	//接收值
	mer_code := c.PostForm("mer_code")
	pay_data := c.PostForm("pay_data")
	/////////////////获得输入的值/////////////////
	c_status, c_msg, pay_map = hook.AuthInputAndMap(mer_code, pay_data)

	if c_status == 200 {
		c_status, c_msg, tpl_param, api_jump, form_param = thread.PayCreate(mer_code, pay_map)
	}
	tpl_name := ""
	//1=表单提交,2=扫码,3=普通window.location.href跳转,4=自动跳转的js代码,5=禁止referrer的跳转
	if c_status == 200 {
		if api_jump == 2 {
			if tpl_param["pay_class"] == "2" && tpl_param["is_mobile"] == "1" {
				tpl_name = "mobile/alipay.tpl"
			} else if tpl_param["pay_class"] == "2" && tpl_param["is_mobile"] == "0" {
				tpl_name = "pc/alipay.tpl"
			} else if tpl_param["pay_class"] == "3" && tpl_param["is_mobile"] == "1" {
				tpl_name = "mobile/tenpay.tpl"
			} else if tpl_param["pay_class"] == "3" && tpl_param["is_mobile"] == "0" {
				tpl_name = "pc/tenpay.tpl"
			} else if tpl_param["pay_class"] == "4" && tpl_param["is_mobile"] == "1" {
				tpl_name = "mobile/wechatpay.tpl"
			} else if tpl_param["pay_class"] == "4" && tpl_param["is_mobile"] == "0" {
				tpl_name = "pc/wechatpay.tpl"
			} else if tpl_param["pay_class"] == "6" && tpl_param["is_mobile"] == "1" {
				tpl_name = "mobile/unionpay.tpl"
			} else if tpl_param["pay_class"] == "6" && tpl_param["is_mobile"] == "0" {
				tpl_name = "pc/unionpay.tpl"
			}
		} else if api_jump == 3 {
			tpl_name = "jump.tpl"
		} else if api_jump == 1 {
			//默认使用表单提交
			tpl_name = "put.tpl"
		} else if api_jump == 4 {
			tpl_name = "auto_jump.tpl"
		} else if api_jump == 5 {
			tpl_name = "new_jump.tpl"
		}
	} else {
		tpl_name = "error.tpl"
	}

	c.HTML(200, tpl_name, gin.H{
		"title":             "接口服务中心的调试",
		"pay_url":           tpl_param["img_url"],
		"amount":            tpl_param["amout"],
		"Api_create_method": tpl_param["api_method"],
		"Api_create_url":    tpl_param["pay_url"],
		"Api_form_param":    form_param,
		"msg":               c_msg,
	})

}

/**
* 支付订单查询
 */
func PayOrder(c *gin.Context) {
	//定义需要输出的结果
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求完成"

	//接收值
	mer_code := c.PostForm("mer_code")
	order_number := c.PostForm("order_number")

	c_status, c_msg, d["order_info"] = thread.PayOrder(mer_code, order_number)

	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
* 查询余额
 */
func PayBalance(c *gin.Context) {
	//定义需要输出的结果
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求完成"

	//接收值
	pay_map := map[string]string{}
	mer_code := c.PostForm("mer_code")
	pay_data := c.PostForm("pay_data")
	///////////////获得输入的值/////////////////
	c_status, c_msg, pay_map = hook.AuthInputAndMap(mer_code, pay_data)

	c_status, c_msg, d["result"] = thread.PayBalance(mer_code, pay_map)

	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
* 支付支持的银行
 */
func Bank(c *gin.Context) {
	//定义需要输出的结果
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求完成"

	//接收值
	is_mobile := c.PostForm("is_mobile")
	pay_id := c.PostForm("pay_id")

	d["is_mobile"] = is_mobile
	d["pay_id"] = pay_id

	c_status, c_msg, d["bank_list"] = thread.Bank(is_mobile, pay_id)

	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*   查询下发的订单
 */
func CashOrder(c *gin.Context) {
	//定义需要输出的结果
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求完成"

	//接收值
	mer_code := c.PostForm("mer_code")
	order_number := c.PostForm("order_number")

	c_status, c_msg, d["order_info"] = thread.CashOrder(mer_code, order_number)

	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*   下发
 */
func PayFor(c *gin.Context) {
	//定义需要输出的结果
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求完成"

	//接收值
	mer_code := c.PostForm("mer_code")
	pay_data := c.PostForm("pay_data")
	pay_map := map[string]string{}
	/////////////////获得输入的值/////////////////
	c_status, c_msg, pay_map = hook.AuthInputAndMap(mer_code, pay_data)

	if c_status == 200 {
		c_status, c_msg, d["result"] = thread.PayFor(mer_code, pay_map)
	}

	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*   下发
 */
func DFpay(c *gin.Context) {
	//定义需要输出的结果
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求完成"

	//接收值
	mer_code := c.PostForm("mer_code")
	pay_data := c.PostForm("pay_data")
	pay_map := map[string]string{}
	/////////////////获得输入的值/////////////////
	c_status, c_msg, pay_map = hook.AuthInputAndMap(mer_code, pay_data)

	if c_status == 200 {
		c_status, c_msg, d["result"] = thread.DFpay(mer_code, pay_map)
	}

	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}
