package thread

import (
	"fmt"
	"TFService/model"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/zhibingzhou/go_public/common"
)

/**
* 支付订单下载
 */
func DownPayList(pay_status, order_number, web_order, pay_code, class_code, is_mobile, start_time, end_time, mer_code, is_agent string, ctx *gin.Context) (int, string, string) {
	t_status := 100
	t_msg := "管理员信息错误"
	file_name := "./download/pay_" + strconv.FormatInt(time.Now().Unix(), 10) + ".xlsx"
	session := sessions.Default(ctx)
	admin_mer := fmt.Sprintf("%v", session.Get("mer_code"))
	p_where := map[string]interface{}{}
	like_sql := ""
	if mer_code != "" {
		p_where["mer_code"] = mer_code
	}
	if pay_status != "" {
		p_where["status"] = pay_status
	}
	if order_number != "" {
		p_where["id"] = order_number
	}
	if web_order != "" {
		p_where["order_number"] = web_order
	}
	if pay_code != "" {
		p_where["pay_code"] = pay_code
	}
	if class_code != "" {
		p_where["class_code"] = class_code
	}
	if is_mobile != "" {
		p_where["is_mobile"] = is_mobile
	}
	if admin_mer != "all" {
		//所有下线
		if is_agent == "2" {
			admin := model.MerInfoRedis(admin_mer)
			like_sql = " and agent_path like '" + admin["agent_path"] + admin["code"] + "_%'"
		} else if is_agent == "1" {
			//直属下线
			admin := model.MerInfoRedis(admin_mer)
			p_where["agent_path"] = admin["agent_path"] + admin["code"] + "_"
		} else {
			//自己
			p_where["mer_code"] = admin_mer
		}
	}

	if start_time == "" {
		start_time = time.Now().Format(format_day) + " 00:00:00"
	} else {
		start_time = start_time + " 00:00:00"
	}

	if end_time == "" {
		end_time = time.Now().Format(format_day) + " 23:59:59"
	} else {
		end_time = end_time + " 23:59:59"
	}
	table_name := "pay_list"
	count_field := "count(0) as num"
	date_field := "create_time"
	total, _ := model.DateListTotal(table_name, date_field, start_time, end_time, like_sql, count_field, p_where)
	if total < 1 {
		return t_status, t_msg, file_name
	}
	field_arr := []string{"mer_code", "id", "order_number", "amount", "real_amount", "create_time", "pay_time", "class_code", "status", "is_mobile", "pay_code"}
	title_arr := []string{"商户号", "系统订单号", "商户订单号", "订单金额", "实际到账金额", "订单创建时间", "订单完成时间", "支付类型", "订单状态", "网站类型", "支付名称"}
	title_arr_1 := [][]string{}
	title_arr_1 = append(title_arr_1[0:], title_arr)
	common.WriteXlsx(file_name, title_arr_1)
	page_size := 200
	pages := total / page_size
	if total%page_size != 0 {
		pages = pages + 1
	}

	offset := 0
	for i := 0; i < pages; i++ {
		offset = i * page_size
		pay_list, _ := model.PageDateList(table_name, date_field, start_time, end_time, like_sql, page_size, offset, field_arr, p_where)
		if len(pay_list) < 1 {
			return t_status, t_msg, file_name
		}
		file_content := [][]string{}
		for _, p_v := range pay_list {
			content_arr := []string{}
			for _, t_val := range field_arr {
				if t_val == "pay_code" {
					p_info := model.PayInfoRedis(p_v["pay_code"])
					content_arr = append(content_arr[0:], p_info["title"])
				} else if t_val == "status" {
					v_val := "处理中"
					if p_v["status"] == "3" {
						v_val = "成功"
					}
					content_arr = append(content_arr[0:], v_val)
				} else if t_val == "class_code" {
					c_info := model.ClassInfoRedis(p_v["class_code"])
					content_arr = append(content_arr[0:], c_info["title"])
				} else {
					content_arr = append(content_arr[0:], p_v[t_val])
				}
			}
			file_content = append(file_content[0:], content_arr)
		}
		common.AppendWriteXlsx(file_name, file_content)
	}

	t_status = 200
	t_msg = "success"
	return t_status, t_msg, file_name
}

