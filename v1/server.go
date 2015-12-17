package sebar

import (
	"errors"
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

func (s *SebarServer) SetURL(rawurl string) error {
	u, e := url.Parse(rawurl)
	if e != nil {
		return errors.New("Unable to parse URL: " + rawurl)
	}
	s.Protocol = u.Scheme
	s.Address = u.Host
	return nil
}
