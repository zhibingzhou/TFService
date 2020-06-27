package model

import (
	"database/sql"
	"fmt"
)

func PageList(table_name, order_by string, page_size, offset int, fields []string, p_where map[string]interface{}) ([]map[string]string, error) {
	records := []map[string]string{}
	if order_by == "" {
		order_by = fields[0] + " desc"
	}
	u_rows, err := gdb.DB.Table(table_name).Select(fields).Where(p_where).Limit(page_size).Order(order_by).Offset(offset).Rows()
	if err != nil {
		return records, err
	}
	records, err = rows2Map(fields, u_rows)
	return records, err
}

func InPageList(table_name, in_field string, page_size, offset int, fields, in_where []string, p_where map[string]interface{}) ([]map[string]string, error) {
	records := []map[string]string{}
	order_by := fields[0] + " desc"
	u_rows, err := gdb.DB.Table(table_name).Select(fields).Where(p_where).Where(in_field+" in (?)", in_where).Limit(page_size).Offset(offset).Order(order_by).
		Rows()
	if err != nil {
		return records, err
	}
	records, err = rows2Map(fields, u_rows)
	return records, err
}

func LikePageList(table_name, like_field, like_where string, page_size, offset int, fields []string, p_where map[string]interface{}) ([]map[string]string, error) {
	records := []map[string]string{}

	u_rows, err := gdb.DB.Table(table_name).Select(fields).Where(like_field+" like ?", like_where).Where(p_where).Limit(page_size).Offset(offset).Rows()
	if err != nil {
		return records, err
	}
	records, err = rows2Map(fields, u_rows)
	return records, err
}

func RLikePageList(table_name, like_field, like_first, like_second string, page_size, offset int, fields []string, p_where map[string]interface{}) ([]map[string]string, error) {
	records := []map[string]string{}

	u_rows, err := gdb.DB.Table(table_name).Select(fields).Where(like_field, like_first, like_second).Where(p_where).Limit(page_size).Offset(offset).Rows()
	if err != nil {
		return records, err
	}
	records, err = rows2Map(fields, u_rows)
	return records, err
}

func LikeListTotal(table_name, like_field, like_where, field string, p_where map[string]interface{}) (int, float64) {
	var c_total CountTotal
	gdb.DB.Table(table_name).Select(field).Where(like_field+" like ?", like_where).Where(p_where).Scan(&c_total)

	return c_total.Num, c_total.Total
}

func ListTotal(table_name, field string, p_where map[string]interface{}) (int, float64) {
	var c_total CountTotal
	gdb.DB.Table(table_name).Select(field).Where(p_where).Scan(&c_total)

	return c_total.Num, c_total.Total
}

/**
*  查询
 */
func CommonTotal(table_name, field string, p_w map[string]interface{}) (int, float64) {
	var c_total CountTotal
	gdb.DB.Table(table_name).Select(field).Where(p_w).Scan(&c_total)

	return c_total.Num, c_total.Total
}

/**
*  查询
 */
func CommonDateTotal(table_name, field, s_time, e_time, time_field string, p_w map[string]interface{}) (int, float64) {
	var c_total CountTotal
	gdb.DB.Table(table_name).Select(field).Where(time_field+">=? and "+time_field+"<?", s_time, e_time).Where(p_w).Scan(&c_total)

	return c_total.Num, c_total.Total
}

func InListTotal(table_name, in_field, field string, in_where []string, p_where map[string]interface{}) (int, float64) {
	var c_total CountTotal
	gdb.DB.Table(table_name).Select(field).Where(p_where).Where(in_field+" in (?)", in_where).Scan(&c_total)

	return c_total.Num, c_total.Total
}

/**
*  日期分页
 */
func PageDateList(table_name, date_field, s_date, e_date, like_sql string, page_size, offset int, fields []string,
	p_where map[string]interface{}) (
	[]map[string]string,
	error) {
	records := []map[string]string{}

	u_rows, err := gdb.DB.Table(table_name).Select(fields).Where(p_where).Where(
		date_field+">=? and "+date_field+"<=? "+like_sql, s_date, e_date).
		Limit(page_size).Order(date_field + " desc").
		Offset(offset).Rows()
	if err != nil {
		return records, err
	}
	records, err = rows2Map(fields, u_rows)
	return records, err
}

/**
*  日期分页 无耐之举再加个方法
 */
func SecondPageDateList(table_name, date_field, s_date, e_date, like_sql, field_in string, typevalue []string, page_size, offset int, fields []string,
	p_where map[string]interface{}) (
	[]map[string]string,
	error) {
	records := []map[string]string{}

	u_rows, err := gdb.DB.Table(table_name).Select(fields).Where(p_where).Where(
		date_field+">=? and "+date_field+"<=? "+like_sql, s_date, e_date).Where(field_in, typevalue).
		Limit(page_size).Order(date_field + " desc").
		Offset(offset).Rows()
	if err != nil {
		return records, err
	}
	records, err = rows2Map(fields, u_rows)

	return records, err
}

func DateListTotal(table_name, date_field, s_date, e_date, like_sql, field string, p_where map[string]interface{}) (int, float64) {
	var c_total CountTotal

	gdb.DB.Table(table_name).Select(field).Where(p_where).Where(date_field+">=? and "+date_field+"<=? "+like_sql, s_date, e_date).Scan(&c_total)

	return c_total.Num, c_total.Total
}

func CommonFieldsRow(table_name string, fields []string, c_w map[string]interface{}) (map[string]string, error) {
	u_row := gdb.DB.Table(table_name).Select(fields).Where(c_w).Row()

	record, err := row2Map(fields, u_row)

	return record, err
}

/*
数据库数据转MAP
*/
func row2Map(fields []string, u_row *sql.Row) (map[string]string, error) {
	record := map[string]string{}

	//创建有效切片
	values := make([]interface{}, len(fields))
	//行扫描，必须复制到这样切片的内存地址中去
	scanArgs := make([]interface{}, len(fields))
	for j := range values {
		scanArgs[j] = &values[j]
	}

	err := u_row.Scan(scanArgs...)

	for i, col := range values {
		if col == nil {
			record[fields[i]] = ""
			continue
		}
		col_s, ok := col.([]byte)
		if ok {
			record[fields[i]] = string(col_s)
		} else {
			record[fields[i]] = fmt.Sprintf("%v", col)
		}
	}

	return record, err
}

/*
数据库数据转MAP
*/
func rows2Map(fields []string, u_rows *sql.Rows) ([]map[string]string, error) {

	records := []map[string]string{}
	var err error

	//创建有效切片
	values := make([]interface{}, len(fields))
	//行扫描，必须复制到这样切片的内存地址中去
	scanArgs := make([]interface{}, len(fields))

	for j := range values {
		scanArgs[j] = &values[j]
	}

	for u_rows.Next() {
		err = u_rows.Scan(scanArgs...)

		if err != nil {
			break
		}
		record := map[string]string{}
		for i, col := range values {
			if col == nil {
				record[fields[i]] = ""
				continue
			}
			col_s, ok := col.([]byte)
			if ok {
				record[fields[i]] = string(col_s)
			} else {
				record[fields[i]] = fmt.Sprintf("%v", col)
			}
		}
		records = append(records[0:], record)
	}
	return records, err
}
