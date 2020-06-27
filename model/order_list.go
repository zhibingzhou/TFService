package model

/**
*  根据系统订单查询
 */
func CashOrderById(order_number string) OrderList {
	var p_list OrderList
	gdb.DB.Model(&OrderList{}).Where("id = ?", order_number).First(&p_list)
	return p_list
}

/**
*  根据web订单查询
 */
func CashOrderByweb(order_number string) OrderList {
	var p_list OrderList
	gdb.DB.Model(&OrderList{}).Where("order_number = ?", order_number).First(&p_list)
	return p_list
}

func UpdatesOrderList(p_list OrderList, p_data map[string]interface{}) error {
	res := gdb.DB.Model(&p_list).UpdateColumns(p_data)
	return res.Error
}
