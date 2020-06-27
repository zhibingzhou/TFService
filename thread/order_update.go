package thread

import (
	"TFService/model"
	"time"

	"github.com/zhibingzhou/go_public/redis"
)

/**
* 支付订单完成
* 接收参数后，对数据做进一步处理
 */
func finished_pay(pay_order, note string, p_list model.PayList) (int, string) {
	order_type := 1
	t_status, t_msg := updateMerCash(p_list.Id, pay_order, note, order_type)
	if t_status == 200 {
		p_l := model.OrderById(p_list.Id)
		go Push(p_l)
	}

	return t_status, t_msg
}

/*
* 更新支付订单
**/
func ThreadUpdatePay(order_number, pay_order, amount, sign, sign_str string, is_cent int) (int, string) {
	t_status := 100
	t_msg := "订单请求频繁"
	lock_key := "lock:UpdatePay" + order_number
	//写入redis锁(10秒)，来排除重复处理的可能
	order_lock := redis.RediGo.StringWriteNx(lock_key, order_number, 15)
	if order_lock == 0 {
		return t_status, t_msg
	}
	defer redis.RediGo.KeyDel(lock_key)

	p_list := model.OrderById(order_number)
	t_msg = "订单号不存在"
	if len(p_list.Id) < 1 {
		return t_status, t_msg
	}

	if p_list.Status != 1 {
		t_msg = "订单状态错误"
		return t_status, t_msg
	}

	t_status, t_msg = payCallBack(amount, sign, sign_str, is_cent, p_list)
	if t_status != 200 {
		return t_status, t_msg
	}
	note := "自动完成"
	t_status, t_msg = finished_pay(pay_order, note, p_list)
	return t_status, t_msg
}

/*
* 更新支付订单
**/
func updatePayStatus(order_number, pay_status string) (int, string) {
	t_status := 100
	t_msg := "订单请求频繁"
	lock_key := "lock:UpdatePay" + order_number
	//写入redis锁(10秒)，来排除重复处理的可能
	order_lock := redis.RediGo.StringWriteNx(lock_key, order_number, 15)
	if order_lock == 0 {
		return t_status, t_msg
	}
	defer redis.RediGo.KeyDel(lock_key)

	p_list := model.OrderById(order_number)
	t_msg = "订单号不存在"
	if len(p_list.Id) < 1 {
		return t_status, t_msg
	}

	if p_list.Status != 1 {
		t_msg = "订单状态错误"
		return t_status, t_msg
	}
	note := "手动完成"
	if pay_status != "3" {
		u_date := map[string]interface{}{}
		u_date["status"] = 9
		u_date["note"] = "手动完成"
		err := model.UpdatesPayList(p_list, u_date)
		if err != nil {
			t_msg = "订单更新失败"
			return t_status, t_msg
		}
		t_status = 200
		t_msg = "success"
		return t_status, t_msg
	}

	pay_order := p_list.Pay_order
	if pay_order == "" {
		pay_order = p_list.Order_number
	}

	t_status, t_msg = finished_pay(pay_order, note, p_list)
	return t_status, t_msg
}

/**
* 支付订单完成
* 接收参数后，对数据做进一步处理
 */
func finished_cash(pay_order, note string, c_list model.CashList) (int, string) {
	order_type := 3
	t_status, t_msg := updateMerCash(c_list.Order_number, pay_order, note, order_type)
	return t_status, t_msg
}

/**
* 支付订单完成
* 接收参数后，对数据做进一步处理
 */
func finished_order(pay_order, note string, c_list model.OrderList) (int, string) {
	order_type := 5
	t_status, t_msg := updateMerCash(c_list.Id, pay_order, note, order_type)
	return t_status, t_msg
}

/*
* 更新下发订单
**/
func ThreadUpdateOrder(order_id, pay_order, amount, order_status, sign, sign_str string, is_cent int) (int, string) {
	t_status := 100
	t_msg := "订单请求频繁"
	lock_key := "lock:ThreadUpdateOrder:" + order_id
	//写入redis锁(10秒)，来排除重复处理的可能
	order_lock := redis.RediGo.StringWriteNx(lock_key, order_id, 15)
	if order_lock == 0 {
		return t_status, t_msg
	}

	defer redis.RediGo.KeyDel(lock_key)

	o_list := model.CashOrderById(order_id)
	t_msg = "订单号不存在"
	if len(o_list.Id) < 1 {
		return t_status, t_msg
	}

	if o_list.Status != 1 {
		t_msg = "订单状态错误"
		return t_status, t_msg
	}

	t_status, t_msg = cashCallBack(amount, sign, sign_str, is_cent, o_list)
	if t_status != 200 {
		return t_status, t_msg
	}
	note := "自动完成"
	t_status, t_msg = updateOrderStatus(order_id, pay_order, order_status, note)
	if t_status != 200 {
		return t_status, t_msg
	}
	//判断是否纯代下发
	if o_list.Order_type == 1 {
		return t_status, t_msg
	}
	t_status, t_msg = updateCashStatus(o_list.Cash_id, pay_order, order_status, note)
	return t_status, t_msg
}

