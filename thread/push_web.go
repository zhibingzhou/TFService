package thread

import (
	"TFService/model"
	"time"
)

func ThreadPush() {
	PushWeb()
	//每1分钟执行一次
	bc_timer := time.NewTicker(time.Duration(1) * time.Minute)
	for {
		select {
		case <-bc_timer.C:
			PushWeb()
		}
	}
}

func PushWeb() {
	//每次100条数据，半个小时之前的数据
	dt, _ := time.ParseDuration("-30m")
	now_time := time.Now().Add(dt).Format(format_date)
	push_res := model.PushPayList(now_time)
	if len(push_res) > 0 {
		for _, p_list := range push_res {
			Push(p_list)
		}
	}
}

func ThreadPushCash() {
	PushWebCash()
	//每1分钟执行一次
	bc_timer := time.NewTicker(time.Duration(1) * time.Minute)
	for {
		select {
		case <-bc_timer.C:
			PushWebCash()
		}
	}
}

func PushWebCash() {
	//每次100条数据，半个小时之前的数据
	dt, _ := time.ParseDuration("-30m")
	now_time := time.Now().Add(dt).Format(format_date)
	push_res := model.PushCashList(now_time)
	if len(push_res) > 0 {
		for _, c_list := range push_res {
			PushCash(c_list)
		}
	}
}

func ThreadNotPush() {
	table_name := "pay_channel"
	p_w := map[string]interface{}{}
	p_w["is_push"] = 0
	fields := []string{"code"}
	pay_list, _ := model.PageList(table_name, "", 1000, 0, fields, p_w)
	if len(pay_list) < 1 {
		return
	}
	for _, p_val := range pay_list {
		go NotPushFor(p_val["code"])
	}

}

func NotPushFor(pay_code string) {
	NotPush(pay_code)
	//每1分钟执行一次
	bc_timer := time.NewTicker(time.Duration(2) * time.Minute)
	for {
		select {
		case <-bc_timer.C:
			NotPush(pay_code)
		}
	}
}

func NotPush(pay_code string) {
	table_name := "cash_list"
	//判断是否有该渠道的下发未完成的订单
	p_where := map[string]interface{}{}
	p_where["pay_code"] = pay_code
	p_where["status"] = 1
	field := []string{"id", "mer_code", "pay_id"}
	count_field := "count(id) as num"
	total, _ := model.ListTotal(table_name, count_field, p_where)
	if total < 1 {
		return
	}
	page_size := 100
	pages := total / page_size
	if total%page_size != 0 {
		pages = pages + 1
	}
	offset := 0
	for i := 0; i < pages; i++ {
		offset = i * page_size
		cash_list, _ := model.PageList(table_name, "", page_size, offset, field, p_where)
		if len(cash_list) < 1 {
			continue
		}
		for _, c_val := range cash_list {
			PayForQuery(c_val)
		}
	}
}
