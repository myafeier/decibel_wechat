package decibel_wechat

import (
	goxorm	"github.com/go-xorm/xorm"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"testing"
	"time"
)

func init() {
	LocalDb, err := goxorm.NewEngine("mysql", "test:test@tcp(localhost:3306)/test?charset=utf8mb4")
	if err != nil {
		log.Println(err)
		return
	}
	LocalDb.SetMaxIdleConns(10)
	LocalDb.SetMaxOpenConns(100)
	LocalDb.SetConnMaxLifetime(100 * time.Second)
	LocalDb.ShowSQL(true)

	config:=&Config{
		&WeChatMicroAppConfig{
			OriginId:"",
			AppId:"",
			Secret:"",
			Base64AESKey:"",
		},
		&WePayConfig{
			MerchantId:"",
			MerchantSecret:"",
			NotifyUrl:"",
		},
	}
	InitWeChatDaemon(true,true,nil,LocalDb,config)
}

func TestGetMicroAppSession(t *testing.T) {

	session,err:=GetMicroAppSession("test")
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v",*session)
}

func TestWePay_UnifiedOrder(t *testing.T) {

	payData,err:=Daemon.NewPay().UnifiedOrder("test","testorder","testinfo","127.0.0.1",1)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%s",payData)
}