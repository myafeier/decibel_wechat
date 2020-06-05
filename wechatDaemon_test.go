package wechat

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/myafeier/log"
	"testing"
	"time"
	goxorm "xorm.io/xorm"
)

func init() {
	log.SetLogLevel(log.DEBUG)
	LocalDb, err := goxorm.NewEngine("mysql", "test:test@tcp(localhost:3306)/test?charset=utf8mb4")
	if err != nil {
		log.Error(err.Error())
		return
	}
	LocalDb.SetMaxIdleConns(10)
	LocalDb.SetMaxOpenConns(100)
	LocalDb.SetConnMaxLifetime(100 * time.Second)
	LocalDb.ShowSQL(true)

	config := &Config{
		&WeChatMicroAppConfig{
			OriginId:     "",
			AppId:        "wx884199c1f98151f3",
			Secret:       "8119bb3228e1b80705f94971551301f7",
			Base64AESKey: "",
		},
		&WePayConfig{
			MerchantId:     "1577820171",
			MerchantSecret: "Uydh7635ysgh89ikojs63526352fhjdk",
			NotifyUrl:      "https://pay.u1200.com",
		},
		&WePayVendorConfig{},
	}
	InitWeChatDaemon(true, true, true, true, nil, LocalDb, config)
}

func TestGetMicroAppSession(t *testing.T) {

	session, err := GetMicroAppSession("test")
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", *session)
}

func TestWePay_UnifiedOrderJSAPI(t *testing.T) {

	payData, err := Daemon.NewPay().UnifiedOrder_JSAPI("test", "testorder", "testinfo", "127.0.0.1", 1)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%s", payData)
}

func TestWePay_UnifiedOrderNative(t *testing.T) {
	payUrl, err := Daemon.NewPay().UnifiedOrder_Native(OrderSourceOfMerchantRecharge, 1, "info", 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("payUrl: %s \n", payUrl)
}

func TestSendSubscribeMsg(t *testing.T) {
	msg := &WxNotifyEntity{}
	msg.OpenId = "oSWz25VYk9p5iYGdybEoC8BfHo2k"
	msg.TemplateId = "fTlB0J83PtcS5DsfpIjfIAj0Oir6sfV_KOOGQ1UyZME"
	msg.Page = "/pages/index/index"
	msg.Data = map[string]DataItem{
		"no":      {Value: "test", Color: "#ff00000"},
		"addtion": {Value: "test", Color: "#ff00000"},
		"remark":  {Value: "test", Color: "#ff00000"},
	}

	err := Daemon.NotifyDaemon.Store(msg)
	if err != nil {
		t.Error(err.Error())
		return
	}
	time.Sleep(2 * time.Minute)

}
