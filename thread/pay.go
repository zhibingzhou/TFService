package thread

import (
	"fmt"
	"TFService/model"
	"TFService/payapi"
	"strconv"
	"time"

	"github.com/zhibingzhou/go_public/common"

	"github.com/zhibingzhou/go_public/redis"
)

/**
* 支付
 */
func PayCreate(mer_code string, p_map map[string]string) (int, string, map[string]string, int, map[string]string) {
	t_status := 100
	t_msg := "商户号错误"
	//1:扫码,2:调转,3:post提交
	api_jump := 3
	img_url := ""
	pay_url := ""
	api_method := ""
	form_param := map[string]string{}

	tpl_param := map[string]string{}

	//获取参数
	is_mobile := p_map["is_mobile"]
	amount := p_map["amount"]
	pay_id := p_map["pay_id"]
	order_number := p_map["order_number"]
	class_code := p_map["class_code"]
	bank_code := p_map["bank_code"]
	push_url := p_map["push_url"]
	ip := p_map["ip"]
	tpl_param["class_code"] = class_code
	tpl_param["is_mobile"] = is_mobile
	tpl_param["amount"] = amount
	tpl_param["img_url"] = ""

	redis_key := "PayCreate:" + order_number

	lock_res := redis.RediGo.StringWriteNx(redis_key, order_number, 5)
	if lock_res < 1 {
		t_msg = "请求频繁"
		return t_status, t_msg, tpl_param, api_jump, form_param
	}

	defer redis.RediGo.KeyDel(redis_key)

	is_exist := model.OrderByWebRedis(order_number)
	if len(is_exist["id"]) > 1 {
		t_msg = "订单已存在"
		return t_status, t_msg, tpl_param, api_jump, form_param
	}

	mer_info := model.MerInfoRedis(mer_code)
	if len(mer_info["id"]) < 1 {
		return t_status, t_msg, tpl_param, api_jump, form_param
	}

	if class_code == "" || bank_code == "" {
		t_msg = "银行编码和渠道类型不能为空"
		return t_status, t_msg, tpl_param, api_jump, form_param
	}

	//转账金额必须大于0
	amount_f, err := strconv.ParseFloat(amount, 64)

	if err != nil {
		t_msg = "支付金额格式错误"
		return t_status, t_msg, tpl_param, api_jump, form_param
	}

	//查看商户可以支付的渠道
	r_list := model.RateList(mer_code, class_code, bank_code)
	if len(r_list) < 1 {
		t_msg = "商户未分配渠道"
		return t_status, t_msg, tpl_param, api_jump, form_param
	}

	code_arr := []string{}
	for _, r_val := range r_list {
		if r_val.Id < 1 {
			continue
		}
		code_arr = append(code_arr[0:], r_val.Pay_code)
	}
	if len(code_arr) < 1 {
		t_msg = "商户未分配渠道"
		return t_status, t_msg, tpl_param, api_jump, form_param
	}

	//判断是否有渠道已满
	rate_list := model.PayRateList(code_arr)
	if len(rate_list) < 1 {
		t_msg = "暂无支付渠道"
		return t_status, t_msg, tpl_param, api_jump, form_param
	}

	pay_code_arr := []string{}
	for _, rate_val := range rate_list {
		if rate_val.Id < 1 {
			continue
		}
		pay_code_arr = append(pay_code_arr[0:], rate_val.Pay_code)
	}
	if len(pay_code_arr) < 1 {
		t_msg = "暂无支付渠道"
		return t_status, t_msg, tpl_param, api_jump, form_param
	}
	pin_field := "pay_code"
	p_fields := []string{"pay_code", "pay_id"}
	pp_where := map[string]interface{}{}
	pp_where["mer_code"] = mer_code
	pp_where["status"] = 1
	p_list, _ := model.InPageList("mer_pay", pin_field, 1000, 0, p_fields, pay_code_arr, pp_where)
	if len(p_list) < 1 {
		t_msg = "商户暂无支付渠道"
		return t_status, t_msg, tpl_param, api_jump, form_param
	}
	id_arr := []string{}
	for _, p_val := range p_list {
		if len(p_val["pay_id"]) < 1 {
			continue
		}
		id_arr = append(id_arr[0:], p_val["pay_id"])
	}
	if len(id_arr) < 1 {
		t_msg = "商户暂无支付渠道"
		return t_status, t_msg, tpl_param, api_jump, form_param
	}

	in_field := "id"
	fields := []string{"pay_code", "id"}
	p_where := map[string]interface{}{}
	p_where["status"] = 1
	conf_list, _ := model.InPageList("pay_config", in_field, 100, 0, fields, id_arr, p_where)
	if len(conf_list) < 1 {
		t_msg = "暂无支付渠道"
		return t_status, t_msg, tpl_param, api_jump, form_param
	}
	p_code_arr := []string{}
	for _, c_val := range conf_list {
		p_code_arr = append(p_code_arr[0:], c_val["pay_code"])
	}

	in_bank := "pay_code"
	b_fields := []string{"pay_bank", "jump_type", "pay_code", "class_code"}
	b_where := map[string]interface{}{}
	b_where["class_code"] = class_code
	b_where["bank_code"] = bank_code
	b_where["is_mobile"] = is_mobile
	bank_list, _ := model.InPageList("pay_bank", in_bank, 100, 0, b_fields, p_code_arr, b_where)
	b_len := len(bank_list)
	if b_len < 1 {
		t_msg = "该渠道未开通此支付方式"
		return t_status, t_msg, tpl_param, api_jump, form_param
	}

	b_index := common.RandomMaxAndMinInt(0, b_len-1)
	pay_bank := bank_list[b_index]
	pay_code := pay_bank["pay_code"]
	if len(pay_bank["pay_bank"]) < 1 {
		t_msg = "该渠道未开通此支付方式"
		return t_status, t_msg, tpl_param, api_jump, form_param
	}
	api_jump, _ = strconv.Atoi(pay_bank["jump_type"])

	for _, c_val := range conf_list {
		if pay_code == c_val["pay_code"] {
			pay_id = c_val["id"]
			break
		}
	}

	for _, r_val := range rate_list {
		if pay_code == r_val.Pay_code && class_code == r_val.Class_code && bank_code == r_val.Bank_code {
			if amount_f < r_val.Min_amount {
				t_msg = "支付的最小金额" + fmt.Sprintf("%.2f", r_val.Min_amount)
				return t_status, t_msg, tpl_param, api_jump, form_param
			}
			if r_val.Max_amount > 0.00 && amount_f > r_val.Max_amount {
				t_msg = "支付的最大金额" + fmt.Sprintf("%.2f", r_val.Max_amount)
				return t_status, t_msg, tpl_param, api_jump, form_param
			}
			break
		}
	}

	real_amout := 0.00
	m_rate := 0.00

	for _, r_val := range r_list {
		if pay_code == r_val.Pay_code && class_code == r_val.Class_code && bank_code == r_val.Bank_code {
			m_rate = r_val.Rate
			break
		}
	}
	if m_rate == 0.00 || m_rate > 1.00 {
		t_msg = "费率异常"
		return t_status, t_msg, tpl_param, api_jump, form_param
	}

	real_amout = (1 - m_rate) * amount_f
	//生成存款订单
	list_id := model.GetKey(17)
	create_time := time.Now().Format(format_date)

	table_name := "pay_list"
	pay_list := map[string]string{}
	pay_list["id"] = list_id
	pay_list["pay_code"] = pay_code
	pay_list["pay_id"] = pay_id
	pay_list["mer_code"] = mer_code
	pay_list["amount"] = amount
	pay_list["real_amount"] = fmt.Sprintf("%.4f", real_amout)
	pay_list["create_time"] = create_time
	pay_list["order_number"] = order_number
	pay_list["class_code"] = class_code
	pay_list["bank_code"] = bank_code
	pay_list["push_url"] = push_url
	pay_list["is_mobile"] = is_mobile
	//pay_list["status"] = "3"
	pay_list["rate"] = fmt.Sprintf("%.4f", m_rate)
	pay_list["agent_path"] = mer_info["agent_path"]

	pay_sql := common.InsertSql(table_name, pay_list)
	err = model.Query(pay_sql)
	if err != nil {
		t_msg = "订单插入失败，订单号重复"
		return t_status, t_msg, tpl_param, api_jump, form_param
	}
	api_pay := apiPayInit(pay_id, pay_code)
	if api_pay == nil {
		t_msg = "支付初始化失败"
		return t_status, t_msg, tpl_param, api_jump, form_param
	}
	t_msg = "支付请求失败"

	//对数据赋值
	pay_data := &payapi.PayData{Amount: amount, Order_number: list_id, Pay_bank: pay_bank["pay_bank"], Is_mobile: is_mobile, Ip: ip, Class_code: class_code}

	t_status, t_msg, api_method, pay_url, img_url, form_param = api_pay.CreatePay(pay_data)
	tpl_param["img_url"] = img_url
	tpl_param["pay_url"] = pay_url
	tpl_param["api_method"] = api_method
	return t_status, t_msg, tpl_param, api_jump, form_param
}

