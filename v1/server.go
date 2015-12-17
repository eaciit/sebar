package sebar

import (
	//"errors"
	"net/url"
)

type SebarServer struct {
	Protocol, Address, Secret string
	url                       string
}

func (s *SebarServer) Start() error {
	return nil
}

func (s *SebarServer) Stop() error {
	return nil
}

func (s *SebarServer) SetURL(rawurl string) *SebarServer {
	u, e := url.Parse(rawurl)
	if e != nil {
		return s
	}
	s.Protocol = u.Scheme
	s.Address = u.Host
	return s
}
