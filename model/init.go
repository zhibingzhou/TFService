package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/zhibingzhou/go_public/common"
	"github.com/zhibingzhou/go_public/redis"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

//Redis缓存的时间
var redis_time = 600
var redis_short_time = 60
var redis_min_time = 10
var redis_data_time = 3600
var redis_long_time = 86400
var redis_max_time = 604800

//配置数据的缓存key集合
var Conf_Redis_Key = "Conf_Redis_Key"

//数据缓存的key集合
var Data_Redis_Key = "Data_Redis_Key"

type Gorm struct {
	DB *gorm.DB
}

type CountTotal struct {
	Total float64
	Num   int
}

func init() {
	ReloadConf("")
}

var gdb Gorm

func ReloadConf(file_name string) {
	if file_name == "" {
		file_name = "./conf/database.json"
	}
	conf_byte, err := common.ReadFile(file_name)
	if err != nil {
		panic(err)
	}
	var json_conf map[string]string
	//解析json格式
	err = json.Unmarshal(conf_byte, &json_conf)
	if err != nil {
		panic(err)
	}
	life_time, _ := time.ParseDuration(json_conf["life_time"])
	max_open, _ := strconv.Atoi(json_conf["max_open"])
	if max_open < 1 {
		max_open = 40
	}
	max_idle, _ := strconv.Atoi(json_conf["max_idle"])
	if max_idle < 1 {
		max_idle = 10
	}
	conn_str := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?charset=utf8", json_conf["user"], json_conf["pwd"], json_conf["network"], json_conf["host"], json_conf["port"], json_conf["db_name"])
	db, err := gorm.Open("mysql", conn_str)
	if err != nil {
		fmt.Println("conn_str->", conn_str)
		panic(err)
	}
	//最大生命周期
	db.DB().SetConnMaxLifetime(life_time)
	//连接池的最大打开连接数
	db.DB().SetMaxOpenConns(max_open)
	//连接池的最大空闲连接数
	db.DB().SetMaxIdleConns(max_idle)
	db.SingularTable(true)
	//启用Logger，显示详细日志
	//db.LogMode(true)

	// 禁用日志记录器，不显示任何日志
	//db.LogMode(false)
	gdb = Gorm{DB: db}
}

func GetKey(length int) string {
	sec := strconv.FormatInt(time.Now().Unix(), 10)
	red_key := "model_get_key:" + sec
	rand_len := length
	ex_time := 1
	pre_id := ""

	if length > 10 {
		rand_len = length - 10
		pre_id = sec
	}
	rand_str := ""
	for i := 0; i < 50; i++ {
		rand_str = common.Random("smallnumber", rand_len)
		red_res := redis.RediGo.Sadd(red_key, rand_str, ex_time)
		if red_res > 0 {
			break
		}
	}

	key_str := pre_id + rand_str
	return key_str
}

func Query(sql_str string) error {
	res := gdb.DB.Exec(sql_str)
	return res.Error
}

func Trans(sql_arr []string) error {
	tx := gdb.DB.Begin()
	// 注意，一旦你在一个事务中，使用tx作为数据库句柄
	var err error
	for _, sql := range sql_arr {
		//更新订单状态
		if err = tx.Exec(sql).Error; err != nil {
			tx.Rollback()
			break
		}
	}

	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}