/**
* 支付回调
 */
func payCallBack(amount, sign, sign_str string, is_cent int, p_list model.PayList) (int, string) {
	t_status := 100
	t_msg := "订单金额错误"
	back_amout, _ := strconv.ParseFloat(amount, 64)
	if is_cent == 1 {
		back_amout = back_amout / 100.00
	}
	//有的支付额度是不定的
	back_amout = back_amout + 1.00
	if back_amout < p_list.Amount {
		return t_status, t_msg
	}
	pay_id := strconv.Itoa(p_list.Pay_id)
	api_pay := apiPayInit(pay_id, p_list.Pay_code)
	if api_pay == nil {
		t_msg = "订单初始化错误"
		return t_status, t_msg
	}
	t_status = api_pay.CallBackPay(sign, sign_str)
	if t_status != 200 {
		t_msg = "验签失败"
	}
	return t_status, t_msg
}

/**
* 下发回调
 */
func cashCallBack(amount, sign, sign_str string, is_cent int, p_list model.OrderList) (int, string) {
	t_status := 100
	t_msg := "订单金额错误"
	back_amout, _ := strconv.ParseFloat(amount, 64)
	if is_cent == 1 {
		back_amout = back_amout / 100.00
	}
	//有的支付额度是不定的
	back_amout = back_amout + 1.00
	if back_amout < p_list.Amount {
		return t_status, t_msg
	}
	pay_id := strconv.Itoa(p_list.Pay_id)
	api_pay := apiPayInit(pay_id, p_list.Pay_code)
	if api_pay == nil {
		t_msg = "订单初始化错误"
		return t_status, t_msg
	}
	t_status = api_pay.CallBackPay(sign, sign_str)
	if t_status != 200 {
		t_msg = "验签失败"
	}
	return t_status, t_msg
}

