package sebar

import (
	"errors"
	"github.com/eaciit/toolkit"
	"time"
)

type Session struct {
	UserID, Secret, URL string
}

type usersession struct {
	Session
	Created time.Time
	Epired  time.Time
}

func (s *Session) ID() string {
	return s.UserID + "||" + s.Secret
}

func Login(url, userid, password string) (*Session, error) {
	//return nil, errors.New("Login is not yet activated")
	if userid != "arief" || password != "darmawan" {
		return nil, errors.New("Authorisation failed")
	}

	loginUrl := MakeUrl(url, "master", "login")
	result, e := toolkit.CallResult(loginUrl, "POST", toolkit.M{}.Set("userid", userid).Set("password", password).ToBytes("json", nil))
	if e != nil {
		return nil, e
	}
	sess := &Session{UserID: userid, Secret: result.Data.(string)}
	sess.URL = url
	return sess, nil
}

func (s *Session) Logout() {
	url := MakeUrl(s.URL, "master", "logout")
	toolkit.CallResult(url, "POST", toolkit.ToBytes(s, "json"))
}

func (s *Session) CallResult(method string, data []byte) *toolkit.Result {
	weburl := s.URL + "/master/" + method
	result, e := toolkit.CallResult(weburl, "POST", data)
	if e != nil {
		result.Data = nil
		return result.SetError(e)
	}
	return result
}

func (s *Session) Write(key string, data interface{}, mode int) error {
	r := s.CallResult("write", toolkit.ToBytes(struct {
		UserID, Secret, Key string
		Data                interface{}
		Mode                int
	}{s.UserID, s.Secret, key, data, mode}, "json"))
	if r.Status == toolkit.Status_NOK {
		return errors.New(r.Message)
	}
	return nil
}

func (s *Session) Read(key string) ([]byte, error) {
	r := s.CallResult("read", toolkit.ToBytes(struct {
		UserID, Secret, Key string
	}{s.UserID, s.Secret, key}, "json"))
	if r.Status == toolkit.Status_NOK {
		return nil, errors.New(r.Message)
	}
	return r.Data.([]byte), nil
}
