package sebar

import (
	"errors"
	"github.com/eaciit/toolkit"
)

type StorageMedia struct {
	Active        bool
	AllocatedSize float64
	UnitSize      float64
	Usage         float64

	datas map[string]*DataPoint
}

func NewStorageMedia() *StorageMedia {
	sm := new(StorageMedia)
	sm.Active = true
	sm.datas = map[string]*DataPoint{}
	return sm
}

func (sm *StorageMedia) Load(path string) error {
	return nil
}

type Storage struct {
	SebarServer
	MemoryStorage   *StorageMedia
	PhysicalStorage *StorageMedia
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

	/*
		Init Storage Data

		##TODO BEGIN
		- initialize storage info (active, size)
		- load storage data from physical folder
		- update the coordinator metadata with metadata from this server
		##TODO END
	*/
	s.MemoryStorage = NewStorageMedia()
	s.PhysicalStorage = NewStorageMedia()

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