/**
* 查询网银支付支持的银行列表
 */
func Bank(is_mobile, pay_id string) (int, string, []map[string]string) {
	t_status := 100
	t_msg := "支付渠道ID错误"
	p_map := []map[string]string{}
	//根据pay_id查询cash_type_code
	p_conf := model.ApiConfigRedis(pay_id)
	if len(p_conf["pay_code"]) < 0 {
		return t_status, t_msg, p_map
	}
	b_list := model.BankList(is_mobile, p_conf["pay_code"])
	if len(b_list) < 1 {
		t_msg = "该支付渠道不支持网银"
		return t_status, t_msg, p_map
	}
	t_status = 200
	t_msg = "success"
	for _, bank_info := range b_list {
		bank_map := map[string]string{}
		bank_map["bank_code"] = bank_info.Bank_code
		bank_map["bank_title"] = bank_info.Bank_title
		p_map = append(p_map[0:], bank_map)
	}
	return t_status, t_msg, p_map
}

/**
*  查询支付订单信息
 */
func PayOrder(mer_code, order_number string) (int, string, map[string]string) {
	t_status := 100
	t_msg := "订单号不能为空"
	order_info := map[string]string{}
	table_name := "pay_list"
	if order_number == "" {
		return t_status, t_msg, order_info
	}
	p_where := map[string]interface{}{}
	p_where["order_number"] = order_number

	field := []string{"order_number", "id", "status", "pay_code", "mer_code", "pay_id", "push_status", "push_num", "amount", "real_amount", "create_time", "pay_time", "class_code", "bank_code", "note", "is_mobile", "rate"}
	order_info, _ = model.CommonFieldsRow(table_name, field, p_where)
	if order_info["id"] == "<nil>" {
		t_msg = "订单号错误"
		return t_status, t_msg, order_info
	}

	t_status = 200
	t_msg = "success"
	return t_status, t_msg, order_info
}

