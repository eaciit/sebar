package sebar_test

import (
	"github.com/eaciit/sebar/v1"
	"strconv"
	"testing"
	"time"
)

var (
	secretMaster = "master"
	userID       = "admin"
	password     = "admin"
	userSecret   = ""
	masterUrl    = "localhost:12345"

	master  *sebar.Master
	nodes   []*sebar.Node
	session *sebar.Session
)

func TestMaster(t *testing.T) {
	master = sebar.NewServer(sebar.RoleMaster, masterUrl).(*sebar.Master)
	master.Secret = secretMaster
	e := master.Start()
	if e != nil {
		t.Error(e)
	}
}

func TestNode(t *testing.T) {
	for i := 0; i < 3; i++ {
		node := sebar.NewServer(sebar.RoleNode, "http://localhost:"+strconv.Itoa(3500+i)).(*sebar.Node)
		node.SetMaster(masterUrl, secretMaster)
		e := node.Start()
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

func TestClose(t *testing.T) {
	if master != nil {
		master.Stop()
	}
}
