package controller

import (
	"encoding/json"
	"net/http"

	"github.com/zhibingzhou/go_public/common"
)

/*
* 定义一个json的返回值类型(首字母必须大写)
 */
type JsonOut struct {
	Status int
	Msg    string
	Data   map[string]interface{}
}

var http_status = http.StatusOK

var log_path string
var down_url string

func init() {
	conf_byte, err := common.ReadFile("./conf/conf.json")

	if err != nil {
		panic(err)
	}
	var json_conf map[string]string
	//解析json格式r
	err = json.Unmarshal(conf_byte, &json_conf)
	if err != nil {
		panic(err)
	}
	log_path = json_conf["log_path"]
	down_url = json_conf["down_url"]
}
