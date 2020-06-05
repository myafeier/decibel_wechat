package wechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/myafeier/log"
	"gopkg.in/chanxuehong/wechat.v2/mp/core"
)

const (
	SubscribeMsgSendUrl = "https://api.weixin.qq.com/cgi-bin/message/subscribe/send?access_token=%s"
)

// 微信通知服务
type WxNotifyService struct {
	AccessTokenServer *core.DefaultAccessTokenServer
	Config            *WeChatMicroAppConfig
}

//微信小程序通知消息
type SubscribeMsg struct {
	TemplateId       string      `json:"template_id"`                 //小程序模板ID
	Data             interface{} `json:"data,omitempty"`              //小程序模板数据
	Page             string      `json:"page,omitempty"`              //小程序页面路径
	MiniprogramState string      `json:"miniprogram_state,omitempty"` //跳转小程序类型：developer为开发版；trial为体验版；formal为正式版；默认为正式版
	Lang             string      `json:"lang,omitempty"`              //进入小程序查看”的语言类型，支持zh_CN(简体中文)、en_US(英文)、zh_HK(繁体中文)、zh_TW(繁体中文)，默认为zh_CN
}

func (s *WxNotifyService) SendSubscribeMsg(toOpenId string, msg *SubscribeMsg) error {
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
