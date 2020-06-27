package thread

import (
	"TFService/model"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/zhibingzhou/go_public/common"

	"github.com/zhibingzhou/go_public/redis"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

/**
*  更新管理员权限
 */
func UpdatePower(account, power_path, secret string, ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "权限错误"
	session := sessions.Default(ctx)
	admin_power := fmt.Sprintf("%v", session.Get("power_path"))
	admin := fmt.Sprintf("%v", session.Get("account"))

	power_arr := strings.Split(power_path, ",")
	if len(power_arr) < 1 {
		return t_status, t_msg
	}

	user := model.AdminInfoRedis(account)
	if len(user["account"]) < 1 {
		t_msg = "账号错误"
		return t_status, t_msg
	}
	mer_code := fmt.Sprintf("%v", session.Get("mer_code"))
	if mer_code != "all" && mer_code != user["mer_code"] {
		t_msg = "不能修改该账号的权限"
		return t_status, t_msg
	}
	admin_info := model.AdminInfoRedis(admin)

	t_status, t_msg = AuthGoogle(secret, admin_info["secret"])
	if t_status != 200 {
		t_msg = "谷歌动态验证码错误"
		return t_status, t_msg
	}

	admin_arr := []string{}
	table_name := "admin_power"
	power_map := []map[string]string{}
	if admin_power != "all" {
		page_size := 1000
		offset := 0
		fields := []string{"power_code"}
		p_where := map[string]interface{}{}
		p_where["account"] = admin
		a_power, _ := model.PageList(table_name, "", page_size, offset, fields, p_where)
		if len(a_power) < 1 {
			t_msg = "权限不足"
			return t_status, t_msg
		}
		for _, p_val := range a_power {
			admin_arr = append(admin_arr[0:], p_val["power_code"])
		}
		for _, power_val := range power_arr {
			if !common.Arr_In(admin_arr, power_val) {
				t_msg = "权限错误"
				return t_status, t_msg
			}
		}
	}

	for _, power_val := range power_arr {
		p_map := map[string]string{}
		p_map["power_code"] = power_val
		p_map["account"] = account
		power_map = append(power_map[0:], p_map)
	}

	del_sql := fmt.Sprintf("delete from admin_power where account='%s';", account)
	up_sql := fmt.Sprintf("update admin_list set power_path='admin' where account='%s';", account)

	in_sql := common.BatchInsertSql(table_name, power_map)
	sql_arr := []string{del_sql, in_sql, up_sql}
	err := model.Trans(sql_arr)
	if err != nil {
		t_msg = "权限修改失败"
		return t_status, t_msg
	}

	t_status = 200
	t_msg = "success"
	return t_status, t_msg
}

/**
*  锁定用户
 */
func EditAdmin(account, is_edit string, ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "请填写完整"
	session := sessions.Default(ctx)
	mer_code := fmt.Sprintf("%v", session.Get("mer_code"))

	if account == "" || is_edit == "" {
		return t_status, t_msg
	}

	user := model.AdminInfo(account)
	if user.Account == "" {
		t_msg = "用户信息错误"
		return t_status, t_msg
	}

	if mer_code != "all" && mer_code != user.Mer_code {
		t_msg = "该账号没有权限"
		return t_status, t_msg
	}

	c_data := map[string]interface{}{}
	if is_edit == "1" {
		c_data["pwd"] = common.HexMd5("123456")
	} else if is_edit == "2" {
		c_data["status"] = 0
	} else if is_edit == "3" {
		c_data["status"] = 1
	} else if is_edit == "4" {
		c_data["secret"] = ""
	}

	err := model.UpdatesAdminList(user, c_data)
	if err != nil {
		t_msg = "操作失败"
		return t_status, t_msg
	}
	redis_key := "admin_list:" + user.Account
	redis.RediGo.KeyDel(redis_key)
	redis_key = "admin_list:session_id:" + user.Session_id
	redis.RediGo.KeyDel(redis_key)
	t_status = 200
	t_msg = "success"
	return t_status, t_msg
}

/**
*  更新管理员密码
 */
func UpdatePwd(old_pwd, new_pwd, con_pwd, secret string, ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "用户信息错误"
	session := sessions.Default(ctx)
	sess_account := session.Get("account")
	if sess_account == nil {
		return t_status, t_msg
	}
	account := fmt.Sprintf("%v", sess_account)

	if new_pwd != con_pwd {
		t_msg = "两次密码不一致"
		return t_status, t_msg
	}
	if new_pwd == "" {
		t_msg = "新密码不能为空"
		return t_status, t_msg
	}
	if old_pwd == "" {
		t_msg = "旧密码不能为空"
		return t_status, t_msg
	}
	a_list := model.AdminInfo(account)
	if a_list.Id < 1 {
		t_msg = "管理员信息异常"
		return t_status, t_msg
	}
	if a_list.Pwd != common.HexMd5(old_pwd) {
		t_msg = "旧密码错误"
		return t_status, t_msg
	}

	t_status, t_msg = AuthGoogle(secret, a_list.Secret)
	if t_status != 200 {
		t_msg = "谷歌动态验证码错误"
		return t_status, t_msg
	}

	c_data := map[string]interface{}{}
	c_data["pwd"] = common.HexMd5(new_pwd)
	err := model.UpdatesAdminList(a_list, c_data)
	if err != nil {
		t_msg = "密码更新失败"
		return t_status, t_msg
	}
	t_status = 200
	t_msg = "success"
	return t_status, t_msg
}

/**
*  管理员列表
 */
func AdminList(account, mer_code, admin_status, page, page_size string, ctx *gin.Context) (int, int, string, []map[string]string) {
	t_status := 200
	t_msg := "success"
	total := 0
	admin_list := []map[string]string{}
	table_name := "admin_list"
	c_w := map[string]interface{}{}
	session := sessions.Default(ctx)
	sess_mer := fmt.Sprintf("%v", session.Get("mer_code"))
	if account != "" {
		c_w["account"] = account
	}

	if mer_code != sess_mer && sess_mer != "all" {
		c_w["mer_code"] = sess_mer
	} else if sess_mer == "all" && mer_code != "" {
		c_w["mer_code"] = mer_code
	}

	count_field := "count(0) as num"

	if admin_status != "" {
		status_int, _ := strconv.Atoi(admin_status)
		c_w["status"] = status_int
	}
	total, _ = model.ListTotal(table_name, count_field, c_w)
	if total < 1 {
		return t_status, total, t_msg, admin_list
	}

	page_int, size_int := ThreadPage(page, page_size)

	offset := (page_int - 1) * size_int
	fields := []string{"account", "status", "mer_code", "login_time", "login_ip"}
	admin_list, _ = model.PageList(table_name, "", size_int, offset, fields, c_w)
	if len(admin_list) < 1 {
		return t_status, total, t_msg, admin_list
	}
	for a_k, a_val := range admin_list {
		mer_info := model.MerInfoRedis(a_val["mer_code"])
		admin_list[a_k]["mer_title"] = mer_info["title"]
	}
	return t_status, total, t_msg, admin_list
}

/**
* 获取用户权限
 */
func AdminPower(account string) (int, string, []map[string]interface{}) {
	t_status := 100
	t_msg := "管理员账户信息错误"
	power_list := []map[string]interface{}{}

	return t_status, t_msg, power_list
}

/**
*  解除管理员的谷歌验证
 */
func DelBind(account string, ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "没有权限"
	session := sessions.Default(ctx)
	mer_code := fmt.Sprintf("%v", session.Get("mer_code"))
	if mer_code != "all" {
		return t_status, t_msg
	}

	a_list := model.AdminInfoRedis(account)
	if len(a_list["account"]) < 1 {
		t_msg = "用户信息异常"
		return t_status, t_msg
	}

	up_sql := fmt.Sprintf("update admin_list set secret='' where id='%s';", a_list["id"])
	err := model.Query(up_sql)
	if err != nil {
		t_msg = "解绑失败"
		return t_status, t_msg
	}
	redis_key := "admin_list:" + a_list["account"]
	redis.RediGo.KeyDel(redis_key)
	redis_key = "admin_list:session_id:" + a_list["session_id"]
	redis.RediGo.KeyDel(redis_key)
	t_status = 200
	t_msg = "success"

	return t_status, t_msg
}

/**
*  清除缓存
 */
func DelCache(cache_status string) (int, string) {
	t_status := 200
	t_msg := "success"
	key_key := Data_Redis_Key
	if cache_status == "config" {
		key_key = Conf_Redis_Key
	}
	key_count := redis.RediGo.Scard(key_key)
	for i := 0; i <= key_count; i++ {
		red_key := redis.RediGo.Spop(key_key)
		if k := redis.RediGo.KeyDel(red_key); k < 1 {
			redis.RediGo.Sadd(key_key, red_key, 0)
		}
	}
	return t_status, t_msg
}

/**
*  新增管理员
 */
func AddAdmin(account, pwd, mer_code, power_code string) (int, string) {
	t_status := 100
	t_msg := "填写完整"
	if account == "" || mer_code == "" || power_code == "" {
		return t_status, t_msg
	}
	if pwd == "" {
		pwd = "123456"
	}
	a_info := model.AdminInfoRedis(account)
	if len(a_info["account"]) > 1 {
		t_msg = "用户名已存在"
		return t_status, t_msg
	}
	table_name := "admin_list"
	a_data := map[string]string{}
	a_data["account"] = account
	a_data["pwd"] = common.HexMd5(pwd)
	a_data["mer_code"] = mer_code

	admin_sql := common.InsertSql(table_name, a_data)
	err := model.Query(admin_sql)
	if err != nil {
		t_msg = "管理员新增失败"
		return t_status, t_msg
	}
	t_status = 200
	t_msg = "success"
	return t_status, t_msg
}

/**
*  更新用户的支付
 */
func PayStatus(order_number, pay_status string, ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "请填写完整"
	if order_number == "" || pay_status == "" {
		return t_status, t_msg
	}
	t_status, t_msg = updatePayStatus(order_number, pay_status)
	return t_status, t_msg
}

/**
*  更新用户的出款订单
 */
func CashStatus(order_number, cash_status string, ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "请填写完整"
	order := true
	if order_number == "" || cash_status == "" {
		return t_status, t_msg
	}

	c_list := model.CashById(order_number)
	t_msg = "订单号不存在"
	if len(c_list.Id) < 1 {
		return t_status, t_msg
	}

	o_list := model.CashOrderByweb(c_list.Order_number)
	if len(o_list.Id) < 1 {
		order = false
	}

	note := "手动完成"
	if cash_status == "1" { //改为处理中

		note = "手动改为处理中"
		//修改渠道订单
		if order {
			updateOrderStatusForRun(o_list.Id, o_list.Id, cash_status, note)
		}
		//修改代付订单
		t_status, t_msg = updateCashStatusForRun(order_number, order_number, cash_status, note)

	} else { //改为成功或者失败
		//修改渠道订单
		if order {
			updateOrderStatus(o_list.Id, o_list.Id, cash_status, note)
		}
		//修改代付订单
		t_status, t_msg = updateCashStatus(order_number, order_number, cash_status, note)
	}

	return t_status, t_msg
}

/**
*  支付订单列表
 */
func PayList(pay_status, order_number, web_order, pay_code, class_code, is_mobile, start_time, end_time, mer_code, is_agent, page, page_size string, ctx *gin.Context) (int, int, string, float64, []map[string]string) {
	t_status := 200
	t_msg := "success"
	total := 0
	pay_list := []map[string]string{}
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
	total, _ = model.DateListTotal(table_name, date_field, start_time, end_time, like_sql, count_field, p_where)
	if total < 1 {
		return t_status, total, t_msg, 0, pay_list
	}
	page_int, size_int := ThreadPage(page, page_size)
	offset := (page_int - 1) * size_int
	fields := []string{"id", "order_number", "status", "pay_code", "mer_code", "push_status", "push_num", "amount", "real_amount", "create_time", "pay_time", "class_code", "bank_code", "note", "is_mobile", "rate"}
	pay_list, _ = model.PageDateList(table_name, date_field, start_time, end_time, like_sql, size_int, offset, fields, p_where)
	if len(pay_list) < 1 {
		return t_status, total, t_msg, 0, pay_list
	}
	field := "sum(amount) as total"
	_, sum_amount := model.DateListTotal(table_name, date_field, start_time, end_time, like_sql, field, p_where)
	for p_k, p_v := range pay_list {
		b_info := model.BankInfoRedis(p_v["bank_code"])
		pay_list[p_k]["bank_title"] = b_info["title"]
		p_info := model.PayInfoRedis(p_v["pay_code"])
		pay_list[p_k]["pay_title"] = p_info["title"]
		c_info := model.ClassInfoRedis(p_v["class_code"])
		pay_list[p_k]["class_title"] = c_info["title"]
	}
	return t_status, total, t_msg, sum_amount, pay_list
}

func CashList(cash_status, order_number, web_order, pay_code, start_time, end_time, mer_code, is_agent, page, page_size string, ctx *gin.Context) (int, int, string, float64, []map[string]string) {
	t_status := 200
	t_msg := "success"
	total := 0
	pay_list := []map[string]string{}
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
	date_field := "create_time"
	count_field := "count(0) as num"
	total, _ = model.DateListTotal(table_name, date_field, start_time, end_time, like_sql, count_field, p_where)
	if total < 1 {
		return t_status, total, t_msg, 0, pay_list
	}
	page_int, size_int := ThreadPage(page, page_size)
	offset := (page_int - 1) * size_int
	fields := []string{"id", "status", "order_number", "pay_code", "mer_code", "push_status", "push_num", "amount", "real_amount", "create_time", "pay_time", "branch", "bank_title", "note", "card_name", "card_number"}
	pay_list, _ = model.PageDateList(table_name, date_field, start_time, end_time, like_sql, size_int, offset, fields, p_where)
	if len(pay_list) < 1 {
		return t_status, total, t_msg, 0, pay_list
	}
	field := "sum(amount) as total"
	_, sum_amount := model.DateListTotal(table_name, date_field, start_time, end_time, like_sql, field, p_where)
	for p_k, p_v := range pay_list {
		p_info := model.PayInfoRedis(p_v["pay_code"])
		pay_list[p_k]["pay_title"] = p_info["title"]
	}
	return t_status, total, t_msg, sum_amount, pay_list
}

/**
*  现金流水
 */
func AmountList(amount_type, order_number, pay_code, start_time, end_time, mer_code, is_agent, page, page_size string, ctx *gin.Context) (int, int, string, []map[string]string) {
	t_status := 200
	t_msg := "success"
	total := 0
	pay_list := []map[string]string{}
	session := sessions.Default(ctx)
	admin_mer := fmt.Sprintf("%v", session.Get("mer_code"))

	p_where := map[string]interface{}{}
	like_sql := ""
	amount_type_in := []string{}

	if amount_type != "" {
		amount_type_in = strings.Split(amount_type, ",")
	}

	if mer_code != "" {
		p_where["mer_code"] = mer_code
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
	date_field := "create_time"
	count_field := "count(0) as num"
	total, _ = model.DateListTotal(table_name, date_field, start_time, end_time, like_sql, count_field, p_where)
	if total < 1 {
		return t_status, total, t_msg, pay_list
	}
	page_int, size_int := ThreadPage(page, page_size)
	offset := (page_int - 1) * size_int
	fields := []string{"id", "amount_type", "pay_code", "mer_code", "before_amount", "after_amount", "amount", "create_time", "order_number", "note"}
	fild_in := "amount_type in (?)"
	if len(amount_type_in) < 1 {
		for i := 1; i < 11; i++ {
			if admin_mer != "all" {
				if i == 9 || i == 10 {
					continue
				}
			}
			amount_type_in = append(amount_type_in, strconv.Itoa(i))
		}
	}
	pay_list, _ = model.SecondPageDateList(table_name, date_field, start_time, end_time, like_sql, fild_in, amount_type_in, size_int, offset, fields, p_where)

	if len(pay_list) < 1 {
		return t_status, total, t_msg, pay_list
	}
	for p_k, p_v := range pay_list {
		p_info := model.PayInfoRedis(p_v["pay_code"])
		pay_list[p_k]["pay_title"] = p_info["title"]
	}
	return t_status, total, t_msg, pay_list
}

/**
*  近期总存款，总出款
 */
func DateTotal(start_date, end_date string, ctx *gin.Context) (int, string, []string, []map[string]interface{}) {
	t_status := 100
	t_msg := "日期格式错误"
	title_list := []string{}
	data_list := []map[string]interface{}{}
	if start_date == "" {
		start_date = time.Now().Format(format_day)
	}
	if end_date == "" {
		end_date = time.Now().Format(format_day)
	}
	//获取当前日期
	day_num := common.DifferDays(start_date, end_date, format_day)
	if day_num < 0 {
		return t_status, t_msg, title_list, data_list
	}
	if day_num > 15 {
		day_num = 15
	}
	data_title := []string{"总收入", "总出款"}
	title_list = common.DateDiff(start_date, format_day, day_num)
	sort.Strings(title_list)

	session := sessions.Default(ctx)
	mer_code := fmt.Sprintf("%v", session.Get("mer_code"))

	for _, t_v := range data_title {
		v_map := map[string]interface{}{}
		v_map["title"] = t_v
		data_list = append(data_list[0:], v_map)
	}

	m_info := model.MerInfoRedis(mer_code)
	if len(m_info["id"]) < 1 {
		return t_status, t_msg, title_list, data_list
	}

	p_w := map[string]interface{}{}
	like_sql := ""
	if mer_code != "all" {
		like_sql = " and (agent_path like '" + m_info["agent_path"] + m_info["code"] + "_%' or mer_code='" + mer_code + "')"
	}

	p_w["status"] = 3

	field := "sum(amount) as total"
	date_field := "create_time"

	pay_total := []float64{}
	cash_total := []float64{}
	for _, title_val := range title_list {
		s_time := title_val + " 00:00:00"
		e_time := title_val + " 23:59:59"
		_, p_total := model.DateListTotal("pay_list", date_field, s_time, e_time, like_sql, field, p_w)
		_, c_total := model.DateListTotal("cash_list", date_field, s_time, e_time, like_sql, field, p_w)
		pay_total = append(pay_total[0:], p_total)
		cash_total = append(cash_total[0:], c_total)
	}
	for d_k, _ := range data_list {
		if d_k == 0 {
			data_list[d_k]["data_list"] = pay_total
		} else if d_k == 1 {
			data_list[d_k]["data_list"] = cash_total
		}
	}
	t_status = 200
	t_msg = "success"
	return t_status, t_msg, title_list, data_list
}

func TotalBalance(ctx *gin.Context) (int, string, map[string]interface{}) {
	t_status := 200
	t_msg := "success"
	session := sessions.Default(ctx)
	mer_code := fmt.Sprintf("%v", session.Get("mer_code"))

	balance_map := map[string]interface{}{}

	balance_map["balance"] = 0.00
	balance_map["group_balance"] = 0.00
	admin_mer := model.MerInfo(mer_code)
	if admin_mer.Id < 1 {
		t_status = 100
		t_msg = "管理员信息异常"
		return t_status, t_msg, balance_map
	}
	table_name := "mer_list"
	field := "sum(amount) as total"
	p_where := map[string]interface{}{}
	p_where["code"] = mer_code
	balance_map["balance"] = admin_mer.Amount
	like_field := "agent_path"
	like_where := admin_mer.Agent_path + admin_mer.Code + "_%"
	_, balance_map["group_balance"] = model.LikeListTotal(table_name, like_field, like_where, field, p_where)

	return t_status, t_msg, balance_map
}

func TodayCount(date_type string, ctx *gin.Context) (int, string, map[string]interface{}) {
	t_status := 100
	t_msg := "管理员信息异常"
	count_info := map[string]interface{}{}
	time_date := time.Now().Format(format_day)
	start_date := time_date + " 00:00:00"
	end_date := time_date + " 23:59:59"

	if date_type == "1" {
		start_date = common.ChangeDate(time_date, format_day, 0, 0, -1) + " 00:00:00"
		end_date = common.ChangeDate(time_date, format_day, 0, 0, -1) + " 23:59:59"
	} else if date_type == "2" {
		start_date = common.ChangeDate(time_date, format_day, 0, 0, -7) + " 00:00:00"
	} else if date_type == "3" {
		start_date = common.ChangeDate(time_date, format_day, 0, 0, -15) + " 00:00:00"
	}

	session := sessions.Default(ctx)
	mer_code := fmt.Sprintf("%v", session.Get("mer_code"))

	m_info := model.MerInfoRedis(mer_code)
	if len(m_info["id"]) < 1 {
		return t_status, t_msg, count_info
	}

	p_w := map[string]interface{}{}
	like_sql := ""
	if mer_code != "all" {
		like_sql = " and (agent_path like '" + m_info["agent_path"] + m_info["code"] + "_%' or mer_code='" + mer_code + "')"
	}

	p_w["status"] = 3

	field := "sum(amount) as total"
	date_field := "create_time"

	_, count_info["in_total"] = model.DateListTotal("pay_list", date_field, start_date, end_date, like_sql, field, p_w)
	_, count_info["out_total"] = model.DateListTotal("cash_list", date_field, start_date, end_date, like_sql, field, p_w)
	t_status = 200
	t_msg = "success"
	return t_status, t_msg, count_info
}

func MerInfo(ctx *gin.Context) (int, string, map[string]string) {
	t_status := 200
	t_msg := "success"
	session := sessions.Default(ctx)
	mer_code := fmt.Sprintf("%v", session.Get("mer_code"))
	mer_info := model.MerInfoRedis(mer_code)
	return t_status, t_msg, mer_info
}

func RateInfo(ctx *gin.Context) (int, string, []map[string]string) {
	t_status := 200
	t_msg := "success"
	session := sessions.Default(ctx)
	mer_code := fmt.Sprintf("%v", session.Get("mer_code"))
	p_where := map[string]interface{}{}
	p_where["mer_code"] = mer_code
	table_name := "mer_rate"
	fields := []string{"mer_code", "pay_code", "class_code", "bank_code", "rate", "limit_amount", "day_amount"}
	rate_info, _ := model.PageList(table_name, "", 1000, 0, fields, p_where)
	if len(rate_info) < 1 {
		return t_status, t_msg, rate_info
	}
	for r_k, r_val := range rate_info {
		pay := model.PayInfoRedis(r_val["pay_code"])
		rate_info[r_k]["pay_title"] = pay["title"]
		class := model.ClassInfoRedis(r_val["class_code"])
		rate_info[r_k]["class_title"] = class["title"]
		sys := model.BankInfoRedis(r_val["bank_code"])
		rate_info[r_k]["bank_title"] = sys["title"]
	}
	return t_status, t_msg, rate_info
}

func CallPay(order_number string) (int, string) {
	t_status := 100
	t_msg := "订单号错误"
	p_list := model.OrderById(order_number)
	if len(p_list.Order_number) < 1 {
		return t_status, t_msg
	}
	go Push(p_list)
	t_status = 200
	t_msg = "success"
	return t_status, t_msg
}

func CallCash(order_number string) (int, string) {
	t_status := 100
	t_msg := "订单号错误"
	p_list := model.CashById(order_number)
	if len(p_list.Order_number) < 1 {
		return t_status, t_msg
	}
	go PushCash(p_list)
	t_status = 200
	t_msg = "success"
	return t_status, t_msg
}

func MerPay(ctx *gin.Context) (int, string, []map[string]string) {
	t_status := 200
	t_msg := "success"
	session := sessions.Default(ctx)
	mer_code := fmt.Sprintf("%v", session.Get("mer_code"))
	mer_info := model.MerInfo(mer_code)
	b_map := map[string]string{}
	b_map["mer_code"] = mer_info.Code
	b_map["amount"] = fmt.Sprintf("%.2f", mer_info.Amount)
	b_map["total_in"] = fmt.Sprintf("%.2f", mer_info.Total_in)
	b_map["total_out"] = fmt.Sprintf("%.2f", mer_info.Total_out)
	b_map["pay_id"] = "0"
	balance_list := []map[string]string{}
	balance_list = append(balance_list[0:], b_map)
	return t_status, t_msg, balance_list
}

/**
*  后台下发
 */
func MerCash(pay_id, bank_id, amount, is_auto, secret string, ctx *gin.Context) (int, string, string) {
	t_status := 100
	t_msg := "请填写完整"
	api_res := "fail"
	if pay_id == "" || bank_id == "" || amount == "" || secret == "" {
		return t_status, t_msg, api_res
	}

	mer_bank := model.MerBankByIdRedis(bank_id)
	if len(mer_bank["bank_code"]) < 1 {
		t_msg = "银行卡ID错误"
		return t_status, t_msg, api_res
	}

	if mer_bank["status"] != "1" {
		t_msg = "银行卡已被锁定"
		return t_status, t_msg, api_res
	}

	session := sessions.Default(ctx)
	mer_code := fmt.Sprintf("%v", session.Get("mer_code"))
	account := fmt.Sprintf("%v", session.Get("account"))

	a_list := model.AdminInfo(account)
	if a_list.Id < 1 {
		t_msg = "管理员信息异常"
		return t_status, t_msg, api_res
	}

	t_status, t_msg = AuthGoogle(secret, a_list.Secret)
	if t_status != 200 {
		t_msg = "谷歌动态验证码错误"
		return t_status, t_msg, api_res
	}
	t_status = 100
	p_map := map[string]string{}
	p_map["order_number"] = model.GetKey(20)
	p_map["bank_code"] = mer_bank["bank_code"]
	p_map["bank_title"] = mer_bank["bank_title"]
	p_map["card_number"] = mer_bank["card_number"]
	p_map["card_name"] = mer_bank["card_name"]
	p_map["pay_id"] = pay_id
	p_map["amount"] = amount
	p_map["bank_branch"] = mer_bank["bank_branch"]
	p_map["phone"] = mer_bank["bank_phone"]
	if is_auto == "1" {
		//查询额度
		if pay_id == "0" {
			//查询额度大于
			mer_table := "mer_pay"
			order_by := "pay_id desc"
			page_size := 100
			offset := 0
			fields := []string{"pay_id"}
			p_where := map[string]interface{}{}
			p_where["mer_code"] = mer_code
			p_where["status"] = 1
			m_list, _ := model.PageList(mer_table, order_by, page_size, offset, fields, p_where)
			if len(m_list) < 1 {
				t_msg = "商户暂时没有渠道"
				return t_status, t_msg, api_res
			}
			//id的数组
			id_arr := []string{}
			for _, m_info := range m_list {
				if len(m_info["pay_id"]) > 0 {
					id_arr = append(id_arr[0:], m_info["pay_id"])
				}
			}
			//
			conf_table := "pay_config"
			in_field := "id"
			in_fields := []string{"amount", "id"}
			c_where := map[string]interface{}{}
			c_where["status"] = 1
			p_size := 1
			conf_list, _ := model.InPageList(conf_table, in_field, p_size, offset, in_fields, id_arr, c_where)
			if len(conf_list) < 1 {
				t_msg = "商户暂时没有渠道"
				return t_status, t_msg, api_res
			}
			amount_f, _ := strconv.ParseFloat(amount, 64)
			conf_amount_f, _ := strconv.ParseFloat(conf_list[0]["amount"], 64)
			if conf_amount_f < amount_f {
				t_msg = "渠道额度不足"
				return t_status, t_msg, api_res
			}
			pay_id = conf_list[0]["id"]
		}
		p_map["pay_id"] = pay_id
		t_status, t_msg, api_res = DFpay(mer_code, p_map)
		return t_status, t_msg, api_res
	}
	t_status, t_msg, api_res = PayFor(mer_code, p_map)
	if api_res == "fail" {
		t_status = 100
	}

	return t_status, t_msg, api_res
}

/**
*  绑定银行卡
 */
func AddBank(bank_code, card_number, card_name, bank_branch, bank_phone, secret string, ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "请填写完整"
	if bank_code == "" || card_number == "" || card_name == "" || bank_branch == "" || bank_phone == "" || secret == "" {
		return t_status, t_msg
	}

	session := sessions.Default(ctx)
	mer_code := fmt.Sprintf("%v", session.Get("mer_code"))
	account := fmt.Sprintf("%v", session.Get("account"))

	is_exist := model.MerBankRedis(card_number, mer_code)
	if len(is_exist["bank_code"]) > 1 {
		t_msg = "银行卡已被绑定"
		return t_status, t_msg
	}

	a_list := model.AdminInfo(account)
	if a_list.Id < 1 {
		t_msg = "管理员信息异常"
		return t_status, t_msg
	}

	t_status, t_msg = AuthGoogle(secret, a_list.Secret)
	if t_status != 200 {
		t_msg = "谷歌动态验证码错误"
		return t_status, t_msg
	}

	bank := model.BankInfoRedis(bank_code)
	if len(bank["title"]) < 1 {
		t_msg = "银行编码错误"
		return t_status, t_msg
	}

	table_name := "mer_bank"
	m_data := map[string]string{}
	m_data["mer_code"] = mer_code
	m_data["bank_code"] = bank_code
	m_data["bank_title"] = bank["title"]
	m_data["card_number"] = card_number
	m_data["card_name"] = card_name
	m_data["bank_branch"] = bank_branch
	m_data["bank_phone"] = bank_phone
	in_sql := common.InsertSql(table_name, m_data)
	err := model.Query(in_sql)
	if err != nil {
		t_msg = "绑定失败"
		return t_status, t_msg
	}
	t_status = 200
	t_msg = "success"
	return t_status, t_msg
}

/**
*  商户绑定的银行卡列表
 */
func MerBank(bank_code, bank_status, card_number, card_name, page, page_size string, ctx *gin.Context) (int, int, string, []map[string]string) {
	t_status := 200
	t_msg := "success"
	mer_bank := []map[string]string{}
	session := sessions.Default(ctx)
	mer_code := fmt.Sprintf("%v", session.Get("mer_code"))
	p_w := map[string]interface{}{}
	p_w["mer_code"] = mer_code

	if bank_status == "0" {
		p_w["status"] = 0
	} else if bank_status == "2" {
	} else {
		p_w["status"] = 1
	}
	if bank_code != "" {
		p_w["bank_code"] = bank_code
	}
	table_name := "mer_bank"
	count_field := "count(id) as num"
	total, _ := model.ListTotal(table_name, count_field, p_w)
	if total < 1 {
		return t_status, total, t_msg, mer_bank
	}

	field := []string{"id", "mer_code", "bank_title", "card_number", "card_name", "bank_branch", "bank_phone", "status"}
	page_int, size_int := ThreadPage(page, page_size)
	offset := (page_int - 1) * size_int
	//由于用户银行卡过多，当卡号,名字不为空时，这里加上模糊搜索，方便用户查找
	if card_number == "" && card_name == "" {
		mer_bank, _ = model.PageList(table_name, "", size_int, offset, field, p_w)
	} else {
		lieke_field := "card_number"
		like_where := "%" + card_number + "%"
		if card_name != "" {
			lieke_field = "card_name"
			like_where = "%" + card_name + "%"
		}
		mer_bank, _ = model.LikePageList(table_name, lieke_field, like_where, size_int, offset, field, p_w)
	}
	return t_status, total, t_msg, mer_bank
}

func LockBank(b_id, secret string, ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "请填写完整"
	if b_id == "" || secret == "" {
		return t_status, t_msg
	}
	is_exist := model.MerBankByIdRedis(b_id)
	if len(is_exist["id"]) < 1 {
		t_msg = "ID错误"
		return t_status, t_msg
	}
	session := sessions.Default(ctx)
	mer_code := fmt.Sprintf("%v", session.Get("mer_code"))
	if mer_code != is_exist["mer_code"] {
		t_msg = "没有权限"
		return t_status, t_msg
	}

	account := fmt.Sprintf("%v", session.Get("account"))
	a_list := model.AdminInfo(account)
	if a_list.Id < 1 {
		t_msg = "管理员信息异常"
		return t_status, t_msg
	}

	t_status, t_msg = AuthGoogle(secret, a_list.Secret)
	if t_status != 200 {
		t_msg = "谷歌动态验证码错误"
		return t_status, t_msg
	}

	up_sql := fmt.Sprintf("update mer_bank set `status`=0 where id='%s';", b_id)
	err := model.Query(up_sql)
	if err != nil {
		t_msg = "锁定失败"
		return t_status, t_msg
	}
	t_status = 200
	t_msg = "success"
	return t_status, t_msg
}

func AgentReport(mer_code, start_date, end_date, page, page_size string, ctx *gin.Context) (int, int, string, map[string]float64, []map[string]string) {
	t_status := 200
	t_msg := "success"
	total := 0
	agent_report := []map[string]string{}
	sum_total := map[string]float64{}
	session := sessions.Default(ctx)
	admin_mer := fmt.Sprintf("%v", session.Get("mer_code"))
	admin := model.MerInfoRedis(admin_mer)
	p_path := admin["agent_path"] + admin_mer + "_"
	if mer_code != "" {
		mer_info := model.MerInfoRedis(mer_code)
		if !strings.Contains(mer_info["agent_path"], p_path) {
			return t_status, total, t_msg, sum_total, agent_report
		}
		p_path = mer_info["agent_path"] + mer_info["code"] + "_"
	}
	if start_date == "" {
		start_date = time.Now().AddDate(0, 0, -1).Format(format_day) + " 00:00:00"
	} else {
		start_date = start_date + " 00:00:00"
	}
	if end_date == "" {
		end_date = time.Now().AddDate(0, 0, -1).Format(format_day) + " 23:59:59"
	} else {
		end_date = end_date + " 23:59:59"
	}
	p_w := map[string]interface{}{}
	p_w["agent_path"] = p_path

	table_name := "mer_report"
	count_field := "count(id) as num"
	date_field := "report_date"
	like_sql := ""
	total, _ = model.DateListTotal(table_name, date_field, start_date, end_date, like_sql, count_field, p_w)
	if total < 1 {
		return t_status, total, t_msg, sum_total, agent_report
	}
	field := []string{"id", "mer_code", "report_date", "total_in", "total_in_rate", "total_out", "total_out_rate", "is_agent"}
	page_int, size_int := ThreadPage(page, page_size)
	offset := (page_int - 1) * size_int
	agent_report, _ = model.PageDateList(table_name, date_field, start_date, end_date, like_sql, size_int, offset, field, p_w)

	//计算入款总和
	count_field = "sum(total_in) as total"
	_, total_in := model.DateListTotal(table_name, date_field, start_date, end_date, like_sql, count_field, p_w)
	sum_total["total_in"] = total_in
	//计算费率总和
	count_field = "sum(total_in_rate) as total"
	_, total_in_rate := model.DateListTotal(table_name, date_field, start_date, end_date, like_sql, count_field, p_w)
	sum_total["total_in_rate"] = total_in_rate
	//计算出歀总和
	count_field = "sum(total_out) as total"
	_, total_out := model.DateListTotal(table_name, date_field, start_date, end_date, like_sql, count_field, p_w)
	sum_total["total_out"] = total_out
	//计算出歀费率总和
	count_field = "sum(total_out_rate) as total"
	_, total_out_rate := model.DateListTotal(table_name, date_field, start_date, end_date, like_sql, count_field, p_w)
	sum_total["total_out_rate"] = total_out_rate

	return t_status, total, t_msg, sum_total, agent_report
}

func MerChannel(mer_code, pay_code, chann_status, page, page_size string, ctx *gin.Context) (int, int, string, []map[string]string) {
	t_status := 200
	t_msg := "success"

	mer_channel := []map[string]string{}
	session := sessions.Default(ctx)
	admin_mer := fmt.Sprintf("%v", session.Get("mer_code"))
	p_where := map[string]interface{}{}
	if admin_mer != "all" {
		p_where["mer_code"] = admin_mer
	} else if admin_mer == "all" && mer_code != "" {
		p_where["mer_code"] = mer_code
	}
	if pay_code != "" {
		p_where["pay_code"] = pay_code
	}
	if chann_status != "" {
		status, err := strconv.Atoi(chann_status)
		if err == nil {
			p_where["status"] = status
		}
	}

	order_by := "pay_id desc"
	table_name := "mer_pay"
	total_field := "count(0) as num"
	total, _ := model.ListTotal(table_name, total_field, p_where)
	if total < 1 {
		return t_status, total, t_msg, mer_channel
	}

	page_int, size_int := ThreadPage(page, page_size)
	offset := (page_int - 1) * size_int
	fields := []string{"id", "pay_code", "pay_id", "status", "mer_code"}
	mer_channel, _ = model.PageList(table_name, order_by, size_int, offset, fields, p_where)
	if len(mer_channel) < 1 {
		return t_status, total, t_msg, mer_channel
	}
	for p_k, p_val := range mer_channel {
		pay := model.PayInfoRedis(p_val["pay_code"])
		mer_channel[p_k]["pay_title"] = pay["title"]
		conf := model.ApiConfigRedis(p_val["pay_id"])
		mer_channel[p_k]["merchant_code"] = conf["merchant_code"]
	}
	return t_status, total, t_msg, mer_channel
}

/**
*  新增商户渠道
 */
func AddMerChannel(pay_id, mer_code string, ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "请填写完整"
	if pay_id == "" || mer_code == "" {
		return t_status, t_msg
	}

	session := sessions.Default(ctx)
	admin_mer := fmt.Sprintf("%v", session.Get("mer_code"))

	if admin_mer != "all" {
		t_msg = "该账号没有权限"
		return t_status, t_msg
	}

	is_exist := model.ApiConfigRedis(pay_id)
	if len(is_exist["pay_code"]) < 1 {
		t_msg = "渠道编码错误"
		return t_status, t_msg
	}

	table_name := "mer_pay"
	m_data := map[string]string{}
	m_data["mer_code"] = mer_code
	m_data["pay_code"] = is_exist["pay_code"]
	m_data["pay_id"] = pay_id
	m_data["status"] = "1"
	in_sql := common.InsertSql(table_name, m_data)
	err := model.Query(in_sql)
	if err != nil {
		t_msg = "添加失败"
		return t_status, t_msg
	}
	t_status = 200
	t_msg = "success"
	return t_status, t_msg
}

/**
*  修改商户渠道
 */
func EditMerChannel(p_id, channel_status string) (int, string) {
	t_status := 100
	t_msg := "请填写完整"
	if p_id == "" || channel_status == "" {
		return t_status, t_msg
	}

	is_exist := model.MerPayById(p_id)
	if is_exist.Id < 1 {
		t_msg = "渠道编码错误"
		return t_status, t_msg
	}
	p_data := map[string]interface{}{}
	p_data["status"] = 1
	if channel_status == "0" {
		p_data["status"] = 0
	}
	err := model.UpdatesMerPay(is_exist, p_data)
	if err != nil {
		t_msg = "修改失败"
		return t_status, t_msg
	}
	t_status = 200
	t_msg = "success"
	return t_status, t_msg
}

func OrderList(order_id, cash_status, order_number, web_order, pay_code, start_time, end_time, page, page_size string, ctx *gin.Context) (int, int, string, []map[string]string) {
	t_status := 200
	t_msg := "success"
	total := 0
	pay_list := []map[string]string{}
	session := sessions.Default(ctx)
	admin_mer := fmt.Sprintf("%v", session.Get("mer_code"))

	if admin_mer != "all" {
		return t_status, total, t_msg, pay_list
	}

	p_where := map[string]interface{}{}
	like_sql := ""
	if order_id != "" {
		p_where["id"] = order_id
	}
	if cash_status != "" {
		p_where["status"] = cash_status
	}
	if order_number != "" {
		p_where["cash_id"] = order_number
	}
	if web_order != "" {
		p_where["order_number"] = web_order
	}
	if pay_code != "" {
		p_where["pay_code"] = pay_code
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
	table_name := "order_list"
	date_field := "create_time"
	count_field := "count(0) as num"
	total, _ = model.DateListTotal(table_name, date_field, start_time, end_time, like_sql, count_field, p_where)
	if total < 1 {
		return t_status, total, t_msg, pay_list
	}
	page_int, size_int := ThreadPage(page, page_size)
	offset := (page_int - 1) * size_int
	fields := []string{"id", "status", "order_number", "pay_code", "amount", "real_amount", "create_time", "pay_time", "branch", "bank_title", "note", "card_name", "card_number", "cash_id", "pay_id"}
	pay_list, _ = model.PageDateList(table_name, date_field, start_time, end_time, like_sql, size_int, offset, fields, p_where)
	if len(pay_list) < 1 {
		return t_status, total, t_msg, pay_list
	}
	for p_k, p_v := range pay_list {
		p_info := model.PayInfoRedis(p_v["pay_code"])
		pay_list[p_k]["pay_title"] = p_info["title"]
	}
	return t_status, total, t_msg, pay_list
}

/**
*  渠道下发
 */
func AddOrder(pay_id, order_number, amount, secret string, ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "请填写完整"
	if pay_id == "" || order_number == "" || amount == "" || secret == "" {
		return t_status, t_msg
	}

	session := sessions.Default(ctx)
	mer_code := fmt.Sprintf("%v", session.Get("mer_code"))
	if mer_code != "all" {
		t_msg = "账号没有权限"
		return t_status, t_msg
	}

	c_list := model.CashById(order_number)
	if len(c_list.Id) < 1 {
		t_msg = "系统订单号错误"
		return t_status, t_msg
	}

	if c_list.Status != 1 {
		t_msg = "订单状态错误"
		return t_status, t_msg
	}

	account := fmt.Sprintf("%v", session.Get("account"))

	a_list := model.AdminInfo(account)
	if a_list.Id < 1 {
		t_msg = "管理员信息异常"
		return t_status, t_msg
	}

	t_status, t_msg = AuthGoogle(secret, a_list.Secret)
	if t_status != 200 {
		t_msg = "谷歌动态验证码错误"
		return t_status, t_msg
	}

	t_status, t_msg = adminOrder(pay_id, amount, c_list)

	return t_status, t_msg
}

/**
*  更新用户的出款订单
 */
func UpdateOrder(order_id, order_status string, ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "请填写完整"
	if order_id == "" || order_status == "" {
		return t_status, t_msg
	}
	session := sessions.Default(ctx)
	mer_code := fmt.Sprintf("%v", session.Get("mer_code"))
	if mer_code != "all" {
		t_msg = "账号没有权限"
		return t_status, t_msg
	}
	account := fmt.Sprintf("%v", session.Get("account"))
	note := account + "后台手动完成"
	t_status, t_msg = updateOrderStatus(order_id, order_id, order_status, note)
	return t_status, t_msg
}

/**
*  新增用户额度
 */
func AddMerAmount(mer_code, amount, secret, pay_id, note string, ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "请填写完整"
	if mer_code == "" || amount == "" || secret == "" || pay_id == "" {
		return t_status, t_msg
	}
	amount_f, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		t_msg = "额度错误"
		return t_status, t_msg
	}
	session := sessions.Default(ctx)
	admin_mer := fmt.Sprintf("%v", session.Get("mer_code"))
	if admin_mer != "all" {
		t_msg = "账号没有权限"
		return t_status, t_msg
	}
	account := fmt.Sprintf("%v", session.Get("account"))
	a_list := model.AdminInfo(account)
	if a_list.Id < 1 {
		t_msg = "管理员信息异常"
		return t_status, t_msg
	}
	t_status, t_msg = AuthGoogle(secret, a_list.Secret)
	if t_status != 200 {
		t_msg = "谷歌动态验证码错误"
		return t_status, t_msg
	}
	mer_info := model.MerInfo(mer_code)
	if mer_info.Id < 1 {
		t_msg = "商户号错误"
		return t_status, t_msg
	}
	pay_conf := model.PayConf(pay_id)
	if pay_conf.Id < 1 {
		t_msg = "渠道编码错误"
		return t_status, t_msg
	}
	if note == "" {
		note = "后台手工添加"
	}

	//新增用户额度
	mer_sql := fmt.Sprintf("update mer_list set amount=amount+%.2f,total_in=total_in+%.2f where code='%s';", amount_f, amount_f, mer_info.Code)
	//更新pay_config记录
	conf_sql := fmt.Sprintf("update pay_config set amount=amount+%.2f,total_in=total_in+%.2f where id=%d;", amount_f, amount_f, pay_conf.Id)
	//新增pay_list记录
	now_time := time.Now().Format(format_date)
	table_name := "pay_list"
	list_id := model.GetKey(17)
	pay_list := map[string]string{}
	pay_list["id"] = list_id
	amount_type := "1"
	pay_list["pay_code"] = pay_conf.Pay_code
	pay_list["pay_id"] = pay_id
	pay_list["mer_code"] = mer_code
	pay_list["amount"] = amount
	pay_list["real_amount"] = amount
	pay_list["create_time"] = now_time
	pay_list["order_number"] = list_id
	pay_list["class_code"] = "bank"
	pay_list["bank_code"] = "bank"
	pay_list["push_url"] = ""
	pay_list["note"] = note
	pay_list["status"] = "3"
	pay_list["agent_path"] = mer_info.Agent_path

	if amount_f < 0.00 {
		amount_type = "5"
	}

	pay_sql := common.InsertSql(table_name, pay_list)

	//新增账变记录
	after_amount := mer_info.Amount + amount_f
	log_table := "amount_list"
	a_data := map[string]string{}
	a_data["id"] = model.GetKey(20)
	a_data["order_number"] = list_id
	a_data["create_time"] = time.Now().Format(format_date)
	a_data["amount_type"] = amount_type
	a_data["amount"] = amount
	a_data["pay_code"] = pay_conf.Pay_code
	a_data["mer_code"] = mer_code
	a_data["pay_id"] = fmt.Sprintf("%d", pay_conf.Id)
	a_data["before_amount"] = fmt.Sprintf("%.4f", mer_info.Amount)
	a_data["after_amount"] = fmt.Sprintf("%.4f", after_amount)
	a_data["agent_path"] = mer_info.Agent_path
	a_data["note"] = note
	log_sql := common.InsertSql(log_table, a_data)
	//采用事务处理
	sql_arr := []string{mer_sql, conf_sql, pay_sql, log_sql}
	err = model.Trans(sql_arr)
	if err != nil {
		t_msg = "充值失败"
		return t_status, t_msg
	}
	t_status = 200
	t_msg = "success"
	return t_status, t_msg
}

/**
*  人工充值
 */
func ManualRecharge(mer_code, amount, secret, pay_id, bank_code, class_code, note string, ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "请填写完整"
	if mer_code == "" || amount == "" || secret == "" || pay_id == "" || bank_code == "" || class_code == "" {
		return t_status, t_msg
	}
	amount_f, err := strconv.ParseFloat(amount, 64)
	if err != nil || amount_f <= 0 {
		t_msg = "额度错误"
		return t_status, t_msg
	}
	session := sessions.Default(ctx)
	admin_mer := fmt.Sprintf("%v", session.Get("mer_code"))
	if admin_mer != "all" {
		t_msg = "账号没有权限"
		return t_status, t_msg
	}
	account := fmt.Sprintf("%v", session.Get("account"))
	a_list := model.AdminInfo(account)
	if a_list.Id < 1 {
		t_msg = "管理员信息异常"
		return t_status, t_msg
	}
	t_status, t_msg = AuthGoogle(secret, a_list.Secret)
	if t_status != 200 {
		t_msg = "谷歌动态验证码错误"
		return t_status, t_msg
	}
	mer_info := model.MerInfo(mer_code)
	if mer_info.Id < 1 {
		t_msg = "商户号错误"
		return t_status, t_msg
	}
	pay_conf := model.PayConf(pay_id)
	if pay_conf.Id < 1 {
		t_msg = "渠道编码错误"
		return t_status, t_msg
	}
	//查看商户可以支付的渠道
	r_list := model.RateList(mer_code, class_code, bank_code)
	if len(r_list) < 1 {
		t_msg = "商户未分配渠道"
		return t_status, t_msg
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
		return t_status, t_msg
	}

	//判断是否有渠道已满
	rate_list := model.PayRateList(code_arr)
	if len(rate_list) < 1 {
		t_msg = "暂无支付渠道"
		return t_status, t_msg
	}

	if note == "" {
		note = "后台手工添加"
	}

	m_rate := 0.00

	for _, r_val := range r_list {
		if pay_conf.Pay_code == r_val.Pay_code && class_code == r_val.Class_code && bank_code == r_val.Bank_code {
			m_rate = r_val.Rate
			break
		}
	}

	if m_rate == 0.00 || m_rate > 1.00 {
		t_msg = "费率异常"
		return t_status, t_msg
	}
	real_amout := (1 - m_rate) * amount_f
	//新增pay_list记录
	now_time := time.Now().Format(format_date)
	table_name := "pay_list"
	list_id := model.GetKey(17)
	pay_list := map[string]string{}
	pay_list["id"] = list_id
	pay_list["pay_code"] = pay_conf.Pay_code
	pay_list["pay_id"] = pay_id
	pay_list["mer_code"] = mer_code
	pay_list["amount"] = amount
	pay_list["real_amount"] = fmt.Sprintf("%.4f", real_amout)
	pay_list["create_time"] = now_time
	pay_list["order_number"] = list_id
	pay_list["class_code"] = class_code
	pay_list["bank_code"] = bank_code
	pay_list["push_url"] = ""
	pay_list["note"] = note
	pay_list["status"] = "1"
	pay_list["agent_path"] = mer_info.Agent_path
	pay_list["rate"] = fmt.Sprintf("%.4f", m_rate)
	pay_sql := common.InsertSql(table_name, pay_list)
	sql_arr := []string{pay_sql}
	err = model.Trans(sql_arr)
	if err != nil {
		t_msg = "充值失败"
		return t_status, t_msg
	}

	p_list := model.OrderById(list_id)

	if len(p_list.Id) < 1 {
		t_msg = "订单号不存在"
		return t_status, t_msg
	}

	t_status, t_msg = finished_pay(list_id, note, p_list)
	return t_status, t_msg
}
