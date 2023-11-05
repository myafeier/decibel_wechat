package wechat

import "xorm.io/xorm"

type WxUserService struct {
	session *xorm.Session
}

func NewWxUserService(session *xorm.Session) *WxUserService {
	return &WxUserService{session: session}
}
func (s *WxUserService) GetOne(id int64) (user *WxUserEntity, has bool, err error) {
	user = new(WxUserEntity)
	has, err = s.session.ID(id).Get(user)
	return
}

func (s *WxUserService) InsertOrUpdateUser(user *WxUserEntity) error {
	exist := new(WxUserEntity)
	if has, err := s.session.Where("open_id=?", user.OpenId).Get(exist); err != nil {
		return err
	} else if has {
		user.Id = exist.Id
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

func (s *WxUserService) Update(user *WxUserEntity, cols ...string) error {
	_, err := s.session.Where("id=?", user.Id).Cols(cols...).Update(user)
	if err != nil {
		return err
	} else {
		return nil
	}
}
