package wechat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"gopkg.in/chanxuehong/wechat.v2/mp/core"
)

type MicroappService struct {
	AccessTokenServer *core.DefaultAccessTokenServer
	WeChatMicroAppConfig
}

var MicroappDaemon *MicroappService

//单例模式
func NewMicroappDaemon(config *WeChatMicroAppConfig) *MicroappService {
	if MicroappDaemon != nil {
		return MicroappDaemon
	}
	MicroappDaemon = new(MicroappService)
	MicroappDaemon.WeChatMicroAppConfig = *config
	MicroappDaemon.AccessTokenServer = core.NewDefaultAccessTokenServer(config.AppId, config.Secret, nil)
	return MicroappDaemon
}

type WxappCodeRequest struct {
	Scene string `json:"scene"`
	Page  string `json:"page"`
	Width int    `json:"width"`
}

func (m *MicroappService) GenerateWxappCode(page, scene string, width int) (code []byte, err error) {
	token, err := m.AccessTokenServer.Token()
	if err != nil {
		return
	}
	req := new(WxappCodeRequest)
	req.Scene = scene
	req.Page = page
	req.Width = width
	data, err := json.Marshal(req)
	if err != nil {
		return
	}
	buf := bytes.NewBuffer(data)
	resp, err := http.Post("https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token="+token, "application/json;charset=utf-8", buf)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	code, err = ioutil.ReadAll(resp.Body)
	return

}

const SESSION_SERVER_URL = "https://api.weixin.qq.com/sns/jscode2session?"

type Session struct {
	OpenID     string `json:"openid"`            //用户唯一标识
	SessionKey string `json:"session_key"`       //会话密钥
	UnionID    string `json:"unionid,omitempty"` //用户在开放平台的唯一标识符。本字段在满足一定条件的情况下才返回。具体参看UnionID机制说明
}

func GetMicroAppSession(code string) (session *Session, err error) {
	url := SESSION_SERVER_URL + "appid=" + Daemon.Config.WeChatMicroAppConfig.AppId + "&secret=" + Daemon.Config.WeChatMicroAppConfig.Secret + "&js_code=" + code + "&grant_type=authorization_code"
	Daemon.Logger.Debug("GetUrl:", url)
	client := new(http.Client)
	response, err := client.Get(url)
	if err != nil {
		Daemon.Logger.Error(err)
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		Daemon.Logger.Error(err)
		return nil, err
	}

	if response.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("response code:%d,body:%s", response.StatusCode, body))
		return nil, err
	}

	session = new(Session)
	err = json.Unmarshal(body, session)
	if err != nil {
		return
	}
	return
}
