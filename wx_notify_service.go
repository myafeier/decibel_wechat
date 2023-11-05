package wechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/myafeier/log"
	"gopkg.in/chanxuehong/wechat.v2/mp/core"
	"xorm.io/xorm"
)

const (
	SubscribeMsgSendUrl = "https://api.weixin.qq.com/cgi-bin/message/subscribe/send?access_token=%s"
)

func NewWxNotifyService(ats core.AccessTokenServer, cfg *WeChatMicroAppConfig, session *xorm.Session) *WxNotifyService {
	return &WxNotifyService{
		AccessTokenServer: ats,
		Config:            cfg,
		session:           session,
		notifyChannel:     make(chan *WxNotifyEntity),
		logMutex:          sync.Mutex{},
	}
}

// 微信通知服务
type WxNotifyService struct {
	AccessTokenServer core.AccessTokenServer
	Config            *WeChatMicroAppConfig
	session           *xorm.Session
	notifyChannel     chan *WxNotifyEntity //通知通道
	logMutex          sync.Mutex
}

// 需要作为后台服务运行
func (s *WxNotifyService) Run() {
	ticket := time.NewTicker(10 * time.Second)
	for {
		select {
		case e := <-s.notifyChannel:
			log.Debug("found mission: %+v", *e)
			go s.send(e)
		case <-ticket.C:
			//go s.dispatchMission()
		}
	}
}

func (s *WxNotifyService) dispatchMission() {
	rows, err := s.session.Where("state=?", NotifyStateOfNotYet).Rows(&WxNotifyEntity{})
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer rows.Close()
	for rows.Next() {
		var row WxNotifyEntity
		err = rows.Scan(&row)
		if err != nil {
			log.Error(err.Error())
			return
		}
		log.Debug("notify: %+v", &row)
		if row.SendAtTime.IsZero() || row.SendAtTime.Before(time.Now()) {
			s.notifyChannel <- &row
		}
	}
	return
}

// 发送通知
func (s *WxNotifyService) send(e *WxNotifyEntity) {
	log.Debug("sending...%d", e.Id)
	msg := new(SubscribeMsg)
	msg.Data = e.Data
	msg.Lang = "zh_CN"
	msg.MiniprogramState = "formal"
	msg.Page = e.Page
	msg.TemplateId = e.TemplateId
	err := s.sendSubscribeMsg(e.OpenId, msg)
	if err != nil {
		e.State = NotifyStateOfFail
		log.Error("send WxNotify To %s error:%s", e.OpenId, err.Error())
	} else {
		e.State = NotifyStateOfSuccess
	}
	err = s.updateLog(e)
	if err != nil {
		log.Error(err.Error())
		return
	}

}
func (s *WxNotifyService) updateLog(e *WxNotifyEntity) (err error) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()
	_, err = s.session.ID(e.Id).Cols("state").Update(e)
	return
}

// 存储通知
func (s *WxNotifyService) Store(e *WxNotifyEntity) (err error) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()
	if e.State == 0 {
		e.State = NotifyStateOfNotYet
	}
	_, err = s.session.Insert(e)
	if err != nil {
		log.Error(err.Error())
	}
	return
}

// 微信小程序通知消息
type SubscribeMsg struct {
	TemplateId       string              `json:"template_id"`                 //小程序模板ID
	Data             map[string]DataItem `json:"data,omitempty"`              //小程序模板数据
	Page             string              `json:"page,omitempty"`              //小程序页面路径
	MiniprogramState string              `json:"miniprogram_state,omitempty"` //跳转小程序类型：developer为开发版；trial为体验版；formal为正式版；默认为正式版
	Lang             string              `json:"lang,omitempty"`              //进入小程序查看”的语言类型，支持zh_CN(简体中文)、en_US(英文)、zh_HK(繁体中文)、zh_TW(繁体中文)，默认为zh_CN
}

func (s *WxNotifyService) sendSubscribeMsg(toOpenId string, msg *SubscribeMsg) error {
	data := &struct {
		Touser string `json:"touser"` //用户openid，可以是小程序的openid，也可以是mp_template_msg.appid对应的公众号的openid
		SubscribeMsg
	}{
		SubscribeMsg: *msg,
		Touser:       toOpenId,
	}

	token, err := s.AccessTokenServer.Token()
	if err != nil {
		log.Error(err.Error())
		return err
	}
	dataByte, err := json.Marshal(data)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Debug("json: %s", dataByte)
	buf := bytes.NewBuffer(dataByte)
	resp, err := http.DefaultClient.Post(fmt.Sprintf(SubscribeMsgSendUrl, token), "application/json", buf)
	if err != nil {
		log.Error(err.Error())
		return err

	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err.Error())
		return err

	}
	var ce core.Error
	err = json.Unmarshal(respData, &ce)
	if err != nil {
		log.Error(err.Error())
		return err

	}
	if ce.ErrCode != core.ErrCodeOK {
		err = fmt.Errorf("code:%d msg:%s", ce.ErrCode, ce.ErrMsg)
		return err
	}
	return nil
}
