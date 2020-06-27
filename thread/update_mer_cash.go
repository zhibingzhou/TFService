package thread

import (
	"TFService/model"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/zhibingzhou/go_public/common"
)

/**
*  更新用户额度
 */
func updateMerCash(order_number, pay_order, note string, order_type int) (int, string) {
	sc := &UpdateCash{Order_number: order_number, Pay_order: pay_order, Note: note, Order_type: order_type}
	workerPool.JobQueue <- sc
	pool_res := <-workerPool.PoolRes
	return pool_res.Status, pool_res.Msg
}

func (uc *UpdateCash) Do() {
	t_status := 100
	t_msg := "额度更新错误"
	//order_type:1=支付,2=下发,3=代付失败返回,4=后台下发到卡,5=上游下发失败返还,6=纯代付下发,7=支付返还
	switch uc.Order_type {
	case 1:
		t_status, t_msg = updatePay(uc.Order_number, uc.Pay_order, uc.Note)
	case 2:
		t_status, t_msg = cashOrder(uc.Order_number, uc.Pay_order, uc.Note)
	case 3:
		t_status, t_msg = updateReturnCash(uc.Order_number, uc.Pay_order, uc.Note)
	case 4:
		t_status, t_msg = order(uc.Order_number, uc.Pay_order, uc.Note)
	case 5:
		t_status, t_msg = updateReturnOrder(uc.Order_number, uc.Pay_order, uc.Note)
	case 6:
		t_status, t_msg = dfOrder(uc.Order_number, uc.Pay_order, uc.Note)
	}
	uc.PoolRes = common.PoolResult{Status: t_status, Msg: t_msg, JsonRes: ""}
}

func (uc *UpdateCash) GetResult() common.PoolResult {
	return uc.PoolRes
}

/**
*  支付完成
 */
func updatePay(order_number, pay_order, note string) (int, string) {
	t_status := 100
	t_msg := "订单号错误"
	p_list := model.OrderById(order_number)
	if len(p_list.Id) < 1 {
		return t_status, t_msg
	}
	agent_path := strings.TrimRight(p_list.Agent_path, "_")
	//判断代理的层级
	agent_arr := strings.Split(agent_path, "_")
	agent_len := len(agent_arr)
	if agent_len < 1 {
		t_msg = "路径异常"
		return t_status, t_msg
	}
	sql_arr := []string{}
	insert_sql := ""
	update_sql := ""
	a_status := 100
	//根据代理计算得到每条sql语句
	for i := 0; i < agent_len; i++ {
		note = "自动完成"
		if i == agent_len-1 {
			a_status, insert_sql, update_sql = payAgentSql(agent_arr[i], p_list.Mer_code, p_list)
		} else {
			a_status, insert_sql, update_sql = payAgentSql(agent_arr[i], agent_arr[i+1], p_list)
		}
		if a_status != 200 {
			note = "费率异常"
		} else if a_status == 200 && insert_sql != "" {
			sql_arr = append(sql_arr[0:], insert_sql)
			sql_arr = append(sql_arr[0:], update_sql)
		}
	}
	pay_time := time.Now().Format(format_date)
	//更新订单状态
	list_update := fmt.Sprintf("update pay_list set `status`=3,note='%s',pay_time='%s' where id='%s';", note, pay_time, p_list.Id)
	//更新余额
	m_pay := model.MerInfo(p_list.Mer_code)
	if m_pay.Id < 1 {
		t_msg = "商户异常"
		return t_status, t_msg
	}
	//更新语句
	up_sql := fmt.Sprintf("update mer_list set total_in=total_in+%.4f,amount=amount+%.4f where id='%d';", p_list.Real_amount, p_list.Real_amount, m_pay.Id)
	after_amount := p_list.Real_amount + m_pay.Amount    //算上费率后的余额
	after_amount_no_rate := p_list.Amount + m_pay.Amount //不算上费率后的余额
	//查询上游的费率
	pay_rate := model.PayRateInfo(p_list.Pay_code, p_list.Class_code, p_list.Bank_code)
	if pay_rate.Id < 1 {
		t_msg = "上游费率异常"
		return t_status, t_msg
	}
	p_real := p_list.Amount * (1 - pay_rate.Rate)
	//更新渠道的额度
	conf_sql := fmt.Sprintf("update pay_config set total_in=total_in+%.4f,amount=amount+%.4f where id='%d';", p_real, p_real, p_list.Pay_id)
	//账变记录语句 支付类型 1.
	table_name := "amount_list"
	a_data := map[string]string{}
	a_data["id"] = model.GetKey(20)
	a_data["order_number"] = p_list.Id
	a_data["create_time"] = time.Now().Format(format_date)
	a_data["amount_type"] = "1"
	a_data["amount"] = fmt.Sprintf("%.4f", p_list.Amount)
	a_data["pay_code"] = p_list.Pay_code
	a_data["mer_code"] = p_list.Mer_code
	a_data["pay_id"] = fmt.Sprintf("%d", p_list.Pay_id)
	a_data["before_amount"] = fmt.Sprintf("%.4f", m_pay.Amount)
	a_data["after_amount"] = fmt.Sprintf("%.4f", after_amount_no_rate)
	a_data["agent_path"] = p_list.Agent_path
	a_data["note"] = "支付"
	in_sql_pay := common.InsertSql(table_name, a_data)

	//支付手续费  加上 -1 是为了体现扣手续费
	rate_amout := -1 * p_list.Rate * p_list.Amount
	//账变记录语句 支付费率类型 6
	r_data := map[string]string{}
	r_data["id"] = model.GetKey(20)
	r_data["order_number"] = p_list.Id
	r_data["create_time"] = time.Now().Format(format_date)
	r_data["amount_type"] = "6"
	r_data["amount"] = fmt.Sprintf("%.4f", rate_amout)
	r_data["pay_code"] = p_list.Pay_code
	r_data["mer_code"] = p_list.Mer_code
	r_data["pay_id"] = fmt.Sprintf("%d", p_list.Pay_id)
	r_data["before_amount"] = fmt.Sprintf("%.4f", after_amount_no_rate)
	r_data["after_amount"] = fmt.Sprintf("%.4f", after_amount)
	r_data["agent_path"] = p_list.Agent_path
	r_data["note"] = "支付手续费"
	in_sql_rate := common.InsertSql(table_name, r_data)

	sql_arr = append(sql_arr[0:], conf_sql)
	sql_arr = append(sql_arr[0:], list_update)
	sql_arr = append(sql_arr[0:], up_sql)
	sql_arr = append(sql_arr[0:], in_sql_pay)
	sql_arr = append(sql_arr[0:], in_sql_rate)

	err := model.Trans(sql_arr)

	t_msg = "订单更新失败"
	//将结果推送到网站去
	if err != nil {
		return t_status, err.Error()
	}
	t_status = 200
	t_msg = "success"
	return t_status, t_msg
}

