package wechat

import (
	"github.com/go-xorm/xorm"
	"gopkg.in/chanxuehong/wechat.v2/mp/core"
)

var Daemon *WeChatDaemon

type WeChatDaemon struct {
	Config              *Config
	Logger              ILogger
	Engine              *xorm.Engine
	AccessTokenServer   *core.DefaultAccessTokenServer
	CallbackServer      *core.Server
	DefaultWepayService *WePayService
}

func InitWeChatDaemon(initAccessTokenServer, initCallBackServer, initDefaultWepayService bool, logger ILogger, dbEngine *xorm.Engine, config *Config) {

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
	err := Daemon.Init(initAccessTokenServer, initCallBackServer, initDefaultWepayService)
	if err != nil {
		panic(err)
	}
}

func (self *WeChatDaemon) Init(initAccessTokenServer bool, initCallBackServer, initDefaultWepayService bool) (err error) {
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
		self.CallbackServer = core.NewServer(self.Config.OriginId, self.Config.WeChatMicroAppConfig.AppId, self.Config.Secret, self.Config.Base64AESKey, mux, nil)
	}

	if initDefaultWepayService {
		self.DefaultWepayService = new(WePayService)
		self.DefaultWepayService.db = self.Engine.NewSession()
		self.DefaultWepayService.AppId = self.Config.WeChatMicroAppConfig.AppId
		self.DefaultWepayService.logger = self.Logger
		self.DefaultWepayService.config = self.Config.WePayConfig
		if self.Config.WePayVendorConfig != nil {
			self.DefaultWepayService.vendor = self.Config.WePayVendorConfig
		}
	}
	return
}

func InitDb(db *xorm.Session) (err error) {
	var tables = []interface{}{&WePaymentEntity{}, &WxUserEntity{}}

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
func (self *WeChatDaemon) NewPay() *WePayService {
	pay := new(WePayService)
	pay.db = self.Engine.NewSession()
	pay.AppId = self.Config.WeChatMicroAppConfig.AppId
	pay.logger = self.Logger
	pay.config = self.Config.WePayConfig
	pay.vendor = self.Config.WePayVendorConfig
	return pay
}