/**
* 代付订单下载
 */
func DownCashList(cash_status, order_number, web_order, pay_code, start_time, end_time, mer_code, is_agent string, ctx *gin.Context) (int, string, string) {
	t_status := 100
	t_msg := "管理员信息错误"
	file_name := "./download/cash_" + strconv.FormatInt(time.Now().Unix(), 10) + ".xlsx"
	session := sessions.Default(ctx)
	admin_mer := fmt.Sprintf("%v", session.Get("mer_code"))
	p_where := map[string]interface{}{}
	like_sql := ""
	if mer_code != "" {
		p_where["mer_code"] = mer_code
	}
	if cash_status != "" {
		p_where["status"] = cash_status
	}
	if order_number != "" {
		p_where["id"] = order_number
	}
	if web_order != "" {
		p_where["order_number"] = web_order
	}
	if pay_code != "" {
		p_where["pay_code"] = pay_code
	}

	if admin_mer != "all" {
		//所有下线
		if is_agent == "2" {
			admin := model.MerInfoRedis(admin_mer)
			like_sql = " and agent_path like '" + admin["agent_path"] + admin["code"] + "_%'"
		} else if is_agent == "1" {
			//直属下线
			admin := model.MerInfoRedis(admin_mer)
			p_where["agent_path"] = admin["agent_path"] + admin["code"] + "_"
		} else {
			//自己
			p_where["mer_code"] = admin_mer
		}
	}

	if start_time == "" {
		start_time = time.Now().Format(format_day) + " 00:00:00"
	} else {
		start_time = start_time + " 00:00:00"
	}

	if end_time == "" {
		end_time = time.Now().Format(format_day) + " 23:59:59"
	} else {
		end_time = end_time + " 23:59:59"
	}
	table_name := "cash_list"
	count_field := "count(0) as num"
	date_field := "create_time"
	total, _ := model.DateListTotal(table_name, date_field, start_time, end_time, like_sql, count_field, p_where)
	if total < 1 {
		return t_status, t_msg, file_name
	}

	field_arr := []string{"mer_code", "id", "order_number", "amount", "real_amount", "create_time", "pay_time", "fee_amount", "status", "pay_code", "bank_title", "card_name", "card_number", "note"}
	title_arr := []string{"商户号", "系统订单号", "商户订单号", "订单金额", "实际到账金额", "订单创建时间", "订单完成时间", "手续费", "订单状态", "支付名称", "银行名称", "持卡人姓名", "卡号", "备注"}
	title_arr_1 := [][]string{}
	title_arr_1 = append(title_arr_1[0:], title_arr)
	common.WriteXlsx(file_name, title_arr_1)
	page_size := 200
	pages := total / page_size
	if total%page_size != 0 {
		pages = pages + 1
	}

	offset := 0
	for i := 0; i < pages; i++ {
		offset = i * page_size
		pay_list, _ := model.PageDateList(table_name, date_field, start_time, end_time, like_sql, page_size, offset, field_arr, p_where)
		if len(pay_list) < 1 {
			return t_status, t_msg, file_name
		}
		file_content := [][]string{}
		for _, p_v := range pay_list {
			content_arr := []string{}
			for _, t_val := range field_arr {
				if t_val == "pay_code" {
					p_info := model.PayInfoRedis(p_v["pay_code"])
					content_arr = append(content_arr[0:], p_info["title"])
				} else if t_val == "status" {
					v_val := "处理中"
					if p_v["status"] == "3" {
						v_val = "成功"
					} else if p_v["status"] == "9" {
						v_val = "失败"
					}
					content_arr = append(content_arr[0:], v_val)
				} else {
					content_arr = append(content_arr[0:], p_v[t_val])
				}
			}
			file_content = append(file_content[0:], content_arr)
		}
		common.AppendWriteXlsx(file_name, file_content)
	}

	t_status = 200
	t_msg = "success"
	return t_status, t_msg, file_name
}