/*
* 更新下发订单
**/
func updateCashStatus(order_number, pay_order, cash_status, note string) (int, string) {
	t_status := 100
	t_msg := "订单请求频繁"
	lock_key := "lock:UpdateCash:" + order_number
	pay_time := time.Now().Format(format_date)
	//写入redis锁(10秒)，来排除重复处理的可能
	order_lock := redis.RediGo.StringWriteNx(lock_key, order_number, 15)
	if order_lock == 0 {
		return t_status, t_msg
	}
	defer redis.RediGo.KeyDel(lock_key)

	c_list := model.CashById(order_number)
	t_msg = "订单号不存在"
	if len(c_list.Id) < 1 {
		return t_status, t_msg
	}

	if c_list.Status != 1 {
		t_msg = "订单状态错误"
		return t_status, t_msg
	}
	if cash_status == "3" {
		u_date := map[string]interface{}{}
		u_date["status"] = 3
		u_date["note"] = note
		u_date["pay_order"] = pay_order
		u_date["pay_time"] = pay_time
		err := model.UpdatesCashList(c_list, u_date)
		if err != nil {
			t_msg = "订单更新失败"
			return t_status, t_msg
		}
		t_status = 200
		t_msg = "success"
		c_list.Status = 3
	} else {
		//下发失败，需要返还额度
		t_status, t_msg = finished_cash(pay_order, note, c_list)
		c_list.Status = 9
	}
	if t_status == 200 {
		go PushCash(c_list)
	}
	return t_status, t_msg
}

/*
* 手动更新下发订单为进行中状态
**/
func updateCashStatusForRun(order_number, pay_order, cash_status, note string) (int, string) {
	t_status := 100
	t_msg := "订单请求频繁"
	lock_key := "lock:UpdateCashRun:" + order_number
	//写入redis锁(10秒)，来排除重复处理的可能
	order_lock := redis.RediGo.StringWriteNx(lock_key, order_number, 15)
	if order_lock == 0 {
		return t_status, t_msg
	}
	defer redis.RediGo.KeyDel(lock_key)

	c_list := model.CashById(order_number)
	t_msg = "订单号不存在"
	if len(c_list.Id) < 1 {
		return t_status, t_msg
	}
	//订单状态为成功，才能改为进行中
	if c_list.Status != 3 {
		t_msg = "订单状态错误"
		return t_status, t_msg
	}

	u_date := map[string]interface{}{}
	u_date["status"] = 1
	u_date["note"] = note
	u_date["pay_order"] = pay_order
	err := model.UpdatesCashList(c_list, u_date)
	if err != nil {
		t_msg = "订单更新失败"
		return t_status, t_msg
	}
	t_status = 200
	t_msg = "success"

	return t_status, t_msg
}

/*
* 更新下发订单
**/
func updateOrderStatus(order_id, pay_order, order_status, note string) (int, string) {
	t_status := 100
	t_msg := "订单请求频繁"
	lock_key := "lock:updateOrderStatus:" + order_id
	//写入redis锁(10秒)，来排除重复处理的可能
	order_lock := redis.RediGo.StringWriteNx(lock_key, order_id, 15)
	if order_lock == 0 {
		return t_status, t_msg
	}
	defer redis.RediGo.KeyDel(lock_key)

	p_list := model.CashOrderById(order_id)
	t_msg = "订单号不存在"
	if len(p_list.Id) < 1 {
		return t_status, t_msg
	}

	if p_list.Status != 1 {
		t_msg = "订单状态错误"
		return t_status, t_msg
	}
	if order_status == "3" {
		u_date := map[string]interface{}{}
		u_date["status"] = 3
		u_date["note"] = note
		u_date["pay_order"] = pay_order
		err := model.UpdatesOrderList(p_list, u_date)
		if err != nil {
			t_msg = "订单更新失败"
			return t_status, t_msg
		}
		t_status = 200
		t_msg = "success"
		return t_status, t_msg
	}

	//下发失败，需要返还额度
	t_status, t_msg = finished_order(pay_order, note, p_list)
	return t_status, t_msg
}

/*
* 手动更新渠道下发订单为进行中状态
**/
func updateOrderStatusForRun(order_id, pay_order, order_status, note string) (int, string) {
	t_status := 100
	t_msg := "订单请求频繁"
	lock_key := "lock:updateOrderStatus:" + order_id
	//写入redis锁(10秒)，来排除重复处理的可能
	order_lock := redis.RediGo.StringWriteNx(lock_key, order_id, 15)
	if order_lock == 0 {
		return t_status, t_msg
	}
	defer redis.RediGo.KeyDel(lock_key)

	p_list := model.CashOrderById(order_id)
	t_msg = "订单号不存在"
	if len(p_list.Id) < 1 {
		return t_status, t_msg
	}

	//订单状态为成功，才能改为进行中
	if p_list.Status != 3 {
		t_msg = "订单状态错误"
		return t_status, t_msg
	}

	u_date := map[string]interface{}{}
	u_date["status"] = 1
	u_date["note"] = note
	u_date["pay_order"] = pay_order
	err := model.UpdatesOrderList(p_list, u_date)
	if err != nil {
		t_msg = "订单更新失败"
		return t_status, t_msg
	}
	t_status = 200
	t_msg = "success"
	return t_status, t_msg
}

func GetBackConfig(order_number string, order_pay bool) (string, string) {
	msg := "success"
	//获取 pay_config 参数
	p_where := map[string]interface{}{}

	if order_pay == true { //支付
		p_list := model.OrderById(order_number)
		if len(p_list.Id) < 1 {
			msg = "未找到此订单"
			return "", msg
		}
		p_where["id"] = p_list.Pay_id
	} else { //代付
		o_list := model.CashOrderById(order_number)
		if len(o_list.Id) < 1 {
			msg = "订单号不存在"
			return "", msg
		}
		p_where["id"] = o_list.Pay_id
	}

	table_name := "pay_config"
	field := []string{"api_conf"}
	api_config, _ := model.CommonFieldsRow(table_name, field, p_where)
	result := api_config["api_conf"]
	return result, msg
}
