package controller

import (
	"image/png"
	"TFService/hook"
	"TFService/thread"

	"github.com/gin-gonic/gin"
	"github.com/zhibingzhou/go_public/common"
)

func Test(c *gin.Context) {
	c.HTML(200, "index.tpl", gin.H{
		"title": "支付接口调试",
	})
}

func TestAdmin(c *gin.Context) {
	c.HTML(200, "admin.tpl", gin.H{
		"title": "后台接口调试",
	})
}

/**
* 测试加密
 */
func TestEncode(c *gin.Context) {
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求完成"
	private_key := c.PostForm("private_key")
	pay_data := c.PostForm("pay_data")

	testAES := common.SetAES(private_key, "", "pkcs5", 16)
	d["result"] = testAES.AesEncryptString(pay_data)
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
* 测试支付
 */
func TestPay(c *gin.Context) {
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求完成"
	pay_map := map[string]string{}

	//接收值
	mer_code := c.PostForm("mer_code")
	pay_data := c.PostForm("pay_data")
	/////////////////获得输入的值/////////////////
	c_status, c_msg, pay_map = hook.AuthInputAndMap(mer_code, pay_data)

	if c_status == 200 {
		c_status, c_msg, d["tpl_param"], d["api_jump"], d["form_param"] = thread.PayCreate(mer_code, pay_map)
	}

	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
* 测试测试查询订单
 */

func TestPayQuery(c *gin.Context) {
	//定义需要输出的结果
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求完成"

	//接收值
	// mer_code := c.PostForm("mer_code")
	// order_number := c.PostForm("order_number")

	//c_status, d["result"], c_msg = thread.PayQuery(merchant_id, ordernumber)

	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

func TestPayFor(c *gin.Context) {
	//定义需要输出的结果
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求完成"

	// //接收值
	// pay_map := map[string]string{}

	// //接收值
	// mer_code := c.PostForm("mer_code")
	// pay_data := c.PostForm("pay_data")
	// /////////////////获得输入的值/////////////////
	// c_status, c_msg, pay_map = hook.AuthInputAndMap(mer_code, pay_data)

	//c_status, c_msg, d["result"] = thread.PayFor(merchant_id, "2", paramsMap)

	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

func TestPayForQuery(c *gin.Context) {
	//定义需要输出的结果
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求完成"

	//接收值
	// mer_code := c.PostForm("mer_code")
	// order_number := c.PostForm("order_number")

	//c_status, d["result"], c_msg = thread.PayForQuery(mer_code, ordernumber)

	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
* 支付支持的银行
 */
func TestPayBank(c *gin.Context) {
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
* 支付支持的银行
 */
func TestCreateQrCode(c *gin.Context) {

	//接收值
	qr_code := c.GetString("qr_code")

	c_img, _ := common.CreateQrCode(qr_code)

	c.Writer.Header().Set("Content-Type", "image/png")
	png.Encode(c.Writer, c_img)
}

func TestEncodeCreateQrCode(c *gin.Context) {
	//接收值
	qr_code := c.GetString("qr_code")
	qr_key := c.GetString("qr_key")

	TestAES := common.SetAES(qr_key, "", "", 16)

	de_qr := TestAES.AesDecryptString(qr_code)

	c_img, _ := common.CreateQrCode(de_qr)

	c.Writer.Header().Set("Content-Type", "image/png")
	png.Encode(c.Writer, c_img)
}

func TestQrCode(c *gin.Context) {
	//接收值
	img_url := c.GetString("qr_code")

	c.HTML(200, "test.tpl", gin.H{
		"img_url": img_url,
	})
}

func TestAliJs(c *gin.Context) {
	money := c.PostForm("money")
	account := c.PostForm("account")
	order_id := c.PostForm("order_id")

	c.HTML(200, "alipay_js.tpl", gin.H{
		"account":  account,
		"money":    money,
		"order_id": order_id,
	})
}

func TestJs(c *gin.Context) {
	money := c.PostForm("money")
	account := c.PostForm("account")
	order_id := c.PostForm("order_id")

	c.HTML(200, "test.tpl", gin.H{
		"account":  account,
		"money":    money,
		"order_id": order_id,
	})
}
