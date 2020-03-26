package wechat

import "time"

type WxUserEntity struct {
	Id        int64     `json:"id"`
	MemberId  int64     `json:"member_id" binding:"-" xorm:"default 0 index"`
	OpenId    string    `json:"open_id" binding:"-" xorm:"varchar(100) default '' index"`
	UnionId   string    `json:"union_id" binding:"-" xorm:"varchar(100) default '' index"`
	NickName  string    `json:"nickName" xorm:"varchar(100) defult ''"`
	AvatarUrl string    `json:"avatarUrl" xorm:"varchar(200) default ''"`
	Gender    int       `json:"gender" xorm:"tinyint(2) default 0"` // 1男，2女， 0未知
	Country   string    `json:"country" xorm:"varchar(100) defult ''"`
	Province  string    `json:"province" xorm:"varchar(100) defult ''"`
	City      string    `json:"city" xorm:"varchar(100) defult ''"`
	Language  string    `json:"language" xorm:"varchar(200) default ''"`
	RefereeId int64     `json:"referee_id" binding:"-" xorm:"default 0 index"`
	Created   time.Time `json:"created" xorm:"created"`
	Update    time.Time `json:"update" xorm:"updated"`
}

func (s *WxUserEntity) TableName() string {
	return "wx_user"
}