/**
*  结算代理的费用
 */
func payAgentSql(p_agent, agent string, p_list model.PayList) (int, string, string) {
	t_status := 100
	in_sql := ""
	up_sql := ""
	var err error

	p_rate_f := 0.00
	p_rate := map[string]string{}
	p_rate = model.MerRateRedis(p_agent, p_list.Pay_code, p_list.Class_code, p_list.Bank_code)
	if len(p_rate["rate"]) < 1 {
		return t_status, in_sql, up_sql
	}
	p_rate_f, err = strconv.ParseFloat(p_rate["rate"], 64)
	if err != nil {
		return t_status, in_sql, up_sql
	}

	if p_rate_f > p_list.Rate {
		return t_status, in_sql, up_sql
	}
	//查询下游费率
	a_rate := model.MerRateRedis(agent, p_list.Pay_code, p_list.Class_code, p_list.Bank_code)
	if len(a_rate["rate"]) < 1 {
		return t_status, in_sql, up_sql
	}

	a_rate_f, err := strconv.ParseFloat(a_rate["rate"], 64)
	//代理的盈利费率 = 下游费率 - 代理费率
	rate_f := a_rate_f - p_rate_f
	if err != nil || rate_f < 0.00 || a_rate_f > p_list.Rate {
		return t_status, in_sql, up_sql
	}
	//查询当前额度
	p_pay := model.MerInfo(p_agent)
	if p_pay.Id < 1 {
		return t_status, in_sql, up_sql
	}
	t_status = 200

	//代理的赢利额度
	amount_f := p_list.Amount * rate_f
	if amount_f < 0.0001 {
		return t_status, in_sql, up_sql
	}

	after_amount := p_pay.Amount + amount_f
	//更新语句
	up_sql = fmt.Sprintf("update mer_list set total_in=total_in+%.4f,amount=amount+%.4f where id='%d';", amount_f, amount_f, p_pay.Id)
	//账变记录语句
	//代理的佣金
	table_name := "amount_list"
	a_data := map[string]string{}
	a_data["id"] = model.GetKey(20)
	a_data["order_number"] = p_list.Id
	a_data["create_time"] = time.Now().Format(format_date)
	a_data["amount_type"] = "3"
	a_data["amount"] = fmt.Sprintf("%.4f", amount_f)
	a_data["pay_code"] = p_list.Pay_code
	a_data["mer_code"] = p_agent
	a_data["pay_id"] = fmt.Sprintf("%d", p_list.Pay_id)
	a_data["before_amount"] = fmt.Sprintf("%.4f", p_pay.Amount)
	a_data["after_amount"] = fmt.Sprintf("%.4f", after_amount)
	a_data["agent_path"] = p_list.Agent_path
	a_data["note"] = "支付的代理佣金"
	in_sql = common.InsertSql(table_name, a_data)
	//额度更新语句和账变记录语句
	return t_status, in_sql, up_sql
}

