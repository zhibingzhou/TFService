package thread

import (
	"fmt"
	"image"
	"TFService/model"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/zhibingzhou/go_public/common"
	"github.com/zhibingzhou/go_public/redis"
)

func Login(account, pwd, secret string, ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "账号密码不能为空"
	if account == "" || pwd == "" {
		return t_status, t_msg
	}
	a_info := model.AdminInfo(account)
	if a_info.Id < 1 {
		t_msg = "账号或密码错误"
		return t_status, t_msg
	}
	if a_info.Status != 1 {
		t_msg = "账号已被锁定或删除"
		return t_status, t_msg
	}
	if common.HexMd5(pwd) != a_info.Pwd {
		t_msg = "账号或密码错误"
		return t_status, t_msg
	}
	if a_info.Secret != "" && secret == "" {
		t_msg = "请填写谷歌动态验证"
		return t_status, t_msg
	}

	if a_info.Secret != "" {
		t_status, t_msg = AuthGoogle(secret, a_info.Secret)
		if t_status != 200 {
			t_msg = "谷歌动态验证码错误"
			return t_status, t_msg
		}
	}

	session := sessions.Default(ctx)

	//session_id := model.GetKey(20)
	session_id := common.Random("smallnumber", 30)
	session.Set("session_id", session_id)
	session.Set("account", a_info.Account)
	session.Set("mer_code", a_info.Mer_code)
	session.Set("power_path", a_info.Power_path)
	session.Save()
	a_data := map[string]interface{}{}
	a_data["session_id"] = session_id
	a_data["login_time"] = time.Now().Format(format_date)
	a_data["login_ip"] = ctx.ClientIP()
	model.UpdatesAdminList(a_info, a_data)
	if a_info.Secret == "" {
		t_status = 500
		t_msg = "请绑定谷歌动态验证"
		return t_status, t_msg
	}

	t_status = 200
	t_msg = "success"

	return t_status, t_msg
}

/**
*  退出
 */
func Logout(ctx *gin.Context) (int, string) {
	t_status := 200
	t_msg := "success"
	session := sessions.Default(ctx)
	session_id := session.Get("session_id")
	if session_id == nil {
		return t_status, t_msg
	}
	sess_id := fmt.Sprintf("%v", session_id)
	session.Clear()
	session.Save()
	sess_key := "admin_list:session_id:" + sess_id
	redis.RediGo.KeyDel(sess_key)
	//查询数据库里的session_id是否是当前的
	sess_info := model.AdminBySession(sess_id)
	if sess_info.Id < 1 {
		return t_status, t_msg
	}
	a_data := map[string]interface{}{}
	a_data["session_id"] = "session_id"
	model.UpdatesAdminList(sess_info, a_data)
	return t_status, t_msg
}

/**
*  是否登录
 */
func ThreadIsLogin(ctx *gin.Context) (int, string) {
	t_status := 600
	t_msg := "未登录"
	session := sessions.Default(ctx)
	session_id := session.Get("session_id")

	if session_id == nil {
		return t_status, t_msg
	}
	//判断唯一登录
	sess_id := fmt.Sprintf("%v", session_id)
	a_info := model.AdminBySessionRedis(sess_id)

	if len(a_info["account"]) < 1 {
		t_msg = "账号已经在其他地方登录"
		return t_status, t_msg
	}

	if len(a_info["secret"]) < 1 {
		t_status = 500
		t_msg = "请绑定谷歌动态验证"
		return t_status, t_msg
	}

	// t_status = 200
	// t_msg = "success"
	t_status, t_msg = IsPower(ctx)

	return t_status, t_msg
}

/**
*  判断用户是否有权限
 */
func IsPower(ctx *gin.Context) (int, string) {
	t_status := 200
	t_msg := "success"
	uri := ctx.Request.RequestURI
	session := sessions.Default(ctx)
	power_path := fmt.Sprintf("%v", session.Get("power_path"))
	account := fmt.Sprintf("%v", session.Get("account"))
	if power_path == "all" {
		return t_status, t_msg
	}
	p_info := model.PowerByUrlRedis(uri)
	if len(p_info) < 1 {
		return t_status, t_msg
	}
	a_pwer := model.UserPowerRedis(account, p_info["code"])
	if len(a_pwer["id"]) < 1 {
		t_status = 403
		t_msg = "没有权限"
		return t_status, t_msg
	}
	return t_status, t_msg
}

/**
*  生成谷歌二维码
 */
func GoogleQr(ctx *gin.Context) (int, image.Image) {
	t_status := 100
	var t_img image.Image
	session := sessions.Default(ctx)
	sess_acc := session.Get("account")
	if sess_acc == nil {
		return t_status, t_img
	}
	account := fmt.Sprintf("%v", sess_acc)
	a_info := model.AdminInfoRedis(account)
	if len(a_info["account"]) < 1 {
		return t_status, t_img
	}
	secret := a_info["secret"]
	if len(secret) < 1 {
		secret = common.GetSecret()
		session.Set("secret", secret)
		session.Save()
	}

	header := map[string]string{}
	header["Content-Type"] = "charset=UTF-8"
	header["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"
	header["Connection"] = "	keep-alive"
	header["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36"
	c_url := common.GetImageUrl(a_info["account"], secret, 150, 150)
	t_status, t_img = common.HttpBodyByImg(c_url, "GET", "", header)
	if t_status != 200 {
		t_status, t_img = common.HttpBodyByImg(c_url, "GET", "", header)
	}

	return t_status, t_img
}

/**
*  绑定/验证谷歌验证码
 */
func CheckGoogleCode(code_val string, ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "验证码错误"
	session := sessions.Default(ctx)
	sess_acc := session.Get("account")
	if sess_acc == nil {
		t_msg = "用户信息异常"
		return t_status, t_msg
	}
	account := fmt.Sprintf("%v", sess_acc)

	secret := fmt.Sprintf("%v", session.Get("secret"))

	is_bind := true
	a_list := model.AdminInfo(account)
	if a_list.Id < 1 {
		t_msg = "用户信息异常"
		return t_status, t_msg
	}

	if len(a_list.Secret) < 1 && len(secret) < 1 {
		t_msg = "未绑定谷歌验证码,请绑定"
		return t_status, t_msg
	} else if len(a_list.Secret) > 1 {
		secret = a_list.Secret
	} else if len(a_list.Secret) < 1 {
		is_bind = false
	}

	t_status, t_msg = AuthGoogle(code_val, secret)
	if t_status != 200 {
		return t_status, t_msg
	}

	if is_bind {
		return t_status, t_msg

	}
	t_status = 100

	up_sql := fmt.Sprintf("update admin_list set secret='%s' where id='%d';", secret, a_list.Id)
	err := model.Query(up_sql)
	if err != nil {
		t_msg = "绑定失败"
		return t_status, t_msg
	}
	redis_key := "admin_list:" + a_list.Account
	redis.RediGo.KeyDel(redis_key)
	redis_key = "admin_list:session_id:" + a_list.Session_id
	redis.RediGo.KeyDel(redis_key)
	t_status = 200
	t_msg = "success"

	return t_status, t_msg
}
