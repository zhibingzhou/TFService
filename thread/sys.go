package thread

import (
	"fmt"
	"TFService/model"
	"strconv"
	"strings"

	"github.com/zhibingzhou/go_public/common"

	//"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

/**
*  用户所能看到的权限列表
 */
func PowerList(account string, ctx *gin.Context) (int, string, []map[string]interface{}) {
	t_status := 100
	t_msg := "该账号没有权限"
	power_list := []map[string]interface{}{}
	session := sessions.Default(ctx)
	admin := fmt.Sprintf("%v", session.Get("account"))
	power_path := fmt.Sprintf("%v", session.Get("power_path"))
	mer_code := fmt.Sprintf("%v", session.Get("mer_code"))
	p_list := []map[string]string{}
	offset := 0
	page_size := 1000
	p_where := map[string]interface{}{}
	in_code := []string{}
	if (account != "" && account != admin) || power_path != "all" {
		if account != "" && account != admin {
			user := model.AdminInfoRedis(account)
			if user["mer_code"] != mer_code && mer_code != "all" {
				return t_status, t_msg, power_list
			}
			p_where["account"] = account
		} else {
			p_where["account"] = admin
		}

		//查询管理员拥有的权限
		fields := []string{"power_code"}

		table_name := "admin_power"
		a_p, _ := model.PageList(table_name, "", page_size, offset, fields, p_where)
		if len(a_p) < 1 {
			t_status = 200
			t_msg = "success"
			return t_status, t_msg, power_list
		}
		for _, a_val := range a_p {
			in_code = append(in_code[0:], a_val["power_code"])
		}
	}

	t_status = 200
	t_msg = "success"

	p_table := "power_list"
	in_where := map[string]interface{}{}
	in_field := []string{"title", "code", "p_code", "power_type"}
	field_in := "code"

	if len(in_code) > 0 {
		p_list, _ = model.InPageList(p_table, field_in, page_size, offset, in_field, in_code, in_where)
	} else {
		p_list, _ = model.PageList(p_table, "", page_size, offset, in_field, in_where)
	}

	//查询权限值
	if len(p_list) < 1 {
		return t_status, t_msg, power_list
	}
	f_map := []map[string]string{}
	s_map := []map[string]string{}
	t_map := []map[string]string{}
	for _, p_map := range p_list {
		if p_map["power_type"] == "0" {
			f_map = append(f_map[0:], p_map)
		} else if p_map["power_type"] == "1" {
			s_map = append(s_map[0:], p_map)
		} else if p_map["power_type"] == "2" {
			t_map = append(t_map[0:], p_map)
		}
	}

	s_list := []map[string]interface{}{}

	if len(f_map) < 1 {
		return t_status, t_msg, power_list
	}

	for _, s_val := range s_map {
		l_map := map[string]interface{}{}
		l_map["name"] = s_val["title"]
		l_map["path"] = strings.Replace(s_val["code"], s_val["p_code"], "", 1)
		l_map["p_code"] = s_val["p_code"]
		v_map := []map[string]string{}
		for _, t_val := range t_map {
			if t_val["p_code"] == s_val["code"] {
				v_v_map := map[string]string{}
				v_v_map["name"] = t_val["title"]
				v_v_map["path"] = strings.Replace(t_val["code"], t_val["p_code"], "", 1)
				v_map = append(v_map[0:], v_v_map)
			}
		}
		l_map["children"] = v_map
		s_list = append(s_list[0:], l_map)
	}

	for _, f_val := range f_map {
		v_list := map[string]interface{}{}
		v_list["name"] = f_val["title"]
		v_list["path"] = f_val["code"]
		f_list := []map[string]interface{}{}
		for _, s_l_val := range s_list {
			if f_val["code"] == fmt.Sprintf("%v", s_l_val["p_code"]) {
				f_list = append(f_list[0:], s_l_val)
			}
		}
		v_list["children"] = f_list
		power_list = append(power_list[0:], v_list)
	}

	return t_status, t_msg, power_list
}

/**
*  新增权限
 */
func AddPower(code, title, url, p_code, power_type string) (int, string) {
	t_status := 100
	t_msg := "请填写完整"
	if title == "" || code == "" || power_type == "" {
		return t_status, t_msg
	}

	if power_type != "0" && p_code == "" {
		t_msg = "请填写上级权限的path"
		return t_status, t_msg
	}

	p_info := model.PowerInfoRedis(code)
	if len(p_info["code"]) > 1 {
		t_msg = "权限path已存在"
		return t_status, t_msg
	}

	table_name := "power_list"
	p_data := map[string]string{}
	p_data["code"] = p_code + code
	p_data["title"] = title
	p_data["url"] = url
	p_data["p_code"] = p_code
	p_data["power_type"] = power_type

	power_sql := common.InsertSql(table_name, p_data)
	err := model.Query(power_sql)
	if err != nil {
		t_msg = "权限新增失败"
		return t_status, t_msg
	}
	t_status = 200
	t_msg = "success"
	return t_status, t_msg
}

/**
*  新增商户
 */
func AddMer(p_agent string, m_map map[string]string, ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "请填写完整"
	if m_map["code"] == "" || m_map["title"] == "" || m_map["domain"] == "" || m_map["is_agent"] == "" {
		return t_status, t_msg
	}
	p_path := "all_"
	session := sessions.Default(ctx)

	admin_code := fmt.Sprintf("%v", session.Get("mer_code"))
	admin_mer := model.MerInfoRedis(admin_code)
	if admin_mer["code"] == "" {
		t_msg = "管理员信息异常"
		return t_status, t_msg
	}
	if admin_mer["is_agent"] != "1" {
		t_msg = "该账号不允许添加商户"
		return t_status, t_msg
	}
	if p_agent != "" {
		p_mer := model.MerInfoRedis(p_agent)
		if len(p_mer["code"]) < 1 {
			t_msg = "上级商户不存在"
			return t_status, t_msg
		}
		if p_mer["is_agent"] != "1" {
			t_msg = "该上级商户不属于代理"
			return t_status, t_msg
		}
		p_path = p_mer["agent_path"] + p_mer["code"] + "_"
	}
	mer_data := m_map
	mer_data["agent_path"] = p_path
	mer_data["private_key"] = common.Random("all", 8) + common.Random("all", 8)
	table_name := "mer_list"
	mer_sql := common.InsertSql(table_name, mer_data)
	err := model.Query(mer_sql)
	if err != nil {
		t_msg = "商户新增失败"
		return t_status, t_msg
	}
	t_status = 200
	t_msg = "success"
	return t_status, t_msg
}

/**
*  商户列表
 */
func MerList(is_under, page, page_size, Ismer_code string, ctx *gin.Context) (int, int, string, []map[string]string) {
	t_status := 100
	total := 0
	t_msg := "管理员信息异常"
	mer_list := []map[string]string{}
	fmer_list := []map[string]string{}
	session := sessions.Default(ctx)
	mer_code := fmt.Sprintf("%v", session.Get("mer_code"))
	power_path := fmt.Sprintf("%v", session.Get("power_path"))
	account := fmt.Sprintf("%v", session.Get("account"))

	t_status = 200
	t_msg = "success"
	table_name := "mer_list"
	p_where := map[string]interface{}{}
	count_field := "count(0) as num"
	fields := []string{"code", "title", "domain", "title", "qq", "skype", "telegram", "phone", "email", "private_key", "amount", "total_in", "total_out", "is_agent", "status"}

	//显示本商户
	user := model.AdminInfoRedis(account)
	if len(user["account"]) < 1 {
		t_msg = "账号错误"
		return t_status, total, t_msg, mer_list
	}
	p_where["code"] = user["mer_code"]
	fmer_list, _ = model.PageList(table_name, "", 1000, 0, fields, p_where)
	delete(p_where, "code")

	total, mer_list = GetMerInfromation(is_under, Ismer_code, mer_code, count_field, page, page_size, "", p_where, fields)

	mer_list = append(fmer_list, mer_list...)

	if len(mer_list) < 1 {
		return t_status, total, t_msg, mer_list
	}

	if power_path != "all" || mer_code != "all" {
		for m_k, _ := range mer_list {
			mer_list[m_k]["qq"] = "*****"
			mer_list[m_k]["qq"] = "*****"
			mer_list[m_k]["skype"] = "*****"
			mer_list[m_k]["telegram"] = "*****"
			mer_list[m_k]["phone"] = "*****"
			mer_list[m_k]["email"] = "*****"
			mer_list[m_k]["private_key"] = "*****"
		}
	}

	return t_status, total, t_msg, mer_list
}

/**
*  获取商户信息
 */
func GetMerInfromation(is_under, Ismer_code, mer_code, count_field, page, page_size, ismer_rate string, p_where map[string]interface{}, fields []string) (int, []map[string]string) {
	table_name := "mer_list"
	total := 0
	mer_list := []map[string]string{}
	admin := model.MerInfoRedis(mer_code)
	if len(admin["id"]) < 1 {
		return total, mer_list
	}

	page_int, size_int := ThreadPage(page, page_size)
	offset := (page_int - 1) * size_int

	if ismer_rate == "rate" {
		size_int = 1000
		offset = 0
	}

	if is_under == "1" {
		p_path := admin["agent_path"] + mer_code + "_"
		p_where["agent_path"] = p_path
		total, _ = model.ListTotal(table_name, count_field, p_where)
		if total < 1 {
			return total, mer_list
		}
		if Ismer_code != "" {
			mer_info := model.MerInfoRedis(Ismer_code)
			if strings.Contains(mer_info["agent_path"], p_path) == false {
				return total, mer_list
			}
			p_where["agent_path"] = mer_info["agent_path"] + mer_info["code"] + "_"
		}
		mer_list, _ = model.PageList(table_name, "", size_int, offset, fields, p_where)
	} else {
		lieke_field := "agent_path"
		like_where := admin["agent_path"] + mer_code + "_%"
		total, _ = model.LikeListTotal(table_name, lieke_field, like_where, count_field, p_where)
		if total < 1 {
			return total, mer_list
		}
		//商户号不为空，以商户号为搜索条件进行模糊搜索
		if Ismer_code != "" {
			lieke_field = "agent_path like ? and code like ?"
			like_first := admin["agent_path"] + mer_code + "_%"
			like_second := Ismer_code + "%"
			mer_list, _ = model.RLikePageList(table_name, lieke_field, like_first, like_second, size_int, offset, fields, p_where)
		} else {
			mer_list, _ = model.LikePageList(table_name, lieke_field, like_where, size_int, offset, fields, p_where)
		}
	}
	return total, mer_list
}

/**
*  修改商户信息
 */
func UpdateMer(mer_code, mer_status string, up_map map[string]string, ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "该账号不能修改商户信息"
	session := sessions.Default(ctx)
	admin_mer := fmt.Sprintf("%v", session.Get("mer_code"))
	if admin_mer != "all" {
		return t_status, t_msg
	}

	up_data := map[string]interface{}{}

	for u_k, u_v := range up_map {
		if u_v != "" {
			up_data[u_k] = u_v
		}
	}

	if mer_status != "" {
		up_data["status"], _ = strconv.Atoi(mer_status)
	}

	if len(up_data) < 1 {
		t_msg = "请填写需要修改的商户信息"
		return t_status, t_msg
	}

	mer_info := model.MerInfo(mer_code)
	if len(mer_info.Code) < 1 {
		t_msg = "商户号错误"
		return t_status, t_msg
	}

	err := model.UpdatesMer(mer_info, up_data)
	if err != nil {
		t_msg = "修改失败"
		return t_status, t_msg
	}
	t_status = 200
	t_msg = "success"
	return t_status, t_msg
}

func PayClass() (int, string, []map[string]string) {
	t_status := 200
	t_msg := "success"
	table_name := "pay_class"
	page_size := 100
	offset := 0
	fields := []string{"code", "title"}
	p_where := map[string]interface{}{}
	class_list, _ := model.PageList(table_name, "", page_size, offset, fields, p_where)
	return t_status, t_msg, class_list
}

/**
*  新增上游支付
 */
func AddPay(code, title, fee_amount, fee_type, is_push string, ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "请填写完整"
	if code == "" || title == "" {
		return t_status, t_msg
	}
	session := sessions.Default(ctx)
	admin_mer := fmt.Sprintf("%v", session.Get("mer_code"))
	if admin_mer != "all" {
		t_msg = "该账号没有权限"
		return t_status, t_msg
	}
	//判断支付code是否存在
	pay := model.PayInfoRedis(code)
	if len(pay["code"]) > 0 {
		t_msg = "支付编码已存在"
		return t_status, t_msg
	}

	data := map[string]string{}
	data["code"] = code
	data["title"] = title
	if fee_amount != "" {
		_, err := strconv.ParseFloat(fee_amount, 64)
		if err == nil {
			data["fee_amount"] = fee_amount
		}
	}
	if fee_type != "" {
		_, err := strconv.Atoi(fee_type)
		if err == nil {
			data["fee_type"] = fee_type
		}
	}
	data["is_push"] = "1"
	if is_push == "0" {
		data["is_push"] = "0"
	}
	table_name := "pay_channel"
	in_sql := common.InsertSql(table_name, data)
	err := model.Query(in_sql)
	if err != nil {
		t_msg = "新增失败"
		return t_status, t_msg
	}
	t_status = 200
	t_msg = "success"
	return t_status, t_msg
}

/**
*  新增上游渠道费率
 */
func AddPayClass(pay_code, class_code, bank_code, rate, min_amount, max_amount, limit_amount string, ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "请填写完整"
	if pay_code == "" || class_code == "" || rate == "" {
		return t_status, t_msg
	}
	session := sessions.Default(ctx)
	admin_mer := fmt.Sprintf("%v", session.Get("mer_code"))
	if admin_mer != "all" {
		t_msg = "该账号没有权限"
		return t_status, t_msg
	}
	var err error
	_, err = strconv.ParseFloat(rate, 64)
	if err != nil {
		t_msg = "费率错误"
		return t_status, t_msg
	}
	if min_amount != "" {
		_, err = strconv.ParseFloat(min_amount, 64)
		if err != nil {
			t_msg = "最小支付额度错误"
			return t_status, t_msg
		}
	} else {
		min_amount = "0.00"
	}

	if max_amount != "" {
		_, err = strconv.ParseFloat(max_amount, 64)
		if err != nil {
			t_msg = "最大支付额度错误"
			return t_status, t_msg
		}
	} else {
		max_amount = "0.00"
	}

	if limit_amount != "" {
		_, err = strconv.ParseFloat(limit_amount, 64)
		if err != nil {
			t_msg = "限制额度错误"
			return t_status, t_msg
		}
	} else {
		limit_amount = "0.00"
	}
	table_name := "pay_rate"
	p_data := map[string]string{}
	p_data["pay_code"] = pay_code
	p_data["class_code"] = class_code
	p_data["bank_code"] = bank_code
	p_data["rate"] = rate
	p_data["min_amount"] = min_amount
	p_data["max_amount"] = max_amount
	p_data["limit_amount"] = limit_amount
	in_sql := common.InsertSql(table_name, p_data)

	m_table := "mer_rate"
	m_data := map[string]string{}
	m_data["pay_code"] = pay_code
	m_data["class_code"] = class_code
	m_data["rate"] = rate
	m_data["mer_code"] = "all"
	m_data["bank_code"] = bank_code
	m_sql := common.InsertSql(m_table, m_data)

	sql_arr := []string{in_sql, m_sql}
	err = model.Trans(sql_arr)

	if err != nil {
		t_msg = "添加失败"
		return t_status, t_msg
	}
	t_status = 200
	t_msg = "success"
	return t_status, t_msg
}

/**
*  支付渠道列表
 */
func ChannelList(class_code, page, page_size string, ctx *gin.Context) (int, int, string, []map[string]string) {
	t_status := 200
	t_msg := "success"
	total := 0
	c_list := []map[string]string{}
	session := sessions.Default(ctx)
	admin_mer := fmt.Sprintf("%v", session.Get("mer_code"))

	count_field := "count(0) as num"
	table_name := "pay_channel"
	p_where := map[string]interface{}{}
	in_field := "code"
	in_where := []string{}
	if class_code != "" || admin_mer != "all" {
		r_where := map[string]interface{}{}
		if class_code != "" {
			r_where["class_code"] = class_code
		}
		if admin_mer != "all" {
			r_where["mer_code"] = admin_mer
		}

		r_field := []string{"pay_code"}
		r_list, _ := model.PageList("mer_rate", "", 1000, 0, r_field, r_where)
		if len(r_list) < 1 {
			return t_status, total, t_msg, c_list
		}
		for _, r_v := range r_list {
			in_where = append(in_where[0:], r_v["pay_code"])
		}
		total, _ = model.InListTotal(table_name, in_field, count_field, in_where, p_where)
	} else {
		total, _ = model.ListTotal(table_name, count_field, p_where)
	}

	if total < 1 {
		return t_status, total, t_msg, c_list
	}

	fields := []string{"code", "title", "fee_amount", "fee_type", "is_push"}

	page_int, size_int := ThreadPage(page, page_size)
	offset := (page_int - 1) * size_int
	if len(in_where) > 0 {
		c_list, _ = model.InPageList(table_name, in_field, size_int, offset, fields, in_where, p_where)
	} else {
		c_list, _ = model.PageList(table_name, "", size_int, offset, fields, p_where)
	}

	return t_status, total, t_msg, c_list
}

/**
*  查询支付渠道详情
 */
func PayDetail(pay_code, class_code, page, page_size string, ctx *gin.Context) (int, int, string, []map[string]string) {
	t_status := 200
	t_msg := "success"
	total := 0
	c_list := []map[string]string{}
	session := sessions.Default(ctx)
	admin_mer := fmt.Sprintf("%v", session.Get("mer_code"))
	if admin_mer != "all" {
		t_msg = "该账号没有权限"
		t_status = 100
		return t_status, total, t_msg, c_list
	}

	table_name := "pay_rate"
	r_where := map[string]interface{}{}
	if pay_code != "" {
		r_where["pay_code"] = pay_code
	}
	if class_code != "" {
		r_where["class_code"] = class_code
	}
	count_field := "count(0) as num"
	total, _ = model.ListTotal(table_name, count_field, r_where)
	if total < 1 {
		return t_status, total, t_msg, c_list
	}
	page_int, size_int := ThreadPage(page, page_size)
	offset := (page_int - 1) * size_int

	r_field := []string{"id", "pay_code", "class_code", "rate", "min_amount", "max_amount", "limit_amount", "day_amount", "bank_code"}
	c_list, _ = model.PageList(table_name, "", size_int, offset, r_field, r_where)
	if len(c_list) < 1 {
		return t_status, total, t_msg, c_list
	}
	for r_k, r_v := range c_list {
		pay_i := model.PayInfoRedis(r_v["pay_code"])
		c_list[r_k]["pay_title"] = pay_i["title"]
		class_i := model.ClassInfoRedis(r_v["class_code"])
		c_list[r_k]["class_title"] = class_i["title"]
		sys := model.BankInfoRedis(r_v["bank_code"])
		c_list[r_k]["class_title"] = sys["title"]
	}
	return t_status, total, t_msg, c_list
}

/**
*  修改详情
 */
func UpdatePay(pay_id, rate, min_amount, max_amount string, ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "该账号没有权限"
	session := sessions.Default(ctx)
	admin_mer := fmt.Sprintf("%v", session.Get("mer_code"))
	if admin_mer != "all" {
		t_msg = "该账号没有权限"
		t_status = 100
		return t_status, t_msg
	}
	p_rate := model.PayRateById(pay_id)
	if p_rate.Id < 1 {
		t_msg = "费率ID错误"
		return t_status, t_msg
	}
	p_data := map[string]interface{}{}

	if rate != "" {
		rate_f, err := strconv.ParseFloat(rate, 64)
		if err != nil {
			t_msg = "费率错误"
			return t_status, t_msg
		}
		p_data["rate"] = rate_f
	}

	if min_amount != "" {
		min_f, err := strconv.ParseFloat(min_amount, 64)
		if err != nil {
			t_msg = "最小支付额度错误"
			return t_status, t_msg
		}
		p_data["min_amount"] = min_f
	}

	if max_amount != "" {
		max_f, err := strconv.ParseFloat(max_amount, 64)
		if err != nil {
			t_msg = "最大支付额度错误"
			return t_status, t_msg
		}
		p_data["max_amount"] = max_f
	}
	if len(p_data) < 1 {
		t_msg = "没有修改"
		return t_status, t_msg
	}

	err := model.UpdatesPayRate(p_rate, p_data)
	if err != nil {
		t_msg = "修改失败"
		return t_status, t_msg
	}
	if rate != "" {
		pay_sql := fmt.Sprintf("update mer_rate set rate='%s' where mer_code='all' and pay_code='%s' and class_code='%s' and bank_code ='%s';", rate, p_rate.Pay_code, p_rate.Class_code, p_rate.Bank_code)
		model.Query(pay_sql)
	}
	t_status = 200
	t_msg = "success"
	return t_status, t_msg
}

/**
*  新增/修改支付费率
 */
func MerRate(rate_id, mer_code, pay_code, class_code, bank_code, rate string, ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "该账号没有权限"
	if rate == "" {
		t_msg = "费率必填"
		return t_status, t_msg
	}
	rate_f, err := strconv.ParseFloat(rate, 64)
	if err != nil {
		t_msg = "费率格式错误"
		return t_status, t_msg
	}
	session := sessions.Default(ctx)
	admin_mer := fmt.Sprintf("%v", session.Get("mer_code"))
	if admin_mer != "all" && mer_code == admin_mer {
		return t_status, t_msg
	}

	admin := model.MerInfoRedis(admin_mer)
	if len(admin["code"]) < 1 {
		t_msg = "管理员信息异常"
		return t_status, t_msg
	}

	is_update := false

	//更新
	if rate_id != "" {
		m_rate := model.MerRateByIdRedis(rate_id)
		if len(m_rate["id"]) < 1 {
			t_msg = "费率ID错误"
			return t_status, t_msg
		}
		is_update = true
	}
	//查询上级费率
	p_rate := model.MerRateInfo(admin["code"], pay_code, class_code, bank_code)
	if p_rate.Rate > rate_f {
		t_msg = "商户费率不能低于上级费率"
		return t_status, t_msg
	}

	rate_mer := model.MerInfoRedis(mer_code)
	if len(rate_mer["code"]) < 1 {
		t_msg = "费率信息异常"
		return t_status, t_msg
	}
	if admin["agent_path"]+admin["code"]+"_" != rate_mer["agent_path"] && admin_mer != "all" {
		return t_status, t_msg
	}
	sql := ""
	table_name := "mer_rate"
	if is_update {
		sql = fmt.Sprintf("update %s set rate='%s',pay_code='%s',bank_code='%s',class_code='%s' where id='%s';", table_name, rate, pay_code, bank_code, class_code, rate_id)
	} else {
		r_data := map[string]string{}
		r_data["mer_code"] = mer_code
		r_data["class_code"] = class_code
		r_data["pay_code"] = pay_code
		r_data["rate"] = rate
		r_data["bank_code"] = bank_code
		sql = common.InsertSql(table_name, r_data)
	}
	err = model.Query(sql)
	if err != nil {
		t_msg = "设置失败"
		return t_status, t_msg
	}
	t_status = 200
	t_msg = "success"
	return t_status, t_msg
}

/**
*  商户及商户下属的费率详情
 */
func MerRateList(Ismer_code, pay_code, class_code, is_under, page, page_size string, ctx *gin.Context) (int, int, string, []map[string]string) {
	t_status := 200
	t_msg := "success"
	total := 0
	r_list := []map[string]string{}
	fr_list := []map[string]string{}
	session := sessions.Default(ctx)
	mer_code := fmt.Sprintf("%v", session.Get("mer_code"))
	account := fmt.Sprintf("%v", session.Get("account"))
	p_where := map[string]interface{}{}
	m_p_where := map[string]interface{}{}
	if pay_code != "" {
		p_where["pay_code"] = pay_code
	}
	if class_code != "" {
		p_where["class_code"] = class_code
	}

	admin := model.MerInfoRedis(mer_code)
	if len(admin["id"]) < 1 {
		return t_status, total, t_msg, r_list
	}
	mer_list := []map[string]string{}
	count_field := "count(0) as num"
	table_name := "mer_rate"
	page_int, size_int := ThreadPage(page, page_size)
	offset := (page_int - 1) * size_int
	fields := []string{"id", "pay_code", "class_code", "mer_code", "rate", "bank_code"}
	mer_fields := []string{"code"}
	//显示本商户费率
	user := model.AdminInfoRedis(account)
	if len(user["account"]) < 1 {
		t_msg = "账号错误"
		return t_status, total, t_msg, mer_list
	}
	if user["mer_code"] == Ismer_code {
		p_where["mer_code"] = user["mer_code"]
		fr_list, _ = model.PageList(table_name, "", size_int, offset, fields, p_where)
		delete(p_where, "mer_code")
	}
	total, mer_list = GetMerInfromation(is_under, Ismer_code, mer_code, count_field, page, page_size, "rate", m_p_where, mer_fields)
	in_field := "mer_code"
	in_where := []string{}
	for _, m_v := range mer_list {
		in_where = append(in_where[0:], m_v["code"])
	}
	total, _ = model.InListTotal(table_name, in_field, count_field, in_where, p_where)
	if total > 0 {
		r_list, _ = model.InPageList(table_name, in_field, size_int, offset, fields, in_where, p_where)
	}
	r_list = append(fr_list, r_list...)
	if len(r_list) < 1 {
		return t_status, total, t_msg, r_list
	}
	for r_k, r_v := range r_list {
		pay_i := model.PayInfoRedis(r_v["pay_code"])
		r_list[r_k]["pay_title"] = pay_i["title"]
		class_i := model.ClassInfoRedis(r_v["class_code"])
		r_list[r_k]["class_title"] = class_i["title"]
	}
	return t_status, total, t_msg, r_list
}

/**
*  删除商户费率
 */
func DelRate(rate_id, mer_code string, ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "该账号没有权限"

	session := sessions.Default(ctx)
	admin_mer := fmt.Sprintf("%v", session.Get("mer_code"))

	if admin_mer != "all" && mer_code == admin_mer {
		return t_status, t_msg
	}

	admin := model.MerInfoRedis(admin_mer)
	if len(admin["code"]) < 1 {
		t_msg = "管理员信息异常"
		return t_status, t_msg
	}

	m_rate := model.MerRateByIdRedis(rate_id)
	if len(m_rate["id"]) < 1 {
		t_msg = "费率ID错误"
		return t_status, t_msg
	}

	rate_mer := model.MerInfoRedis(mer_code)
	if len(rate_mer["code"]) < 1 {
		t_msg = "费率信息异常"
		return t_status, t_msg
	}

	if admin["agent_path"]+admin["code"]+"_" != rate_mer["agent_path"] && admin_mer != "all" {
		return t_status, t_msg
	}

	table_name := "mer_rate"
	sql := fmt.Sprintf("delete from %s where id='%s';", table_name, rate_id)
	err := model.Query(sql)
	if err != nil {
		t_msg = "删除失败"
		return t_status, t_msg
	}

	t_status = 200
	t_msg = "success"
	return t_status, t_msg
}

/**
*  白名单列表
 */
func IpList(mer_code, ip, page, page_size string, ctx *gin.Context) (int, int, string, []map[string]string) {
	t_status := 100
	t_msg := "该账号没有权限"
	total := 0
	r_list := []map[string]string{}
	session := sessions.Default(ctx)
	admin_mer := fmt.Sprintf("%v", session.Get("mer_code"))
	if admin_mer != "all" {
		return t_status, total, t_msg, r_list
	}
	p_where := map[string]interface{}{}
	if mer_code != "" {
		p_where["mer_code"] = mer_code
	}
	if ip != "" {
		p_where["ip"] = ip
	}
	table_name := "mer_ip"
	t_status = 200
	t_msg = "success"
	page_int, size_int := ThreadPage(page, page_size)
	offset := (page_int - 1) * size_int
	fields := []string{"mer_code", "ip"}
	count_field := "count(0) as num"
	total, _ = model.ListTotal(table_name, count_field, p_where)
	if total < 1 {
		return t_status, total, t_msg, r_list
	}
	r_list, _ = model.PageList(table_name, "", size_int, offset, fields, p_where)
	return t_status, total, t_msg, r_list
}

/**
*  新增IP白名单
 */
func AddIp(mer_code, ip string, ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "该账号没有权限"
	session := sessions.Default(ctx)
	admin_mer := fmt.Sprintf("%v", session.Get("mer_code"))
	if admin_mer != "all" {
		return t_status, t_msg
	}
	if ip == "" || mer_code == "" {
		t_msg = "请填写完整"
		return t_status, t_msg
	}
	is_exist := model.MerIpAuthRedis(mer_code, ip)
	if len(is_exist["mer_code"]) > 0 {
		t_msg = "iP已存在,请勿重复添加"
		return t_status, t_msg
	}
	table_name := "mer_ip"
	ip_data := map[string]string{}
	ip_data["ip"] = ip
	ip_data["mer_code"] = mer_code
	ip_sql := common.InsertSql(table_name, ip_data)
	err := model.Query(ip_sql)
	if err != nil {
		t_msg = "白名单添加失败"
		return t_status, t_msg
	}
	t_status = 200
	t_msg = "success"
	return t_status, t_msg
}

/**
*  删除IP白名单
 */
func DelIp(ip string, ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "该账号没有权限"
	session := sessions.Default(ctx)
	admin_mer := fmt.Sprintf("%v", session.Get("mer_code"))
	if admin_mer != "all" {
		return t_status, t_msg
	}
	is_exist := model.MerIpByIpRedis(ip)
	if len(is_exist["mer_code"]) < 1 {
		t_msg = "ip不存在"
		return t_status, t_msg
	}
	table_name := "mer_ip"
	ip_sql := fmt.Sprintf("delete from %s where ip='%s';", table_name, ip)
	err := model.Query(ip_sql)
	if err != nil {
		t_msg = "白名单删除失败"
		return t_status, t_msg
	}
	t_status = 200
	t_msg = "success"
	return t_status, t_msg
}

func NoticeList(is_all, page, page_size string) (int, int, string, []map[string]string) {
	t_status := 200
	t_msg := "success"
	n_list := []map[string]string{}
	page_int, size_int := ThreadPage(page, page_size)
	offset := (page_int - 1) * size_int
	fields := []string{"id", "title", "content", "create_time"}
	p_where := map[string]interface{}{}
	if is_all == "0" {
		p_where["status"] = 1
	}

	order_by := "create_time desc"
	table_name := "note_list"
	total_field := "count(0) as num"
	total, _ := model.ListTotal(table_name, total_field, p_where)
	if total < 1 {
		return t_status, total, t_msg, n_list
	}
	n_list, _ = model.PageList(table_name, order_by, size_int, offset, fields, p_where)
	return t_status, total, t_msg, n_list
}

func AddNotice(title, content string, ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "请填写完整"
	if title == "" || content == "" {
		return t_status, t_msg
	}
	session := sessions.Default(ctx)
	admin_mer := fmt.Sprintf("%v", session.Get("mer_code"))
	if admin_mer != "all" {
		t_msg = "该账号没有权限"
		return t_status, t_msg
	}
	table_name := "note_list"
	n_data := map[string]string{}
	n_data["id"] = model.GetKey(16)
	n_data["title"] = title
	n_data["content"] = content
	n_sql := common.InsertSql(table_name, n_data)
	err := model.Query(n_sql)
	if err != nil {
		t_msg = "新增失败"
		return t_status, t_msg
	}
	t_status = 200
	t_msg = "success"
	return t_status, t_msg
}

func UpdateNotice(n_id string, ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "请填写完整"
	if n_id == "" {
		return t_status, t_msg
	}

	session := sessions.Default(ctx)
	admin_mer := fmt.Sprintf("%v", session.Get("mer_code"))
	if admin_mer != "all" {
		t_msg = "该账号没有权限"
		return t_status, t_msg
	}

	n_sql := fmt.Sprintf("update note_list set `status`='%s' where id='%s';", "0", n_id)
	err := model.Query(n_sql)
	if err != nil {
		t_msg = "修改失败"
		return t_status, t_msg
	}
	t_status = 200
	t_msg = "success"
	return t_status, t_msg
}

/**
*  系统银行列表
 */
func SysBank() (int, string, []map[string]string) {
	t_status := 200
	t_msg := "success"

	p_w := map[string]interface{}{}
	table_name := "sys_bank"
	field := []string{"code", "title"}
	sys_bank, _ := model.PageList(table_name, "", 1000, 0, field, p_w)
	return t_status, t_msg, sys_bank
}

func PayConf(pay_code, pay_status, page, page_size string) (int, int, string, []map[string]string) {
	t_status := 200
	t_msg := "success"
	p_list := []map[string]string{}
	page_int, size_int := ThreadPage(page, page_size)
	offset := (page_int - 1) * size_int
	fields := []string{"id", "pay_code", "note", "status", "merchant_code", "amount", "total_in", "total_out"}
	p_where := map[string]interface{}{}
	if pay_code != "" {
		p_where["pay_code"] = pay_code
	}
	if pay_status != "" {
		p_where["status"] = pay_status
	}

	order_by := "id desc"
	table_name := "pay_config"
	total_field := "count(0) as num"
	total, _ := model.ListTotal(table_name, total_field, p_where)
	if total < 1 {
		return t_status, total, t_msg, p_list
	}
	p_list, _ = model.PageList(table_name, order_by, size_int, offset, fields, p_where)
	if len(p_list) < 1 {
		return t_status, total, t_msg, p_list
	}
	for p_k, p_val := range p_list {
		pay := model.PayInfoRedis(p_val["pay_code"])
		p_list[p_k]["title"] = pay["title"]
	}
	return t_status, total, t_msg, p_list
}

func PayBank(pay_code, class_code, page, page_size string) (int, int, string, []map[string]string) {
	t_status := 200
	t_msg := "success"
	p_list := []map[string]string{}
	page_int, size_int := ThreadPage(page, page_size)
	offset := (page_int - 1) * size_int

	p_where := map[string]interface{}{}
	if pay_code != "" {
		p_where["pay_code"] = pay_code
	}
	if class_code != "" {
		p_where["class_code"] = class_code
	}

	order_by := "id desc"
	table_name := "pay_bank"
	total_field := "count(0) as num"
	total, _ := model.ListTotal(table_name, total_field, p_where)
	if total < 1 {
		return t_status, total, t_msg, p_list
	}

	fields := []string{"pay_code", "is_mobile", "class_code", "bank_code", "bank_title", "jump_type", "pay_bank"}
	p_list, _ = model.PageList(table_name, order_by, size_int, offset, fields, p_where)
	if len(p_list) < 1 {
		return t_status, total, t_msg, p_list
	}
	for p_k, p_val := range p_list {
		pay := model.PayInfoRedis(p_val["pay_code"])
		p_list[p_k]["pay_title"] = pay["title"]
		class := model.ClassInfoRedis(p_val["class_code"])
		p_list[p_k]["class_title"] = class["title"]
	}
	return t_status, total, t_msg, p_list
}

func AddPayBank(p_data map[string]string) (int, string) {
	t_status := 100
	t_msg := "请填写完整"
	if len(p_data) < 1 {
		return t_status, t_msg
	}
	for _, p_val := range p_data {
		if p_val == "" {
			return t_status, t_msg
		}
	}
	sys_bank := model.BankInfoRedis(p_data["bank_code"])
	p_data["bank_title"] = sys_bank["title"]
	table_name := "pay_bank"
	n_sql := common.InsertSql(table_name, p_data)
	err := model.Query(n_sql)
	if err != nil {
		t_msg = "新增失败"
		return t_status, t_msg
	}
	t_status = 200
	t_msg = "success"
	return t_status, t_msg
}

/**
*  查询商户支付渠道详情
 */
func PayMerDetail(is_under, page, page_size string, ctx *gin.Context) (int, int, string, []map[string]string) {
	t_status := 200
	t_msg := "success"
	total := 0
	c_list := []map[string]string{}
	r_list := []map[string]string{}
	session := sessions.Default(ctx)
	admin_mer := fmt.Sprintf("%v", session.Get("mer_code"))
	admin := model.MerInfoRedis(admin_mer)
	if len(admin["code"]) < 1 {
		t_msg = "管理员信息异常"
		return t_status, total, t_msg, c_list
	}
	p_where := map[string]interface{}{}
	mer_list := []map[string]string{}
	count_field := "count(0) as num"
	table_name := "mer_rate"
	page_int, size_int := ThreadPage(page, page_size)
	offset := (page_int - 1) * size_int
	fields := []string{"id", "pay_code", "class_code", "mer_code", "rate", "bank_code"}
	//当前商户
	p_where["mer_code"] = admin["code"]
	total, _ = model.ListTotal(table_name, count_field, p_where)
	if total < 1 {
		return t_status, total, t_msg, c_list
	}
	mer_list, _ = model.PageList(table_name, "", size_int, offset, fields, p_where)

	table_name = "pay_rate"
	r_where := map[string]interface{}{}
	r_field := []string{"id", "pay_code", "class_code", "rate", "min_amount", "max_amount", "limit_amount", "day_amount", "bank_code"}
	for r_k, r_v := range mer_list {

		r_where["pay_code"] = mer_list[r_k]["pay_code"]
		r_where["class_code"] = mer_list[r_k]["class_code"]
		r_where["bank_code"] = mer_list[r_k]["bank_code"]

		totals, _ := model.ListTotal(table_name, count_field, r_where)
		if totals < 1 {
			return t_status, totals, t_msg, c_list
		}
		c_list, _ = model.PageList(table_name, "", 1, 0, r_field, r_where)
		class_i := model.ClassInfoRedis(r_v["class_code"])
		sys := model.BankInfoRedis(r_v["bank_code"])
		c_list[0]["class_title"] = class_i["title"] + "_" + sys["title"] //支付宝_个人扫码
		c_list[0]["rate"] = r_v["rate"]
		r_list = append(c_list, r_list...)
	}
	return t_status, total, t_msg, r_list
}
