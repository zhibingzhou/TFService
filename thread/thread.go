package thread

import (
	"fmt"
	"TFService/model"
	"time"

	"github.com/zhibingzhou/go_public/common"
	"github.com/zhibingzhou/go_public/redis"
)

func ThreadReport() {
	ThreadReportFor()
	//每60分钟执行一次
	bc_timer := time.NewTicker(time.Duration(60) * time.Minute)
	for {
		select {
		case <-bc_timer.C:
			ThreadReportFor()
		}
	}

}

func ThreadReportFor() {
	fields := []string{"code", "is_agent", "agent_path"}
	p_where := map[string]interface{}{}
	p_where["status"] = 1
	mer_list, _ := model.PageList("mer_list", "", 100, 0, fields, p_where)
	if len(mer_list) < 1 {
		return
	}
	for _, mer_info := range mer_list {
		merReport(mer_info)
	}
}

func merReport(mer_info map[string]string) {
	redis_key := "lock:merReport:" + mer_info["code"]

	lock_res := redis.RediGo.StringWriteNx(redis_key, mer_info["code"], 15)
	if lock_res < 1 {
		return
	}
	defer redis.RediGo.KeyDel(redis_key)

	time_now := time.Now()
	if time_now.Hour() < 3 {
		time_now = time_now.AddDate(0, 0, -1)
	}

	time_date := time_now.Format(format_day)
	start_date := time_date + " 00:00:00"
	end_date := time_date + " 23:59:59"
	in_table := "amount_list"
	date_field := "create_time"
	like_sql := ""
	field := "sum(amount) as total"
	//充值
	p_where := map[string]interface{}{}
	p_where["amount_type"] = 1
	p_where["mer_code"] = mer_info["code"]
	_, total_in := model.DateListTotal(in_table, date_field, start_date, end_date, like_sql, field, p_where)
	//代理佣金
	dp_where := map[string]interface{}{}
	dp_where["amount_type"] = 3
	dp_where["mer_code"] = mer_info["code"]
	_, dtotal_in := model.DateListTotal(in_table, date_field, start_date, end_date, like_sql, field, dp_where)
	//调整额度
	tp_where := map[string]interface{}{}
	tp_where["amount_type"] = 5
	tp_where["mer_code"] = mer_info["code"]
	_, ttotal_in := model.DateListTotal(in_table, date_field, start_date, end_date, like_sql, field, tp_where)
	//实际充值 =  充值 +  代理佣金 +  调整额度
	total_in = total_in + dtotal_in + ttotal_in
	//下发
	out_where := map[string]interface{}{}
	out_where["amount_type"] = 2
	out_where["mer_code"] = mer_info["code"]
	_, total_out := model.DateListTotal(in_table, date_field, start_date, end_date, like_sql, field, out_where)
	//下发失败返还
	fout_where := map[string]interface{}{}
	fout_where["amount_type"] = 4
	fout_where["mer_code"] = mer_info["code"]
	_, ftotal_out := model.DateListTotal(in_table, date_field, start_date, end_date, like_sql, field, fout_where)
	// 下发 + 下发失败返还 = 实际下发
	total_out = total_out + ftotal_out
	//充值手续费
	p_where_rate := map[string]interface{}{}
	p_where_rate["amount_type"] = 6
	p_where_rate["mer_code"] = mer_info["code"]
	_, total_in_rate := model.DateListTotal(in_table, date_field, start_date, end_date, like_sql, field, p_where_rate)
	//下发手续费
	out_where_rate := map[string]interface{}{}
	out_where_rate["amount_type"] = 7
	out_where_rate["mer_code"] = mer_info["code"]
	_, total_out_rate := model.DateListTotal(in_table, date_field, start_date, end_date, like_sql, field, out_where_rate)

	c_w := map[string]interface{}{}
	c_w["report_date"] = start_date
	c_w["mer_code"] = mer_info["code"]
	fields := []string{"id"}
	table_name := "mer_report"
	//查询是否存在
	is_exist, _ := model.CommonFieldsRow(table_name, fields, c_w)
	sql := ""
	if len(is_exist["id"]) > 0 && is_exist["id"] != "<nil>" {
		sql = fmt.Sprintf("update %s set total_in=%.2f,total_out=%.2f,total_in_rate=%.2f,total_out_rate=%.2f where id='%s';", table_name, total_in, total_out, total_in_rate, total_out_rate, is_exist["id"])
	} else {
		m_data := map[string]string{}
		m_data["id"] = model.GetKey(18)
		m_data["mer_code"] = mer_info["code"]
		m_data["is_agent"] = mer_info["is_agent"]
		m_data["agent_path"] = mer_info["agent_path"]
		m_data["report_date"] = start_date
		m_data["total_out"] = fmt.Sprintf("%.2f", total_out)
		m_data["total_in"] = fmt.Sprintf("%.2f", total_in)
		m_data["total_in_rate"] = fmt.Sprintf("%.2f", total_in_rate)
		m_data["total_out_rate"] = fmt.Sprintf("%.2f", total_out_rate)
		sql = common.InsertSql(table_name, m_data)
	}
	model.Query(sql)
	return
}