/**
*  创建代付订单
*  需要将一笔订单拆分成几笔
 */
func cashOrder(order_number, pay_order, note string) (int, string) {
	t_status := 100
	t_msg := "信息异常"
	c_list := model.CashByWeb(order_number)
	if len(c_list.Id) < 1 {
		return t_status, t_msg
	}
	if c_list.Status != -1 {
		t_msg = "订单状态异常"
		return t_status, t_msg
	}
	//查询额度
	m_info := model.MerInfo(c_list.Mer_code)
	if m_info.Id < 1 {
		return t_status, t_msg
	}

	fee_amount := 2.00
	order_amount := c_list.Amount + fee_amount
	real_amount := c_list.Amount
	if m_info.Amount < order_amount {
		t_msg = "额度不足"
		return t_status, t_msg
	}

	//更新语句
	update_sql := fmt.Sprintf("update cash_list set `status`=1,real_amount=%.2f,fee_amount=%.2f,order_amount=%.2f where id='%s';", real_amount, fee_amount, order_amount, c_list.Id)
	//更新用户额度
	mer_sql := fmt.Sprintf("update mer_list set amount=amount-%.2f,total_out=total_out+%.2f where id='%d';", order_amount, order_amount, m_info.Id)
	//新增账变记录
	after_amount := m_info.Amount - order_amount
	table_name := "amount_list"
	a_data := map[string]string{}
	a_data["id"] = model.GetKey(20)
	a_data["order_number"] = c_list.Id
	a_data["create_time"] = c_list.Create_time
	a_data["amount_type"] = "2"
	a_data["amount"] = fmt.Sprintf("%.4f", -1*order_amount)
	a_data["pay_code"] = c_list.Pay_code
	a_data["mer_code"] = c_list.Mer_code
	a_data["pay_id"] = fmt.Sprintf("%d", c_list.Pay_id)
	a_data["before_amount"] = fmt.Sprintf("%.4f", m_info.Amount)
	a_data["after_amount"] = fmt.Sprintf("%.4f", after_amount)
	a_data["agent_path"] = c_list.Agent_path
	a_data["note"] = "支付的代理费率"
	in_sql := common.InsertSql(table_name, a_data)
	sql_arr := []string{update_sql, mer_sql, in_sql}
	err := model.Trans(sql_arr)
	if err != nil {
		t_msg = "下发订单处理失败"
		return t_status, t_msg
	}
	t_status = 200
	t_msg = "success"
	//新增账变记录
	return t_status, t_msg
}

/**
*  创建代付订单
 */
