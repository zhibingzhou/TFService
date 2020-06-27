package thread

import (
	"fmt"
	"TFService/hook"
	"TFService/model"
	"strconv"
	"strings"
	"time"

	"github.com/zhibingzhou/go_public/common"

	"github.com/zhibingzhou/go_public/redis"
)

func init() {
}

/**
* 订单推送的进程
* @param	access_code	接入商编号
* @param	web_ordernumber	网站订单号
* @param	ordernumber	系统订单号
* @param	amout			订单金额
* @param	status			订单的状态
 */
func Push(p_list model.PayList) {
	if len(p_list.Push_url) < 1 {
		return
	}

	var result = fmt.Sprintf(`{"order_number":"%s","pay_id":"%d","class_code":"%s","pay_order":"%s","amount":"%.2f","status":"%d","real_amount":"%.2f"}`, p_list.Order_number, p_list.Pay_id, p_list.Class_code, p_list.Id, p_list.Amount, p_list.Status, p_list.Real_amount)
	//第三部，加密需要的数据
	aes_res := hook.HookAesEncrypt(p_list.Mer_code, result)
	aes_res = strings.Replace(aes_res, "+", "%2B", -1)
	strResult := fmt.Sprintf("mer_code=%s&order_number=%s&result=%s", p_list.Mer_code, p_list.Order_number, aes_res)
	//写入redis锁(15秒)，来排除并发的可能
	push_key := "lock:push_" + p_list.Order_number
	rd_lock := redis.RediGo.StringWriteNx(push_key, p_list.Order_number, 30)
	if rd_lock == 0 {
		//重复提单
		return
	} else {
		pushPool(strResult, p_list.Push_url, p_list)
	}
}

/**
* 结果推送，并记录结果
* @param	str	需要推送的内容
* @param	ordernumber	系统订单
 */
func pushPool(str, push_url string, p_list model.PayList) {
	//将推送的结果写入数据库
	p_data := map[string]interface{}{}
	push_status := 1
	status, msg_b := common.HttpBody(push_url, "POST", str, push_header)
	msg := string(msg_b)
	p_num := p_list.Push_num + 1
	common.LogsWithFileName(log_path, "push_", "ordernumber->"+p_list.Order_number+"\nmsg->"+msg+" \n"+str+"\npush_url->"+push_url)
	if status == 200 && msg == "success" {
		push_status = 3
		p_data["push_status"] = push_status
		p_data["push_num"] = p_num
		model.UpdatesPayList(p_list, p_data)
		return
	}

	push_status = 9

	//继续推送,连续10秒,推送12次，如果都失败了，放弃任务，并记录
	for i := 1; i < 10; i++ {
		p_num = p_num + 1
		time.Sleep(30 * time.Second)
		status, msg_b = common.HttpBody(push_url, "POST", str, push_header)
		common.LogsWithFileName(log_path, "push_", "ordernumber->"+p_list.Order_number+"\ntimes->"+strconv.Itoa(i)+"\nrequest->"+str+"\nmsg->"+msg+"\npush_url->"+push_url)
		msg = string(msg_b)
		if status == 200 && msg == "success" {
			//成功后跳出循环
			push_status = 3
			break
		}
	}
	p_data["push_status"] = push_status
	p_data["push_num"] = p_num
	model.UpdatesPayList(p_list, p_data)
}

func PushCash(p_list model.CashList) {
	mer_info := model.MerInfoRedis(p_list.Mer_code)
	if len(mer_info["code"]) < 1 {
		return
	}

	if len(p_list.Push_url) < 1 {
		return
	}

	var result = fmt.Sprintf(`{"order_number":"%s","pay_id":"%d","pay_order":"%s","amount":"%.2f","status":"%d","real_amount":"%.2f"}`, p_list.Order_number, p_list.Pay_id, p_list.Id, p_list.Amount, p_list.Status, p_list.Real_amount)
	//第三部，加密需要的数据
	aes_res := hook.HookAesEncrypt(mer_info["code"], result)
	aes_res = strings.Replace(aes_res, "+", "%2B", -1)
	strResult := fmt.Sprintf("mer_code=%s&order_number=%s&result=%s", p_list.Mer_code, p_list.Order_number, aes_res)
	//写入redis锁(15秒)，来排除并发的可能
	push_key := "lock:push_" + p_list.Order_number
	rd_lock := redis.RediGo.StringWriteNx(push_key, p_list.Order_number, 30)
	if rd_lock == 0 {
		//重复提单
		return
	} else {
		pushCashPool(strResult, p_list.Push_url, p_list)
	}
}

/**
* 结果推送，并记录结果
* @param	str	需要推送的内容
* @param	ordernumber	系统订单
 */
func pushCashPool(str, push_url string, p_list model.CashList) {
	//将推送的结果写入数据库
	p_data := map[string]interface{}{}
	push_status := 1
	status, msg_b := common.HttpBody(push_url, "POST", str, push_header)
	msg := string(msg_b)
	p_num := p_list.Push_num + 1
	common.LogsWithFileName(log_path, "push_", "ordernumber->"+p_list.Order_number+"\nmsg->"+msg+" \n"+str+"\npush_url->"+push_url)
	if status == 200 && msg == "success" {
		push_status = 3
		p_data["push_status"] = push_status
		p_data["push_num"] = p_num
		model.UpdatesCashList(p_list, p_data)
		return
	}

	push_status = 9

	//继续推送,连续10秒,推送12次，如果都失败了，放弃任务，并记录
	for i := 1; i < 10; i++ {
		p_num = p_num + 1
		time.Sleep(30 * time.Second)
		status, msg_b = common.HttpBody(push_url, "POST", str, push_header)
		common.LogsWithFileName(log_path, "push_", "ordernumber->"+p_list.Order_number+"\ntimes->"+strconv.Itoa(i)+"\nrequest->"+str+"\nmsg->"+msg+"\npush_url->"+push_url)
		msg = string(msg_b)
		if status == 200 && msg == "success" {
			//成功后跳出循环
			push_status = 3
			break
		}
	}
	p_data["push_status"] = push_status
	p_data["push_num"] = p_num
	model.UpdatesCashList(p_list, p_data)
}
