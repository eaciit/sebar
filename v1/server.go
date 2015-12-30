package sebar

import (
	"errors"
	"github.com/eaciit/appserver/v1"
)

type IServer interface {
	Start() error
	Stop() error
}

type SebarServer struct {
	appserver.Server
	Protocol, Address, Secret string

	//_rpcAddress string
	//_urlAddress string
}

func (s *SebarServer) Start() error {
	var e error
	secret := "sebar"

	e = s.Register(s)
	if e != nil {
		return errors.New("Unable to register RPC: " + e.Error())
	}

	s.Server.SetSecret(secret)
	e = s.Server.Start(s.Address)
	return nil
}

func (s *SebarServer) Stop() error {
	return nil
}

/*
func (s *SebarServer) SetURL(rawurl string) *SebarServer {
	u, e := url.Parse(rawurl)
	if e != nil {
		return s
	}
	s.Protocol = u.Scheme
	s.Address = u.Host
	return s
}
*/
