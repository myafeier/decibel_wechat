package wechat

import "github.com/go-xorm/xorm"

type WxUserService struct {
	session *xorm.Session
}

func NewWxUserService(session *xorm.Session) *WxUserService {
	return &WxUserService{session: session}
}

func (s *WxUserService) InsertOrUpdateUser(user *WxUserEntity) error {
	exist := new(WxUserEntity)
	if has, err := s.session.Where("open_id=?", user.OpenId).Get(exist); err != nil {
		return err
	} else if has {
		if _, err := s.session.ID(exist.Id).Update(user); err != nil {
			return err
		}
	} else {
		if _, err := s.session.Insert(user); err != nil {
			return err
		}
	}
	return nil
}
