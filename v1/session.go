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