func CashOrder(mer_code, order_number string) (int, string, map[string]string) {
	t_status := 100
	t_msg := "订单号不能为空"
	order_info := map[string]string{}
	table_name := "cash_list"
	if order_number == "" {
		return t_status, t_msg, order_info
	}
	p_where := map[string]interface{}{}
	p_where["order_number"] = order_number

	field := []string{"order_number", "id", "status", "pay_code", "mer_code", "pay_id", "push_status", "push_num", "amount", "real_amount", "create_time", "pay_time", "bank_title", "bank_code", "note"}
	order_info, _ = model.CommonFieldsRow(table_name, field, p_where)
	if order_info["id"] == "<nil>" {
		t_msg = "订单号错误"
		return t_status, t_msg, order_info
	}

	t_status = 200
	t_msg = "success"
	return t_status, t_msg, order_info
}

/**
*  查询余额信息
 */
func PayBalance(mer_code string, pay_map map[string]string) (int, string, map[string]interface{}) {
	t_status := 100
	t_msg := "支付渠道不存在"
	mer_info := map[string]interface{}{}
	m_pay := model.MerInfo(mer_code)
	if m_pay.Id < 1 {
		return t_status, t_msg, mer_info
	}
	t_status = 200
	t_msg = "success"
	mer_info["amount"] = m_pay.Amount
	mer_info["total_in"] = m_pay.Total_in
	mer_info["total_out"] = m_pay.Total_out
	return t_status, t_msg, mer_info
}

/**
* 代付
 */