func order(order_id, pay_order, note string) (int, string) {
	t_status := 100
	t_msg := "信息异常"
	o_list := model.CashOrderById(order_id)
	if len(o_list.Id) < 1 {
		return t_status, t_msg
	}
	if o_list.Status != -1 {
		t_msg = "订单状态异常"
		return t_status, t_msg
	}
	pay_id := fmt.Sprintf("%d", o_list.Pay_id)
	pay_conf := model.PayConf(pay_id)
	order_amount, real_amount := ThreadEveDraw(o_list.Pay_code, o_list.Amount)
	if pay_conf.Amount < order_amount {
		t_msg = "额度不足"
		return t_status, t_msg
	}

	fee_amount := order_amount - real_amount

	//更新语句
	update_sql := fmt.Sprintf("update order_list set `status`=1,real_amount=%.2f,fee_amount=%.2f,order_amount=%.2f where id='%s';", real_amount, fee_amount, order_amount, o_list.Id)
	//更新用户额度
	mer_sql := fmt.Sprintf("update pay_config set amount=amount-%.2f,total_out=total_out+%.2f where id='%s';", order_amount, order_amount, pay_id)
	//新增账变记录
	after_amount := pay_conf.Amount - order_amount
	table_name := "amount_list"
	a_data := map[string]string{}
	a_data["id"] = model.GetKey(20)
	a_data["order_number"] = o_list.Id
	a_data["create_time"] = o_list.Create_time
	a_data["amount_type"] = "9"
	a_data["amount"] = fmt.Sprintf("%.4f", -1*order_amount)
	a_data["pay_code"] = o_list.Pay_code
	a_data["mer_code"] = "all"
	a_data["pay_id"] = pay_id
	a_data["before_amount"] = fmt.Sprintf("%.4f", pay_conf.Amount)
	a_data["after_amount"] = fmt.Sprintf("%.4f", after_amount)
	a_data["agent_path"] = ""
	a_data["note"] = "上游支付下发"
	in_sql := common.InsertSql(table_name, a_data)
	sql_arr := []string{update_sql, mer_sql, in_sql}
	err := model.Trans(sql_arr)
	if err != nil {
		t_msg = "下发订单处理失败"
		return t_status, t_msg
	}
	t_status = 200
	t_msg = "success"
	//新增账变记录
	return t_status, t_msg
}

func updateReturnCash(order_number, pay_order, note string) (int, string) {
	t_status := 100
	t_msg := "信息异常"
	c_list := model.CashByWeb(order_number)
	if len(c_list.Id) < 1 {
		return t_status, t_msg
	}
	//查询额度
	m_info := model.MerInfo(c_list.Mer_code)
	if m_info.Id < 1 {
		return t_status, t_msg
	}

	if c_list.Status != 1 {
		t_msg = "订单状态异常"
		return t_status, t_msg
	}

	//更新语句
	update_sql := fmt.Sprintf("update cash_list set `status`=9,`note`='取消' where id='%s';", c_list.Id)
	//更新用户额度
	mer_sql := fmt.Sprintf("update mer_list set amount=amount+%.2f,total_out=total_out-%.2f where id='%d';", c_list.Order_amount, c_list.Order_amount, m_info.Id)
	//新增账变记录
	after_amount_no_rate := m_info.Amount + c_list.Amount
	after_amount := m_info.Amount + c_list.Order_amount
	table_name := "amount_list"

	//下发失败返回 4.
	a_data := map[string]string{}
	a_data["id"] = model.GetKey(20)
	a_data["order_number"] = c_list.Id
	a_data["create_time"] = c_list.Create_time
	a_data["amount_type"] = "4"
	a_data["amount"] = fmt.Sprintf("%.4f", c_list.Amount)
	a_data["pay_code"] = c_list.Pay_code
	a_data["mer_code"] = c_list.Mer_code
	a_data["pay_id"] = fmt.Sprintf("%d", c_list.Pay_id)
	a_data["before_amount"] = fmt.Sprintf("%.4f", m_info.Amount)
	a_data["after_amount"] = fmt.Sprintf("%.4f", after_amount_no_rate)
	a_data["agent_path"] = c_list.Agent_path
	a_data["note"] = "下发失败返回"
	in_sql := common.InsertSql(table_name, a_data)

	//下发手续费 7.
	r_data := map[string]string{}
	r_data["id"] = model.GetKey(20)
	r_data["order_number"] = c_list.Id
	r_data["create_time"] = c_list.Create_time
	r_data["amount_type"] = "7"
	r_data["amount"] = fmt.Sprintf("%.4f", c_list.Fee_amount)
	r_data["pay_code"] = c_list.Pay_code
	r_data["mer_code"] = c_list.Mer_code
	r_data["pay_id"] = fmt.Sprintf("%d", c_list.Pay_id)
	r_data["before_amount"] = fmt.Sprintf("%.4f", after_amount_no_rate)
	r_data["after_amount"] = fmt.Sprintf("%.4f", after_amount)
	r_data["agent_path"] = c_list.Agent_path
	r_data["note"] = "下发失败手续费返回"
	rate_in_sql := common.InsertSql(table_name, r_data)

	sql_arr := []string{update_sql, mer_sql, in_sql, rate_in_sql}
	err := model.Trans(sql_arr)
	if err != nil {
		t_msg = "下发订单处理失败"
		return t_status, t_msg
	}
	t_status = 200
	t_msg = "success"
	//新增账变记录
	return t_status, t_msg
}

