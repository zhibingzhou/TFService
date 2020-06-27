package controller

import (
	"image/png"
	"TFService/thread"

	"github.com/gin-gonic/gin"
	"github.com/zhibingzhou/go_public/common"
)

/**
* 管理员登录
 */
func Login(c *gin.Context) {
	//定义需要输出的结果
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求完成"

	//接收值
	account := c.PostForm("account")
	pwd := c.PostForm("pwd")
	secret := c.PostForm("secret")
	d["account"] = account
	c_status, c_msg = thread.Login(account, pwd, secret, c)

	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  退出
 */
func Logout(c *gin.Context) {
	data := map[string]interface{}{}
	c_status, c_msg := thread.Logout(c)
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
* 支付支持的银行
 */
func PaySuccess(c *gin.Context) {
	//定义需要输出的结果
	host := c.Request.Host
	c.HTML(200, "success.tpl", gin.H{
		"title": "支付接口调试",
		"host":  host,
	})
}

/**
* 生成二维码
 */
func CreateQrCode(c *gin.Context) {
	//接收值
	qr_code := c.GetString("qr_code")

	c_img, _ := common.CreateQrCode(qr_code)
	c.Writer.Header().Set("Content-Type", "image/png")
	png.Encode(c.Writer, c_img)
}

/**
*  生成谷歌验证码
 */
func GoogleQr(c *gin.Context) {
	_, c_img := thread.GoogleQr(c)
	c.Writer.Header().Set("Content-Type", "image/png")
	png.Encode(c.Writer, c_img)
}
