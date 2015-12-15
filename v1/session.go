package sebar

import (
	"errors"
)

type Session struct {
	UserID, Secret string
}

func Login(url, userid, password string) (*Session, error) {
	return nil, errors.New("Login is not yet activated")
}

func Logout(userid, secret string) {
}
