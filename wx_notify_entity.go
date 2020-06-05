package wechat

import (
	"time"
)

type NotifyState int8

const (
	NotifyStateOfNotYet  NotifyState = 1  //尚未发送
	NotifyStateOfSuccess NotifyState = 2  //发送成功
	NotifyStateOfFail    NotifyState = -1 //发送失败
)

type WxNotifyEntity struct {
	Id         int64               `json:"id"`
	OpenId     string              `json:"open_id" xorm:"varchar(100) default ''"`     //发送对象id
	TemplateId string              `json:"template_id" xorm:"varchar(100) default ''"` //发送模版ID
	SendAtTime time.Time           `json:"send_at_time" xorm:""`                       //计划发送时间,0时为即时发送，否则为定时发送
	State      NotifyState         `json:"state" xorm:"tinyint(2) default 0 index"`    //发送状态
	Data       map[string]DataItem `json:"data" xorm:"json"`                           //数据
	Page       string              `json:"page" xorm:"varchar(100) default ''"`        //跳转页面
	Created    time.Time           `xorm:"created"`                                    //生成数据
	Updated    time.Time           `xorm:"updated"`                                    //发送时间
}

func (s *WxNotifyEntity) TableName() string {
	return "wx_notify"
}

type DataItem struct {
	Value string `json:"value"`
	Color string `json:"color,omitempty"`
}
