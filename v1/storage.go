package sebar

import (
	"errors"
	"github.com/eaciit/toolkit"
)

type StorageTypeEnum string

const (
	StorageTypeMemory StorageTypeEnum = "MEM"
	StorageTypeDisk   StorageTypeEnum = "DSK"
)

var defaultMemoryStorageSize, defaultDiskStorageSize float64

func SetDefaultStorageSize(size float64, destination StorageTypeEnum) {
	if destination == StorageTypeMemory {
		defaultMemoryStorageSize = size
	} else if destination == StorageTypeDisk {
		defaultDiskStorageSize = size
	}
}

func DefaultStorageSize(storageType StorageTypeEnum) float64 {
	ret := float64(1.5 * 1024 * 1024 * 1024)
	if storageType == StorageTypeMemory {
		if defaultMemoryStorageSize == 0 {
			defaultMemoryStorageSize = ret
		}
		return defaultMemoryStorageSize
	} else if storageType == StorageTypeDisk {
		if defaultDiskStorageSize == 0 {
			defaultDiskStorageSize = ret
		}
		return defaultDiskStorageSize
	}
	return 0
}

type StorageInfo struct {
	Active        bool
	AllocatedSize float64
	UnitSize      float64
	Usage         float64
}

type StorageMedia struct {
	StorageInfo
	datas map[string]*DataPoint
}

func (m *StorageMedia) Available() float64 {
	return m.AllocatedSize - m.Usage
}

func NewStorageMedia(size float64) *StorageMedia {
	sm := new(StorageMedia)
	sm.Active = true
	sm.AllocatedSize = size
	sm.datas = map[string]*DataPoint{}
	return sm
}

func (sm *StorageMedia) Load(path string) error {
	return nil
}

type Storage struct {
	SebarServer
	MemoryStorage *StorageMedia
	DiskStorage   *StorageMedia
}

func (s *Storage) StopServer(in toolkit.M) *toolkit.Result {
	r := toolkit.NewResult()
	s.SebarServer.Stop()
	return r
}

func (s *Storage) StorageStatus(in toolkit.M) *toolkit.Result {
	r := toolkit.NewResult()
	r.SetBytes(struct {
		Memory StorageInfo
		Disk   StorageInfo
	}{s.MemoryStorage.StorageInfo, s.DiskStorage.StorageInfo}, "")
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
		- DONE - initialize storage info (active, size)
		- load storage data from physical folder
		- update the coordinator metadata with metadata from this server
		##TODO END
	*/
	s.MemoryStorage = NewStorageMedia(DefaultStorageSize(StorageTypeMemory))
	s.DiskStorage = NewStorageMedia(DefaultStorageSize(StorageTypeDisk))

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