func updateReturnOrder(order_id, pay_order, note string) (int, string) {
	t_status := 100
	t_msg := "信息异常"
	o_list := model.CashOrderById(order_id)
	if len(o_list.Id) < 1 {
		return t_status, t_msg
	}
	pay_id := fmt.Sprintf("%d", o_list.Pay_id)
	//查询额度
	m_info := model.PayConf(pay_id)
	if m_info.Id < 1 {
		return t_status, t_msg
	}

	if o_list.Status != 1 {
		t_msg = "订单状态异常"
		return t_status, t_msg
	}

	//更新语句
	update_sql := fmt.Sprintf("update order_list set `status`=9 where id='%s';", o_list.Id)
	//更新用户额度
	mer_sql := fmt.Sprintf("update pay_config set amount=amount+%.2f,total_out=total_out-%.2f where id='%s';", o_list.Order_amount, o_list.Order_amount, pay_id)
	//新增账变记录
	after_amount := m_info.Amount + o_list.Order_amount
	table_name := "amount_list"
	a_data := map[string]string{}
	a_data["id"] = model.GetKey(20)
	a_data["order_number"] = o_list.Id
	a_data["create_time"] = o_list.Create_time
	a_data["amount_type"] = "10"
	a_data["amount"] = fmt.Sprintf("%.4f", o_list.Order_amount)
	a_data["pay_code"] = o_list.Pay_code
	a_data["mer_code"] = "all"
	a_data["pay_id"] = pay_id
	a_data["before_amount"] = fmt.Sprintf("%.4f", m_info.Amount)
	a_data["after_amount"] = fmt.Sprintf("%.4f", after_amount)
	a_data["agent_path"] = ""
	a_data["note"] = "上游渠道下发失败返回"
	in_sql := common.InsertSql(table_name, a_data)
	sql_arr := []string{update_sql, mer_sql, in_sql}
	err := model.Trans(sql_arr)
	if err != nil {
		t_msg = "下发订单处理失败"
		return t_status, t_msg
	}
	t_status = 200
	t_msg = "success"
	//新增账变记录
	return t_status, t_msg
}

/**
*  创建代付订单
*  需要将一笔订单拆分成几笔
 */