func PayFor(mer_code string, p_map map[string]string) (int, string, string) {
	t_status := 100

	t_msg := "订单号为空"
	re_result := "fail"

	web_order := p_map["order_number"]
	if len(web_order) < 1 {
		return t_status, t_msg, re_result
	}

	redis_key := "lock:PayFor:" + web_order

	lock_res := redis.RediGo.StringWriteNx(redis_key, web_order, 5)
	if lock_res < 1 {
		t_msg = "请求频繁"
		return t_status, t_msg, re_result
	}

	defer redis.RediGo.KeyDel(redis_key)

	is_exist := model.CashByWebRedis(web_order)
	if len(is_exist["id"]) > 1 {
		t_msg = "订单已存在"
		return t_status, t_msg, re_result
	}

	mer_info := model.MerInfo(mer_code)
	if mer_info.Id < 1 {
		t_msg = "商户号错误"
		return t_status, t_msg, re_result
	}

	amout_f, _ := strconv.ParseFloat(p_map["amount"], 64)
	if amout_f < 1.00 {
		t_msg = "代付金额错误"
		return t_status, t_msg, re_result
	}

	pay_conf := model.ApiConfigRedis(p_map["pay_id"])
	if pay_conf["status"] == "-1" {
		t_msg = "支付渠道维护中"
		return t_status, t_msg, re_result
	}

	bank_info := model.BankInfoRedis(p_map["bank_code"])
	if len(bank_info["title"]) < 1 {
		t_msg = "银行编码错误"
		return t_status, t_msg, re_result
	}

	//查询额度
	if mer_info.Amount < amout_f {
		t_msg = "额度不足"
		return t_status, t_msg, re_result
	}

	if amout_f < 10000 || amout_f > 50000 {
		t_msg = "下发金额范围在 10000-50000"
		return t_status, t_msg, re_result
	}

	create_time := time.Now().Format(format_date)
	//订单生成
	c_list := map[string]string{}
	c_id := model.GetKey(20)
	c_list["id"] = c_id
	c_list["status"] = "-1"
	c_list["bank_code"] = p_map["bank_code"]
	c_list["bank_title"] = bank_info["title"]
	c_list["card_number"] = p_map["card_number"]
	c_list["order_number"] = web_order
	c_list["mer_code"] = mer_code
	c_list["card_name"] = p_map["card_name"]
	c_list["pay_id"] = "0"
	c_list["push_url"] = p_map["push_url"]
	c_list["amount"] = p_map["amount"]
	c_list["create_time"] = create_time
	c_list["agent_path"] = mer_info.Agent_path
	c_list["branch"] = p_map["branch"]
	c_list["phone"] = p_map["phone"]
	c_list["pay_code"] = "all"
	table_name := "cash_list"
	c_sql := common.InsertSql(table_name, c_list)
	err := model.Query(c_sql)
	if err != nil {
		t_msg = "订单生成失败"
		return t_status, t_msg, re_result
	}
	note := ""
	order_type := 2
	t_status, t_msg = updateMerCash(web_order, "", note, order_type)
	if t_status != 200 {
		return t_status, t_msg, re_result
	}

	t_status = 200
	t_msg = "success"
	re_result = "processing"

	return t_status, t_msg, re_result
}

/**
*  查询下发订单
 */
func PayForQuery(cash_map map[string]string) (int, string, string) {
	t_status := 100
	t_msg := "请求频繁"
	re_result := "processing"

	redis_key := "lock:PayForQuery:" + cash_map["id"]

	lock_res := redis.RediGo.StringWriteNx(redis_key, cash_map["id"], 5)
	if lock_res < 1 {
		return t_status, t_msg, re_result
	}

	defer redis.RediGo.KeyDel(redis_key)

	api_pay := apiPayInit(cash_map["pay_id"], cash_map["pay_code"])

	//初始化api的配置
	if api_pay == nil {
		t_msg = "初始化失败"
		return t_status, t_msg, re_result
	}
	pay_data := map[string]string{}
	pay_data["order_number"] = cash_map["id"]
	t_status, t_msg, re_result = api_pay.PayQuery(pay_data)
	return t_status, t_msg, re_result
}

/**
* 代付
 */
