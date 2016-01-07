package sebar

import (
	"errors"
	"github.com/eaciit/toolkit"
)

type StorageMediaInfo struct {
	Active        bool
	AllocatedSize float64
	UnitSize      float64
	Usage         float64
}

type Storage struct {
	SebarServer
	MemoryStorage   *StorageMediaInfo
	PhysicalStorage *StorageMediaInfo
}

func (s *Storage) StopServer(in toolkit.M) *toolkit.Result {
	r := toolkit.NewResult()
	s.SebarServer.Stop()
	return r
}

func (s *Storage) Write(in toolkit.M) *toolkit.Result {
	r := toolkit.NewResult()
	r.SetErrorTxt("Storage.Write is not yet implemented")
	return r
}

func (s *Storage) Start() error {
	errorPrefix := "Unable to start storage server " + s.Address + " : "

	//--- validate
	if s.Coordinator == "" {
		return errors.New(errorPrefix + "No coordinator has been specified")
	}

	s.SebarServer.Server.RegisterRPCFunctions(s)
	/*
		for _, v := range s.Functions() {
			v.AuthType = ""
		}
	*/

	s.AddUser(s.CoordinatorUserID, s.CoordinatorSecret)
	e := s.SebarServer.Start()
	if e != nil {
		return errors.New(errorPrefix + e.Error())
	}

	return nil
}

func (s *Storage) Stop() error {
	s.StopServer(nil)
	return nil
}
