package decibel_wechat

import "time"

type WePayStat int8

const (
	WePayStatToPay      WePayStat = 1  //待支付
	WePayStatPaySuccess WePayStat = 2  //支付成功
	WePayStatPayFail    WePayStat = -1 //支付失败
)

//微信支付订单
type WePayOrder struct {
	Id             int       `json:"id"`
	OutTradeNo     string    `json:"out_trade_no" xorm:"varchar(200) default '' index"`
	Openid         string    `json:"openid" xorm:"varchar(200) default '' index"`
	NonceStr       string    `json:"nonce_str" xorm:"varchar(200) default ''"`
	Body           string    `json:"body" xorm:"varchar(200) default ''"`
	TotalFee       int       `json:"total_fee" xorm:"int(11) default 0"`
	SpbillCreateIp string    `json:"spbill_create_ip" xorm:"varchar(200) default ''"`
	NotifyUrl      string    `json:"notify_url" xorm:"varchar(200) default ''"`
	TradeType      string    `json:"trade_type" xorm:"varchar(200) default ''"`
	SignType       string    `json:"sign_type" xorm:"varchar(200) default ''"`
	ReturnCode     string    `json:"return_code" xorm:"varchar(200) default ''"`
	PrepayId       string    `json:"prepay_id" xorm:"varchar(200) default ''"`
	TransactionId  string    `json:"transaction_id"  xorm:"varchar(200) default ''"`
	Stat           WePayStat `json:"stat" xorm:"tinyint(2) default 0 index"`
	WatcherResult  string    `json:"watcher_result" xorm:"varchar(200) default 'unknown'"`
	Created        time.Time `json:"created" xorm:"created"`
	Updated        time.Time `json:"updated" xorm:""`
}
