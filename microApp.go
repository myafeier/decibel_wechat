package wechat

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

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
