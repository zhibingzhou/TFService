package thread

import (
	"TFService/model"
	"TFService/payapi"
	"encoding/json"
	"strconv"

	"github.com/zhibingzhou/go_public/common"
)

//创建映射结构体
type PAYAPI interface {
	//初始化,加载配置
	Init(string)
	//创建支付订单
	CreatePay(*payapi.PayData) (int, string, string, string, string, map[string]string)
	//回调验证
	CallBackPay(string, string) int
	//代付
	PayFor(map[string]string) (int, string, string)
	//下发订单查询
	PayQuery(map[string]string) (int, string, string)
}

var format_date string = "2006-01-02 15:04:05"
var format_day string = "2006-01-02"

type UpdateCash struct {
	Order_number string
	Pay_order    string
	Note         string
	Order_type   int
	PoolRes      common.PoolResult
}

//配置数据的缓存key集合
var Conf_Redis_Key = "Conf_Redis_Key"

//数据缓存的key集合
var Data_Redis_Key = "Data_Redis_Key"

var push_header map[string]string
var workerPool *common.WorkerPool

var log_path string

func init() {
	conf_byte, err := common.ReadFile("./conf/conf.json")

	if err != nil {
		panic(err)
	}
	var json_conf map[string]string
	//解析json格式r
	err = json.Unmarshal(conf_byte, &json_conf)
	if err != nil {
		panic(err)
	}
	log_path = json_conf["log_path"]
	push_header = make(map[string]string)
	push_header["Content-Type"] = "application/x-www-form-urlencoded; charset=UTF-8"
	push_header["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36"
	workerPool = common.NewWorkerPool(1)
	workerPool.Run()
}

/**
* 验证谷歌验证码
 */
func AuthGoogle(code_val, secret string) (int, string) {
	t_status := 100
	t_msg := "验证码不能为空"

	if len(code_val) < 1 || len(secret) < 1 {
		return t_status, t_msg
	}

	google := common.SetGoogleAuth(30, 6)
	t_code_int, _ := strconv.ParseInt(code_val, 10, 64)
	t_status, t_msg = google.CheckGoogleCode(secret, t_code_int)
	return t_status, t_msg
}

/**
*  处理分页
 */
func ThreadPage(page, page_size string) (int, int) {
	page_int, _ := strconv.Atoi(page)
	if page_int < 1 {
		page_int = 1
	}
	size_int, _ := strconv.Atoi(page_size)
	if size_int < 1 {
		size_int = 20
	} else if size_int > 100 {
		size_int = 100
	}
	return page_int, size_int
}

/**
* 初始化支付API
 */
func apiPayInit(pay_id, pay_code string) PAYAPI {
	//对结构体实例化
	var api_Pay PAYAPI

	switch pay_code {
	case "zf66":
		api_Pay = new(payapi.ZF66)
		api_Pay.Init(pay_id)
	case "ipay":
		api_Pay = new(payapi.IPAY)
		api_Pay.Init(pay_id)
	case "day88":
		api_Pay = new(payapi.DAY88)
		api_Pay.Init(pay_id)
	case "dfpay":
		api_Pay = new(payapi.DFPAY)
		api_Pay.Init(pay_id)
	case "hytpay":
		api_Pay = new(payapi.HYTPAY)
		api_Pay.Init(pay_id)
	case "newpay":
		api_Pay = new(payapi.NEWPAY)
		api_Pay.Init(pay_id)
	case "jxpay":
		api_Pay = new(payapi.JXPAY)
		api_Pay.Init(pay_id)
	case "rlpay":
		api_Pay = new(payapi.RLPAY)
		api_Pay.Init(pay_id)
	case "yeopay":
		api_Pay = new(payapi.YEOPAY)
		api_Pay.Init(pay_id)
	case "stpay":
		api_Pay = new(payapi.STPAY)
		api_Pay.Init(pay_id)
	case "bfpay":
		api_Pay = new(payapi.BFPAY)
		api_Pay.Init(pay_id)
	case "yfpay":
		api_Pay = new(payapi.YFPAY)
		api_Pay.Init(pay_id)
	case "jpay":
		api_Pay = new(payapi.JPAY)
		api_Pay.Init(pay_id)
	case "sxpay":
		api_Pay = new(payapi.SXPAY)
		api_Pay.Init(pay_id)
	case "pgpay": //同jpay
		api_Pay = new(payapi.JPAY)
		api_Pay.Init(pay_id)
	case "thpay":
		api_Pay = new(payapi.THPAY)
		api_Pay.Init(pay_id)
	case "ggpay":
		api_Pay = new(payapi.GGPAY)
		api_Pay.Init(pay_id)
	case "xxfpay":
		api_Pay = new(payapi.XXFPAY)
		api_Pay.Init(pay_id)
	case "zofpay":
		api_Pay = new(payapi.ZOFPAY)
		api_Pay.Init(pay_id)
	case "cfpay":
		api_Pay = new(payapi.CFPAY)
		api_Pay.Init(pay_id)
	case "c2cpay":
		api_Pay = new(payapi.C2CPAY)
		api_Pay.Init(pay_id)
	case "xpay":
		api_Pay = new(payapi.XPAY)
		api_Pay.Init(pay_id)
	case "xwspay":
		api_Pay = new(payapi.XWSPAY)
		api_Pay.Init(pay_id)
	case "mbpay":
		api_Pay = new(payapi.MBPAY)
		api_Pay.Init(pay_id)

	}
	return api_Pay
}

/**
*  每一个渠道的下发的费率计算
 */
func ThreadEveDraw(pay_code string, amount float64) (float64, float64) {
	// fee_amount := 5.00
	// amount_f := amount + fee_amount
	// real_amount := amount
	// return amount_f, real_amount
	real_amount := 0.00
	pay_chan := model.PayInfoRedis(pay_code)
	amount_f := 0.00
	if len(pay_chan["code"]) < 1 {
		return amount_f, real_amount
	}
	fee_amount, _ := strconv.ParseFloat(pay_chan["fee_amount"], 64)
	real_amount = amount
	amount_f = amount
	//下发扣费规则
	switch pay_chan["fee_type"] {
	case "1":
		amount_f = amount + fee_amount
	case "2":
		amount_f = amount
		real_amount = amount - fee_amount
	}

	return amount_f, real_amount
}
