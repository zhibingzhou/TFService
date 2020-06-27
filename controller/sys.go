package controller

import (
	"TFService/thread"

	"github.com/gin-gonic/gin"
)

/**
*  权限列表
 */
func PowerList(c *gin.Context) {
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求完成"
	account := c.PostForm("account")
	c_status, c_msg, d["power_list"] = thread.PowerList(account, c)
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  新增权限
 */
func AddPower(c *gin.Context) {
	d := map[string]interface{}{}

	path := c.PostForm("path")
	name := c.PostForm("name")
	url := c.PostForm("url")
	p_path := c.PostForm("p_path")
	power_type := c.PostForm("power_type")

	c_status, c_msg := thread.AddPower(path, name, url, p_path, power_type)
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  新增商户
 */
func AddMer(c *gin.Context) {
	d := map[string]interface{}{}
	m_map := map[string]string{}
	m_map["code"] = c.PostForm("code")
	m_map["title"] = c.PostForm("title")
	m_map["domain"] = c.PostForm("domain")
	m_map["qq"] = c.PostForm("qq")
	m_map["skype"] = c.PostForm("skype")
	m_map["telegram"] = c.PostForm("telegram")
	m_map["phone"] = c.PostForm("phone")
	m_map["email"] = c.PostForm("email")
	m_map["is_agent"] = c.PostForm("is_agent")
	p_agent := c.PostForm("p_agent")

	c_status, c_msg := thread.AddMer(p_agent, m_map, c)
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  商户列表
 */
func MerList(c *gin.Context) {
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求错误"
	is_under := c.PostForm("is_under")
	page := c.PostForm("page")
	page_size := c.PostForm("page_size")
	Ismer_code := c.PostForm("mer_code")

	c_status, d["total"], c_msg, d["mer_list"] = thread.MerList(is_under, page, page_size, Ismer_code, c)
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  修改商户信息
 */
func UpdateMer(c *gin.Context) {
	d := map[string]interface{}{}

	up_map := map[string]string{}
	mer_code := c.PostForm("mer_code")
	up_map["qq"] = c.PostForm("qq")
	up_map["domain"] = c.PostForm("domain")
	up_map["skype"] = c.PostForm("skype")
	up_map["telegram"] = c.PostForm("telegram")
	up_map["phone"] = c.PostForm("phone")
	up_map["email"] = c.PostForm("email")
	mer_status := c.PostForm("mer_status")

	c_status, c_msg := thread.UpdateMer(mer_code, mer_status, up_map, c)
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  支付类型列表
 */
func PayClass(c *gin.Context) {
	c_status := 100
	c_msg := "请求错误"
	d := map[string]interface{}{}

	c_status, c_msg, d["class_list"] = thread.PayClass()
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  新增上游渠道
 */
func AddPay(c *gin.Context) {
	d := map[string]interface{}{}

	code := c.PostForm("code")
	title := c.PostForm("title")
	fee_amount := c.PostForm("fee_amount")
	fee_type := c.PostForm("fee_type")
	is_push := c.PostForm("is_push")

	c_status, c_msg := thread.AddPay(code, title, fee_amount, fee_type, is_push, c)
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  设置上游渠道费率
 */
func AddPayClass(c *gin.Context) {
	d := map[string]interface{}{}

	class_code := c.PostForm("class_code")
	rate := c.PostForm("rate")
	pay_code := c.PostForm("pay_code")
	bank_code := c.PostForm("bank_code")
	min_amount := c.PostForm("min_amount")
	max_amount := c.PostForm("max_amount")
	limit_amount := c.PostForm("limit_amount")

	c_status, c_msg := thread.AddPayClass(pay_code, class_code, bank_code, rate, min_amount, max_amount, limit_amount, c)
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  上游渠道列表
 */
func ChannelList(c *gin.Context) {
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求错误"
	class_code := c.PostForm("class_code")
	page := c.PostForm("page")
	page_size := c.PostForm("page_size")

	c_status, d["total"], c_msg, d["list"] = thread.ChannelList(class_code, page, page_size, c)
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  查询上游渠道支付类型详情
 */
func PayDetail(c *gin.Context) {
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求错误"
	pay_code := c.PostForm("pay_code")
	class_code := c.PostForm("class_code")
	page := c.PostForm("page")
	page_size := c.PostForm("page_size")

	c_status, d["total"], c_msg, d["list"] = thread.PayDetail(pay_code, class_code, page, page_size, c)
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  修改上游支付类型详情
 */
func UpdatePay(c *gin.Context) {
	d := map[string]interface{}{}
	pay_id := c.PostForm("pay_id")
	rate := c.PostForm("rate")
	min_amount := c.PostForm("min_amount")
	max_amount := c.PostForm("max_amount")

	c_status, c_msg := thread.UpdatePay(pay_id, rate, min_amount, max_amount, c)
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  新增/修改用户的费率
 */
func MerRate(c *gin.Context) {
	d := map[string]interface{}{}
	rate_id := c.PostForm("rate_id")
	mer_code := c.PostForm("mer_code")
	class_code := c.PostForm("class_code")
	rate := c.PostForm("rate")
	pay_code := c.PostForm("pay_code")
	bank_code := c.PostForm("bank_code")

	c_status, c_msg := thread.MerRate(rate_id, mer_code, pay_code, class_code, bank_code, rate, c)
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  商户的费率
 */
func MerRateList(c *gin.Context) {
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求错误"
	mer_code := c.PostForm("mer_code")
	pay_code := c.PostForm("pay_code")
	class_code := c.PostForm("class_code")
	is_under := c.PostForm("is_under")
	page := c.PostForm("page")
	page_size := c.PostForm("page_size")

	c_status, d["total"], c_msg, d["list"] = thread.MerRateList(mer_code, pay_code, class_code, is_under, page, page_size, c)
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  删除商户的费率
 */
func DelMerRate(c *gin.Context) {
	d := map[string]interface{}{}
	rate_id := c.PostForm("rate_id")
	mer_code := c.PostForm("mer_code")
	c_status, c_msg := thread.DelRate(rate_id, mer_code, c)
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  首页查询商户信息
 */
func PayMerDetail(c *gin.Context) {
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求错误"
	is_under := c.PostForm("is_under")
	page := c.PostForm("page")
	page_size := c.PostForm("page_size")

	c_status, d["total"], c_msg, d["list"] = thread.PayMerDetail(is_under, page, page_size, c)
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  商户IP白名单
 */
func IpList(c *gin.Context) {
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求错误"
	mer_code := c.PostForm("mer_code")
	ip := c.PostForm("ip")
	page := c.PostForm("page")
	page_size := c.PostForm("page_size")

	c_status, d["total"], c_msg, d["list"] = thread.IpList(mer_code, ip, page, page_size, c)
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  添加白名单
 */
func AddIp(c *gin.Context) {
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求错误"
	mer_code := c.PostForm("mer_code")
	ip := c.PostForm("ip")

	c_status, c_msg = thread.AddIp(mer_code, ip, c)
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  删除白名单
 */
func DelIp(c *gin.Context) {
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求错误"
	ip := c.PostForm("ip")

	c_status, c_msg = thread.DelIp(ip, c)
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

func NoticeList(c *gin.Context) {
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求完成"
	is_all := c.PostForm("is_all")
	page := c.PostForm("page")
	page_size := c.PostForm("page_size")

	c_status, d["total"], c_msg, d["notice_list"] = thread.NoticeList(is_all, page, page_size)
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

func AddNotice(c *gin.Context) {
	d := map[string]interface{}{}
	title := c.PostForm("title")
	content := c.PostForm("content")

	c_status, c_msg := thread.AddNotice(title, content, c)
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

func UpdateNotice(c *gin.Context) {
	d := map[string]interface{}{}
	n_id := c.PostForm("n_id")

	c_status, c_msg := thread.UpdateNotice(n_id, c)
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

func SysBank(c *gin.Context) {
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求完成"
	c_status, c_msg, d["sys_bank"] = thread.SysBank()
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

func PayConf(c *gin.Context) {
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求完成"
	pay_code := c.PostForm("pay_code")
	pay_status := c.PostForm("pay_status")
	page := c.PostForm("page")
	page_size := c.PostForm("page_size")

	c_status, d["total"], c_msg, d["pay_conf"] = thread.PayConf(pay_code, pay_status, page, page_size)
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

func PayBank(c *gin.Context) {
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求完成"
	pay_code := c.PostForm("pay_code")
	class_code := c.PostForm("class_code")
	page := c.PostForm("page")
	page_size := c.PostForm("page_size")

	c_status, d["total"], c_msg, d["pay_bank"] = thread.PayBank(pay_code, class_code, page, page_size)
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

func AddPayBank(c *gin.Context) {
	d := map[string]interface{}{}
	p_data := map[string]string{}
	p_data["pay_code"] = c.PostForm("pay_code")
	p_data["class_code"] = c.PostForm("class_code")
	p_data["is_mobile"] = c.PostForm("is_mobile")
	p_data["bank_code"] = c.PostForm("bank_code")
	p_data["pay_bank"] = c.PostForm("pay_bank")
	p_data["jump_type"] = c.PostForm("jump_type")

	c_status, c_msg := thread.AddPayBank(p_data)
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}
