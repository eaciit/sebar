package sebar_test

import (
	"github.com/eaciit/sebar/v1"
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

/*
func TestNode(t *testing.T) {
	for i := 0; i < 3; i++ {
		node := sebar.NewServer(sebar.RoleWorker, "http://localhost:"+strconv.Itoa(3500+i)).(*sebar.Node)
		//node.SetMaster(masterUrl, secretMaster)
		e := node.Star()
		if e == nil {
			nodes = append(nodes, node)
		}
	}
}

func TestLogin(t *testing.T) {
	var e error
	session, e = sebar.Login(masterUrl, userID, password)
	if e != nil {
		t.Errorf(e.Error())
		return
	}
	t.Logf("User %s logged in. Secret: %s", userID, session.Secret)
}

func TestWrite(t *testing.T) {
	session.Write("Public:Sales:Orders",
		struct {
			OrderID, Customer string
			OrderDate         time.Time
			Amount            float64
		}{"ORD01", "Shell USA", time.Now(), 5000},
		sebar.WriteMemory+sebar.WriteStorage)
}
*/

func TestClose(t *testing.T) {
	var e error
	if coordinator != nil {
		e = coordinator.Stop()
		if e != nil {
			t.Error(e)
			return
		}
	}
}
