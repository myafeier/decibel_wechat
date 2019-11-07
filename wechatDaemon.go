package decibel_wechat

import (
	"github.com/go-xorm/xorm"
	"gopkg.in/chanxuehong/wechat.v2/mp/core"
)

var Daemon *WeChatDaemon

type WeChatDaemon struct {
	Config            *Config
	Logger            ILogger
	Engine            *xorm.Engine
	AccessTokenServer *core.DefaultAccessTokenServer
	CallbackServer    *core.Server
}

func InitWeChatDaemon(initAccessTokenServer, initCallBackServer bool, logger ILogger, dbEngine *xorm.Engine, config *Config) {

	if Daemon == nil {
		Daemon = new(WeChatDaemon)
	}
	Daemon.Engine = dbEngine
	Daemon.Config = config
	if logger == nil {
		Daemon.Logger = NewDefaultLogger()
	} else {
		Daemon.Logger = logger
	}
	err := Daemon.Init(initAccessTokenServer, initCallBackServer)
	if err != nil {
		panic(err)
	}
}

func (self *WeChatDaemon) Init(initAccessTokenServer bool, initCallBackServer bool) (err error) {
	err = InitDb(self.Engine.NewSession())
	if err != nil {
		self.Logger.Error(err)
		return
	}

	if initAccessTokenServer {
		self.AccessTokenServer = core.NewDefaultAccessTokenServer(self.Config.WeChatMicroAppConfig.AppId, self.Config.WeChatMicroAppConfig.Secret, nil)
		_, err = self.AccessTokenServer.Token()
		if err != nil {
			self.Logger.Error(err)
			return
		}
	}

	if initCallBackServer {
		self.CallbackServer = core.NewServer(self.Config.OriginId, self.Config.AppId, self.Config.Secret, self.Config.Base64AESKey, mux, nil)
	}

	return
}

func InitDb(db *xorm.Session) (err error) {
	var tables = []interface{}{WePayment{}}

	var isExist bool
	for _, v := range tables {
		isExist, err = db.IsTableExist(v)
		if err != nil {
			return
		}
		if !isExist {
			err = db.CreateTable(v)
			if err != nil {
				return
			}
			err = db.CreateIndexes(v)
			if err != nil {
				return
			}
		} else {
			err = db.Sync2(v)
			if err != nil {
				return
			}
		}
	}
	return
}

//得到支付接口
func (self *WeChatDaemon) NewPay() *WePay {
	pay := new(WePay)
	pay.db = self.Engine.NewSession()
	pay.AppId = self.Config.WeChatMicroAppConfig.AppId
	pay.logger = self.Logger
	pay.config = self.Config.WePayConfig
	return pay
}

//直接调用Pay
func NewPay(appId string,config *WePayConfig,db *xorm.Session,logger ILogger,)*WePay{
	pay := new(WePay)
	pay.db = db
	pay.AppId = appId
	pay.logger = logger
	pay.config = config
	if logger==nil{
		pay.logger=NewDefaultLogger()
	}
	return pay
}
