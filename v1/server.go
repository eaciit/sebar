package sebar

import (
	"errors"
	//"net/url"
	"github.com/eaciit/appserver"
)

type IServer interface {
}

type SebarServer struct {
	appserver.AppServer
	Protocol, Address, Secret string

	//_rpcAddress string
	//_urlAddress string
}

func (s *SebarServer) Start() error {
	var e error

	e := s.Register(s)
	if e != nil {
		return errors.New("Unable to register RPC: " + e.Error())
	}

	e := s.AppServer.Start(false)
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