func dfOrder(order_number, pay_order, note string) (int, string) {
	t_status := 100
	t_msg := "信息异常"
	c_list := model.CashByWeb(order_number)
	if len(c_list.Id) < 1 {
		return t_status, t_msg
	}
	if c_list.Status != -1 {
		t_msg = "订单状态异常"
		return t_status, t_msg
	}
	//查询额度
	m_info := model.MerInfo(c_list.Mer_code)
	if m_info.Id < 1 {
		return t_status, t_msg
	}

	if m_info.Status != 1 {
		t_msg = "商户已被锁定"
		return t_status, t_msg
	}

	order_amount, real_amount := ThreadEveDraw(c_list.Pay_code, c_list.Amount)
	if m_info.Amount < order_amount {
		t_msg = "额度不足"
		return t_status, t_msg
	}
	pay_id := fmt.Sprintf("%d", c_list.Pay_id)
	pay_conf := model.PayConf(pay_id)
	if pay_conf.Id < 1 {
		t_msg = "渠道异常"
		return t_status, t_msg
	}

	if pay_conf.Amount < real_amount {
		t_msg = "额度不足"
		return t_status, t_msg
	}
	fee_amount := order_amount - real_amount
	//扣除渠道额度
	conf_up_sql := fmt.Sprintf("update pay_config set amount=amount-%.2f,total_out=total_out+%.2f where id=%d;", order_amount, order_amount, pay_conf.Id)

	//新增支付记录
	c_data := map[string]string{}
	c_id := model.GetKey(20)
	c_data["id"] = c_id
	c_data["status"] = "1"
	c_data["bank_code"] = c_list.Bank_code
	c_data["bank_title"] = c_list.Bank_title
	c_data["card_number"] = c_list.Card_number
	c_data["order_number"] = c_list.Order_number
	c_data["card_name"] = c_list.Card_name
	c_data["pay_id"] = pay_id
	c_data["amount"] = fmt.Sprintf("%.2f", real_amount)
	c_data["order_amount"] = fmt.Sprintf("%.2f", order_amount)
	c_data["real_amount"] = fmt.Sprintf("%.2f", real_amount)
	c_data["fee_amount"] = fmt.Sprintf("%.2f", fee_amount)
	c_data["create_time"] = c_list.Create_time
	c_data["cash_id"] = c_list.Id
	c_data["branch"] = c_list.Branch
	c_data["phone"] = c_list.Phone
	c_data["mer_code"] = c_list.Mer_code
	c_data["pay_code"] = pay_conf.Pay_code
	c_data["order_type"] = "2"
	order_table := "order_list"
	c_sql := common.InsertSql(order_table, c_data)
	//新增账变记录
	//上游支付下发
	conf_after := pay_conf.Amount - order_amount
	amount_table := "amount_list"
	o_data := map[string]string{}
	o_data["id"] = model.GetKey(20)
	o_data["order_number"] = c_id
	o_data["create_time"] = c_list.Create_time
	o_data["amount_type"] = "9"
	o_data["amount"] = fmt.Sprintf("%.4f", -1*order_amount)
	o_data["pay_code"] = c_list.Pay_code
	o_data["mer_code"] = "all"
	o_data["pay_id"] = pay_id
	o_data["before_amount"] = fmt.Sprintf("%.4f", pay_conf.Amount)
	o_data["after_amount"] = fmt.Sprintf("%.4f", conf_after)
	o_data["agent_path"] = ""
	o_data["note"] = "上游支付下发"
	amount_sql := common.InsertSql(amount_table, o_data)

	//更新语句
	update_sql := fmt.Sprintf("update cash_list set `status`=1,real_amount=%.2f,fee_amount=%.2f,order_amount=%.2f where id='%s';", real_amount, fee_amount, order_amount, c_list.Id)
	//更新用户额度
	mer_sql := fmt.Sprintf("update mer_list set amount=amount-%.2f,total_out=total_out+%.2f where id='%d';", order_amount, order_amount, m_info.Id)
	//新增账变记录
	//下发不扣手续费后 余额
	after_amount_no_rate := m_info.Amount - real_amount
	//下发扣除手续费后 余额
	after_amount := m_info.Amount - order_amount
	//下发 2.
	a_data := map[string]string{}
	a_data["id"] = model.GetKey(20)
	a_data["order_number"] = c_list.Id
	a_data["create_time"] = c_list.Create_time
	a_data["amount_type"] = "2"
	a_data["amount"] = fmt.Sprintf("%.4f", -1*real_amount)
	a_data["pay_code"] = c_list.Pay_code
	a_data["mer_code"] = c_list.Mer_code
	a_data["pay_id"] = pay_id
	a_data["before_amount"] = fmt.Sprintf("%.4f", m_info.Amount)
	a_data["after_amount"] = fmt.Sprintf("%.4f", after_amount_no_rate)
	a_data["agent_path"] = c_list.Agent_path
	a_data["note"] = "下发"
	in_sql := common.InsertSql(amount_table, a_data)
	//下发手续费 7.
	r_data := map[string]string{}
	r_data["id"] = model.GetKey(20)
	r_data["order_number"] = c_list.Id
	r_data["create_time"] = c_list.Create_time
	r_data["amount_type"] = "7"
	r_data["amount"] = fmt.Sprintf("%.4f", -1*fee_amount)
	r_data["pay_code"] = c_list.Pay_code
	r_data["mer_code"] = c_list.Mer_code
	r_data["pay_id"] = pay_id
	r_data["before_amount"] = fmt.Sprintf("%.4f", after_amount_no_rate)
	r_data["after_amount"] = fmt.Sprintf("%.4f", after_amount)
	r_data["agent_path"] = c_list.Agent_path
	r_data["note"] = "下发手续费"
	rate_in_sql := common.InsertSql(amount_table, r_data)
	sql_arr := []string{update_sql, mer_sql, in_sql, amount_sql, c_sql, conf_up_sql, rate_in_sql}
	err := model.Trans(sql_arr)
	if err != nil {
		t_msg = "下发订单处理失败"
		return t_status, t_msg
	}
	t_status = 200
	t_msg = c_id
	//新增账变记录
	return t_status, t_msg
}
