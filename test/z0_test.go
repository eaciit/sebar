package sebar_test

import (
	"github.com/eaciit/sebar/client.v1"
	"github.com/eaciit/sebar/v1"
	"github.com/eaciit/toolkit"
	//"strconv"
	"fmt"
	"testing"
	//"time"
)

var (
	secretMaster = "master"
	userID       = "admin"
	password     = "admin"
	userSecret   = ""
	masterUrl    = "localhost:6789"

	coordinator *sebar.Coordinator
	nodes       []*sebar.Node
	session     *sebar.Session

	client *sebarclient.Client
)

func TestMaster(t *testing.T) {
	coordinator = sebar.NewServer(sebar.RoleCoordinator, masterUrl).(*sebar.Coordinator)
	coordinator.AddUser(userID, password)
	coordinator.AllowMultiLogin = true
	e := coordinator.Start()
	if e != nil {
		t.Error(e)
	}
}

func TestStorageNode(t *testing.T) {
	for i := 0; i < 5; i++ {
		var storage *sebar.Storage
		storage = sebar.NewServer(sebar.RoleStorage, fmt.Sprintf("localhost:%d", 9601+i)).(*sebar.Storage)
		storage.Coordinator = coordinator.Address
		storage.CoordinatorUserID = userID
		storage.CoordinatorSecret = password
		e := storage.Start()
		if e != nil {
			t.Error(e)
		}
	}
}

func TestStorageSave(t *testing.T) {
	client = new(sebarclient.Client)
	client.Host = coordinator.Address
	client.UserID = userID
	client.Secret = password
	e := client.Connect()
	if e != nil {
		t.Error(e.Error())
	}

	toolkit.Printf("Writing Data:\n")
	for i := 0; i < 200; i++ {
		dataku := toolkit.RandInt(1000)
		toolkit.Printf("%d ", dataku)
	}
	toolkit.Println("")
}

func TestClose(t *testing.T) {
	var e error

	if client != nil {
		client.Close()
	}

	if coordinator != nil {
		e = coordinator.Stop()
		if e != nil {
			t.Error(e)
			return
		}
	}
}
