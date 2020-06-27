package model

type MerIp struct {
	Id       int
	Ip       string
	Mer_code string
}

type MerList struct {
	Id          int
	Code        string
	Status      int
	Domain      string
	Title       string
	Qq          string
	Skype       string
	Telegram    string
	Phone       string
	Email       string
	Private_key string
	Is_agent    int
	Agent_path  string
	Amount      float64
	Total_in    float64
	Total_out   float64
}

type PayBank struct {
	Id         int `orm:"pk"`
	Is_mobile  int
	Pay_code   string
	Class_code string
	Bank_code  string
	Bank_title string
	Jump_type  int
	Pay_bank   string
}

type SysBank struct {
	Id    int
	Code  string
	Title string
}

type PayConfig struct {
	Id            int `orm:"pk"`
	Status        int
	Pay_code      string
	Merchant_code string
	Api_conf      string
	Amount        float64
	Total_in      float64
	Total_out     float64
	Note          string
}

type PayList struct {
	Id           string `orm:"pk"`
	Status       int    //订单状态(1=处理中,3=完成,9=拒绝)
	Pay_code     string
	Pay_id       int
	Mer_code     string
	Push_status  int
	Push_num     int
	Amount       float64
	Real_amount  float64
	Create_time  string
	Pay_time     string
	Order_number string
	Pay_order    string
	Class_code   string
	Bank_code    string
	Push_url     string
	Note         string
	Is_mobile    int
	Rate         float64
	Agent_path   string
}

type PayChannel struct {
	Id         int `orm:"pk"`
	Code       string
	Title      string
	Fee_amount float64 //下发的手续费
	Fee_type   int     //手续费的收取类型
	Is_push    int     //是否有下发推送
}

type PayClass struct {
	Id    int
	Code  string
	Title string
}

type MerPay struct {
	Id       int
	Mer_code string
	Pay_code string
	Pay_id   int
	Status   int
}

type MerRate struct {
	Id           int
	Mer_code     string
	Pay_code     string
	Class_code   string
	Bank_code    string
	Rate         float64
	Limit_amount float64
	Day_amount   float64
}

type PayRate struct {
	Id           int
	Pay_code     string
	Class_code   string
	Bank_code    string
	Rate         float64
	Min_amount   float64
	Max_amount   float64
	Limit_amount float64
	Day_amount   float64
}

type CashList struct {
	Id           string `orm:"pk"`
	Status       int    //订单状态(1=处理中,3=完成,9=拒绝,-1=未扣款的废单)
	Pay_code     string
	Pay_id       int
	Mer_code     string
	Push_status  int
	Push_num     int
	Amount       float64
	Real_amount  float64
	Create_time  string
	Pay_time     string
	Order_number string
	Pay_order    string
	Bank_code    string
	Push_url     string
	Note         string
	Bank_title   string
	Branch       string
	Card_name    string
	Card_number  string
	Phone        string
	Fee_amount   float64
	Order_amount float64
	Agent_path   string
}

type OrderList struct {
	Id           string `orm:"pk"`
	Status       int    //订单状态(1=处理中,3=完成,9=拒绝,-1=未扣款的废单)
	Pay_code     string
	Pay_id       int
	Amount       float64
	Real_amount  float64
	Create_time  string
	Pay_time     string
	Order_number string
	Cash_id      string
	Pay_order    string
	Bank_code    string
	Note         string
	Bank_title   string
	Branch       string
	Card_name    string
	Card_number  string
	Phone        string
	Fee_amount   float64
	Order_amount float64
	Order_type   int //是否纯代付下发:1=代收下发,2=纯代付下发
}

type AmountList struct {
	Id            string `orm:"pk"`
	Amount_type   int    //账变类型(1=支付,2=下发,3=代理收入,4=下发失败返还额度,5=调整额度,6.支付的手续费,7.下发的手续费,8.代理佣金返还,9.上游支付下发,10.上游支付下发失败返回)
	Pay_code      string
	Pay_id        int
	Mer_code      string
	Amount        float64
	Before_amount float64
	After_amount  float64
	Create_time   string
	Order_number  string
	Note          string
	Agent_path    string
}

type AdminList struct {
	Id         int `orm:"pk"`
	Account    string
	Pwd        string
	Status     int //状态(1=正常,0=锁定,-1=删除)
	Mer_code   string
	Power_path string
	Secret     string
	Login_time string
	Login_ip   string
	Session_id string
}

type PowerList struct {
	Id         int `orm:"pk"`
	Path       string
	Code       string
	Power_type int //状态(1=正常,0=锁定,-1=删除)
	Url        string
	P_code     string
}

type AdminPower struct {
	Id         int `orm:"pk"`
	Account    string
	Power_code string
}

type NoteList struct {
	Id          string `orm:"pk"`
	Title       string
	Content     string
	Create_time string
	Status      int
}

type MerBank struct {
	Id          int `orm:"pk"`
	Mer_code    string
	Bank_code   string
	Bank_title  string
	Card_number string
	Card_name   string
	Bank_branch string
	Bank_phone  string
	Status      int
}

type MerReport struct {
	Id             string
	Mer_code       string
	Report_date    string
	Total_in       float64
	Total_in_rate  float64
	Total_out      float64
	Total_out_rate float64
	Is_agent       int
	Agent_path     string
}

type CashBank struct {
	Id         int
	Bank_code  string
	Bank_title string
	Cash_bank  string
	Pay_code   string
}
