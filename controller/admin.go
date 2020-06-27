package controller

import (
	"TFService/thread"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

/**
*  绑定谷歌验证码
 */
func CheckGoogleCode(c *gin.Context) {
	c_status := 100
	c_msg := "验证码必须填写"
	//定义需要输出的数据格式
	d := map[string]interface{}{}
	code_val := c.PostForm("code_val")

	if len(code_val) > 0 {
		c_status, c_msg = thread.CheckGoogleCode(code_val, c)
	}
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

func EditAdmin(c *gin.Context) {
	c_status := 100
	c_msg := "验证码必须填写"
	//定义需要输出的数据格式
	d := map[string]interface{}{}
	account := c.PostForm("account")
	is_edit := c.PostForm("is_edit")
	c_status, c_msg = thread.EditAdmin(account, is_edit, c)
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  更新用户的权限
 */
func UpdatePower(c *gin.Context) {
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求完成"

	power_path := c.PostForm("power_path")
	account := c.PostForm("account")
	secret := c.PostForm("secret")

	c_status, c_msg = thread.UpdatePower(account, power_path, secret, c)
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  更新管理员密码
 */
func UpdatePwd(c *gin.Context) {
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求完成"

	old_pwd := c.PostForm("old_pwd")
	new_pwd := c.PostForm("new_pwd")
	con_pwd := c.PostForm("con_pwd")
	secret := c.PostForm("secret")

	c_status, c_msg = thread.UpdatePwd(old_pwd, new_pwd, con_pwd, secret, c)
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  管理员列表
 */
func AdminList(c *gin.Context) {
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求完成"

	account := c.PostForm("account")
	mer_code := c.PostForm("mer_code")
	admin_status := c.PostForm("admin_status")
	page := c.PostForm("page")
	page_size := c.PostForm("page_size")

	c_status, d["total"], c_msg, d["admin_list"] = thread.AdminList(account, mer_code, admin_status, page, page_size, c)
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  新增管理员
 */
func AddAdmin(c *gin.Context) {
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求完成"

	account := c.PostForm("account")
	pwd := c.PostForm("pwd")
	mer_code := c.PostForm("mer_code")
	power_code := c.PostForm("power_code")

	c_status, c_msg = thread.AddAdmin(account, pwd, mer_code, power_code)
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  取消管理员的绑定
 */
func DelBind(c *gin.Context) {
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求完成"

	account := c.PostForm("account")

	c_status, c_msg = thread.DelBind(account, c)
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  清除缓存
 */
func DelCache(c *gin.Context) {
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求完成"

	cache_status := c.PostForm("cache_status")

	c_status, c_msg = thread.DelCache(cache_status)
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  支付订单
 */
func PayList(c *gin.Context) {
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求完成"
	sum_amount := 0.00
	pay_status := c.PostForm("pay_status")
	order_number := c.PostForm("order_number")
	web_order := c.PostForm("web_order")
	pay_code := c.PostForm("pay_code")
	class_code := c.PostForm("class_code")
	is_mobile := c.PostForm("is_mobile")
	start_time := c.PostForm("start_time")
	end_time := c.PostForm("end_time")
	mer_code := c.PostForm("mer_code")
	is_agent := c.PostForm("is_agent")
	page := c.PostForm("page")
	page_size := c.PostForm("page_size")

	c_status, d["total"], c_msg, sum_amount, d["pay_list"] = thread.PayList(pay_status, order_number, web_order, pay_code, class_code, is_mobile, start_time, end_time, mer_code, is_agent, page, page_size, c)
	d["sum_amount"] = strconv.FormatFloat(float64(sum_amount), 'f', 2, 64)
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  下发订单
 */
func CashList(c *gin.Context) {
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求完成"
	cash_status := c.PostForm("cash_status")
	order_number := c.PostForm("order_number")
	web_order := c.PostForm("web_order")
	pay_code := c.PostForm("pay_code")
	start_time := c.PostForm("start_time")
	end_time := c.PostForm("end_time")
	mer_code := c.PostForm("mer_code")
	is_agent := c.PostForm("is_agent")
	page := c.PostForm("page")
	page_size := c.PostForm("page_size")

	c_status, d["total"], c_msg, d["sum_amount"], d["cash_list"] = thread.CashList(cash_status, order_number, web_order, pay_code, start_time, end_time, mer_code, is_agent, page, page_size, c)

	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  商户现金流水
 */
func AmountList(c *gin.Context) {
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求完成"

	amount_type := c.PostForm("amount_type")
	order_number := c.PostForm("order_number")
	pay_code := c.PostForm("pay_code")
	start_time := c.PostForm("start_time")
	end_time := c.PostForm("end_time")
	mer_code := c.PostForm("mer_code")
	is_agent := c.PostForm("is_agent")
	page := c.PostForm("page")
	page_size := c.PostForm("page_size")

	c_status, d["total"], c_msg, d["cash_list"] = thread.AmountList(amount_type, order_number, pay_code, start_time, end_time, mer_code, is_agent, page, page_size, c)
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  一段时间内的总出入款
 */
func DateTotal(c *gin.Context) {
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求完成"

	start_date := c.PostForm("start_time")
	end_date := c.PostForm("end_time")

	c_status, c_msg, d["title_list"], d["data_list"] = thread.DateTotal(start_date, end_date, c)
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

func TotalBalance(c *gin.Context) {
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求完成"

	c_status, c_msg, d["balance_info"] = thread.TotalBalance(c)
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

func TodayCount(c *gin.Context) {
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求完成"
	date_type := c.PostForm("date_type")
	c_status, c_msg, d["count_list"] = thread.TodayCount(date_type, c)
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

func MerInfo(c *gin.Context) {
	c_status := 100
	c_msg := "验证码必须填写"
	//定义需要输出的数据格式
	d := map[string]interface{}{}
	c_status, c_msg, d["mer_info"] = thread.MerInfo(c)
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

func RateInfo(c *gin.Context) {
	c_status := 100
	c_msg := "验证码必须填写"
	//定义需要输出的数据格式
	d := map[string]interface{}{}
	c_status, c_msg, d["rate_info"] = thread.RateInfo(c)
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  支付订单
 */
func DownPayList(c *gin.Context) {
	d := map[string]interface{}{}

	pay_status := c.PostForm("pay_status")
	order_number := c.PostForm("order_number")
	web_order := c.PostForm("web_order")
	pay_code := c.PostForm("pay_code")
	class_code := c.PostForm("class_code")
	is_mobile := c.PostForm("is_mobile")
	start_time := c.PostForm("start_time")
	end_time := c.PostForm("end_time")
	mer_code := c.PostForm("mer_code")
	is_agent := c.PostForm("is_agent")
	c_status, c_msg, down_path := thread.DownPayList(pay_status, order_number, web_order, pay_code, class_code, is_mobile, start_time, end_time, mer_code, is_agent, c)
	down_path = strings.Replace(down_path, "./", "", 1)

	d["down_url"] = down_url + "/" + down_path
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  下发订单
 */
func DownCashList(c *gin.Context) {
	d := map[string]interface{}{}

	cash_status := c.PostForm("cash_status")
	order_number := c.PostForm("order_number")
	web_order := c.PostForm("web_order")
	pay_code := c.PostForm("pay_code")
	start_time := c.PostForm("start_time")
	end_time := c.PostForm("end_time")
	mer_code := c.PostForm("mer_code")
	is_agent := c.PostForm("is_agent")

	c_status, c_msg, down_path := thread.DownCashList(cash_status, order_number, web_order, pay_code, start_time, end_time, mer_code, is_agent, c)
	down_path = strings.Replace(down_path, "./", "", 1)

	d["down_url"] = down_url + "/" + down_path
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  商户现金流水
 */
func DownAmountList(c *gin.Context) {
	d := map[string]interface{}{}
	amount_type := c.PostForm("amount_type")
	order_number := c.PostForm("order_number")
	pay_code := c.PostForm("pay_code")
	start_time := c.PostForm("start_time")
	end_time := c.PostForm("end_time")
	mer_code := c.PostForm("mer_code")
	is_agent := c.PostForm("is_agent")

	c_status, c_msg, down_path := thread.DownAmountList(amount_type, order_number, pay_code, start_time, end_time, mer_code, is_agent, c)
	down_path = strings.Replace(down_path, "./", "", 1)

	d["down_url"] = down_url + "/" + down_path
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  支付回调
 */
func CallPay(c *gin.Context) {
	d := map[string]interface{}{}

	order_number := c.PostForm("order_number")

	c_status, c_msg := thread.CallPay(order_number)

	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  下发回调
 */
func CallCash(c *gin.Context) {
	d := map[string]interface{}{}

	order_number := c.PostForm("order_number")

	c_status, c_msg := thread.CallCash(order_number)

	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  支付订单的成功或者失败
 */
func PayStatus(c *gin.Context) {
	d := map[string]interface{}{}

	order_number := c.PostForm("p_id")
	pay_status := c.PostForm("pay_status")

	c_status, c_msg := thread.PayStatus(order_number, pay_status, c)

	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  下发订单的成功或者失败
 */
func CashStatus(c *gin.Context) {
	d := map[string]interface{}{}

	order_number := c.PostForm("c_id")
	cash_status := c.PostForm("cash_status")

	c_status, c_msg := thread.CashStatus(order_number, cash_status, c)

	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  下发订单的成功或者失败
 */
func MerPay(c *gin.Context) {
	d := map[string]interface{}{}

	c_status := 100
	c_msg := "请求错误"

	c_status, c_msg, d["balance_list"] = thread.MerPay(c)

	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  后台下发
 */
func MerCash(c *gin.Context) {
	d := map[string]interface{}{}

	c_status := 100
	c_msg := "验证码必须填写"
	pay_id := c.PostForm("pay_id")
	bank_id := c.PostForm("bank_id")
	amount := c.PostForm("amount")
	is_auto := c.PostForm("is_auto")
	secret := c.PostForm("secret")
	c_status, c_msg, d["result"] = thread.MerCash(pay_id, bank_id, amount, is_auto, secret, c)

	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  绑定银行卡
 */
func AddBank(c *gin.Context) {
	d := map[string]interface{}{}

	bank_code := c.PostForm("bank_code")
	card_number := c.PostForm("card_number")
	card_name := c.PostForm("card_name")
	bank_branch := c.PostForm("bank_branch")
	bank_phone := c.PostForm("bank_phone")
	secret := c.PostForm("secret")
	c_status, c_msg := thread.AddBank(bank_code, card_number, card_name, bank_branch, bank_phone, secret, c)

	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  商户绑定的银行卡
 */
func MerBank(c *gin.Context) {
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求完成"
	bank_code := c.PostForm("bank_code")
	bank_status := c.PostForm("bank_status")
	page := c.PostForm("page")
	page_size := c.PostForm("page_size")
	card_number := c.PostForm("card_number")
	card_name := c.PostForm("card_name")
	c_status, d["total"], c_msg, d["mer_bank"] = thread.MerBank(bank_code, bank_status, card_number, card_name, page, page_size, c)

	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  锁定银行卡
 */
func LockBank(c *gin.Context) {
	d := map[string]interface{}{}

	b_id := c.PostForm("b_id")
	secret := c.PostForm("secret")
	c_status, c_msg := thread.LockBank(b_id, secret, c)

	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

func AgentReport(c *gin.Context) {
	d := map[string]interface{}{}
	sum_total := map[string]float64{} //总和
	c_status := 100
	c_msg := "请求完成"
	mer_code := c.PostForm("mer_code")
	start_date := c.PostForm("start_date")
	end_date := c.PostForm("end_date")
	page := c.PostForm("page")
	page_size := c.PostForm("page_size")
	c_status, d["total"], c_msg, sum_total, d["agent_report"] = thread.AgentReport(mer_code, start_date, end_date, page, page_size, c)
	for key, value := range sum_total {
		d[key] = value
	}
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

func MerChannel(c *gin.Context) {
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求完成"
	mer_code := c.PostForm("mer_code")
	chann_status := c.PostForm("chann_status")
	pay_code := c.PostForm("pay_code")
	page := c.PostForm("page")
	page_size := c.PostForm("page_size")
	c_status, d["total"], c_msg, d["mer_channel"] = thread.MerChannel(mer_code, pay_code, chann_status, page, page_size, c)

	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  新增商户渠道
 */
func AddMerChannel(c *gin.Context) {
	d := map[string]interface{}{}

	mer_code := c.PostForm("mer_code")
	pay_id := c.PostForm("pay_id")
	c_status, c_msg := thread.AddMerChannel(pay_id, mer_code, c)

	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

func EditMerChannel(c *gin.Context) {
	d := map[string]interface{}{}

	p_id := c.PostForm("p_id")
	channel_status := c.PostForm("channel_status")
	c_status, c_msg := thread.EditMerChannel(p_id, channel_status)

	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  下发订单
 */
func OrderList(c *gin.Context) {
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求完成"

	order_id := c.PostForm("order_id")
	order_status := c.PostForm("order_status")
	order_number := c.PostForm("order_number")
	web_order := c.PostForm("web_order")
	pay_code := c.PostForm("pay_code")
	start_time := c.PostForm("start_time")
	end_time := c.PostForm("end_time")
	page := c.PostForm("page")
	page_size := c.PostForm("page_size")

	c_status, d["total"], c_msg, d["order_list"] = thread.OrderList(order_id, order_status, order_number, web_order, pay_code, start_time, end_time, page, page_size, c)
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  后台渠道下发
 */
func AddOrder(c *gin.Context) {
	d := map[string]interface{}{}

	pay_id := c.PostForm("pay_id")
	order_number := c.PostForm("order_number")
	amount := c.PostForm("amount")
	secret := c.PostForm("secret")
	c_status, c_msg := thread.AddOrder(pay_id, order_number, amount, secret, c)

	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  下发订单的成功或者失败
 */
func UpdateOrder(c *gin.Context) {
	d := map[string]interface{}{}

	order_id := c.PostForm("order_id")
	order_status := c.PostForm("order_status")

	c_status, c_msg := thread.UpdateOrder(order_id, order_status, c)

	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  给用户新增额度
 */
func AddMerAmount(c *gin.Context) {
	d := map[string]interface{}{}

	mer_code := c.PostForm("mer_code")
	amount := c.PostForm("amount")
	secret := c.PostForm("secret")
	pay_id := c.PostForm("pay_id")
	note := c.PostForm("note")

	c_status, c_msg := thread.AddMerAmount(mer_code, amount, secret, pay_id, note, c)

	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  人工充值
 */
func ManualRecharge(c *gin.Context) {
	d := map[string]interface{}{}

	mer_code := c.PostForm("mer_code")
	amount := c.PostForm("amount")
	secret := c.PostForm("secret")
	pay_id := c.PostForm("pay_id")
	bank_code := c.PostForm("bank_code")
	class_code := c.PostForm("class_code")
	note := c.PostForm("note")

	c_status, c_msg := thread.ManualRecharge(mer_code, amount, secret, pay_id, bank_code, class_code, note, c)

	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}
