package payapi

import (
	"encoding/json"

	"github.com/zhibingzhou/go_public/common"
)

/**
* 定义需要传递到api的数据结构
 */
type PayData struct {
	Amount       string //订单金额
	Order_number string //订单编号
	Class_code   string //支付类型
	Pay_bank     string //选择的银行或第三方支付平台
	Is_mobile    string //是否手机版 0是网页  1是手机
	Ip           string //客户的IP
}

var date_format string = "2006-01-02 15:04:05"

var day_f string = "20060102150405"

var day_format string = "2006-01-02"

var month_format string = "2006-01"

var s_format string = "2006-01-02 15:04:05"

var log_path string

var http_header map[string]string

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
	http_header = make(map[string]string)
	http_header["Content-type"] = "application/x-www-form-urlencoded; charset=UTF-8"
	http_header["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.79 Safari/537.36"
}