func adminOrder(pay_id, amount string, c_list model.CashList) (int, string) {
	t_status := 100

	t_msg := "订单号为空"
	redis_key := "lock:adminOrder:" + pay_id + ":" + c_list.Order_number

	lock_res := redis.RediGo.StringWriteNx(redis_key, c_list.Order_number, 5)
	if lock_res < 1 {
		t_msg = "请求频繁"
		return t_status, t_msg
	}

	defer redis.RediGo.KeyDel(redis_key)

	amout_f, _ := strconv.ParseFloat(amount, 64)
	if amout_f < 1.00 {
		t_msg = "下发金额错误"
		return t_status, t_msg
	}

	pay_conf := model.ApiConfigRedis(pay_id)

	create_time := time.Now().Format(format_date)
	//订单生成
	c_data := map[string]string{}

	c_id := model.GetKey(20)
	c_data["id"] = c_id
	c_data["status"] = "-1"
	c_data["bank_code"] = c_list.Bank_code
	c_data["bank_title"] = c_list.Bank_title
	c_data["card_number"] = c_list.Card_number
	c_data["order_number"] = c_list.Order_number
	c_data["card_name"] = c_list.Card_name
	c_data["mer_code"] = c_list.Mer_code
	c_data["pay_id"] = pay_id
	c_data["amount"] = amount
	c_data["create_time"] = create_time
	c_data["cash_id"] = c_list.Id
	c_data["branch"] = c_list.Branch
	c_data["phone"] = c_list.Phone
	c_data["pay_code"] = pay_conf["pay_code"]
	table_name := "order_list"
	c_sql := common.InsertSql(table_name, c_data)
	err := model.Query(c_sql)
	if err != nil {
		t_msg = "订单生成失败"
		return t_status, t_msg
	}
	note := ""
	order_type := 4
	t_status, t_msg = updateMerCash(c_id, "", note, order_type)
	if t_status != 200 {
		return t_status, t_msg
	}

	api_pay := apiPayInit(pay_id, pay_conf["pay_code"])
	if api_pay == nil {
		t_msg = "支付初始化失败"
		t_status = 100
		return t_status, t_msg
	}
	bank_code := c_list.Bank_code
	bank_title := c_list.Bank_title

	cash_bank := model.CashBankRedis(pay_conf["pay_code"], bank_code)
	if len(cash_bank["bank_code"]) > 0 && cash_bank["bank_code"] != "nil" {
		bank_code = cash_bank["cash_bank"]
		bank_title = cash_bank["bank_title"]
	}

	pay_data := map[string]string{}
	pay_data["amount"] = amount                  //代付金额
	pay_data["bank_code"] = bank_code            //银行编码
	pay_data["card_name"] = c_list.Card_name     //持卡人姓名
	pay_data["card_number"] = c_list.Card_number //卡号
	pay_data["bank_branch"] = c_list.Branch      //支行信息
	pay_data["order_number"] = c_id              //代付单号
	pay_data["bank_phone"] = c_list.Phone        //手机号码
	pay_data["bank_title"] = bank_title          //银行名称
	pay_data["pay_time"] = create_time

	re_result := "error"
	t_status, t_msg, re_result = api_pay.PayFor(pay_data)

	if re_result == "fail" {
		order_status := "9"
		note := "下发失败,自动完成"
		updateOrderStatus(c_id, c_id, order_status, note)
	}

	t_status = 200
	t_msg = "success"

	return t_status, t_msg
}

/**
*  代付接口下发
 */