/**
* 账变记录下载
 */
func DownAmountList(amount_type, order_number, pay_code, start_time, end_time, mer_code, is_agent string, ctx *gin.Context) (int, string, string) {
	t_status := 100
	t_msg := "管理员信息错误"
	file_name := "./download/amount_" + strconv.FormatInt(time.Now().Unix(), 10) + ".xlsx"
	session := sessions.Default(ctx)
	admin_mer := fmt.Sprintf("%v", session.Get("mer_code"))
	p_where := map[string]interface{}{}
	like_sql := ""
	if mer_code != "" {
		p_where["mer_code"] = mer_code
	}
	if amount_type != "" {
		p_where["amount_type"] = amount_type
	}
	if order_number != "" {
		p_where["order_number"] = order_number
	}
	if pay_code != "" {
		p_where["pay_code"] = pay_code
	}
	if admin_mer != "all" {
		//所有下线
		if is_agent == "2" {
			admin := model.MerInfoRedis(admin_mer)
			like_sql = " and agent_path like '" + admin["agent_path"] + admin["code"] + "_%'"
		} else if is_agent == "1" {
			//直属下线
			admin := model.MerInfoRedis(admin_mer)
			p_where["agent_path"] = admin["agent_path"] + admin["code"] + "_"
		} else {
			//自己
			p_where["mer_code"] = admin_mer
		}
	}

	if start_time == "" {
		start_time = time.Now().Format(format_day) + " 00:00:00"
	} else {
		start_time = start_time + " 00:00:00"
	}

	if end_time == "" {
		end_time = time.Now().Format(format_day) + " 23:59:59"
	} else {
		end_time = end_time + " 23:59:59"
	}
	table_name := "amount_list"
	count_field := "count(0) as num"
	date_field := "create_time"
	total, _ := model.DateListTotal(table_name, date_field, start_time, end_time, like_sql, count_field, p_where)
	if total < 1 {
		return t_status, t_msg, file_name
	}
	field_arr := []string{"mer_code", "amount_type", "order_number", "amount", "before_amount", "create_time", "after_amount", "note"}
	title_arr := []string{"商户号", "账变记录类型", "系统订单号", "订单金额", "账变前额度", "订单创建时间", "账变后额度"}
	title_arr_1 := [][]string{}
	title_arr_1 = append(title_arr_1[0:], title_arr)
	common.WriteXlsx(file_name, title_arr_1)
	page_size := 200
	pages := total / page_size
	if total%page_size != 0 {
		pages = pages + 1
	}

	offset := 0
	for i := 0; i < pages; i++ {
		offset = i * page_size
		pay_list, _ := model.PageDateList(table_name, date_field, start_time, end_time, like_sql, page_size, offset, field_arr, p_where)
		if len(pay_list) < 1 {
			return t_status, t_msg, file_name
		}
		file_content := [][]string{}
		for _, p_v := range pay_list {
			content_arr := []string{}
			for _, t_val := range field_arr {
				if t_val == "pay_code" {
					p_info := model.PayInfoRedis(p_v["pay_code"])
					content_arr = append(content_arr[0:], p_info["title"])
				} else if t_val == "amount_type" {
					v_val := "支付"
					switch p_v["status"] {
					case "2":
						v_val = "下发"
					case "3":
						v_val = "代理收入"
					case "4":
						v_val = "下发失败返还额度"
					case "5":
						v_val = "调整额度"
					case "6":
						v_val = "支付的手续费"
					case "7":
						v_val = "下发的手续费"
					case "8":
						v_val = "代理佣金返还"
					case "9":
						v_val = "上游支付下发"
					case "10":
						v_val = "上游支付下发失败返回"
					}
					content_arr = append(content_arr[0:], v_val)
				} else {
					content_arr = append(content_arr[0:], p_v[t_val])
				}
			}
			file_content = append(file_content[0:], content_arr)
		}
		common.AppendWriteXlsx(file_name, file_content)
	}

	t_status = 200
	t_msg = "success"
	return t_status, t_msg, file_name
}
