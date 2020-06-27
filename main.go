package main

import (
	"TFService/router"
	"TFService/thread"
	"encoding/json"

	"github.com/zhibingzhou/go_public/common"
)

//
func main() {
	conf_byte, err := common.ReadFile("./conf/conf.json")

	if err != nil {
		panic(err)
	}
	var json_conf map[string]string
	//解析json格式
	err = json.Unmarshal(conf_byte, &json_conf)
	if err != nil {
		panic(err)
	}
	go thread.ThreadPush()
	go thread.ThreadPushCash()
	//查询是否有不推送的下发
	go thread.ThreadNotPush()
	//代理报表
	go thread.ThreadReport()
	_ = router.Router.Run(json_conf["port"])

}