func DFpay(mer_code string, p_map map[string]string) (int, string, string) {
	t_status := 100

	t_msg := "订单号为空"
	t_result := "fail"

	web_order := p_map["order_number"]
	if len(web_order) < 1 {
		return t_status, t_msg, t_result
	}

	redis_key := "lock:PayFor:" + web_order

	lock_res := redis.RediGo.StringWriteNx(redis_key, web_order, 5)
	if lock_res < 1 {
		t_msg = "请求频繁"
		return t_status, t_msg, t_result
	}

	defer redis.RediGo.KeyDel(redis_key)

	is_exist := model.CashByWebRedis(web_order)
	if len(is_exist["id"]) > 1 {
		t_msg = "订单已存在"
		return t_status, t_msg, t_result
	}

	mer_info := model.MerInfo(mer_code)
	if mer_info.Id < 1 {
		t_msg = "商户号错误"
		return t_status, t_msg, t_result
	}

	amout_f, _ := strconv.ParseFloat(p_map["amount"], 64)
	if amout_f < 1.00 {
		t_msg = "代付金额错误"
		return t_status, t_msg, t_result
	}

	pay_conf := model.ApiConfigRedis(p_map["pay_id"])
	if len(pay_conf["id"]) < 1 {
		t_msg = "渠道编码错误"
		return t_status, t_msg, t_result
	}
	if pay_conf["status"] == "-1" {
		t_msg = "支付渠道维护中"
		return t_status, t_msg, t_result
	}

	bank_info := model.BankInfoRedis(p_map["bank_code"])
	if len(bank_info["title"]) < 1 {
		t_msg = "银行编码错误"
		return t_status, t_msg, t_result
	}

	pay_rate := model.PayRateInfo(pay_conf["pay_code"], "bank", "CARD") //暂时限定代付的class_code 和 bank_code 为bank 和 CARD
	if pay_rate.Id < 1 {
		t_msg = "上游费率异常"
		return t_status, t_msg, t_result
	}

	if amout_f < pay_rate.Min_amount || amout_f > pay_rate.Max_amount {
		t_msg = "代付金额不在允许范围内"
		return t_status, t_msg, t_result
	}

	mer_rate := model.MerRateInfo(mer_code, pay_conf["pay_code"], "bank", "CARD")
	if mer_rate.Id < 1 {
		t_msg = "未配置代付费率"
		return t_status, t_msg, t_result
	}

	//查询额度
	if mer_info.Amount < amout_f {
		t_msg = "额度不足"
		return t_status, t_msg, t_result
	}

	create_time := time.Now().Format(format_date)
	//订单生成
	c_list := map[string]string{}
	c_id := model.GetKey(20)
	c_list["id"] = c_id
	c_list["status"] = "-1"
	c_list["bank_code"] = p_map["bank_code"]
	c_list["bank_title"] = bank_info["title"]
	c_list["card_number"] = p_map["card_number"]
	c_list["order_number"] = web_order
	c_list["mer_code"] = mer_code
	c_list["card_name"] = p_map["card_name"]
	c_list["pay_id"] = p_map["pay_id"]
	c_list["push_url"] = p_map["push_url"]
	c_list["amount"] = p_map["amount"]
	c_list["create_time"] = create_time
	c_list["agent_path"] = mer_info.Agent_path
	c_list["branch"] = p_map["branch"]
	c_list["phone"] = p_map["phone"]
	c_list["pay_code"] = pay_conf["pay_code"]
	table_name := "cash_list"
	c_sql := common.InsertSql(table_name, c_list)
	err := model.Query(c_sql)
	if err != nil {
		t_msg = "订单生成失败"
		return t_status, t_msg, t_result
	}
	note := ""
	order_type := 6
	t_status, t_msg = updateMerCash(web_order, "", note, order_type)
	if t_status != 200 {
		return t_status, t_msg, t_result
	}
	t_status = 100
	bank_code := p_map["bank_code"]
	//查询银行编码信息
	cash_bank := model.CashBankRedis(pay_conf["pay_code"], p_map["bank_code"])
	if len(cash_bank["cash_bank"]) > 0 {
		bank_code = cash_bank["cash_bank"]
	}

	_, real_amount := ThreadEveDraw(pay_conf["pay_code"], amout_f)

	pay_data := map[string]string{}
	pay_data["amount"] = fmt.Sprintf("%.2f", real_amount) //代付金额
	pay_data["bank_code"] = bank_code                     //银行编码
	pay_data["card_name"] = p_map["card_name"]            //持卡人姓名
	pay_data["card_number"] = p_map["card_number"]        //卡号
	pay_data["bank_branch"] = p_map["bank_branch"]        //支行信息
	pay_data["order_number"] = t_msg                      //代付单号
	pay_data["bank_ext"] = p_map["bank_ext"]              //身份证
	pay_data["bank_phone"] = p_map["bank_phone"]          //手机号码
	pay_data["bank_province"] = p_map["bank_province"]    //省份
	pay_data["bank_city"] = p_map["bank_city"]            //城市
	pay_data["bank_area"] = p_map["card_name"]            //区
	pay_data["bank_cnaps"] = p_map["bank_cnaps"]          //联行号
	pay_data["bank_title"] = bank_info["title"]           //银行名称
	pay_data["pay_time"] = create_time

	api_pay := apiPayInit(p_map["pay_id"], pay_conf["pay_code"])
	if api_pay == nil {
		t_msg = "支付初始化失败"
		return t_status, t_msg, t_result
	}

	t_status, t_msg, t_result = api_pay.PayFor(pay_data)
	return t_status, t_msg, t_result
}
