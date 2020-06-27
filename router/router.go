package router

import (
	"TFService/controller"
	"TFService/hook"
	"TFService/thread"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/zhibingzhou/go_public/common"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

type JsonOut struct {
	Status int
	Msg    string
	Data   map[string]interface{}
}

//判断域名是否合法
func HookAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		mer_code := ctx.PostForm("mer_code")
		//真实IP
		ip := ctx.ClientIP()
		h_status, h_msg := hook.AuthIp(mer_code, ip)
		if h_status != 200 && ip != "127.0.0.1" {
			ctx.Abort()
			ctx.JSON(http.StatusOK, gin.H{"Status": 403, "Msg": h_msg})
		}
	}
}

//判断是否登录
func LoginAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//判断是否登录
		r_status, r_msg := thread.ThreadIsLogin(ctx)
		if r_status != 200 {
			ctx.Abort()
			ctx.JSON(http.StatusOK, gin.H{"Status": r_status, "Msg": r_msg})
		}
	}
}

var Router *gin.Engine

func init() {
	Router = gin.New()

	file_name := "./conf/redis.json"

	conf_byte, err := common.ReadFile(file_name)
	if err != nil {
		panic(err)
	}
	var json_conf map[string]string
	//解析json格式
	err = json.Unmarshal(conf_byte, &json_conf)
	if err != nil {
		panic(err)
	}
	max_idle, _ := strconv.Atoi(json_conf["max_idle"])
	if max_idle < 1 {
		max_idle = 10
	}

	addr := fmt.Sprintf("%s:%s", json_conf["host"], json_conf["port"])

	store, err := redis.NewStore(max_idle, "tcp", addr, json_conf["auth"], []byte(json_conf["pre_key"]))

	Router.Use(sessions.Sessions(json_conf["pre_key"], store))

	Router.Static("/static", "./assets")

	//静态文件路径，一定需要
	Router.LoadHTMLGlob("view/**/*")
	Router.LoadHTMLFiles("./view/index.tpl", "./view/admin.tpl", "./view/success.tpl", "./view/jump.tpl", "./view/put.tpl", "./view/auto_jump.tpl", "./view/new_jump.tpl", "./view/error.tpl")
	Router.GET("/test/test.do", controller.Test)
	Router.GET("/test/test_admin.do", controller.TestAdmin)
	Router.POST("/test/encode.do", controller.TestEncode)

	//API接口
	pay_group := Router.Group("/pay")
	//支付接口
	pay_group.POST("/create.do", controller.PayCreate)
	//支持的银行列表
	pay_group.POST("/bank.do", controller.Bank)
	//余额及总存款和总出款
	pay_group.POST("/balance.do", controller.PayBalance)
	//支付的订单查询
	pay_group.POST("/pay_order.do", controller.PayOrder)
	//下发的订单查询
	pay_group.POST("/cash_order.do", controller.CashOrder)
	//代收下发
	pay_group.POST("/pay_for.do", controller.PayFor)
	//纯代付系统
	pay_group.POST("/df_pay.do", controller.DFpay)

	//支付回调接口
	call_group := Router.Group("/call")
	call_group.POST("/dypay.do", controller.CallDypay)
	call_group.GET("/zf66.do", controller.CallZf66)
	call_group.GET("/day88.do", controller.CallDay88)
	call_group.POST("/hytpay.do", controller.CallHytpay)
	call_group.POST("/newpay.do", controller.CallNewpay)
	call_group.POST("/jxpay.do", controller.CallJxpay)
	call_group.POST("/rlpay.do", controller.CallRlpay)
	call_group.POST("/yeopay.do", controller.CallYeopay)
	call_group.POST("/stpay.do", controller.CallStpay)
	call_group.POST("/bfpay.do", controller.CallBfpay)
	call_group.GET("/yfpay.do", controller.CallYfpay)
	call_group.POST("/jpay.do", controller.CallJpay)
	call_group.POST("/sxpay.do", controller.CallSxpay)
	call_group.POST("/thpay.do", controller.Callthpay)
	call_group.POST("/ggpay.do", controller.Callggpay)
	call_group.POST("/xpay.do", controller.Callxpay)
	call_group.POST("/xxfpay.do", controller.Callxxfpay)
	call_group.POST("/zofpay.do", controller.Callzofpay)
	call_group.POST("/cfpay.do", controller.Callcfpay)
	call_group.POST("/xwspay.do", controller.Callxwspay)

	//下发回调接口
	back_group := Router.Group("/back")
	back_group.POST("/newpay.do", controller.BackNewpay)
	back_group.POST("/c2cpay.do", controller.BackC2Cpay)
	back_group.POST("/bfpay.do", controller.Backbfpay)
	back_group.POST("/xpay.do", controller.Backxpay)
	back_group.POST("/mbpay.do", controller.Backmbpay)

	//公共接口public
	public_group := Router.Group("/public")
	//退出登录
	public_group.POST("/logout.do", controller.Logout)
	//管理员登录
	public_group.POST("/login.do", controller.Login)
	//支付成功的返回页面
	public_group.GET("/success.do", controller.PaySuccess)
	//生成谷歌动态验证二维码
	public_group.GET("/google_qr.do", controller.GoogleQr)

	//验证类型接口
	check_group := Router.Group("/check")
	//验证谷歌动态
	check_group.POST("/check_google.do", controller.CheckGoogleCode)

	//管理员的接口
	admin_group := Router.Group("/admin", LoginAuth())
	//管理员更新密码
	admin_group.POST("/update_pwd.do", controller.UpdatePwd)
	//管理员列表
	admin_group.POST("/admin_list.do", controller.AdminList)
	//取消绑定
	admin_group.POST("/del_bind.do", controller.DelBind)
	//新增管理员
	admin_group.POST("/admin_add.do", controller.AddAdmin)
	//新增/修改管理员权限
	admin_group.POST("/update_power.do", controller.UpdatePower)
	//操作管理员
	admin_group.POST("/edit_admin.do", controller.EditAdmin)
	//清除缓存
	admin_group.POST("/del_cache.do", controller.DelCache)
	//查询支付订单
	admin_group.POST("/pay_list.do", controller.PayList)
	//查询下发订单
	admin_group.POST("/cash_list.do", controller.CashList)
	//查询商户流水
	admin_group.POST("/amount_list.do", controller.AmountList)
	//一段时间内的总收入,总出款
	admin_group.POST("/date_total.do", controller.DateTotal)
	//可提现额度，团队可提现额度
	admin_group.POST("/total_balance.do", controller.TotalBalance)
	//一段时间内的总收入,总出款
	admin_group.POST("/today_count.do", controller.TodayCount)
	//商户信息
	admin_group.POST("/mer_info.do", controller.MerInfo)
	//商户费率信息
	admin_group.POST("/rate_info.do", controller.RateInfo)
	//下载支付顶顶那
	admin_group.POST("/down_pay.do", controller.DownPayList)
	//商户费率信息
	admin_group.POST("/down_cash.do", controller.DownCashList)
	//商户费率信息
	admin_group.POST("/down_amount.do", controller.DownAmountList)
	//支付回调
	admin_group.POST("/call_pay.do", controller.CallPay)
	//下发回调
	admin_group.POST("/call_cash.do", controller.CallCash)
	//支付的状态修改
	admin_group.POST("/pay_status.do", controller.PayStatus)
	//下发的状态修改
	admin_group.POST("/cash_status.do", controller.CashStatus)
	//商户的渠道额度
	admin_group.POST("/mer_pay.do", controller.MerPay)
	//后台下发
	admin_group.POST("/mer_cash.do", controller.MerCash)
	//后台绑定的银行卡列表
	admin_group.POST("/mer_bank.do", controller.MerBank)
	//后台绑定银行卡
	admin_group.POST("/add_bank.do", controller.AddBank)
	//后台锁定银行卡
	admin_group.POST("/lock_bank.do", controller.LockBank)
	//代理报表
	admin_group.POST("/agent_report.do", controller.AgentReport)
	//商户渠道配置列表
	admin_group.POST("/mer_channel.do", controller.MerChannel)
	//新增商户渠道
	admin_group.POST("/add_mer_channel.do", controller.AddMerChannel)
	//修改商户渠道
	admin_group.POST("/edit_mer_channel.do", controller.EditMerChannel)
	//渠道下发订单列表
	admin_group.POST("/order_list.do", controller.OrderList)
	//新增渠道下发订单
	admin_group.POST("/add_order.do", controller.AddOrder)
	//更新渠道订单状态
	admin_group.POST("/update_order.do", controller.UpdateOrder)
	//给商户新增额度
	admin_group.POST("/add_mer_amount.do", controller.AddMerAmount)
	//人工充值
	admin_group.POST("/manual_recharge.do", controller.ManualRecharge)
	//首页商户信息
	admin_group.POST("/pay_mer_detail.do", controller.PayMerDetail)

	//系统设置的接口
	sys_group := Router.Group("/sys", LoginAuth())
	//权限列表
	sys_group.POST("/user_power.do", controller.PowerList)
	//权限列表
	sys_group.POST("/power_list.do", controller.PowerList)
	//新增权限
	sys_group.POST("/add_power.do", controller.AddPower)
	//新增商户
	sys_group.POST("/add_mer.do", controller.AddMer)
	//商户列表
	sys_group.POST("/mer_list.do", controller.MerList)
	//修改商户信息
	sys_group.POST("/update_mer.do", controller.UpdateMer)
	//新增上游支付渠道
	sys_group.POST("/add_pay.do", controller.AddPay)
	//支付类型列表
	sys_group.POST("/pay_class.do", controller.PayClass)
	//新增上游支付渠道信息
	sys_group.POST("/add_pay_class.do", controller.AddPayClass)
	//上游支付列表
	sys_group.POST("/channel_list.do", controller.ChannelList)
	//上游渠道费率详情
	sys_group.POST("/pay_detail.do", controller.PayDetail)
	//修改上游渠道费率详情
	sys_group.POST("/update_pay.do", controller.UpdatePay)
	//新增/修改商户的支付费率
	sys_group.POST("/mer_rate.do", controller.MerRate)
	//删除商户的支付费率
	sys_group.POST("/del_mer_rate.do", controller.DelMerRate)
	//商户的支付渠道及费率列表
	sys_group.POST("/mer_rate_list.do", controller.MerRateList)
	//新增商户IP白名单
	sys_group.POST("/add_ip.do", controller.AddIp)
	//白名单列表
	sys_group.POST("/ip_list.do", controller.IpList)
	//删除IP白名单
	sys_group.POST("/del_ip.do", controller.DelIp)
	//公告列表
	sys_group.POST("/notice_list.do", controller.NoticeList)
	//新增公告
	sys_group.POST("/add_notice.do", controller.AddNotice)
	//关闭公告
	sys_group.POST("/update_notice.do", controller.UpdateNotice)
	//系统银行列表
	sys_group.POST("/sys_bank.do", controller.SysBank)
	//支付配置列表
	sys_group.POST("/pay_conf.do", controller.PayConf)
	//上游渠道编码表
	sys_group.POST("/pay_bank.do", controller.PayBank)
	//新增上游渠道编码
	sys_group.POST("/add_pay_bank.do", controller.AddPayBank)
}
